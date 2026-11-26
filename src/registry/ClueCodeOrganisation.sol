// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";
import {Initializable} from "@openzeppelin/contracts/proxy/utils/Initializable.sol";
import {Clones} from "@openzeppelin/contracts/proxy/Clones.sol";
import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import {IERC721Receiver} from "@openzeppelin/contracts/token/ERC721/IERC721Receiver.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {IClueCode} from "../code/IClueCode.sol";


abstract contract ClueCodeOrganisation is IERC721Receiver, ReentrancyGuard{
    using EnumerableSet for EnumerableSet.UintSet;

    enum CodeStatus { Open, Stateless, Active, Archived }
    enum OrgStatus { None, Active, Inactive, Archive }

    /// @notice ClueCode ERC-721 contract whose tokens this vault holds.
    IClueCode public immutable clueCode;

    struct OrgRoleData {
        mapping(address account => bool) hasRole;
        bytes32 adminRole;
    }

    /// @notice On-chain organisation record, mirroring the backend Organisation model.
    struct Organisation {
        uint64 id;                   // json:"id"
        string name;                 // json:"name"
        OrgStatus status;            // json:"status"
        string website;              // json:"website"
        bytes32 signer;              // json:"signer"; first Ed25519 public key (backwards-compatible view)
        uint32 verificationFlag;     // json:"verificationFlag"
    }

    uint256 internal _nextOrgId;
    uint248 public totalOrganisation;

    /// @dev Code => organisation that owns it. Zero means the code is not vaulted.
    mapping(uint32 code => uint256 organisation) internal _codeOrg;

    /// @dev org_id => role => role data
    mapping(uint256 => mapping(bytes32 => OrgRoleData)) private _roles;

    /// @dev org_id => organisation record. Replaces the old initiator-only mapping.
    mapping(uint256 orgId => Organisation) internal _organisations;

    /// @dev org_id => Ed25519 signers. Append-only: an organisation cannot remove a signer.
    mapping(uint256 orgId => bytes32[] signers) internal _organisationSigners;

    mapping(uint64 code => CodeStatus) internal code_status;

    /// @dev Role admin that cannot be held, used to lock roles from public administration.
    bytes32 internal constant NO_ROLE_ADMIN = bytes32(type(uint256).max);

    bytes32 public constant ORG_ADMIN_ROLE = 0x00;
    bytes32 public constant ORG_CONTROLLER_ROLE = keccak256("ORG_CONTROLLER_ROLE");
    bytes32 public constant ORG_UPDATER_ROLE = keccak256("ORG_UPDATER_ROLE");

    // --- Custom Errors ---
    error OrganisationNotFound(uint256 orgId);
    error CodeAlreadyVaulted(uint32 code, uint256 existingOrgId);
    error CodeNotVaulted(uint32 code);
    error NotOrgInitiator(uint32 code, address caller);
    error UnauthorizedOrgRole(uint256 orgId, bytes32 role, address account);
    error ZeroAddress();
    error InvalidSigner();

    error InvalidOrganisation(uint256 orgId);
    error RoleLocked(uint256 orgId, bytes32 role);
    error MissingOrgRole(uint256 orgId, bytes32 adminRole, address account);
    error CodeNotStateless(uint32 code, CodeStatus current);
    error CannotSetOpenViaStatus(uint32 code);
    error ReceiverZeroAddress();
    error InvalidERC721Sender(address sender);
    error InvalidOrganisationSigner(uint256 orgId, bytes32 signer);
    error OrganisationSignerAlreadyAdded(uint256 orgId, bytes32 signer);



    // --- Events ---
    event OrganisationCreated(
        uint256 indexed orgId,
        address indexed initiator,
        string name,
        string website,
        bytes32 signer
    );
    event OrganisationSignerAdded(uint256 indexed orgId, bytes32 indexed signer, address indexed sender);
    event CodeVaulted(uint32 indexed code, uint256 indexed orgId, address indexed initiator);
    event CodeReleased(uint32 indexed code, uint256 indexed orgId, address indexed receiver);
    event CodeStatusSet(uint32 indexed code, CodeStatus newStatus);
    event OrgRoleAdminChanged(uint256 indexed orgId, bytes32 indexed role, bytes32 newAdminRole);
    event OrgRoleRevoked(uint256 indexed orgId, bytes32 indexed role, address indexed account, address sender);
    event OrgRoleGranted(
        uint256 indexed orgId,
        bytes32 indexed role,
        address indexed account,
        address sender
    );

    // --- Modifiers ---
    modifier onlyExistingOrg(uint256 orgId) {
        if (orgId == 0 || orgId > _nextOrgId) {
            revert OrganisationNotFound(orgId);
        }
        _;
    }

    modifier onlyOrgRole(uint256 orgId, bytes32 role) {
        if (!hasOrgRole(orgId, role, msg.sender)) {
            revert UnauthorizedOrgRole(orgId, role, msg.sender);
        }
        _;
    }

    ///  Code Vaulting

    /// @notice Transfers a ClueCode ERC-721 token from the caller into this contract and assigns it to `orgId`.
    /// @param code The uint32 token ID to vault.
    /// @param orgId The target organisation ID that will own the vaulted code.
    function activateVaultCode(uint32 code, uint256 orgId) external {
        if (orgId == 0 || orgId > _nextOrgId) {
            revert InvalidOrganisation(orgId);
        }

        if (!_roles[orgId][ORG_UPDATER_ROLE].hasRole[msg.sender]) {
            revert MissingOrgRole(orgId, ORG_UPDATER_ROLE, msg.sender);
        }

        _vaultCode(code, orgId);
    }

    /// @dev Vaults `code` into `orgId` with no role check: the caller is `msg.sender`, who
    ///      must own the token (the ERC-721 transfer enforces it). Callers must have
    ///      validated `orgId` and authorised `msg.sender` themselves.
    function _vaultCode(uint32 code, uint256 orgId) internal {
        uint256 currentOrg = _codeOrg[code];
        if (currentOrg != 0) {
            revert CodeAlreadyVaulted(code, currentOrg);
        }

        _codeOrg[code] = orgId;
        emit CodeVaulted(code, orgId, msg.sender);

        clueCode.safeTransferFrom(msg.sender, address(this), uint256(code));
    }

    /// @notice Releases a vaulted code back to `receiver`, resetting its organisation assignment.
    /// @dev Caller must hold any role within the code's owning organisation.
    ///      The code's status must be `Stateless`; an `Active` code (one that still
    ///      holds Clue IDs or link records) cannot be released.
    ///      Resets the code's status to `Open` and clears its organisation mapping.
    /// @param code     The uint32 token ID to release.
    /// @param receiver The address that will receive the ERC-721 token.
    function releaseVaultCode(uint32 code, address receiver) external nonReentrant {
        if (receiver == address(0)) revert ReceiverZeroAddress();

        uint256 orgId = _codeOrg[code];
        if (orgId == 0) revert CodeNotVaulted(code);

        // Caller must be a member of the owning organisation (any role suffices).
        if (
            !_roles[orgId][ORG_ADMIN_ROLE].hasRole[msg.sender] &&
            !_roles[orgId][ORG_CONTROLLER_ROLE].hasRole[msg.sender] &&
            !_roles[orgId][ORG_UPDATER_ROLE].hasRole[msg.sender]
        ) {
            revert UnauthorizedOrgRole(orgId, bytes32(0), msg.sender);
        }

        // Code must be Stateless — Active codes still hold data and cannot be released.
        CodeStatus current = code_status[uint64(code)];
        if (current != CodeStatus.Stateless) {
            revert CodeNotStateless(code, current);
        }

        // Clean up provider/validator/link state before clearing core vault mappings.
        _cleanupCodeOnRelease(code, orgId);

        // Clear mappings and reset status to Open.
        delete _codeOrg[code];
        code_status[uint64(code)] = CodeStatus.Open;

        emit CodeReleased(code, orgId, receiver);

        clueCode.safeTransferFrom(address(this), receiver, uint256(code));
    }

    /// @dev Hook for inheriting contracts to clean up code-related state on release.
    ///      Called before core vault mappings are cleared.
    function _cleanupCodeOnRelease(uint32 code, uint256 orgId) internal virtual {}

    /// @dev Sets the status of a vaulted code.
    ///      Intended for the inheriting registry:
    ///      - `Stateless` → `Active` on the first Clue ID.
    ///      - `Active` → `Stateless` when the code holds no Clue ID and,
    ///        for a Route code, no link records.
    ///      `Open` cannot be set here; use `releaseVaultCode` instead.
    /// @param code      The code whose status is being updated.
    /// @param newStatus The desired target status (must not be `Open`).
    function _setCodeStatus(uint32 code, CodeStatus newStatus) internal {
        if (newStatus == CodeStatus.Open) revert CannotSetOpenViaStatus(code);

        code_status[uint64(code)] = newStatus;
        emit CodeStatusSet(code, newStatus);
    }



    // Organisation

    /// @notice Creates a new organisation. The caller is both its admin and initiator.
    /// @dev `signer` is the organisation's first Ed25519 public key; a zero key is rejected.
    ///      More signers can be added later by an {ORG_ADMIN_ROLE} holder via
    ///      {addOrganisationSigner}; signers are append-only and can never be removed.
    ///      The organisation starts {OrgStatus.Active} with a zero verification flag.
    /// @param name    Organisation name.
    /// @param website Organisation website.
    /// @param signer  Ed25519 public key of the organisation's Profile Server.
    /// @return orgId  The id of the new organisation.
    function createOrganisation(
        string calldata name,
        string calldata website,
        bytes32 signer
    ) public virtual returns (uint256 orgId) {
        if (signer == bytes32(0)) revert InvalidSigner();
        if (_nextOrgId >= type(uint64).max) revert InvalidOrganisation(_nextOrgId);

        orgId = ++_nextOrgId;
        totalOrganisation += 1;

        _organisations[orgId] = Organisation({
            id: uint64(orgId),
            name: name,
            status: OrgStatus.Active,
            website: website,
            signer: signer,
            verificationFlag: 0
        });
        _organisationSigners[orgId].push(signer);

        _roles[orgId][ORG_ADMIN_ROLE].hasRole[msg.sender] = true;
        _roles[orgId][ORG_CONTROLLER_ROLE].adminRole = ORG_ADMIN_ROLE;
        _roles[orgId][ORG_UPDATER_ROLE].adminRole = ORG_ADMIN_ROLE;

        emit OrganisationCreated(orgId, msg.sender, name, website, signer);
    }

    /// @notice Adds an Ed25519 signer public key to `orgId`.
    /// @dev Only an {ORG_ADMIN_ROLE} holder of `orgId` can add signers. Signers are
    ///      append-only: there is deliberately no remove/replace path, so a compromised
    ///      key cannot invalidate existing signatures or existing registry data.
    ///      Zero and duplicate signers are rejected.
    /// @param orgId  Organisation receiving the signer.
    /// @param signer Ed25519 public key to add.
    function addOrganisationSigner(uint256 orgId, bytes32 signer)
        external
        onlyExistingOrg(orgId)
        onlyOrgRole(orgId, ORG_ADMIN_ROLE)
    {
        if (signer == bytes32(0)) revert InvalidOrganisationSigner(orgId, signer);

        bytes32[] storage signers = _organisationSigners[orgId];
        uint256 len = signers.length;
        for (uint256 i = 0; i < len; ++i) {
            if (signers[i] == signer) revert OrganisationSignerAlreadyAdded(orgId, signer);
        }

        signers.push(signer);
        emit OrganisationSignerAdded(orgId, signer, msg.sender);
    }

    /// @notice Returns every Ed25519 signer public key registered for `orgId`.
    function organisationSignersOf(uint256 orgId)
        external
        view
        onlyExistingOrg(orgId)
        returns (bytes32[] memory)
    {
        return _organisationSigners[orgId];
    }

    /// @notice Returns the stored organisation record for `orgId`.
    function organisationOf(uint256 orgId)
        external
        view
        onlyExistingOrg(orgId)
        returns (Organisation memory)
    {
        return _organisations[orgId];
    }

    // ROLES

        /// @notice Returns true if `account` holds `role` within `orgId`.
    function hasOrgRole(uint256 orgId, bytes32 role, address account)
        public
        view
        virtual
        returns (bool)
    {
        return _roles[orgId][role].hasRole[account];
    }

    /// @notice Grants `role` on `orgId` to `account` if `msg.sender` holds the role's admin role for that organisation.
    /// @param orgId The organisation ID scope for the role assignment.
    /// @param role The role identifier to grant.
    /// @param account The address receiving the role.
    function grantOrgRole(uint256 orgId, bytes32 role, address account) internal {
        if (orgId == 0 || orgId > _nextOrgId) {
            revert InvalidOrganisation(orgId);
        }

        OrgRoleData storage roleData = _roles[orgId][role];
        bytes32 adminRole = roleData.adminRole;

        if (adminRole == NO_ROLE_ADMIN) {
            revert RoleLocked(orgId, role);
        }

        bytes32 effectiveAdmin = adminRole == bytes32(0) ? ORG_ADMIN_ROLE : adminRole;
        if (!_roles[orgId][effectiveAdmin].hasRole[msg.sender]) {
            revert MissingOrgRole(orgId, effectiveAdmin, msg.sender);
        }

        _grantOrgRole(orgId, role, account);
    }

    /// @dev Internal helper to grant an organisation role without admin checks (useful during organisation initialization).
    function _grantOrgRole(uint256 orgId, bytes32 role, address account) internal {
        OrgRoleData storage roleData = _roles[orgId][role];
        if (!roleData.hasRole[account]) {
            roleData.hasRole[account] = true;
            emit OrgRoleGranted(orgId, role, account, msg.sender);
        }
    }

    /// @notice Revokes `role` on `orgId` from `account` if `msg.sender` holds the role's admin role for that organisation.
    /// @param orgId The organisation ID scope for the role revocation.
    /// @param role The role identifier to revoke.
    /// @param account The address losing the role.
    function revokeOrgRole(uint256 orgId, bytes32 role, address account) internal {
        if (orgId == 0 || orgId > _nextOrgId) {
            revert InvalidOrganisation(orgId);
        }

        OrgRoleData storage roleData = _roles[orgId][role];
        bytes32 adminRole = roleData.adminRole;

        if (adminRole == NO_ROLE_ADMIN) {
            revert RoleLocked(orgId, role);
        }

        bytes32 effectiveAdmin = adminRole == bytes32(0) ? ORG_ADMIN_ROLE : adminRole;
        if (!_roles[orgId][effectiveAdmin].hasRole[msg.sender]) {
            revert MissingOrgRole(orgId, effectiveAdmin, msg.sender);
        }

        _revokeOrgRole(orgId, role, account);
    }

    /// @dev Internal helper to revoke an organisation role without admin checks.
    function _revokeOrgRole(uint256 orgId, bytes32 role, address account) internal {
        OrgRoleData storage roleData = _roles[orgId][role];
        if (roleData.hasRole[account]) {
            roleData.hasRole[account] = false;
            emit OrgRoleRevoked(orgId, role, account, msg.sender);
        }
    }

    /// @dev Reverts if `msg.sender` does not hold `role` in `orgId`.
    /// @param orgId Organisation whose role is checked.
    /// @param role  The role identifier to check.
    function _checkRole(uint256 orgId, bytes32 role) internal view {
        if (!_roles[orgId][role].hasRole[msg.sender]) {
            revert MissingOrgRole(orgId, role, msg.sender);
        }
    }

        /// @notice Accepts ClueCode ERC-721 safe transfers into this vault.
    /// @dev Restricts mints/transfers to the configured `clueCode` contract.
    function onERC721Received(address, address, uint256, bytes calldata)
        external
        view
        override
        returns (bytes4)
    {
        if (msg.sender != address(clueCode)) revert InvalidERC721Sender(msg.sender);
        return IERC721Receiver.onERC721Received.selector;
    }

}