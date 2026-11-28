// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {IERC721Receiver} from "@openzeppelin/contracts/token/ERC721/IERC721Receiver.sol";
import {IClueProviderCode, ClueOrgRoles} from "./IClueProviderCode.sol";

/**
 * @title ClueCodeOrganisation
 * @notice Organisation/vault state and registry-driven mutations for Provider Codes.
 *
 * This module is kept separate from the code (ERC-721, grade, canonical form and
 * metadata) implemented by {ClueCode}. The trusted registry holds a typed reference
 * to this contract and drives the state-changing functions on behalf of users,
 * forwarding the original caller explicitly.
 *
 * @dev Abstract: the concrete contract supplies the default-admin role check and
 *      the ERC-721 ownership/transfer primitives used by the vault.
 *
 * NOTE: reference implementation, not audited.
 */
abstract contract ClueCodeOrganisation is IClueProviderCode {
    /// @notice Registry allowed to drive organisation/vault state on behalf of users.
    address public registry;

    // --- Organisation / vault state ---

    struct OrgRoleData {
        mapping(address account => bool) hasRole;
        bytes32 adminRole;
    }

    uint256 internal _nextOrgId;
    uint248 internal _totalOrganisation;

    /// @dev Code => organisation that owns it. Zero means the code is not vaulted.
    mapping(uint32 code => uint256 organisation) internal _codeOrg;

    /// @dev org_id => role => role data.
    mapping(uint256 => mapping(bytes32 => OrgRoleData)) private _roles;

    /// @dev org_id => organisation record.
    mapping(uint256 orgId => Organisation) internal _organisations;

    /// @dev org_id => Ed25519 signers. Append-only: a signer can never be removed.
    mapping(uint256 orgId => bytes32[] signers) internal _organisationSigners;

    /// @dev Code => lifecycle status.
    mapping(uint32 code => CodeStatus) internal _codeStatus;

    // ------------------------------------------------------------------
    // Errors and events
    // ------------------------------------------------------------------

    // All errors and events emitted here are declared once in {IClueProviderCode},
    // so the Solidity interface and the deployed ABI cannot drift apart.

    // ------------------------------------------------------------------
    // Hooks supplied by the concrete code contract
    // ------------------------------------------------------------------

    /// @dev Reverts unless `msg.sender` holds the contract's default admin role.
    function _checkDefaultAdminRole() internal view virtual;

    /// @dev Returns the owner of the Provider Code token.
    function _codeOwner(uint256 code) internal view virtual returns (address);

    /// @dev Moves the Provider Code token between accounts, as the vault requires.
    function _transferCode(address from, address to, uint256 code) internal virtual;

    // ------------------------------------------------------------------
    // Modifiers and registry wiring
    // ------------------------------------------------------------------

    modifier onlyRegistry() {
        if (msg.sender != registry) revert NotRegistry(msg.sender);
        _;
    }

    modifier onlyExistingOrg(uint256 orgId) {
        if (orgId == 0 || orgId > _nextOrgId) revert OrganisationNotFound(orgId);
        _;
    }

    modifier onlyAdmin() {
        _checkDefaultAdminRole();
        _;
    }

    /// @notice Sets the registry contract that may drive organisation/vault state.
    /// @dev Admin only. Set once after the registry is deployed.
    function setRegistry(address newRegistry) external onlyAdmin {
        if (newRegistry == address(0)) revert InvalidRegistry();
        registry = newRegistry;
        emit RegistrySet(newRegistry);
    }

    // ------------------------------------------------------------------
    // Open views
    // ------------------------------------------------------------------

    function codeOrg(uint32 code) external view override returns (uint256) {
        return _codeOrg[code];
    }

    function codeStatus(uint32 code) external view override returns (CodeStatus) {
        return _codeStatus[code];
    }

    function totalOrganisation() external view override returns (uint248) {
        return _totalOrganisation;
    }

    function organisationOf(uint256 orgId)
        external
        view
        override
        onlyExistingOrg(orgId)
        returns (Organisation memory)
    {
        return _organisations[orgId];
    }

    function organisationSignersOf(uint256 orgId)
        external
        view
        override
        onlyExistingOrg(orgId)
        returns (bytes32[] memory)
    {
        return _organisationSigners[orgId];
    }

    function organisationSigner(uint256 orgId)
        public
        view
        override
        onlyExistingOrg(orgId)
        returns (bytes32)
    {
        return _organisations[orgId].signer;
    }

    function hasOrgRole(uint256 orgId, bytes32 role, address account)
        public
        view
        override
        returns (bool)
    {
        return _roles[orgId][role].hasRole[account];
    }

    // ------------------------------------------------------------------
    // Registry-only mutations
    // ------------------------------------------------------------------

    /// @dev Creates an organisation with `initiator` as its first admin.
    function createOrganisationFor(
        address initiator,
        string calldata name,
        string calldata website,
        bytes32 signer
    ) external override onlyRegistry returns (uint256 orgId) {
        if (signer == bytes32(0)) revert InvalidSigner();
        if (_nextOrgId >= type(uint64).max) revert InvalidOrganisation(_nextOrgId);

        orgId = ++_nextOrgId;
        _totalOrganisation += 1;

        _organisations[orgId] = Organisation({
            id: uint64(orgId),
            name: name,
            status: OrgStatus.Active,
            website: website,
            signer: signer,
            verificationFlag: 0
        });
        _organisationSigners[orgId].push(signer);

        // The creator is the organisation's admin and also holds the two operational roles.
        // Granting them here avoids extra registry -> ClueCode round trips at vault time.
        _roles[orgId][ClueOrgRoles.ORG_ADMIN_ROLE].hasRole[initiator] = true;
        _roles[orgId][ClueOrgRoles.ORG_CONTROLLER_ROLE].hasRole[initiator] = true;
        _roles[orgId][ClueOrgRoles.ORG_UPDATER_ROLE].hasRole[initiator] = true;
        _roles[orgId][ClueOrgRoles.ORG_CONTROLLER_ROLE].adminRole = ClueOrgRoles.ORG_ADMIN_ROLE;
        _roles[orgId][ClueOrgRoles.ORG_UPDATER_ROLE].adminRole = ClueOrgRoles.ORG_ADMIN_ROLE;

        emit OrganisationCreated(orgId, initiator, name, website, signer);
    }

    /// @dev Adds an Ed25519 signer. Only an admin of `orgId` may do so.
    function addOrganisationSignerFor(uint256 orgId, bytes32 signer, address sender)
        external
        override
        onlyRegistry
        onlyExistingOrg(orgId)
    {
        if (!_roles[orgId][ClueOrgRoles.ORG_ADMIN_ROLE].hasRole[sender]) {
            revert UnauthorizedOrgRole(orgId, ClueOrgRoles.ORG_ADMIN_ROLE, sender);
        }
        if (signer == bytes32(0)) revert InvalidOrganisationSigner(orgId, signer);

        bytes32[] storage signers = _organisationSigners[orgId];
        uint256 len = signers.length;
        for (uint256 i = 0; i < len; ++i) {
            if (signers[i] == signer) revert OrganisationSignerAlreadyAdded(orgId, signer);
        }
        signers.push(signer);
        emit OrganisationSignerAdded(orgId, signer, sender);
    }

    /// @dev Sets the verification flag of `orgId`. Only an org admin may do so. The flag
    ///      mirrors the backend Organisation.VerificationFlag field and is informational.
    function setOrganisationVerificationFlag(uint256 orgId, uint32 flag, address sender)
        external
        override
        onlyRegistry
        onlyExistingOrg(orgId)
    {
        if (!_roles[orgId][ClueOrgRoles.ORG_ADMIN_ROLE].hasRole[sender]) {
            revert UnauthorizedOrgRole(orgId, ClueOrgRoles.ORG_ADMIN_ROLE, sender);
        }
        _organisations[orgId].verificationFlag = flag;
        emit OrganisationVerificationFlagSet(orgId, flag, sender);
    }

    /// @dev Sets the admin role of `role` in `orgId`. An org admin may delegate or lock a
    ///      role. The admin role itself cannot be changed, and a locked role stays locked.
    function setOrgRoleAdmin(uint256 orgId, bytes32 role, bytes32 adminRole, address sender)
        external
        override
        onlyRegistry
        onlyExistingOrg(orgId)
    {
        if (!_roles[orgId][ClueOrgRoles.ORG_ADMIN_ROLE].hasRole[sender]) {
            revert UnauthorizedOrgRole(orgId, ClueOrgRoles.ORG_ADMIN_ROLE, sender);
        }
        if (role == ClueOrgRoles.ORG_ADMIN_ROLE) revert RoleLocked(orgId, role);

        OrgRoleData storage roleData = _roles[orgId][role];
        if (roleData.adminRole == ClueOrgRoles.NO_ROLE_ADMIN) revert RoleLocked(orgId, role);

        roleData.adminRole = adminRole;
        emit OrgRoleAdminSet(orgId, role, adminRole, sender);
    }

    /// @dev Grants `role` in `orgId` to `account`, with `sender` as the authorising caller.
    function grantOrgRole(uint256 orgId, bytes32 role, address account, address sender)
        external
        override
        onlyRegistry
    {
        if (orgId == 0 || orgId > _nextOrgId) revert InvalidOrganisation(orgId);

        OrgRoleData storage roleData = _roles[orgId][role];
        bytes32 adminRole = roleData.adminRole;
        if (adminRole == ClueOrgRoles.NO_ROLE_ADMIN) revert RoleLocked(orgId, role);

        bytes32 effectiveAdmin = adminRole == bytes32(0) ? ClueOrgRoles.ORG_ADMIN_ROLE : adminRole;
        if (!_roles[orgId][effectiveAdmin].hasRole[sender]) revert MissingOrgRole(orgId, effectiveAdmin, sender);

        _grantOrgRole(orgId, role, account, sender);
    }

    /// @dev Revokes `role` in `orgId` from `account`, with `sender` as the authorising caller.
    function revokeOrgRole(uint256 orgId, bytes32 role, address account, address sender)
        external
        override
        onlyRegistry
    {
        if (orgId == 0 || orgId > _nextOrgId) revert InvalidOrganisation(orgId);

        OrgRoleData storage roleData = _roles[orgId][role];
        bytes32 adminRole = roleData.adminRole;
        if (adminRole == ClueOrgRoles.NO_ROLE_ADMIN) revert RoleLocked(orgId, role);

        bytes32 effectiveAdmin = adminRole == bytes32(0) ? ClueOrgRoles.ORG_ADMIN_ROLE : adminRole;
        if (!_roles[orgId][effectiveAdmin].hasRole[sender]) revert MissingOrgRole(orgId, effectiveAdmin, sender);

        OrgRoleData storage target = _roles[orgId][role];
        if (target.hasRole[account]) {
            target.hasRole[account] = false;
            emit OrgRoleRevoked(orgId, role, account, sender);
        }
    }

    function _grantOrgRole(uint256 orgId, bytes32 role, address account, address sender) private {
        OrgRoleData storage roleData = _roles[orgId][role];
        if (!roleData.hasRole[account]) {
            roleData.hasRole[account] = true;
            emit OrgRoleGranted(orgId, role, account, sender);
        }
    }

    /// @dev Vaults `code` into `orgId` on behalf of `initiator`, who must own the token.
    ///      The token is transferred internally to this contract; no separate approval is
    ///      required because the caller of the registry is the token owner.
    function vaultCode(uint32 code, uint256 orgId, address initiator) external override onlyRegistry {
        if (orgId == 0 || orgId > _nextOrgId) revert InvalidOrganisation(orgId);
        if (!_roles[orgId][ClueOrgRoles.ORG_UPDATER_ROLE].hasRole[initiator]) {
            revert MissingOrgRole(orgId, ClueOrgRoles.ORG_UPDATER_ROLE, initiator);
        }

        uint256 currentOrg = _codeOrg[code];
        if (currentOrg != 0) revert CodeAlreadyVaulted(code, currentOrg);
        if (_codeOwner(code) != initiator) revert NotCodeOwner(code, initiator);

        _codeOrg[code] = orgId;
        emit CodeVaulted(code, orgId, initiator);

        _transferCode(initiator, address(this), uint256(code));
    }

    /// @dev Releases a Stateless code to `receiver`. The registry must have cleaned up its
    ///      own validator/provider state before calling this function.
    function releaseVaultCode(uint32 code, address receiver, address sender) external override onlyRegistry {
        if (receiver == address(0)) revert ReceiverZeroAddress();

        uint256 orgId = _codeOrg[code];
        if (orgId == 0) revert CodeNotVaulted(code);

        if (
            !_roles[orgId][ClueOrgRoles.ORG_ADMIN_ROLE].hasRole[sender]
                && !_roles[orgId][ClueOrgRoles.ORG_CONTROLLER_ROLE].hasRole[sender]
                && !_roles[orgId][ClueOrgRoles.ORG_UPDATER_ROLE].hasRole[sender]
        ) revert UnauthorizedOrgRole(orgId, bytes32(0), sender);

        CodeStatus current = _codeStatus[code];
        if (current != CodeStatus.Stateless) revert CodeNotStateless(code, current);

        delete _codeOrg[code];
        _codeStatus[code] = CodeStatus.Open;
        emit CodeReleased(code, orgId, receiver);

        _transferCode(address(this), receiver, uint256(code));
    }

    /// @notice Sets the status of a vaulted code on behalf of the registry.
    function setCodeStatus(uint32 code, CodeStatus newStatus) external override onlyRegistry {
        _setCodeStatus(code, newStatus);
    }

    /// @dev Shared status write. Open is reachable only via release. Callers are
    ///      gated by the registry or by {ClueConstrict} in the concrete contract.
    function _setCodeStatus(uint32 code, CodeStatus newStatus) internal {
        if (newStatus == CodeStatus.Open) revert CannotSetOpenViaStatus(code);
        _codeStatus[code] = newStatus;
        emit CodeStatusSet(code, newStatus);
    }

    /// @notice Accepts safe transfers of this contract's own tokens back into the vault.
    function onERC721Received(address, address, uint256, bytes calldata)
        external
        view
        override
        returns (bytes4)
    {
        if (msg.sender != address(this)) revert InvalidERC721Sender(msg.sender);
        return IERC721Receiver.onERC721Received.selector;
    }
}
