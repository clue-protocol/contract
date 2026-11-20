// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import {IClueCode} from "../code/IClueCode.sol";
import {ClueCodeVault} from "./ClueCodeVault.sol";

/**
 * @title ClueProvider
 * @notice ClueRegistry component that holds the on-chain record of each vaulted Provider
 *         Code: its name, target, Profile Server key, its Route links and partitions, and
 *         its live serial count.
 * @dev Builds on ClueCodeVault. A code is first vaulted with `activateCode`, then
 *      configured here as Static or Route. Status moves with the code's state:
 *
 *        - Route code: `_lock` on its first link, `_unlock` when its last link is removed.
 *        - Static code: `_lock` when it first holds a Clue ID (direct or linked),
 *          `_unlock` when it holds none.
 *
 *      Keys: owning a code gives no signing authority. The Ed25519 Profile Server key is
 *      registered here and is unaffected by NFT transfers, because a vaulted NFT cannot
 *      be traded.
 *
 *      Serial counts are fed by the registry's root-update path through
 *      `_applyDirectCount` and `_applyLinkedCount`. This contract only enforces the
 *      ceilings; the validated batch decides what was created or deleted.
 *
 *      NOTE: reference implementation, not audited.
 */
abstract contract ClueProvider is ClueCodeVault {
    using EnumerableSet for EnumerableSet.UintSet;

    // Types

    /// @notice Partition type of a link. None marks a link that does not exist.
    enum PartitionType { None, Domain, Alloc }

    /// @notice Configuration of a vaulted Provider Code.
    struct Provider {
        bool    configured;
        bytes32 profileKey;    // Ed25519 public key of the Profile Server
        bytes32 nameHash;      // keccak256 of the normalized resolved name
        bytes32 networkHash;   // keccak256 of the normalized network (Static only)
        string  name;          // Resolved name, e.g. "USDC Clue Ethereum" or "USDC Clue Global"
        string  network;       // Static only
        string  token;         // Static only; empty if the code targets the network itself
        bytes   tokenAddress;  // Static only; empty if token is empty
    }

    /// @notice Link record (routeCode, staticCode, D, R_link) and its partition.
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

    /// @notice Partition requested when linking.
    /// @dev Domain: set domainStart and domainEnd, allocLimit must be 0.
    ///      Alloc: set allocLimit, domainStart and domainEnd must be 0.
    struct PartitionSpec {
        PartitionType partitionType;
        uint32 domainStart;
        uint32 domainEnd;
        uint32 allocLimit;
    }

    // Constants

    uint32 public constant DIRECT_CEILING_S = 500_000_000;
    uint32 public constant DIRECT_CEILING_A = 300_000_000;
    uint32 public constant DIRECT_CEILING_B = 150_000_000;

    uint32 public constant ROUTE_BUDGET_S = 350_000_000;
    uint32 public constant ROUTE_BUDGET_A = 150_000_000;

    uint256 public constant MAX_NAME_LENGTH = 48;

    // Storage

    /// @dev Code => configuration.
    mapping(uint32 => Provider) internal _provider;

    /// @dev Normalized resolved-name hash => code that owns it. Names are unique across codes.
    mapping(bytes32 => uint32) internal _nameOwner;

    /// @notice Static code => Route code it is linked to, or 0.
    mapping(uint32 => uint32) public routeOf;

    /// @notice Static code => Route code its controller has approved for linking, or 0.
    mapping(uint32 => uint32) public approvedRoute;

    /// @dev Route code => (Static code => link).
    mapping(uint32 => mapping(uint32 => Link)) internal _links;

    /// @dev Route code => linked Static codes.
    mapping(uint32 => EnumerableSet.UintSet) internal _linkedStatics;

    /// @dev Route code => network hash => already used by a link of this Route.
    mapping(uint32 => mapping(bytes32 => bool)) internal _networkUsed;

    /// @notice Route code => sum of R_link over its links.
    mapping(uint32 => uint32) public routeAllocated;

    /// @notice Static code => live Clue IDs in its direct zone.
    mapping(uint32 => uint32) public directCount;

    /// @notice Static code => live Linked Serials across its routed zone.
    mapping(uint32 => uint32) public linkedTotal;

    // Errors

    error NotConfigured(uint32 code);
    error AlreadyConfigured(uint32 code);
    error WrongProviderType(uint32 code, ProviderType expected);
    error InvalidName();
    error InvalidNetwork();
    error InvalidToken();
    error InvalidKey();
    error NameTaken(string name);
    error NotApproved(uint32 staticCode, uint32 routeCode);
    error AlreadyLinked(uint32 staticCode);
    error NetworkAlreadyLinked(uint32 routeCode);
    error NotLinked(uint32 routeCode, uint32 staticCode);
    error InvalidLinkShare();
    error RouteBudgetExceeded(uint32 routeCode, uint256 requested, uint256 budget);
    error InvalidPartition();
    error LinkedSerialsRemain(uint32 routeCode, uint32 staticCode, uint32 count);
    error StillLinked(uint32 code);
    error CountUnderflow(uint32 code);
    error DirectCeilingExceeded(uint32 code, uint256 requested, uint32 ceiling);
    error PartitionCapacityExceeded(uint32 routeCode, uint32 staticCode, uint256 requested, uint32 capacity);

    // Events

    event ProviderConfigured(uint32 indexed code, ProviderType providerType, string name, bytes32 profileKey);
    event ProfileKeyRotated(uint32 indexed code, bytes32 oldKey, bytes32 newKey);
    event RouteApproved(uint32 indexed staticCode, uint32 indexed routeCode);
    event Linked(uint32 indexed routeCode, uint32 indexed staticCode, Link record);
    event Unlinked(uint32 indexed routeCode, uint32 indexed staticCode);
    event DirectCountUpdated(uint32 indexed code, uint32 count);
    event LinkedCountUpdated(uint32 indexed routeCode, uint32 indexed staticCode, uint32 count);

    // Constructor

    constructor(IClueCode clueCode_) ClueCodeVault(clueCode_) {}

    // Configuration

    /**
     * @notice Configures a vaulted Static code: its name, its single target and its
     *         Profile Server key.
     * @dev The resolved name is `baseName + " " + network`, so the network is always part
     *      of a Static name. A token target needs a token address; a network target has
     *      none. Names are ASCII, case-insensitive and unique across all codes.
     * @param code Provider Code, vaulted as Static. Caller must be its controller.
     * @param baseName Provider name without the network, e.g. "USDC Clue".
     * @param network Blockchain network, e.g. "Ethereum".
     * @param token Token under the network, or empty for the network itself.
     * @param tokenAddress Token contract address bytes, or empty if `token` is empty.
     * @param profileKey Ed25519 public key of the Profile Server.
     */
    function configureStatic(
        uint32 code,
        string calldata baseName,
        string calldata network,
        string calldata token,
        bytes calldata tokenAddress,
        bytes32 profileKey
    ) external {
        _requireConfigurable(code, ProviderType.Static);
        if (profileKey == bytes32(0)) revert InvalidKey();

        bytes32 networkHash = keccak256(_normalize(network));
        if (networkHash == keccak256("global")) revert InvalidNetwork();

        bool hasToken = bytes(token).length != 0;
        if (hasToken) {
            _normalize(token); // validates characters only
            if (tokenAddress.length == 0) revert InvalidToken();
        } else if (tokenAddress.length != 0) {
            revert InvalidToken();
        }

        string memory resolved = string.concat(baseName, " ", network);
        bytes32 nameHash = _claimName(code, resolved);

        _provider[code] = Provider({
            configured: true,
            profileKey: profileKey,
            nameHash: nameHash,
            networkHash: networkHash,
            name: resolved,
            network: network,
            token: token,
            tokenAddress: tokenAddress
        });

        emit ProviderConfigured(code, ProviderType.Static, resolved, profileKey);
    }

    /**
     * @notice Configures a vaulted Route code: its name and its Profile Server key.
     * @dev The resolved name is `baseName + " Global"`. A Route code has exactly one
     *      delegated Profile Server, which is the holder of `profileKey`.
     * @param code Provider Code, vaulted as Route. Caller must be its controller.
     * @param baseName Provider name without "Global", e.g. "USDC Clue".
     * @param profileKey Ed25519 public key of the delegated Profile Server.
     */
    function configureRoute(uint32 code, string calldata baseName, bytes32 profileKey) external {
        _requireConfigurable(code, ProviderType.Route);
        if (profileKey == bytes32(0)) revert InvalidKey();

        string memory resolved = string.concat(baseName, " Global");
        bytes32 nameHash = _claimName(code, resolved);

        Provider storage p = _provider[code];
        p.configured = true;
        p.profileKey = profileKey;
        p.nameHash = nameHash;
        p.name = resolved;

        emit ProviderConfigured(code, ProviderType.Route, resolved, profileKey);
    }

    /**
     * @notice Replaces the Profile Server key of a configured code.
     * @dev Controller only. The old key stops being the registered key immediately.
     */
    function rotateProfileKey(uint32 code, bytes32 newKey) external {
        _requireController(code);
        Provider storage p = _provider[code];
        if (!p.configured) revert NotConfigured(code);
        if (newKey == bytes32(0)) revert InvalidKey();

        bytes32 oldKey = p.profileKey;
        p.profileKey = newKey;
        emit ProfileKeyRotated(code, oldKey, newKey);
    }

    // Linking

    /**
     * @notice The controller of a Static code approves one Route code to link it.
     * @dev A link takes over the routed zone of the Static code, so the Static controller
     *      must consent. Approving 0 cancels. The approval is consumed by `link`.
     * @param staticCode Configured Static code. Caller must be its controller.
     * @param routeCode Route code to approve, or 0 to cancel.
     */
    function approveRoute(uint32 staticCode, uint32 routeCode) external {
        _requireController(staticCode);
        _requireConfiguredAs(staticCode, ProviderType.Static);

        approvedRoute[staticCode] = routeCode;
        emit RouteApproved(staticCode, routeCode);
    }

    /**
     * @notice Links a Static code to a Route code and delegates a partition of it.
     * @dev Caller must be the Route code's controller. The Static code need not be empty:
     *      its routed zone [D, D + rLink - 1] is empty by construction. The first link
     *      locks the Route code (Stateless to Active). The Static code's status does not
     *      change. The partition type cannot be changed afterwards.
     * @param routeCode Configured Route code.
     * @param staticCode Configured Static code on a network not yet linked to this Route.
     *                   Its controller must have called `approveRoute` for `routeCode`.
     * @param rLink Serials of the Route budget given to this link.
     * @param spec Domain partition (interval inside the routed zone) or Alloc partition
     *             (explicit limit up to `rLink`).
     */
    function link(uint32 routeCode, uint32 staticCode, uint32 rLink, PartitionSpec calldata spec) external {
        _requireController(routeCode);
        _requireConfiguredAs(routeCode, ProviderType.Route);
        _requireConfiguredAs(staticCode, ProviderType.Static);

        if (approvedRoute[staticCode] != routeCode) revert NotApproved(staticCode, routeCode);
        if (routeOf[staticCode] != 0) revert AlreadyLinked(staticCode);

        bytes32 networkHash = _provider[staticCode].networkHash;
        if (_networkUsed[routeCode][networkHash]) revert NetworkAlreadyLinked(routeCode);

        if (rLink == 0) revert InvalidLinkShare();
        uint256 newTotal = uint256(routeAllocated[routeCode]) + rLink;
        if (newTotal > routeBudget(routeCode)) {
            revert RouteBudgetExceeded(routeCode, newTotal, routeBudget(routeCode));
        }

        uint32 d = directCeiling(staticCode);
        Link memory l = Link({
            partitionType: spec.partitionType,
            d: d,
            rLink: rLink,
            domainStart: spec.domainStart,
            domainEnd: spec.domainEnd,
            allocLimit: spec.allocLimit,
            capacity: _partitionCapacity(d, rLink, spec),
            linkedCount: 0
        });

        // Effects.
        delete approvedRoute[staticCode];
        routeOf[staticCode] = routeCode;
        _networkUsed[routeCode][networkHash] = true;
        routeAllocated[routeCode] = uint32(newTotal);
        _links[routeCode][staticCode] = l;
        _linkedStatics[routeCode].add(uint256(staticCode));

        // The first link makes the Route code Active.
        if (_vault[routeCode].status == CodeStatus.Stateless) _lock(routeCode);

        emit Linked(routeCode, staticCode, l);
    }

    /// @dev Validates a partition against the routed zone [d, d + rLink - 1] and returns
    ///      its capacity: the Domain size, or the Alloc limit.
    function _partitionCapacity(uint32 d, uint32 rLink, PartitionSpec calldata spec)
        private
        pure
        returns (uint32)
    {
        if (spec.partitionType == PartitionType.Domain) {
            if (spec.allocLimit != 0) revert InvalidPartition();
            if (spec.domainStart > spec.domainEnd) revert InvalidPartition();
            if (spec.domainStart < d) revert InvalidPartition();
            if (uint256(spec.domainEnd) > uint256(d) + rLink - 1) revert InvalidPartition();
            return spec.domainEnd - spec.domainStart + 1;
        }
        if (spec.partitionType == PartitionType.Alloc) {
            if (spec.domainStart != 0 || spec.domainEnd != 0) revert InvalidPartition();
            if (spec.allocLimit == 0 || spec.allocLimit > rLink) revert InvalidPartition();
            return spec.allocLimit;
        }
        revert InvalidPartition();
    }

    /**
     * @notice Removes a link once no Linked Serial remains.
     * @dev Route controller only: the Static code cannot delete Linked Serials, so the
     *      Route side always has to act first. All Linked Serials must already be
     *      deleted (in the next Provider Version), so the link has no state. The budget
     *      is returned to the Route, the network is freed, and when the Route code has
     *      no link left it is unlocked (Active to Stateless).
     */
    function unlink(uint32 routeCode, uint32 staticCode) external {
        _requireController(routeCode);
        Link storage l = _links[routeCode][staticCode];
        if (l.partitionType == PartitionType.None) revert NotLinked(routeCode, staticCode);

        if (l.linkedCount != 0) revert LinkedSerialsRemain(routeCode, staticCode, l.linkedCount);

        routeAllocated[routeCode] -= l.rLink;
        _networkUsed[routeCode][_provider[staticCode].networkHash] = false;
        delete routeOf[staticCode];
        delete _links[routeCode][staticCode];
        _linkedStatics[routeCode].remove(uint256(staticCode));

        // No link left: the Route code returns to Stateless.
        if (_linkedStatics[routeCode].length() == 0 && _vault[routeCode].status == CodeStatus.Active) {
            _unlock(routeCode);
        }

        emit Unlinked(routeCode, staticCode);
    }

    // Serial counts (called by the registry's root-update path)

    /**
     * @dev Applies a batch's net effect on the direct zone of a Static code: `created`
     *      CREATEs and `deleted` DELETEs. Rejects a count below zero or above the grade
     *      ceiling, then syncs the code's status.
     */
    function _applyDirectCount(uint32 code, uint32 created, uint32 deleted) internal {
        _requireConfiguredAs(code, ProviderType.Static);

        uint256 current = directCount[code];
        if (deleted > current + created) revert CountUnderflow(code);
        uint256 next = current + created - deleted;

        uint32 ceiling = directCeiling(code);
        if (next > ceiling) revert DirectCeilingExceeded(code, next, ceiling);

        directCount[code] = uint32(next);
        emit DirectCountUpdated(code, uint32(next));

        _syncStaticStatus(code);
    }

    /**
     * @dev Applies a batch's net effect on the Linked Serials of one link. The count is
     *      held against the partition capacity (Domain size or Alloc limit). It also
     *      updates the Static code's linked total and status, because a Linked Serial
     *      appears under the Static code too.
     */
    function _applyLinkedCount(uint32 routeCode, uint32 staticCode, uint32 created, uint32 deleted) internal {
        Link storage l = _links[routeCode][staticCode];
        if (l.partitionType == PartitionType.None) revert NotLinked(routeCode, staticCode);

        uint256 current = l.linkedCount;
        if (deleted > current + created) revert CountUnderflow(staticCode);
        uint256 next = current + created - deleted;
        if (next > l.capacity) revert PartitionCapacityExceeded(routeCode, staticCode, next, l.capacity);

        l.linkedCount = uint32(next);
        linkedTotal[staticCode] = uint32(uint256(linkedTotal[staticCode]) + created - deleted);
        emit LinkedCountUpdated(routeCode, staticCode, uint32(next));

        _syncStaticStatus(staticCode);
    }

    /// @dev Static code is Active while it holds any Clue ID (direct or linked) and
    ///      Stateless otherwise.
    function _syncStaticStatus(uint32 code) private {
        bool holdsState = directCount[code] != 0 || linkedTotal[code] != 0;
        CodeStatus status = _vault[code].status;

        if (holdsState && status == CodeStatus.Stateless) _lock(code);
        else if (!holdsState && status == CodeStatus.Active) _unlock(code);
    }

    // Release hook

    /// @dev A linked code cannot leave the vault. On release the configuration is
    ///      cleared and the name is freed. A Stateless Route code has no links by
    ///      construction; a Stateless Static code may still be linked, so it is checked.
    function _beforeRelease(uint32 code) internal override {
        if (routeOf[code] != 0 || _linkedStatics[code].length() != 0) revert StillLinked(code);

        Provider storage p = _provider[code];
        if (p.configured) {
            delete _nameOwner[p.nameHash];
            delete _provider[code];
        }
        delete approvedRoute[code];
    }

    // Internal helpers

    /// @dev Caller is the controller, the code has the expected type and is not configured.
    function _requireConfigurable(uint32 code, ProviderType expected) private view {
        _requireController(code);
        if (_vault[code].providerType != expected) revert WrongProviderType(code, expected);
        if (_provider[code].configured) revert AlreadyConfigured(code);
    }

    /// @dev The code is vaulted, configured and of the expected type.
    function _requireConfiguredAs(uint32 code, ProviderType expected) internal view {
        _requireControlled(code);
        if (_vault[code].providerType != expected) revert WrongProviderType(code, expected);
        if (!_provider[code].configured) revert NotConfigured(code);
    }

    /// @dev Reserves the normalized resolved name for `code`. Reverts if another code has it.
    function _claimName(uint32 code, string memory resolved) private returns (bytes32 nameHash) {
        nameHash = keccak256(_normalize(resolved));
        if (_nameOwner[nameHash] != 0) revert NameTaken(resolved);
        _nameOwner[nameHash] = code;
    }

    /**
     * @dev Normalizes a name: ASCII letters fold to lowercase; digits and single inner
     *      spaces are kept; anything else reverts. A name has 1 to MAX_NAME_LENGTH bytes
     *      and no leading, trailing or doubled space. The rules are a starting point and
     *      may change before launch.
     */
    function _normalize(string memory s) internal pure returns (bytes memory out) {
        bytes memory b = bytes(s);
        uint256 len = b.length;
        if (len == 0 || len > MAX_NAME_LENGTH) revert InvalidName();

        out = new bytes(len);
        bool prevSpace = true; // true at the start, so a leading space is rejected
        for (uint256 i; i < len; ++i) {
            bytes1 c = b[i];
            if (c >= 0x41 && c <= 0x5A) c = bytes1(uint8(c) + 32); // A-Z to a-z

            if (c == 0x20) {
                if (prevSpace) revert InvalidName();
                prevSpace = true;
            } else if ((c >= 0x61 && c <= 0x7A) || (c >= 0x30 && c <= 0x39)) {
                prevSpace = false;
            } else {
                revert InvalidName();
            }
            out[i] = c;
        }
        if (prevSpace) revert InvalidName(); // trailing space
    }

    // Views

    /// @notice Static direct ceiling D for the grade of `code`.
    function directCeiling(uint32 code) public view returns (uint32) {
        IClueCode.Grade g = clueCode.gradeOf(uint256(code));
        if (g == IClueCode.Grade.S) return DIRECT_CEILING_S;
        if (g == IClueCode.Grade.A) return DIRECT_CEILING_A;
        return DIRECT_CEILING_B;
    }

    /// @notice Route budget R for the grade of `code`. Zero for grade B.
    function routeBudget(uint32 code) public view returns (uint32) {
        IClueCode.Grade g = clueCode.gradeOf(uint256(code));
        if (g == IClueCode.Grade.S) return ROUTE_BUDGET_S;
        if (g == IClueCode.Grade.A) return ROUTE_BUDGET_A;
        return 0;
    }

    /// @notice Returns true if `code` has been configured.
    function isConfigured(uint32 code) external view returns (bool) {
        return _provider[code].configured;
    }

    /// @notice Returns the resolved name of `code`. Reverts if it is not configured.
    function nameOf(uint32 code) external view returns (string memory) {
        if (!_provider[code].configured) revert NotConfigured(code);
        return _provider[code].name;
    }

    /// @notice Returns the registered Profile Server key. Reverts if not configured.
    function profileKeyOf(uint32 code) external view returns (bytes32) {
        if (!_provider[code].configured) revert NotConfigured(code);
        return _provider[code].profileKey;
    }

    /// @notice Returns the target of a Static code. Reverts if it is not a configured Static code.
    function targetOf(uint32 code)
        external
        view
        returns (string memory network, string memory token, bytes memory tokenAddress)
    {
        _requireConfiguredAs(code, ProviderType.Static);
        Provider storage p = _provider[code];
        return (p.network, p.token, p.tokenAddress);
    }

    /// @notice Returns the code that owns a resolved name, or 0. Matching ignores case.
    function codeOfName(string calldata resolvedName) external view returns (uint32) {
        return _nameOwner[keccak256(_normalize(resolvedName))];
    }

    /// @notice Returns the link record between a Route code and a Static code.
    function linkOf(uint32 routeCode, uint32 staticCode) external view returns (Link memory) {
        Link memory l = _links[routeCode][staticCode];
        if (l.partitionType == PartitionType.None) revert NotLinked(routeCode, staticCode);
        return l;
    }

    /// @notice Returns the Static codes linked to a Route code.
    function linkedStaticsOf(uint32 routeCode) external view returns (uint256[] memory) {
        return _linkedStatics[routeCode].values();
    }

    /// @notice Returns how much of the Route budget is still unallocated.
    function routeBudgetLeft(uint32 routeCode) external view returns (uint32) {
        return routeBudget(routeCode) - routeAllocated[routeCode];
    }
}