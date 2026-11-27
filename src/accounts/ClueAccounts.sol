// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {ClueAuthorizer} from "../libs/authorizers/ClueAuthorizer.sol";

/**
 * @title ClueAccounts
 * @notice Standalone contract that records which account leaf hashes have been registered
 *         into the Clue Protocol and checks that a leaf, or a list of leaves, is valid:
 *         registered, and a member of the Sparse Merkle Tree (SMT) with a given root.
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
 *      2. The record. One map, `leaf => registered`, append-only. Reads are open to
 *         everyone; writes are restricted to the set of contract authorizers:
 *           - A leaf is added with {registerLeaf}, {registerLeaves} or {registerAccounts}.
 *           - There is no function that removes or clears a leaf, so a registered leaf stays
 *             registered forever.
 *           - Registration is open only to an address marked as an authorizer; every read is
 *             open. The zero hash can never be registered.
 *
 *      The authorizer set is provided by {ClueAuthorizer}: the deployer is the initial owner
 *      and authorizer, and the owner can call {addAuthorizer} and {revokeAuthorizer}. An
 *      authorizer cannot transfer the role, and an unrelated account cannot take it over.
 *      Following {ClueValidators}, the authorizers are contract-level accounts; they are
 *      not per-code validator authorizers.
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
contract ClueAccounts is ClueAuthorizer {
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

    /// @notice One account to register: its hash and its position in the tree. The leaf is
    ///         derived from the two, so it is not carried. No siblings: registration trusts
    ///         the batch, it does not prove membership.
    struct AccountLeaf {
        bytes32 accountHash;   // Hash of the account (Account.Hash)
        uint64 position;       // Position of the account in the SMT (Account.TreePosition)
    }

    // Storage

    /// @dev Leaf hash => true once registered. Never set back to false.
    mapping(bytes32 => bool) internal _registered;

    // Errors

    /// @dev The zero hash cannot be registered.
    error ZeroLeaf();

    /// @dev An account with a zero hash cannot be registered.
    error ZeroAccountHash();

    /// @dev A leaf is not registered. Reverted by {_requireRegistered}.
    error LeafNotRegistered(bytes32 leaf);

    /// @dev A single leaf failed its check.
    error InvalidLeaf(LeafCheck reason);

    /// @dev A batch with no leaves is never valid.
    error EmptyLeafBatch();

    // Events

    /// @notice A leaf was added to the record. Emitted once per leaf, never for a leaf that
    ///         was already registered.
    event LeafRegistered(bytes32 indexed leaf);

    // Public views: the leaf hash

    /// @notice The leaf hash of an account: keccak256 of `position` then `accountHash`.
    ///         Use it to build the leaf a client expects to find in the record.
    function computeLeafHash(uint64 position, bytes32 accountHash) external pure returns (bytes32) {
        return _hashLeaf(position, accountHash);
    }

    // Public views: the record

    /// @notice Returns true if `leaf` has been registered.
    function isRegistered(bytes32 leaf) external view returns (bool) {
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
        external
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
    function computeRoot(LeafProof memory proof) external pure returns (bytes32) {
        return _computeRoot(proof);
    }

    // Public views: full validation (registered and proven)

    /// @notice Checks one leaf and returns the first failing check, or Valid. Full check:
    ///         account hash not zero, leaf registered, proof hashes up to `root`. An unknown
    ///         leaf fails with NotRegistered before any tree hashing.
    function validateLeaf(bytes32 root, LeafProof memory proof) external view returns (LeafCheck) {
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
        external
        view
        returns (LeafCheck[] memory results, bool allValid)
    {
        return _checkLeaves(root, proofs);
    }

    // External: checks that revert

    /// @notice Reverts with {LeafNotRegistered} unless `leaf` is registered.
    function requireRegistered(bytes32 leaf) external view {
        _requireRegistered(leaf);
    }

    // External: registration (Authorizer only)

    /// @notice Registers one leaf. Reverts on the zero hash. A leaf that is already
    ///         registered is left as it is and emits nothing.
    function registerLeaf(bytes32 leaf) external onlyAuthorizer {
        _registerLeaf(leaf);
    }

    /// @notice Registers a batch of leaves the authorizer already trusts. Reverts on an
    ///         empty batch or on the zero hash; duplicates are skipped.
    function registerLeaves(bytes32[] calldata leaves) external onlyAuthorizer {
        _registerLeaves(leaves);
    }

    /// @notice Registers a batch of accounts the authorizer already trusts. The leaf of each
    ///         is derived from its position and account hash. No proof or root is needed.
    function registerAccounts(AccountLeaf[] calldata accounts) external onlyAuthorizer {
        _registerAccounts(accounts);
    }

    // Internal: checks that revert

    /// @dev Reverts with {LeafNotRegistered} unless `leaf` is registered.
    function _requireRegistered(bytes32 leaf) internal view {
        if (!_registered[leaf]) revert LeafNotRegistered(leaf);
    }

    // Internal: registration

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
     * @dev Registers a batch of accounts without checking tree membership. The leaf of each
     *      is derived from its position and account hash. The batch is trusted: the registry
     *      calls this only for a root the validators already accepted, so no root or siblings
     *      are needed and nothing here proves anything. Reverts on an empty batch or a zero
     *      account hash; an account whose leaf is already registered is skipped, so a repeat
     *      never blocks a batch.
     */
    function _registerAccounts(AccountLeaf[] memory accounts) internal {
        uint256 n = accounts.length;
        _requireNonEmpty(n);
        for (uint256 i; i < n; ++i) {
            if (accounts[i].accountHash == bytes32(0)) revert ZeroAccountHash();
            _registerLeaf(_hashLeaf(accounts[i].position, accounts[i].accountHash));
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
    ///      closed. Called by {_checkProof} after its registration check.
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

    /// @dev Reverts with {EmptyLeafBatch} if a batch has no items.
    function _requireNonEmpty(uint256 n) private pure {
        if (n == 0) revert EmptyLeafBatch();
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