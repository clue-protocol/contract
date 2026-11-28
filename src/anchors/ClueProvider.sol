// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {ClueOrgRoles} from "../code/IClueProviderCode.sol";
import {ClueCodeOrgRef} from "./ClueCodeOrgRef.sol";

/// @title ClueProvider
/// @notice Stores the *specification* of each vaulted Provider Code: what the code is, not what
///         it is doing. Everything here is written once at registration and read many times.
///
/// @dev    Scope, following the Provider Code Specification:
///
///         STORED HERE (static specification, mirrors pack.MetaData)
///           - provider type            Static | Route            (Section 2, never both)
///           - grade                    S | A | B                 (Section 5, derived from digit length)
///           - resolved provider name   unique across all codes   (Section 4)
///           - canonical namespace      unique Clue-ID suffix, e.g. "binance" or "binance.eth"
///           - static target            network / symbol / token / token address (Section 4)
///           - createdAt
///           - grade ceilings D and R   pure functions of grade   (Section 3.5.1)
///
///         NOT STORED HERE (live / active state lives in ClueCode, the registry only holds a reference)
///           - Ed25519 public key / Profile Server key -> Organisation
///           - code status (Open / Stateless / Active)  -> ClueCode
///           - code -> organisation, roles, initiator   -> ClueCode
///           - links (routeCode, staticCode, D, R_link), partitions, domain intervals,
///             allocLimit, live occupancy counts, Σ R_link
///           - version, SMT / Name SMT positions, commits, Counter, Controller
///
///         Mapping to ClueCode: specs are keyed by the same `uint32 code`, a spec can
///         only be written while `_codeOrg(code) != 0` (the code is vaulted), and writes are gated
///         by that organisation's roles.
abstract contract ClueProvider is ClueCodeOrgRef {
    // ---------------------------------------------------------------------
    // Types (values match the pack enumerations in Section 9.1.7)
    // ---------------------------------------------------------------------

    enum ProviderType { Unknown, Static, Route }

    enum Grade { None, S, A, B }

    /// @notice The single target a Static code serves (Section 4: one network, or one token under one network).
    /// @dev    For a Route code every field is empty.
    struct StaticTarget {
        string network;        // e.g. "Ethereum"
        string networkSymbol;  // e.g. "eth"  (namespace suffix: binance.eth)
        string token;          // e.g. "USDC"; empty when the code targets the network itself
        bytes tokenAddress;    // token contract address; empty iff `token` is empty
    }

    struct ProviderSpec {
        ProviderType providerType;
        Grade grade;
        uint64 createdAt;
        string name;             // resolved display name, e.g. "binance global"
        string namespace;        // canonical Clue-ID suffix, e.g. "binance" or "binance.eth"
        StaticTarget target;     // empty for Route
    }

    // ---------------------------------------------------------------------
    // Constants: grade ceilings (Section 3.5.1). Fixed by grade, never increased.
    // ---------------------------------------------------------------------

    uint32 internal constant S_DIRECT_CEILING = 500_000_000;
    uint32 internal constant A_DIRECT_CEILING = 300_000_000;
    uint32 internal constant B_DIRECT_CEILING = 150_000_000;

    uint32 internal constant S_ROUTE_BUDGET = 350_000_000;
    uint32 internal constant A_ROUTE_BUDGET = 150_000_000;
    // Grade B cannot be a Route code: budget is zero.

    uint256 internal constant MAX_NAME_LENGTH = 64;
    uint256 internal constant MAX_NAMESPACE_LENGTH = 32;
    uint256 internal constant MAX_SYMBOL_LENGTH = 10;
    uint256 internal constant MAX_TARGET_LENGTH = 64;

    // ---------------------------------------------------------------------
    // Storage
    // ---------------------------------------------------------------------

    /// @dev code => specification.
    mapping(uint32 code => ProviderSpec) internal _providers;

    /// @dev keccak256(namespace) => code. One canonical Clue-ID suffix per Provider Code.
    mapping(bytes32 namespaceHash => uint32 code) internal _codeByNamespace;

    /// @dev keccak256(name) => code. Enforces the "unique across all Provider Codes" rule.
    mapping(bytes32 nameHash => uint32 code) internal _codeByName;

    // ---------------------------------------------------------------------
    // Errors
    // ---------------------------------------------------------------------

    error InvalidProviderCode(uint32 code);
    error InvalidProviderType();
    error ProviderAlreadyRegistered(uint32 code);
    error ProviderNotRegistered(uint32 code);
    error RouteNotAllowedForGrade(uint32 code, Grade grade);
    error NotProviderAuthorised(uint256 orgId, address account);
    error InvalidProviderName();
    error InvalidNamespace();
    error NamespaceTaken(string namespace_, uint32 existingCode);
    error NameTaken(string name, uint32 existingCode);
    error InvalidStaticTarget();
    error RouteMustHaveEmptyTarget();
    error NamespaceMustEndWithNetworkSymbol();

    // ---------------------------------------------------------------------
    // Events
    // ---------------------------------------------------------------------

    event ProviderRegistered(
        uint32 indexed code,
        uint256 indexed orgId,
        ProviderType providerType,
        Grade grade,
        string name_,
        string namespace_
    );
    event ProviderCleared(uint32 indexed code, uint256 indexed orgId);

    // ---------------------------------------------------------------------
    // Modifiers
    // ---------------------------------------------------------------------

    modifier onlyRegistered(uint32 code) {
        if (_providers[code].providerType == ProviderType.Unknown) {
            revert ProviderNotRegistered(code);
        }
        _;
    }

    // ---------------------------------------------------------------------
    // Registration
    // ---------------------------------------------------------------------

    /// @notice Writes the specification of a vaulted Provider Code. Callable once per vaulting.
    /// @dev    Caller must hold ORG_ADMIN_ROLE or ORG_CONTROLLER_ROLE in the organisation that
    ///         vaulted `code`. Name rules that cannot be checked cheaply on-chain (English,
    ///         "Static name includes its network name") are validated by the Profile Server.
    /// @param code          Provider Code (uint32 token id).
    /// @param providerType  Static or Route. Route is allowed for Grades S and A only.
    /// @param name_         Resolved display name, e.g. "binance global" for code 02.
    /// @param namespace_    Canonical Clue-ID suffix: "binance" for code 02, "binance.eth" for code 65601.
    /// @param target        Static target; must be entirely empty for a Route code.
    function registerProvider(
        uint32 code,
        ProviderType providerType,
        string calldata name_,
        string calldata namespace_,
        StaticTarget calldata target
    ) external virtual {
        uint256 orgId = _requireProviderAuth(code);

        if (_providers[code].providerType != ProviderType.Unknown) {
            revert ProviderAlreadyRegistered(code);
        }
        if (providerType == ProviderType.Unknown) revert InvalidProviderType();

        Grade grade = _gradeOf(code);
        if (grade == Grade.None) revert InvalidProviderCode(code);
        if (providerType == ProviderType.Route && grade == Grade.B) {
            revert RouteNotAllowedForGrade(code, grade);
        }

        _validateName(name_);
        _validateNamespace(namespace_);
        _reserveIdentifiers(code, name_, namespace_);

        if (providerType == ProviderType.Static) {
            _validateStaticTarget(target, namespace_);
        } else {
            _requireEmptyTarget(target);
        }

        ProviderSpec storage spec = _providers[code];
        spec.providerType = providerType;
        spec.grade = grade;
        spec.createdAt = uint64(block.timestamp);
        spec.name = name_;
        spec.namespace = namespace_;
        // Field-wise copy: calldata struct with dynamic members -> storage.
        spec.target.network = target.network;
        spec.target.networkSymbol = target.networkSymbol;
        spec.target.token = target.token;
        spec.target.tokenAddress = target.tokenAddress;

        emit ProviderRegistered(code, orgId, providerType, grade, name_, namespace_);
    }

    /// @dev Removes a specification and frees its name and namespace. The inheriting
    ///      contract MUST call this from the release path (`releaseVaultCode`) so that a released
    ///      code does not keep a stale spec or lock its identifiers. Never call it for a code
    ///      that still holds state.
    function _clearProvider(uint32 code) internal virtual {
        ProviderSpec storage spec = _providers[code];
        if (spec.providerType == ProviderType.Unknown) return;

        delete _codeByNamespace[keccak256(bytes(spec.namespace))];
        delete _codeByName[keccak256(bytes(spec.name))];
        delete _providers[code];

        emit ProviderCleared(code, _codeOrg(code));
    }

    /// @dev Clears provider spec on code release.
    function _cleanupCodeOnRelease(uint32 code, uint256 orgId) internal virtual override {
        _clearProvider(code);
        super._cleanupCodeOnRelease(code, orgId);
    }

    // ---------------------------------------------------------------------
    // Views
    // ---------------------------------------------------------------------

    function isProviderRegistered(uint32 code) public view returns (bool) {
        return _providers[code].providerType != ProviderType.Unknown;
    }

    /// @notice Full specification of `code`.
    function getProvider(uint32 code) external view onlyRegistered(code) returns (ProviderSpec memory) {
        return _providers[code];
    }

    function providerTypeOf(uint32 code) external view returns (ProviderType) {
        return _providers[code].providerType;
    }

    function staticTargetOf(uint32 code) external view onlyRegistered(code) returns (StaticTarget memory) {
        return _providers[code].target;
    }

    /// @notice Canonical Clue-ID namespace lookup: the code that owns `namespace_`, or 0 if none.
    /// @dev Clients use this to resolve the suffix of `123456789@binance` or `123456789@binance.eth`
    ///      before validating the Damm check digit against the returned Provider Code.
    function codeByNamespace(string calldata namespace_) external view returns (uint32) {
        return _codeByNamespace[keccak256(bytes(namespace_))];
    }

    /// @notice Resolved display name of `code`, e.g. "binance global" for code 02.
    /// @dev Stored at registration; see {registerProvider}.
    function resolvedNameOf(uint32 code) external view onlyRegistered(code) returns (string memory) {
        return _providers[code].name;
    }

    /// @notice Grade of `code`, from digit length alone (Section 5). `None` for an invalid code.
    function gradeOf(uint32 code) external pure returns (Grade) {
        return _gradeOf(code);
    }

    /// @notice Static direct ceiling D for the code's grade (Section 3.5.1).
    function directCeilingOf(uint32 code) public pure returns (uint32) {
        Grade g = _gradeOf(code);
        if (g == Grade.S) return S_DIRECT_CEILING;
        if (g == Grade.A) return A_DIRECT_CEILING;
        if (g == Grade.B) return B_DIRECT_CEILING;
        revert InvalidProviderCode(code);
    }

    /// @notice Route budget R for the code's grade; zero for Grade B (Section 3.5.1).
    function routeBudgetOf(uint32 code) public pure returns (uint32) {
        Grade g = _gradeOf(code);
        if (g == Grade.S) return S_ROUTE_BUDGET;
        if (g == Grade.A) return A_ROUTE_BUDGET;
        if (g == Grade.B) return 0;
        revert InvalidProviderCode(code);
    }

    /// @notice The one canonical decimal form of a code: zero-padded to grade width (2, 3 or 5 digits).
    ///         Off-chain clients use this string for the Damm input so Go and TypeScript agree.
    function canonicalCodeOf(uint32 code) external pure returns (string memory) {
        Grade g = _gradeOf(code);
        if (g == Grade.None) revert InvalidProviderCode(code);
        uint256 width = g == Grade.S ? 2 : (g == Grade.A ? 3 : 5);

        bytes memory out = new bytes(width);
        uint256 n = code;
        for (uint256 i = width; i > 0; --i) {
            out[i - 1] = bytes1(uint8(48 + (n % 10)));
            n /= 10;
        }
        return string(out);
    }

    // ---------------------------------------------------------------------
    // Internal: grade
    // ---------------------------------------------------------------------

    /// @dev S: 01-99, A: 100-999, B: 10000-99999. Everything else (0, 1000-9999, >99999) is invalid.
    function _gradeOf(uint32 code) internal pure returns (Grade) {
        if (code >= 1 && code <= 99) return Grade.S;
        if (code >= 100 && code <= 999) return Grade.A;
        if (code >= 10_000 && code <= 99_999) return Grade.B;
        return Grade.None;
    }

    // ---------------------------------------------------------------------
    // Internal: authorisation
    // ---------------------------------------------------------------------

    /// @dev Code must be vaulted; caller must be admin or controller of the vaulting organisation.
    function _requireProviderAuth(uint32 code) internal view returns (uint256 orgId) {
        orgId = _codeOrg(code);
        if (orgId == 0) revert CodeNotVaulted(code);
        if (
            !_clueCode.hasOrgRole(orgId, ClueOrgRoles.ORG_ADMIN_ROLE, msg.sender)
                && !_clueCode.hasOrgRole(orgId, ClueOrgRoles.ORG_CONTROLLER_ROLE, msg.sender)
        ) {
            revert NotProviderAuthorised(orgId, msg.sender);
        }
    }

    // ---------------------------------------------------------------------
    // Internal: validation
    // ---------------------------------------------------------------------

    /// @dev Printable ASCII (0x20-0x7E), 1..MAX_NAME_LENGTH, no leading/trailing space.
    ///      "English" is a Profile Server check; ASCII-only blocks look-alike scripts on-chain.
    function _validateName(string calldata name_) internal pure {
        bytes calldata b = bytes(name_);
        uint256 len = b.length;
        if (len == 0 || len > MAX_NAME_LENGTH) revert InvalidProviderName();
        if (b[0] == 0x20 || b[len - 1] == 0x20) revert InvalidProviderName();
        for (uint256 i = 0; i < len; ++i) {
            bytes1 c = b[i];
            if (c < 0x20 || c > 0x7E) revert InvalidProviderName();
        }
    }

    /// @dev Canonical Clue-ID suffix: lowercase a-z, digits and '.'; starts with a
    ///      letter, ends with a letter/digit, and has no adjacent dots. Lowercase-only
    ///      means one spelling per namespace, so hashing is canonical:
    ///      "binance" (Route), "binance.eth" (Static).
    function _validateNamespace(string calldata namespace_) internal pure {
        bytes calldata b = bytes(namespace_);
        uint256 len = b.length;
        if (len < 2 || len > MAX_NAMESPACE_LENGTH) revert InvalidNamespace();
        if (!_isLower(b[0])) revert InvalidNamespace();
        if (b[len - 1] == 0x2E) revert InvalidNamespace();

        for (uint256 i = 0; i < len; ++i) {
            bytes1 c = b[i];
            bool isDot = c == 0x2E;
            if (!_isLower(c) && !_isDigit(c) && !isDot) revert InvalidNamespace();
            if (i != 0 && isDot && b[i - 1] == 0x2E) revert InvalidNamespace();
        }
    }

    function _reserveIdentifiers(
        uint32 code,
        string calldata name_,
        string calldata namespace_
    ) internal {
        bytes32 namespaceHash = keccak256(bytes(namespace_));
        uint32 existing = _codeByNamespace[namespaceHash];
        if (existing != 0) revert NamespaceTaken(namespace_, existing);

        bytes32 nameHash = keccak256(bytes(name_));
        existing = _codeByName[nameHash];
        if (existing != 0) revert NameTaken(name_, existing);

        _codeByNamespace[namespaceHash] = code;
        _codeByName[nameHash] = code;
    }

    /// @dev Static: network and symbol required; token and tokenAddress both present or both absent
    ///      (one network OR one token under it); namespace must end with ".<symbol>".
    function _validateStaticTarget(StaticTarget calldata t, string calldata namespace_) internal pure {
        bytes calldata network = bytes(t.network);
        bytes calldata symbol = bytes(t.networkSymbol);
        bytes calldata token = bytes(t.token);

        if (network.length == 0 || network.length > MAX_TARGET_LENGTH) revert InvalidStaticTarget();
        if (symbol.length == 0 || symbol.length > MAX_SYMBOL_LENGTH) revert InvalidStaticTarget();
        if (token.length > MAX_TARGET_LENGTH) revert InvalidStaticTarget();

        for (uint256 i = 0; i < symbol.length; ++i) {
            if (!_isLower(symbol[i]) && !_isDigit(symbol[i])) revert InvalidStaticTarget();
        }

        // Exactly one target: a token needs its contract address, and a bare network has neither.
        if ((token.length == 0) != (t.tokenAddress.length == 0)) revert InvalidStaticTarget();

        // Namespace suffix: "...<separator><symbol>", with at least one character of base before it.
        bytes calldata n = bytes(namespace_);
        uint256 suffixLen = symbol.length + 1;
        if (n.length <= suffixLen) revert NamespaceMustEndWithNetworkSymbol();
        if (n[n.length - suffixLen] != 0x2E) revert NamespaceMustEndWithNetworkSymbol();
        for (uint256 i = 0; i < symbol.length; ++i) {
            if (n[n.length - symbol.length + i] != symbol[i]) revert NamespaceMustEndWithNetworkSymbol();
        }
    }

    function _requireEmptyTarget(StaticTarget calldata t) internal pure {
        if (
            bytes(t.network).length != 0 ||
            bytes(t.networkSymbol).length != 0 ||
            bytes(t.token).length != 0 ||
            t.tokenAddress.length != 0
        ) revert RouteMustHaveEmptyTarget();
    }

    function _isLower(bytes1 c) private pure returns (bool) {
        return c >= 0x61 && c <= 0x7A;
    }

    function _isDigit(bytes1 c) private pure returns (bool) {
        return c >= 0x30 && c <= 0x39;
    }
}