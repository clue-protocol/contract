// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {IERC721Receiver} from "@openzeppelin/contracts/token/ERC721/IERC721Receiver.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {IClueCode} from "../code/IClueCode.sol";

/**
 * @title ClueCodeVault
 * @notice ClueRegistry component that holds ClueCode Provider Code tokens on behalf of
 *         their depositors.
 * @dev A vaulted code can be released back to its depositor only while it is Stateless,
 *      meaning it holds no Clue ID and no link. The registry that inherits this vault
 *      moves a code between Stateless and Active with `_setStatus`.
 */
abstract contract ClueCodeVault is IERC721Receiver, ReentrancyGuard {

    // Types

    /// @notice Provider Code type, chosen by the depositor at activation.
    enum ProviderType { Static, Route }

    /**
     * @notice On-chain status of a Provider Code.
     * @dev Open is the zero value, so a code that is not vaulted reads as Open.
     *      INTERNAL-SL appears on-chain as Stateless and has no value of its own.
     *      READ-ONLY is not tracked by the vault.
     */
    enum CodeStatus { Open, Stateless, Active }

    /// @notice Record of a vaulted Provider Code.
    struct VaultRecord {
        address      depositor;     // Address that activated the code
        ProviderType providerType;  // Type chosen at activation
        CodeStatus   status;        // Stateless on activation, then set by the registry
    }

    // Storage

    /// @notice ClueCode ERC-721 contract whose tokens this vault holds.
    IClueCode public immutable clueCode;

    /// @dev Provider Code => vault record. A zero depositor means the code is not vaulted.
    mapping(uint32 => VaultRecord) internal _vault;

    // Errors

    error InvalidParameter();
    error InvalidCode(uint32 code);
    error InvalidStatus(CodeStatus status);
    error AlreadyVaulted(uint32 code);
    error NotVaulted(uint32 code);
    error NotDepositor(uint32 code, address caller);
    error NotStateless(uint32 code, CodeStatus status);
    error RouteGradeNotAllowed(uint32 code);
    error CustodyNotTaken(uint32 code);
    error ReleaseFailed(uint32 code);
    error DirectTransferNotAllowed();

    // Events

    event CodeActivated(uint32 indexed code, address indexed depositor, ProviderType providerType);
    event CodeReleased(uint32 indexed code, address indexed depositor);
    event CodeStatusChanged(uint32 indexed code, CodeStatus from, CodeStatus to);

    // Constructor

    /// @param clueCode_ Address of the ClueCode ERC-721 contract.
    constructor(IClueCode clueCode_) {
        if (address(clueCode_) == address(0)) revert InvalidParameter();
        clueCode = clueCode_;
    }

    // Activation and release

    /**
     * @notice Places a Provider Code in the vault. The code moves from Open to Stateless.
     * @dev The caller must own the token and must first approve this contract on ClueCode,
     *      either for the token (`approve`) or for all tokens (`setApprovalForAll`).
     *      ClueCode itself reverts if the token is not minted, the caller is not its owner,
     *      or this contract is not approved.
     * @param code Provider Code to activate.
     * @param providerType Static or Route. Route is allowed only for grades S and A.
     */
    function activateCode(uint32 code, ProviderType providerType) external nonReentrant {
        _requireValidCode(code);
        _requireNotVaulted(code);
        _requireProviderTypeAllowed(code, providerType);

        _vault[code] = VaultRecord({
            depositor: msg.sender,
            providerType: providerType,
            status: CodeStatus.Stateless
        });

        clueCode.safeTransferFrom(msg.sender, address(this), uint256(code));
        _requireCustody(code);

        emit CodeActivated(code, msg.sender, providerType);
        emit CodeStatusChanged(code, CodeStatus.Open, CodeStatus.Stateless);
    }

    /**
     * @notice Returns a Provider Code to its depositor. The code moves back to Open and
     *         its vault record, including the chosen provider type, is deleted.
     * @dev Only the depositor can release, and only while the code is Stateless. An Active
     *      code must first have its state undone: all Clue IDs deleted and, for a Route
     *      code, all links removed. The token is returned with `transferFrom`, so a
     *      depositor that is a contract is not called back.
     * @param code Provider Code to release.
     */
    function releaseCode(uint32 code) external nonReentrant {
        VaultRecord memory record = _vault[code];

        if (record.depositor == address(0)) revert NotVaulted(code);
        if (record.depositor != msg.sender) revert NotDepositor(code, msg.sender);
        if (record.status != CodeStatus.Stateless) revert NotStateless(code, record.status);

        // Update state before the external call.
        delete _vault[code];
 
        clueCode.transferFrom(address(this), record.depositor, uint256(code));
        if (clueCode.ownerOf(uint256(code)) != record.depositor) revert ReleaseFailed(code);

        emit CodeStatusChanged(code, CodeStatus.Stateless, CodeStatus.Open);
        emit CodeReleased(code, record.depositor);
    }

    // Status

    /**
     * @dev Sets the status of a vaulted code. Intended for the inheriting registry:
     *      Stateless to Active on the first Clue ID, and back to Stateless when the code
     *      holds no Clue ID and, for a Route code, no link records.
     *      Open cannot be set here; use `releaseCode`.
     */
    function _setStatus(uint32 code, CodeStatus newStatus) internal {
        if (newStatus == CodeStatus.Open) revert InvalidStatus(newStatus);
        if (_vault[code].depositor == address(0)) revert NotVaulted(code);

        CodeStatus current = _vault[code].status;
        if (current != newStatus) {
            _vault[code].status = newStatus;
            emit CodeStatusChanged(code, current, newStatus);
        }
    }

    // Checks

    /// @dev Reverts unless ClueCode recognizes `code` as a valid Provider Code.
    function _requireValidCode(uint32 code) internal view {
        if (!clueCode.isValidCode(uint256(code))) revert InvalidCode(code);
    }

    /// @dev Reverts if the vault already holds `code`.
    function _requireNotVaulted(uint32 code) internal view {
        if (_vault[code].depositor != address(0)) revert AlreadyVaulted(code);
    }

    /// @dev Reverts if Route is requested for a grade that cannot be a Route.
    ///      Only grades S and A may be Route codes; grade B is Static only.
    function _requireProviderTypeAllowed(uint32 code, ProviderType providerType) internal view {
        if (providerType == ProviderType.Route) {
            IClueCode.Grade grade = clueCode.gradeOf(uint256(code));
            if (grade != IClueCode.Grade.S && grade != IClueCode.Grade.A) {
                revert RouteGradeNotAllowed(code);
            }
        }
    }

    /// @dev Reverts unless this contract holds the token.
    function _requireCustody(uint32 code) internal view {
        if (clueCode.ownerOf(uint256(code)) != address(this)) revert CustodyNotTaken(code);
    }

    // Views

    /// @notice Returns true if the vault holds `code`.
    function isVaulted(uint32 code) public view returns (bool) {
        return _vault[code].depositor != address(0);
    }

    /// @notice Returns the address that activated `code`, or zero if it is not vaulted.
    function depositorOf(uint32 code) external view returns (address) {
        return _vault[code].depositor;
    }

    /// @notice Returns the status of `code`. A code that is not vaulted is Open.
    function statusOf(uint32 code) external view returns (CodeStatus) {
        return _vault[code].status;
    }

    /// @notice Returns true if `code` is vaulted and Stateless, so it can be released.
    function isReleasable(uint32 code) external view returns (bool) {
        return _vault[code].depositor != address(0) && _vault[code].status == CodeStatus.Stateless;
    }

    /// @notice Returns the provider type chosen at activation. Reverts if `code` is not vaulted.
    function providerTypeOf(uint32 code) external view returns (ProviderType) {
        if (_vault[code].depositor == address(0)) revert NotVaulted(code);
        return _vault[code].providerType;
    }

    // ERC-721 receiver

    /// @dev Accepts a token only when this vault itself sends it to itself during `activateCode`.
    ///      Any other `safeTransferFrom` to the vault reverts, so a token cannot enter
    ///      without a depositor record.
    function onERC721Received(address operator, address, uint256, bytes calldata)
        external
        view
        override
        returns (bytes4)
    {
        if (msg.sender != address(clueCode) || operator != address(this)) {
            revert DirectTransferNotAllowed();
        }
        return IERC721Receiver.onERC721Received.selector;
    }
}
