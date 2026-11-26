// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {EIP712} from "@openzeppelin/contracts/utils/cryptography/EIP712.sol";
import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import {ClueValidators} from "./ClueValidators.sol";

/**
 * @title ClueMultiSig
 * @notice Validator approval gateway for the ClueRegistry.
 * @dev Stores only two things per Provider Code:
 *        - the version currently being updated, and
 *        - the validators that have signed it.
 *      It does not store the proposal. The ClueRegistry holds the data and passes in
 *      `updateHash`, a hash of its stored proposal. The validator signs that hash,
 *      so the registry controls what is being approved and this contract only checks
 *      who approved it. If the signature is valid and the signer is eligible the call
 *      returns; otherwise it reverts.
 *
 *      `updateHash` MUST commit to the previous root, the proposed root and the batch
 *      id (and may cover more). The EIP-712 domain binds the chain id and this
 *      contract, which is the registry, so approvals cannot be replayed elsewhere.
 *
 *      Intended use in the registry:
 *
 *        function validatorSignUpdate(uint32 code, bytes calldata signature) external {
 *            Proposal storage p = _proposal[code];                     // registry data
 *            bool done = _validatorSignUpdate(code, p.version, _updateHash(p), signature);
 *            if (done) _confirm(code);       // apply the root, then call _closeUpdate(code)
 *        }
 *
 *      On upload the registry calls _openUpdate(code, version); on confirm or discard
 *      it calls _closeUpdate(code).
 *
 *      NOTE: draft, not compiled, not audited. Assumes OpenZeppelin 5.x.
 */
abstract contract ClueMultiSig is ClueValidators, EIP712 {
    using EnumerableSet for EnumerableSet.UintSet;

    // Constants

    bytes32 public constant APPROVAL_TYPEHASH =
        keccak256("Approval(uint32 code,uint32 version,bytes32 updateHash)");

    // Storage

    /// @dev Provider Code => version currently being updated. 0 = no open update.
    mapping(uint32 => uint32) internal _updatingVersion;

    /// @dev Provider Code => ids of the validators that signed the open update.
    mapping(uint32 => EnumerableSet.UintSet) internal _signers;

    // Errors

    error InvalidUpdateVersion();
    error UpdateAlreadyOpen(uint32 code, uint32 openVersion);
    error NoUpdateOpen(uint32 code);
    error WrongUpdateVersion(uint32 code, uint32 expected, uint32 got);
    error InvalidSignature();
    error NotCodeValidator(uint32 code, address account);
    error AlreadySigned(uint32 code, uint256 validatorId);

    // Events

    event UpdateOpened(uint32 indexed code, uint32 indexed version);
    event ValidatorSigned(uint32 indexed code, uint32 indexed version, uint256 indexed validatorId, uint256 signedCount);
    event UpdateClosed(uint32 indexed code, uint32 indexed version);

    // Constructor

    constructor() EIP712("ClueRegistry", "1") {}

    // Views

    /// @notice Version of `code` currently being updated, 0 if none.
    function updatingVersionOf(uint32 code) external view returns (uint32) {
        return _updatingVersion[code];
    }

    /// @notice Ids of the validators that have signed the open update of `code`.
    function signersOf(uint32 code) external view returns (uint256[] memory) {
        return _signers[code].values();
    }

    /// @notice True if `validatorId` has signed the open update of `code`.
    function hasSigned(uint32 code, uint256 validatorId) external view returns (bool) {
        return _signers[code].contains(validatorId);
    }

    /// @notice The EIP-712 digest a validator signs to approve `updateHash` for `version`.
    function approvalDigest(uint32 code, uint32 version, bytes32 updateHash) public view returns (bytes32) {
        return _hashTypedDataV4(keccak256(abi.encode(APPROVAL_TYPEHASH, code, version, updateHash)));
    }

    /// @notice Number of signers that are still active and still linked to `code`.
    function approvalsOf(uint32 code) public view returns (uint256) {
        return _approvals(code);
    }

    // Registry-facing (internal)

    /// @dev Starts collecting approvals for `version` of `code`. One update at a time.
    function _openUpdate(uint32 code, uint32 version) internal {
        if (version == 0) revert InvalidUpdateVersion();
        uint32 open = _updatingVersion[code];
        if (open != 0) revert UpdateAlreadyOpen(code, open);

        _updatingVersion[code] = version;
        emit UpdateOpened(code, version);
    }

    /// @dev Ends the open update (confirmed or discarded) and forgets its signers, so
    ///      the same version can be opened again after a discard.
    function _closeUpdate(uint32 code) internal {
        uint32 version = _updatingVersion[code];
        if (version == 0) revert NoUpdateOpen(code);

        EnumerableSet.UintSet storage signers = _signers[code];
        while (signers.length() != 0) {
            signers.remove(signers.at(0));
        }
        delete _updatingVersion[code];

        emit UpdateClosed(code, version);
    }

    /**
     * @dev The gateway. Reverts unless `signature` is a valid EIP-712 approval of
     *      (`code`, `version`, `updateHash`) by an active validator linked to `code`
     *      that has not signed this update yet. On success it records the signer and
     *      returns whether the approval quorum is now reached.
     * @param code Provider Code.
     * @param version Version the registry says is being updated; must be the open one.
     * @param updateHash Hash of the registry's stored proposal.
     * @param signature ECDSA signature over {approvalDigest}. Anyone may relay it.
     * @return quorumReached True when enough linked validators have signed.
     */
    function _validatorSignUpdate(
        uint32 code,
        uint32 version,
        bytes32 updateHash,
        bytes calldata signature
    ) internal returns (bool quorumReached) {
        uint32 open = _updatingVersion[code];
        if (open == 0) revert NoUpdateOpen(code);
        if (version != open) revert WrongUpdateVersion(code, open, version);

        (address signer, ECDSA.RecoverError err, ) =
            ECDSA.tryRecover(approvalDigest(code, version, updateHash), signature);
        if (err != ECDSA.RecoverError.NoError) revert InvalidSignature();

        uint256 validatorId = _validatorIdByAccount[signer];
        if (
            validatorId == 0 ||
            !_validatorActive[validatorId] ||
            !_codeValidators[code].contains(validatorId)
        ) revert NotCodeValidator(code, signer);

        if (!_signers[code].add(validatorId)) revert AlreadySigned(code, validatorId);

        // Emit the *effective* approval count (active + still linked), not raw set
        // length. If the roster changed mid-update, `_signers[code].length()` would
        // overstate progress and mislead subgraphs/frontends. Reuse this value for
        // the quorum check to avoid iterating the set twice.
        uint256 approvals = _approvals(code);
        emit ValidatorSigned(code, version, validatorId, approvals);

        return approvals >= _requiredApprovals(code);
    }

    // Policy

    /**
     * @dev Approvals needed: two thirds of the code's linked validators plus one,
     *      never more than all of them. Override for another quorum.
     *      Examples: 1 validator needs 1, 3 need 3, 4 need 3, 10 need 7.
     */
    function _requiredApprovals(uint32 code) internal view virtual returns (uint256) {
        uint256 n = _codeValidators[code].length();
        if (n == 0) return 1;
        uint256 need = (n * 2) / 3 + 1;
        return need > n ? n : need;
    }

    // Internal

    /// @dev Counts only signers that are still active and still linked, so a validator
    ///      unlinked or disabled after signing no longer counts toward the quorum.
    function _approvals(uint32 code) internal view returns (uint256 n) {
        EnumerableSet.UintSet storage s = _signers[code];
        uint256 len = s.length();
        for (uint256 i = 0; i < len; ++i) {
            uint256 id = s.at(i);
            if (_validatorActive[id] && _codeValidators[code].contains(id)) ++n;
        }
    }
}