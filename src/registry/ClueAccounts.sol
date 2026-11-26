// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/**
 * @title ClueAccounts
 * @notice ClueRegistry component that records which account leaf hashes have been
 *         registered into the Clue Protocol and checks that a leaf, or a list of leaves, is
 *         valid: registered, and a member of the Sparse Merkle Tree (SMT) with a given root.
 * @dev Three parts that work together.
 *
 *      1. The leaf hash. The leaf of an account is the hash of its position and its
 *         account hash:
 *
 *             leaf = keccak256(abi.encodePacked(uint64 position, bytes32 accountHash))
 *
 *         so 8 bytes of position (big endian) followed by the 32 bytes of the account hash.
 *         The account hash is the hash of the account data and is computed off-chain; this
 *         contract takes it as given. Binding the position into the leaf means the same
 *         account hash at another position is a different leaf. See {_hashLeaf}.
 *
 *      2. The record. One map, `leaf => registered`, append-only and read-only to the
 *         outside:
 *           - A leaf is added with {_registerLeaf}, {_registerLeaves} or
 *             {_registerProvenLeaves}. All are internal, so only the registry can call them.
 *           - There is no function that removes or clears a leaf, so a registered leaf stays
 *             registered forever.
 *           - Everything public is a view.
 *         The zero hash can never be registered.
 *
 *      3. The compute check. Given the account hash, its position and its siblings, derive
 *         the leaf, hash it up the tree and compare the result with the SMT root. Pure
 *         math, no storage.
 *
 *      A leaf is VALID (see {validateLeaf}) when, in this order:
 *
 *        1. its account hash is not zero,
 *        2. its leaf (derived from position and account hash) is registered, which is a
 *           cheap map read, so an unknown leaf costs no tree hashing,
 *        3. its proof has exactly {_smtDepth} siblings,
 *        4. its position fits in that depth, and
 *        5. hashing the leaf up the tree with its siblings gives the root.
 *
 *      Registered alone only says the leaf was part of an accepted update at some point. A
 *      proof alone only says the leaf is in the tree with that root. Together, checked
 *      against the CURRENT root, they say the leaf is the account's current state: after an
 *      UPDATE or a DELETE the old leaf is still registered (append-only) but no longer
 *      hashes to the current root, so it fails with RootMismatch. Checked against an older
 *      accepted root, a valid result means the leaf was current at that root.
 *
 *      The root is NOT stored here. The registry already holds it and passes it in, and the
 *      result only means "in the tree with THIS root", so the caller must pass a root it
 *      trusts: the registry's accepted root. A zero root never matches.
 *
 *      Path rule: sibling i is at level i, counted from the leaf. Bit i of `position`
 *      (bit 0 is the least significant) is 0 when the node is the left child at that level,
 *      so the parent is hash(node, sibling), and 1 when it is the right child, so the parent
 *      is hash(sibling, node). A fixed depth stops an inner node from being passed off as a
 *      leaf. As a second guard, a leaf is hashed from 40 bytes and an inner node from 64,
 *      so the two never share a preimage.
 *
 *      Not decided by the spec, so they are virtual and must match the Go backend and the
 *      TypeScript clients: the leaf hash ({_hashLeaf}, as above), the node hash
 *      ({_hashNode}, default keccak256 of left || right), the depth ({_smtDepth}, default
 *      64, the width of `position`) and the bit order above.
 *
 *      NOTE: reference implementation, not audited.
 */
abstract contract ClueAccounts {
    // Types

    /// @notice Result of checking one leaf. Valid is 1, never 0, so an unset result can
    ///         never be read as valid.
    enum LeafCheck {
        None,
        Valid,
        EmptyLeaf,
        NotRegistered,
        InvalidDepth,
        InvalidPosition,
        RootMismatch
    }

    /// @notice One account to check and the proof that places it in the tree. The leaf is
    ///         not carried: it is derived from `position` and `accountHash`.
    struct LeafProof {
        bytes32 accountHash;   // Hash of the account (Account.Hash)
        uint64 position;       // Position of the account in the SMT (Account.TreePosition)
        bytes32[] siblings;    // Siblings from the leaf level up to just below the root
    }

    // Storage

    /// @dev Leaf hash => true once registered. Never set back to false.
    mapping(bytes32 => bool) internal _registered;

    // Errors

    /// @dev The zero hash cannot be registered.
    error ZeroLeaf();

    /// @dev A leaf is not registered. Reverted by {_requireRegistered}.
    error LeafNotRegistered(bytes32 leaf);

    /// @dev One leaf of a batch is not registered. `index` is its position in the batch.
    error LeafNotRegisteredAt(uint256 index, bytes32 leaf);

    /// @dev A single leaf failed its check.
    error InvalidLeaf(LeafCheck reason);

    /// @dev One leaf of a batch failed its check. `index` is its position in the batch.
    error InvalidLeafAt(uint256 index, LeafCheck reason);

    /// @dev A batch with no leaves is never valid.
    error EmptyLeafBatch();

    /// @dev Parallel arrays passed to a batch function have different lengths.
    error LengthMismatch();

    // Events

    /// @notice A leaf was added to the record. Emitted once per leaf, never for a leaf that
    ///         was already registered.
    event LeafRegistered(bytes32 indexed leaf);

    // Public views: the leaf hash

    /// @notice The leaf hash of an account: keccak256 of `position` then `accountHash`.
    ///         Use it to build the leaf a client expects to find in the record.
    function computeLeafHash(uint64 position, bytes32 accountHash) public pure returns (bytes32) {
        return _hashLeaf(position, accountHash);
    }

    // Public views: the record

    /// @notice Returns true if `leaf` has been registered.
    function isRegistered(bytes32 leaf) public view returns (bool) {
        return _registered[leaf];
    }

    /**
     * @notice Checks a list of leaves against the record only, and returns one result per
     *         leaf, in order. No proof is checked.
     * @dev Never reverts on an unregistered leaf. An empty list returns an empty result and
     *      `allRegistered` false.
     * @return results True for each leaf that is registered.
     * @return allRegistered True only if the list is non-empty and every leaf is registered.
     */
    function areRegistered(bytes32[] memory leaves)
        public
        view
        returns (bool[] memory results, bool allRegistered)
    {
        uint256 n = leaves.length;
        results = new bool[](n);
        allRegistered = n != 0;
        for (uint256 i; i < n; ++i) {
            bool ok = _registered[leaves[i]];
            results[i] = ok;
            if (!ok) allRegistered = false;
        }
    }

    // Public views: the compute check

    /// @notice The root that the account's leaf and its siblings hash up to. Does not check
    ///         it against anything. Useful for building test vectors.
    function computeRoot(LeafProof memory proof) public pure returns (bytes32) {
        return _computeRoot(proof);
    }

    // Public views: full validation (registered and proven)

    /// @notice Returns true if the leaf is registered and a valid leaf of the tree with root
    ///         `root`.
    function isValidLeaf(bytes32 root, LeafProof memory proof) public view returns (bool) {
        return _checkProof(root, proof) == LeafCheck.Valid;
    }

    /// @notice Checks one leaf and returns the first failing check, or Valid. Full check:
    ///         account hash not zero, leaf registered, proof hashes up to `root`. An unknown
    ///         leaf fails with NotRegistered before any tree hashing.
    function validateLeaf(bytes32 root, LeafProof memory proof) public view returns (LeafCheck) {
        return _checkProof(root, proof);
    }

    /**
     * @notice Batch full check. Same rules as {validateLeaf}, one result per proof, in order.
     * @dev Never reverts on an invalid leaf, so a caller sees every failure at once. An
     *      empty list returns an empty result and `allValid` false.
     * @param root SMT root every leaf is checked against.
     * @param proofs Accounts with their positions and siblings.
     * @return results Result of each leaf; Valid or its first failing check.
     * @return allValid True only if the list is non-empty and every leaf is Valid.
     */
    function validateLeafProofs(bytes32 root, LeafProof[] memory proofs)
        public
        view
        returns (LeafCheck[] memory results, bool allValid)
    {
        return _checkLeaves(root, proofs);
    }

    /**
     * @notice Checks a batch of (accountHash, position) pairs against the record, and
     *         returns one result per pair, in order.
     * @dev Registration-only: no siblings are supplied, so membership in the tree with
     *      `root` is NOT verified. `root` is used only to fail closed on a zero root (every
     *      result is RootMismatch). For membership use {validateLeaf} or
     *      {validateLeafProofs}. Per-item rules are those of {_checkRegistration}.
     *      Reverts with {LengthMismatch} if the arrays differ in length. An empty list
     *      returns an empty result and `allValid` false.
     * @param root SMT root; a zero root fails every item closed. Not otherwise checked.
     * @param accountHashes Account hashes, parallel to `positions`.
     * @param positions     Tree positions, parallel to `accountHashes`.
     * @return results Result of each pair; Valid or its first failing check.
     * @return allValid True only if the list is non-empty and every pair is Valid.
     */
    function validateLeaves(
        bytes32 root,
        bytes32[] memory accountHashes,
        uint256[] memory positions
    ) public view returns (LeafCheck[] memory results, bool allValid) {
        uint256 n = accountHashes.length;
        if (positions.length != n) revert LengthMismatch();

        if (root == bytes32(0)) return (_allRootMismatch(n), false);

        results = new LeafCheck[](n);
        for (uint256 i; i < n; ++i) {
            results[i] = _checkRegistration(accountHashes[i], positions[i]);
        }
        allValid = _allValid(results);
    }

    /**
     * @notice Checks a batch of already-computed leaf hashes against the record. `must`
     *         controls whether registration is required.
     * @dev Registration-only: leaf hashes are taken as given, so membership in the tree with
     *      `root` is NOT verified. `root` is used only to fail closed on a zero root (every
     *      result is RootMismatch). Per-item rules are those of {_checkLeafHash}; when `must`
     *      is false only the zero leaf hash is rejected. An empty list returns an empty
     *      result and `allValid` false.
     * @param root       SMT root; a zero root fails every item closed. Not otherwise checked.
     * @param leafHashes Leaf hashes to check.
     * @param must       True to require each leafHash to be registered.
     * @return results Result of each leaf hash; Valid or its first failing check.
     * @return allValid True only if the list is non-empty and every leaf hash is Valid.
     */
    function validateLeafHash(
        bytes32 root,
        bytes32[] memory leafHashes,
        bool must
    ) public view returns (LeafCheck[] memory results, bool allValid) {
        uint256 n = leafHashes.length;
        if (root == bytes32(0)) return (_allRootMismatch(n), false);

        results = new LeafCheck[](n);
        for (uint256 i; i < n; ++i) {
            results[i] = _checkLeafHash(leafHashes[i], must);
        }
        allValid = _allValid(results);
    }

    // Registry-facing (internal): checks that revert

    /// @dev Reverts with {LeafNotRegistered} unless `leaf` is registered.
    function _requireRegistered(bytes32 leaf) internal view {
        if (!_registered[leaf]) revert LeafNotRegistered(leaf);
    }

    /// @dev Reverts on the first unregistered leaf with {LeafNotRegisteredAt}, or with
    ///      {EmptyLeafBatch} if there are none. No proof is checked.
    function _requireRegisteredLeaves(bytes32[] memory leaves) internal view {
        uint256 n = leaves.length;
        _requireNonEmpty(n);
        for (uint256 i; i < n; ++i) {
            if (!_registered[leaves[i]]) revert LeafNotRegisteredAt(i, leaves[i]);
        }
    }

    /// @dev Reverts with {InvalidLeaf} unless the leaf is registered and valid against `root`.
    function _requireValidLeaf(bytes32 root, LeafProof memory proof) internal view {
        LeafCheck r = _checkProof(root, proof);
        if (r != LeafCheck.Valid) revert InvalidLeaf(r);
    }

    /// @dev Reverts on the first invalid leaf with {InvalidLeafAt}, or with
    ///      {EmptyLeafBatch} if there are none. One bad leaf rejects the whole batch.
    function _requireValidLeaves(bytes32 root, LeafProof[] memory proofs) internal view {
        uint256 n = proofs.length;
        _requireNonEmpty(n);
        for (uint256 i; i < n; ++i) {
            LeafCheck r = _checkProof(root, proofs[i]);
            if (r != LeafCheck.Valid) revert InvalidLeafAt(i, r);
        }
    }

    // Registry-facing (internal): registration

    /// @dev Registers one leaf without a proof. Reverts on the zero hash. A leaf that is
    ///      already registered is left as it is and emits nothing. Build the leaf with
    ///      {_hashLeaf}.
    /// @return added True if the leaf was new.
    function _registerLeaf(bytes32 leaf) internal returns (bool added) {
        if (leaf == bytes32(0)) revert ZeroLeaf();
        if (_registered[leaf]) return false;
        _registered[leaf] = true;
        emit LeafRegistered(leaf);
        return true;
    }

    /// @dev Registers a batch of leaves without proofs, for a registry that already trusts
    ///      the batch (validator quorum reached). Reverts on an empty batch or on the zero
    ///      hash; a leaf that is already registered is skipped, so a repeat never blocks a
    ///      batch.
    /// @return added Number of leaves that were new.
    function _registerLeaves(bytes32[] memory leaves) internal returns (uint256 added) {
        uint256 n = leaves.length;
        _requireNonEmpty(n);
        for (uint256 i; i < n; ++i) {
            if (_registerLeaf(leaves[i])) ++added;
        }
    }

    /**
     * @dev Registers a batch of accounts only if each one proves into `root`, the root just
     *      accepted. The leaf of each is derived from its position and account hash. Nothing
     *      is registered unless every proof passes: one bad proof reverts the whole batch
     *      with {InvalidLeafAt}. Costs one hash per sibling per account, so it is much
     *      heavier than {_registerLeaves}.
     *
     *      Proofs are against the FINAL root of the batch, so only a leaf that is in that
     *      root can be registered here. If a serial changed twice in one Version, only its
     *      last leaf proves; the intermediate leaf does not and is never registered.
     *
     *      Uses {_checkProofOnly}, not {_checkProof}, because the leaves being registered
     *      are by definition not in the record yet.
     * @return added Number of leaves that were new.
     */
    function _registerProvenLeaves(bytes32 root, LeafProof[] memory proofs) internal returns (uint256 added) {
        uint256 n = proofs.length;
        _requireNonEmpty(n);
        for (uint256 i; i < n; ++i) {
            LeafCheck r = _checkProofOnly(root, proofs[i]);
            if (r != LeafCheck.Valid) revert InvalidLeafAt(i, r);
            if (_registerLeaf(_hashLeaf(proofs[i].position, proofs[i].accountHash))) ++added;
        }
    }

    // Core checks

    /// @dev Non-reverting batch full check with proofs. See {validateLeafProofs}.
    function _checkLeaves(bytes32 root, LeafProof[] memory proofs)
        internal
        view
        returns (LeafCheck[] memory results, bool allValid)
    {
        uint256 n = proofs.length;
        results = new LeafCheck[](n);
        for (uint256 i; i < n; ++i) {
            results[i] = _checkProof(root, proofs[i]);
        }
        allValid = _allValid(results);
    }

    /// @dev Full check: account hash not zero, leaf registered, then the proof. The map is
    ///      read before any tree hashing, so an unknown leaf fails cheaply.
    function _checkProof(bytes32 root, LeafProof memory p) internal view virtual returns (LeafCheck) {
        if (p.accountHash == bytes32(0)) return LeafCheck.EmptyLeaf;
        if (!_registered[_hashLeaf(p.position, p.accountHash)]) return LeafCheck.NotRegistered;
        return _checkProofOnly(root, p);
    }

    /// @dev Proof check only, without the registration check: account hash not zero, depth,
    ///      position, then the compute. A zero `root` never matches, so a missing root fails
    ///      closed. Also used by {_registerProvenLeaves} for leaves not yet registered.
    function _checkProofOnly(bytes32 root, LeafProof memory p) internal pure virtual returns (LeafCheck) {
        if (p.accountHash == bytes32(0)) return LeafCheck.EmptyLeaf;

        // A depth of 0 or above 64 cannot be addressed by a uint64 position, so a misconfigured
        // override fails closed here rather than accepting proofs it cannot bind to a position.
        uint256 depth = _smtDepth();
        if (depth == 0 || depth > 64) return LeafCheck.InvalidDepth;
        if (p.siblings.length != depth) return LeafCheck.InvalidDepth;
        if (depth < 64 && (uint256(p.position) >> depth) != 0) return LeafCheck.InvalidPosition;

        if (root == bytes32(0) || _computeRoot(p) != root) return LeafCheck.RootMismatch;

        return LeafCheck.Valid;
    }

    /// @dev Registration-only check for one (accountHash, position) pair, used by
    ///      {validateLeaves}. Cannot verify against a root without siblings. Rules, in order:
    ///      zero account hash -> EmptyLeaf; position over uint64 -> InvalidPosition; leaf not
    ///      registered -> NotRegistered; otherwise Valid.
    function _checkRegistration(bytes32 accountHash, uint256 position) internal view returns (LeafCheck) {
        if (accountHash == bytes32(0)) return LeafCheck.EmptyLeaf;
        if (position > type(uint64).max) return LeafCheck.InvalidPosition;
        return _checkLeafHash(_hashLeaf(uint64(position), accountHash), true);
    }

    /// @dev Check for one already-computed leaf hash: zero -> EmptyLeaf; `must` and not
    ///      registered -> NotRegistered; otherwise Valid. With `must` false the record is
    ///      not read.
    function _checkLeafHash(bytes32 leaf, bool must) internal view returns (LeafCheck) {
        if (leaf == bytes32(0)) return LeafCheck.EmptyLeaf;
        if (must && !_registered[leaf]) return LeafCheck.NotRegistered;
        return LeafCheck.Valid;
    }

    /// @dev Reverts with {EmptyLeafBatch} if a batch has no items.
    function _requireNonEmpty(uint256 n) private pure {
        if (n == 0) revert EmptyLeafBatch();
    }

    /// @dev `n` results, all RootMismatch. Used to fail a whole batch closed on a zero root.
    function _allRootMismatch(uint256 n) private pure returns (LeafCheck[] memory results) {
        results = new LeafCheck[](n);
        for (uint256 i; i < n; ++i) results[i] = LeafCheck.RootMismatch;
    }

    /// @dev True only if `results` is non-empty and every entry is Valid.
    function _allValid(LeafCheck[] memory results) private pure returns (bool) {
        uint256 n = results.length;
        if (n == 0) return false;
        for (uint256 i; i < n; ++i) {
            if (results[i] != LeafCheck.Valid) return false;
        }
        return true;
    }

    /// @dev Derives the leaf from position and account hash, then hashes it up the tree
    ///      with its siblings, following the path rule of the contract header. Uses
    ///      `siblings.length` as the depth, so call it only after the depth check, or for
    ///      test vectors.
    function _computeRoot(LeafProof memory p) internal pure returns (bytes32 node) {
        node = _hashLeaf(p.position, p.accountHash);
        uint256 position = p.position;
        bytes32[] memory siblings = p.siblings;
        uint256 n = siblings.length;
        for (uint256 i; i < n; ++i) {
            bytes32 sibling = siblings[i];
            node = ((position >> i) & 1) == 0 ? _hashNode(node, sibling) : _hashNode(sibling, node);
        }
    }

    // Tree parameters

    /// @dev The leaf hash: keccak256 of the 8-byte position followed by the 32-byte account
    ///      hash. Override to match the SMT's leaf rule.
    ///      Same value as `keccak256(abi.encodePacked(position, accountHash))`, hashed in
    ///      scratch space so the 64 hashes of a proof do not grow memory.
    function _hashLeaf(uint64 position, bytes32 accountHash) internal pure virtual returns (bytes32 h) {
        assembly ("memory-safe") {
            mstore(0x00, shl(192, position)) // 8 bytes, big endian, at 0x00..0x07
            mstore(0x08, accountHash)        // 32 bytes at 0x08..0x27
            h := keccak256(0x00, 0x28)
        }
    }

    /// @dev Hash of two child nodes. Override to match the SMT's hash function. Same value
    ///      as `keccak256(abi.encodePacked(left, right))`, hashed in scratch space.
    function _hashNode(bytes32 left, bytes32 right) internal pure virtual returns (bytes32 h) {
        assembly ("memory-safe") {
            mstore(0x00, left)
            mstore(0x20, right)
            h := keccak256(0x00, 0x40)
        }
    }

    /// @dev Number of siblings every proof must carry. Override to match the SMT depth.
    ///      At most 64, because `position` is a uint64.
    function _smtDepth() internal pure virtual returns (uint256) {
        return 64;
    }
}