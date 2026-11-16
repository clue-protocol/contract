// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {IERC721Receiver} from "@openzeppelin/contracts/token/ERC721/IERC721Receiver.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {IClueCode} from "../code/IClueCode.sol";

/**
 * @title ClueCodeVault - ClueRegistry Component
 * @notice Holds ClueCode Provider Code tokens on behalf of their depositors.
 */
abstract contract ClueCodeVault is IERC721Receiver, ReentrancyGuard {
    
    /// @notice Provider Code type, chosen by the depositor at activation.
    enum ProviderType { Static, Route }

    /// @notice Record of a vaulted Provider Code.
    struct VaultRecord {
        address      depositor;     // Address that activated the code.
        ProviderType providerType;  // Type chosen at activation.
    }

    // State variables
    /// @notice ClueCode ERC-721 contract whose tokens are held by this vault.
    IClueCode public immutable clueCode;

    /// @dev Provider Code => vault record. A zero depositor means the code is not vaulted.
    mapping(uint32 => VaultRecord) internal _vault;

    // Errors
    error InvalidParameter();
    error InvalidCode(uint32 code);
    error AlreadyVaulted(uint32 code);
    error NotVaulted(uint32 code);
    error NotMinted(uint32 code);
    error NotOwner(uint32 code);
    error NotApproved(uint32 code);
    error RouteGradeNotAllowed(uint32 code);
    error CustodyNotTaken(uint32 code);
    error DirectTransferNotAllowed();

    // Events
    event CodeActivated(uint32 indexed code, address indexed depositor, ProviderType providerType);

    // Constructor

    constructor(IClueCode clueCode_) {
        if (address(clueCode_) == address(0)) revert InvalidParameter();
        clueCode = clueCode_;
    }

    // Activation

    /**
     * @notice Activates a Provider Code. Called by the holder of the ClueCode NFT
     *         to place the code in the vault.
     * @dev The caller must first approve this contract on ClueCode, either for the
     *      token (`approve`) or for all tokens (`setApprovalForAll`).
     * @param code The Provider Code to activate.
     * @param providerType Static or Route. A Route code is only allowed for grades S and A.
     */
    function activateCode(uint32 code, ProviderType providerType) external nonReentrant {
        _requireValidCode(code);
        _requireNotVaulted(code);
        _requireProviderTypeAllowed(code, providerType);
        _requireTokenOwner(code, msg.sender);
        _requireApproved(code, msg.sender);

        _vault[code] = VaultRecord({depositor: msg.sender, providerType: providerType});

        clueCode.transferFrom(msg.sender, address(this), uint256(code));
        _requireCustody(code);

        emit CodeActivated(code, msg.sender, providerType);
    }

    /// @dev Reverts unless ClueCode recognizes `code` as a valid Provider Code.
    function _requireValidCode(uint32 code) internal view {
        if (!clueCode.isValidCode(uint256(code))) revert InvalidCode(code);
    }

    /// @dev Reverts if `code` is already held by the vault.
    function _requireNotVaulted(uint32 code) internal view {
        if (_vault[code].depositor != address(0)) revert AlreadyVaulted(code);
    }

    /// @dev Reverts if a Route code is requested for a grade that cannot be a Route.
    ///      Only grades S and A may be Route codes; grade B is Static only.
    function _requireProviderTypeAllowed(uint32 code, ProviderType providerType) internal view {
        if (providerType == ProviderType.Route) {
            IClueCode.Grade grade = clueCode.gradeOf(uint256(code));
            if (grade != IClueCode.Grade.S && grade != IClueCode.Grade.A) {
                revert RouteGradeNotAllowed(code);
            }
        }
    }

    /// @dev Reverts unless `account` owns the token. A valid code that has not been minted
    ///      is reported as NotMinted.
    function _requireTokenOwner(uint32 code, address account) internal view {
        try clueCode.ownerOf(uint256(code)) returns (address owner) {
            if (owner != account) revert NotOwner(code);
        } catch {
            revert NotMinted(code);
        }
    }

    /// @dev Reverts unless this contract is approved to transfer the token.
    function _requireApproved(uint32 code, address owner) internal view {
        bool approved = clueCode.getApproved(uint256(code)) == address(this) ||
            clueCode.isApprovedForAll(owner, address(this));
        if (!approved) revert NotApproved(code);
    }

    /// @dev Reverts unless this contract holds the token after the transfer.
    function _requireCustody(uint32 code) internal view {
        if (clueCode.ownerOf(uint256(code)) != address(this)) revert CustodyNotTaken(code);
    }

    // ---------------------------------------------------------------------
    // Views
    // ---------------------------------------------------------------------

    function isVaulted(uint32 code) public view returns (bool) {
        return _vault[code].depositor != address(0);
    }

    function depositorOf(uint32 code) external view returns (address) {
        return _vault[code].depositor;
    }

    /// @notice Returns the provider type chosen at activation. Reverts if the code is not vaulted.
    function providerTypeOf(uint32 code) external view returns (ProviderType) {
        if (_vault[code].depositor == address(0)) revert NotVaulted(code);
        return _vault[code].providerType;
    }

    // ---------------------------------------------------------------------
    // ERC-721 receiver
    // ---------------------------------------------------------------------

    /// @dev Always reverts, so that a token sent with `safeTransferFrom` cannot be held
    ///      without a depositor. Activation uses `transferFrom` and does not call this function.
    function onERC721Received(address, address, uint256, bytes calldata)
        external
        pure
        override
        returns (bytes4)
    {
        revert DirectTransferNotAllowed();
    }
}