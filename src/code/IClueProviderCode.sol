// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {IERC721Receiver} from "@openzeppelin/contracts/token/ERC721/IERC721Receiver.sol";
import {IClueCode} from "./IClueCode.sol";
import {IClueConstrict} from "../zlibs/constrict/IClueConstrict.sol";

/**
 * @title IClueProviderCode
 * @notice Single ERC-165 interface for the Provider Code contract (ClueCode).
 *         It keeps the NFT surface via {IClueCode}, declares the constrict layer via
 *         {IClueConstrict}, and carries the organisation/vault surface inline. Use
 *         `type(IClueProviderCode).interfaceId` for ERC-165 checks and for
 *         {ClueConstrict} role registrations/grants.
 */
interface IClueProviderCode is IClueCode, IClueConstrict, IERC721Receiver {
    enum CodeStatus { Open, Stateless, Active, Archived }
    enum OrgStatus { None, Active }

    /// @notice On-chain organisation record, mirroring the backend Organisation model.
    struct Organisation {
        uint64 id;                   // json:"id"
        string name;                 // json:"name"
        OrgStatus status;            // json:"status"
        string website;              // json:"website"
        bytes32 signer;              // json:"signer"; first Ed25519 public key
        uint32 verificationFlag;     // json:"verificationFlag"
    }

    // Events

    event RegistrySet(address indexed registry);
    event OrganisationCreated(
        uint256 indexed orgId,
        address indexed initiator,
        string name,
        string website,
        bytes32 signer
    );
    event OrganisationSignerAdded(uint256 indexed orgId, bytes32 indexed signer, address indexed sender);
    event OrganisationVerificationFlagSet(uint256 indexed orgId, uint32 flag, address indexed sender);
    event OrgRoleGranted(uint256 indexed orgId, bytes32 indexed role, address indexed account, address sender);
    event OrgRoleRevoked(uint256 indexed orgId, bytes32 indexed role, address indexed account, address sender);
    event OrgRoleAdminSet(uint256 indexed orgId, bytes32 indexed role, bytes32 indexed adminRole, address sender);
    event CodeVaulted(uint32 indexed code, uint256 indexed orgId, address indexed initiator);
    event CodeReleased(uint32 indexed code, uint256 indexed orgId, address indexed receiver);
    event CodeStatusSet(uint32 indexed code, CodeStatus newStatus);

    // Errors

    error NotRegistry(address caller);
    error InvalidRegistry();
    error NotCodeOwner(uint32 code, address caller);
    error OrganisationNotFound(uint256 orgId);
    error InvalidOrganisation(uint256 orgId);
    error CodeAlreadyVaulted(uint32 code, uint256 existingOrgId);
    error CodeNotVaulted(uint32 code);
    error UnauthorizedOrgRole(uint256 orgId, bytes32 role, address account);
    error RoleLocked(uint256 orgId, bytes32 role);
    error MissingOrgRole(uint256 orgId, bytes32 adminRole, address account);
    error InvalidSigner();
    error InvalidOrganisationSigner(uint256 orgId, bytes32 signer);
    error OrganisationSignerAlreadyAdded(uint256 orgId, bytes32 signer);
    error CodeNotStateless(uint32 code, CodeStatus current);
    error CannotSetOpenViaStatus(uint32 code);
    error ReceiverZeroAddress();
    error InvalidERC721Sender(address sender);

    // Open views

    /// @notice Address allowed to drive organisation/vault state on behalf of users.
    function registry() external view returns (address);

    /// @notice Organisation that owns `code`; zero when the code is not vaulted.
    function codeOrg(uint32 code) external view returns (uint256);

    /// @notice Lifecycle status of `code`.
    function codeStatus(uint32 code) external view returns (CodeStatus);

    /// @notice Number of organisations created so far.
    function totalOrganisation() external view returns (uint248);

    /// @notice Stored record of `orgId`.
    function organisationOf(uint256 orgId) external view returns (Organisation memory);

    /// @notice All Ed25519 signers of `orgId`, in insertion order.
    function organisationSignersOf(uint256 orgId) external view returns (bytes32[] memory);

    /// @notice First Ed25519 signer of `orgId` (the backwards-compatible signer field).
    function organisationSigner(uint256 orgId) external view returns (bytes32);

    /// @notice True when `account` holds `role` in `orgId`.
    function hasOrgRole(uint256 orgId, bytes32 role, address account) external view returns (bool);

    // Administration

    /// @notice Sets the registry that may drive organisation/vault state. Admin only.
    function setRegistry(address newRegistry) external;

    /// @notice Grants an account the constrict role for a registered interface. Admin only.
    function grantConstrictRole(bytes4 interfaceId, address account) external;

    // Registry-only mutations. `sender`/`initiator` is the user who called the registry.

    function createOrganisationFor(
        address initiator,
        string calldata name,
        string calldata website,
        bytes32 signer
    ) external returns (uint256 orgId);

    function addOrganisationSignerFor(uint256 orgId, bytes32 signer, address sender) external;

    function setOrganisationVerificationFlag(uint256 orgId, uint32 flag, address sender) external;

    function grantOrgRole(uint256 orgId, bytes32 role, address account, address sender) external;

    function revokeOrgRole(uint256 orgId, bytes32 role, address account, address sender) external;

    /// @notice Sets the admin role of `role`, or locks it with NO_ROLE_ADMIN.
    function setOrgRoleAdmin(uint256 orgId, bytes32 role, bytes32 adminRole, address sender) external;

    function vaultCode(uint32 code, uint256 orgId, address initiator) external;

    function releaseVaultCode(uint32 code, address receiver, address sender) external;

    function setCodeStatus(uint32 code, CodeStatus newStatus) external;

    /// @notice Sets the status of `code` from a granted sibling contract.
    /// @dev Gated by {ClueConstrict.hasConstrictRole} in the implementation.
    function setProviderCodeStatus(uint32 code, CodeStatus newStatus) external;
}

/**
 * @title ClueOrgRoles
 * @notice Role identifiers used by {IClueProviderCode}. Compile-time constants, shared
 *         by {ClueCode} and the registry components so the values cannot drift apart.
 */
library ClueOrgRoles {
    /// @dev The default/admin role: adminRole == 0 resolves to this role.
    bytes32 internal constant ORG_ADMIN_ROLE = bytes32(0);

    /// @dev May manage validators and release a vaulted code.
    bytes32 internal constant ORG_CONTROLLER_ROLE = keccak256("ORG_CONTROLLER_ROLE");

    /// @dev May vault a code into its organisation.
    bytes32 internal constant ORG_UPDATER_ROLE = keccak256("ORG_UPDATER_ROLE");

    /// @dev Sentinel admin role: a role with this admin can never be granted publicly.
    bytes32 internal constant NO_ROLE_ADMIN = bytes32(type(uint256).max);
}
