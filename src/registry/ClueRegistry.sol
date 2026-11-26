// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import {EIP712} from "@openzeppelin/contracts/utils/cryptography/EIP712.sol";

import {ClueCodeOrganisation} from "./ClueCodeOrganisation.sol";
import {ClueValidators} from "./ClueValidators.sol";
import {ClueMultiSig} from "./ClueMultiSig.sol";
import {ClueAccounts} from "./ClueAccounts.sol";
import {ClueProvider} from "./ClueProvider.sol";

/**
 * @title ClueRegistry
 * @notice The Clue Protocol registry. Combines code vaulting and organisations
 *         ({ClueCodeOrganisation}), validators ({ClueValidators}), validator approval
 *         ({ClueMultiSig}) and the account leaf record ({ClueAccounts}).
 * @dev No constructor yet. Marked `abstract`
 *      until `onERC721Received` and the constructor (which must set the immutable
 *      `clueCode`) are written; then drop `abstract`.
 *
 *      Inheritance (linearized): ClueMultiSig already inherits ClueValidators, which
 *      inherits ClueCodeOrganisation; ClueProvider also inherits ClueCodeOrganisation.
 */
abstract contract ClueRegistry is ClueMultiSig, ClueProvider, ClueAccounts, AccessControl {
    using EnumerableSet for EnumerableSet.Bytes32Set;

    // Types

    enum PartitionType { None, Domain, Alloc }

    struct Link {
        PartitionType partitionType;
        uint32 d;            // Static direct ceiling at link time
        uint32 rLink;        // Share of the Route budget given to this link
        uint32 domainStart;  // Domain partition only
        uint32 domainEnd;    // Domain partition only
        uint32 allocLimit;   // Alloc partition only
        uint32 capacity;     // Domain size, or allocLimit: the most Linked Serials allowed
        uint32 linkedCount;  // Live Linked Serials
    }

    struct UpdateAccount {
        bytes32 totalHash; // hash(leaf1+leaf2+....)
        uint16 count;
        bytes[] accounts;
    }

    struct StaticProposalV {
        bytes root;
        bytes prevRoot;
        uint32 prevVersion;
        uint32 updateVersion;
        uint32 providerCode;
        uint32 updatedSerialCount;
        uint64 commitSequence;
        uint64 commitSequenceEnd;
        UpdateAccount accounts;
    }

    struct RouteProposalV {
        bytes root;
        bytes prevRoot;
        uint32 prevVersion;
        uint32 updateVersion;
        uint32 updatedSerialCount;
        uint32 providerCode;
        uint64 commitSequence;
        uint64 commitSequenceEnd;
        UpdateAccount account;
    }

    struct StaticProposalLink {
        bytes root;
        bytes prevRoot;
        uint32 prevVersion;
        uint32 updateVersion;
        uint32 updatedSerialCount;
        uint32 providerCode;
        uint64 commitSequence;
        uint64 commitSequenceEnd;
        Link staticLinks;
        UpdateAccount account;
    }

    struct Links {
        StaticProposalLink[] statics;
    }

    // Constants

    /// @dev Role allowed to run registry-wide administration (separate from the per-
    ///      organisation roles in {ClueCodeOrganisation}).
    bytes32 public constant REGISTRY_ADMIN_ROLE = keccak256("REGISTRY_ADMIN_ROLE");

    error CodeNotActive(uint32 code);
    error NoConfirmedRoot(uint32 code);

    struct ResolveResult {
        bytes32 root;               // Current OnChainSuccess root
        uint64  version;            // Current OnChainSuccess version
        uint8   providerType;       // 0 - Static, 1 - Route
        bytes32 organisationSigner; // Organisation signer
    }

    struct StaticResolveResult {
        bytes32 root;               // Current Version
        uint64  version;            // Current OnChainSuccess version
        string  network;            // ethereum, base etc - if Provider Code is Route the network should be "Network-Route"
        uint8   providerType;       // 0 - Static, 1 - Route
        bytes32 organisationSigner; // Organisation signer
    }

    struct RouteResolveResult {
        bytes32 root;               // Current Version
        uint64  version;            // Current OnChainSuccess version
        uint8   providerType;       // 0 - Static, 1 - Route
        bytes32 organisationSigner; // Organisation signer
        uint32[] linkedCode;
    }

    // Storage

    /// @dev Provider Code => the update currently open, if any.
    mapping(uint32 code => StaticProposalV) internal _proposal;

    /// @dev Provider Code => SMT root accepted by the validators. Zero until the first
    ///      update is confirmed. This is the root passed to {ClueAccounts} checks.
    mapping(uint32 code => bytes32 root) internal _acceptedRoot;

    /// @dev Provider Code => last confirmed version.
    mapping(uint32 code => uint32 version) internal _confirmedVersion;

    /// @dev Provider Code => network identifier, set once at activation. Empty for Route.
    mapping(uint32 code => string) internal _network;

    /// @dev Route Provider Code => linked static Provider Codes.
    mapping(uint32 code => uint32[]) internal _linkedStaticCodes;

    // Organisations

    /// @notice Creates an organisation, vaults `code` into it, registers a validator,
    ///         and links that validator to the code in one transaction.
    /// @dev The caller becomes the organisation's admin/initiator and the validator's
    ///      account/authorizer. The organisation name is also used as the validator name.
    ///      The caller must own the ClueCode token `code` and have approved this contract
    ///      for it; the ERC-721 transfer reverts otherwise.
    /// @param code                 The uint32 token ID to vault.
    /// @param orgName              Organisation name; also used as the validator name.
    /// @param website              Organisation website.
    /// @param signer               Ed25519 public key of the organisation's Profile Server.
    /// @param validatorDescription Validator profile/summary, may be empty.
    /// @return orgId       The id of the new organisation.
    /// @return validatorId The id of the newly registered validator.
    function createOrganisationVaultCodeAndRegisterValidator(
        uint32 code,
        string calldata orgName,
        string calldata website,
        bytes32 signer,
        string calldata validatorDescription
    ) external nonReentrant returns (uint256 orgId, uint256 validatorId) {
        orgId = createOrganisation(orgName, website, signer);
        _grantOrgRole(orgId, ORG_CONTROLLER_ROLE, msg.sender);
        _vaultCode(code, orgId);

        validatorId = _registerValidator(orgName, validatorDescription, msg.sender);
        _linkValidatorToCode(code, validatorId);
    }

    /// @notice Returns the trusted on-chain anchor for `code`, for off-chain resolve comparison.
    /// @dev View, no inclusion proof and no Ed25519 verification. Reverts on any failure.
    ///      Leaf authenticity/binding is verified client-side: the organisation Ed25519 signer
    ///      signs account data, AccountHash commits to that signed account data, and
    ///      leaf = hash(position, AccountHash). On-chain resolve therefore only checks that the
    ///      leaf/hash is registered in ClueAccounts; it does not prove code/root membership.
    /// @param code Provider Code to resolve.
    /// @param leaf Raw account leaf payload or 32-byte leaf hash. Must be registered.
    ///             Format-dependent checks (code/serial/zone/Damm) are TODO until leaf layout is frozen.
    function resolve(uint32 code, bytes calldata leaf)
        external
        view
        returns (ResolveResult memory)
    {
        (uint256 orgId, uint32 version) = _resolveAnchor(code, leaf);
        return ResolveResult({
            root: _acceptedRoot[code],
            version: uint64(version),
            providerType: uint8(_providers[code].providerType),
            organisationSigner: _organisations[orgId].signer
        });
    }

    /// @notice Returns the on-chain resolve anchor for a Static Provider Code.
    /// @dev Same trust model as {resolve}: on-chain checks anchor + registered leaf/hash;
    ///      client verifies Ed25519-signed account data binding off-chain.
    function resolveStaticCode(uint32 code, bytes calldata leaf)
        external
        view
        returns (StaticResolveResult memory)
    {
        (uint256 orgId, uint32 version) = _resolveAnchor(code, leaf);
        ProviderType pType = _providers[code].providerType;
        if (pType != ProviderType.Static) revert InvalidProviderType();

        return StaticResolveResult({
            root: _acceptedRoot[code],
            version: uint64(version),
            network: _network[code],
            providerType: uint8(pType),
            organisationSigner: _organisations[orgId].signer
        });
    }

    /// @notice Returns the on-chain resolve anchor for a Route Provider Code.
    /// @dev Same trust model as {resolve}: on-chain checks anchor + registered leaf/hash;
    ///      client verifies Ed25519-signed account data binding off-chain.
    function resolveRouteCode(uint32 code, bytes calldata leaf)
        external
        view
        returns (RouteResolveResult memory)
    {
        (uint256 orgId, uint32 version) = _resolveAnchor(code, leaf);
        ProviderType pType = _providers[code].providerType;
        if (pType != ProviderType.Route) revert InvalidProviderType();

        return RouteResolveResult({
            root: _acceptedRoot[code],
            version: uint64(version),
            providerType: uint8(pType),
            organisationSigner: _organisations[orgId].signer,
            linkedCode: _linkedStaticCodes[code]
        });
    }

    /// @dev Resolves diamond override for code-release cleanup.
    function _cleanupCodeOnRelease(uint32 code, uint256 orgId) internal virtual override(ClueValidators, ClueProvider) {
        super._cleanupCodeOnRelease(code, orgId);
    }

    /// @dev Shared resolve checks. Returns orgId and confirmed version.
    function _resolveAnchor(uint32 code, bytes calldata leaf) internal view returns (uint256 orgId, uint32 version) {
        // 1. Code is vaulted.
        orgId = _codeOrg[code];
        if (orgId == 0) revert CodeNotVaulted(code);

        // 2. Code is Active or Archived (Archived: active but no longer updates).
        CodeStatus status = code_status[uint64(code)];
        if (status != CodeStatus.Active && status != CodeStatus.Archived) revert CodeNotActive(code);

        // 3. A confirmed root exists.
        version = _confirmedVersion[code];
        if (version == 0) revert NoConfirmedRoot(code);

        // 4. Leaf checks (cheap, no tree work). Client verifies Ed25519-signed account data
        //    binding off-chain; on-chain only requires the leaf/hash to be registered.
        //    Format-dependent checks (code/serial/zone/Damm) require a frozen leaf layout.
        if (leaf.length == 0) revert InvalidLeaf(LeafCheck.EmptyLeaf);
        bytes32 leafHash;
        if (leaf.length == 32) {
            bytes calldata l = leaf;
            assembly {
                leafHash := calldataload(l.offset)
            }
        } else {
            leafHash = keccak256(leaf);
        }
        if (!_registered[leafHash]) revert LeafNotRegistered(leafHash);
    }
}