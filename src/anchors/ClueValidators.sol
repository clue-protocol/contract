// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import {ClueOrgRoles} from "../code/IClueProviderCode.sol";
import {ClueCodeOrgRef} from "./ClueCodeOrgRef.sol";

/**
 * @title ClueValidators
 * @notice Registry of validators and per-code validator assignments.
 * @dev Validators are registered by name/id with their details. Each Provider Code
 *      maintains its own set of validator assignments that reference global validators.
 *
 *      Registration is open: the `msg.sender` of {registerValidator} becomes that
 *      validator's `account` and is its Validator Authorizer. Linking is a two-step
 *      handshake: the organisation that owns a Provider Code requests a link, and
 *      only that validator's Validator Authorizer can accept it, delete its request
 *      messages and view them. Unlinking is one-sided: either the validator's account
 *      or the organisation's controller may revoke a link.
 *
 *      The separate contract-level `Authorizer` (set at deployment) has no say in
 *      linking: it can only delete validators via {resolveValidator}.
 */
abstract contract ClueValidators is ClueCodeOrgRef {
    using EnumerableSet for EnumerableSet.UintSet;

    // Constants

    uint256 internal constant MIN_VALIDATOR_NAME_LENGTH = 3;
    uint256 internal constant MAX_VALIDATOR_NAME_LENGTH = 32;
    uint256 internal constant MAX_VALIDATOR_DESCRIPTION_LENGTH = 256;
    uint256 internal constant MAX_LINK_MESSAGE_LENGTH = 256;

    // Types

    struct Validator {
        uint256 id;           // Incremental validator id (never zero)
        string  name;         // Registered name (globally unique)
        string  description;  // Profile/summary
        address account;      // Registrant (Validator Authorizer); also signs EIP-712 approvals
    }

    /// @notice A pending request from a Provider Code's organisation to link a validator.
    struct ValidatorLinkRequest {
        uint32  code;        // Provider Code the organisation wants the validator linked to
        uint256 validatorId; // Validator the request is addressed to
        string  message;     // Free-form note from the organisation, capped by MAX_LINK_MESSAGE_LENGTH
    }

    // Storage

    /// @dev Contract-level Authorizer, set to the deployer at construction. Its only
    ///      power is deleting validators through {resolveValidator}; it cannot link,
    ///      accept, delete requests or view them on a validator's behalf.
    address public Authorizer;

    /// @dev Next incremental validator id. First registered validator gets 1.
    uint256 internal _nextValidatorId;

    /// @dev Global validators by incremental id. Existence: `id != 0`.
    mapping(uint256 => Validator) internal _validators;

    /// @dev Globally enabled flag by validator id.
    mapping(uint256 => bool) internal _validatorActive;

    /// @dev keccak256(normalized name) => validator id. Enforces global name uniqueness.
    mapping(bytes32 => uint256) internal _validatorIdByName;

    /// @dev Validator account => validator id. One account maps to one validator,
    ///      and is the reverse lookup used when recovering EIP-712 signers.
    mapping(address => uint256) internal _validatorIdByAccount;

    /// @dev All registered validator ids, for enumeration.
    EnumerableSet.UintSet internal _validatorIds;

    /// @dev Provider Code => set of assigned validator ids.
    mapping(uint32 => EnumerableSet.UintSet) internal _codeValidators;

    /// @dev Validator id => number of Provider Codes it is linked to.
    ///      Incremented on link, decremented on unlink; must be zero to resolve a validator.
    mapping(uint256 => uint256) public validatorLinkCount;

    /// @dev Pending link requests addressed to each validator, in insertion order.
    mapping(uint256 validatorId => ValidatorLinkRequest[]) internal _validatorLinkRequests;

    /// @dev 1-based index of the pending request for (validatorId, code) in
    ///      `_validatorLinkRequests[validatorId]`; zero means no pending request.
    mapping(uint256 validatorId => mapping(uint32 code => uint256 index)) internal _validatorRequestIndex;

    // Errors

    error InvalidValidator();
    error InvalidValidatorName();
    error InvalidValidatorDescription();
    error InvalidValidatorAccount();
    error ValidatorAccountMustBeCaller();
    error InvalidValidatorAuthorizer();
    error NotValidatorAuthorizer(address caller);
    error ValidatorNotFound(uint256 validatorId);
    error ValidatorNameTaken(string name);
    error ValidatorAccountTaken(address account);
    error ValidatorNotActive(uint256 validatorId);
    error ValidatorStillLinked(uint256 validatorId, uint256 count);
    error CodeValidatorAlreadyAdded(uint32 code, uint256 validatorId);
    error CodeValidatorNotFound(uint32 code, uint256 validatorId);
    error NotValidatorAccount(uint256 validatorId, address caller);
    // CodeNotVaulted and MissingOrgRole are inherited from {ClueCodeOrgRef}.
    error InvalidLinkMessage();
    error ValidatorLinkAlreadyRequested(uint32 code, uint256 validatorId);
    error ValidatorLinkRequestNotFound(uint32 code, uint256 validatorId);
    error ValidatorStillRequested(uint256 validatorId, uint256 count);

    // Events

    event ValidatorAuthorizerSet(address indexed previousAuthorizer, address indexed newAuthorizer);
    event ValidatorRegistered(uint256 indexed validatorId, address indexed account, string name);
    event ValidatorResolved(uint256 indexed validatorId);
    event ValidatorUpdated(uint256 indexed validatorId, address indexed account, string name);
    event ValidatorActiveSet(uint256 indexed validatorId, bool active);
    event CodeValidatorAdded(uint32 indexed code, uint256 indexed validatorId, address indexed account);
    event CodeValidatorRemoved(uint32 indexed code, uint256 indexed validatorId, address indexed account);
    event ValidatorLinkRequested(uint32 indexed code, uint256 indexed validatorId, string message);
    event ValidatorLinkRequestAccepted(uint32 indexed code, uint256 indexed validatorId);
    event ValidatorLinkRequestDeleted(uint32 indexed code, uint256 indexed validatorId);

    // Constructor

    /// @dev Sets the deploying account as the initial Authorizer.
    constructor() {
        Authorizer = msg.sender;
    }

    // Modifiers

    /// @dev Contract-level Authorizer: deletion of validators only.
    modifier onlyAuthorizer() {
        if (msg.sender != Authorizer) revert NotValidatorAuthorizer(msg.sender);
        _;
    }

    /// @dev Validator Authorizer: the account that registered `validatorId`.
    modifier onlyValidatorAuthorizer(uint256 validatorId) {
        if (_validatorIdByAccount[msg.sender] != validatorId) revert NotValidatorAccount(validatorId, msg.sender);
        _;
    }

    // Global validator administration

    /// @notice Transfers the Authorizer role to `newAuthorizer`.
    /// @dev Current Authorizer only. The zero address is rejected.
    function transferAuthorizer(address newAuthorizer) internal onlyAuthorizer {
        if (newAuthorizer == address(0)) revert InvalidValidatorAuthorizer();
        address previous = Authorizer;
        Authorizer = newAuthorizer;
        emit ValidatorAuthorizerSet(previous, newAuthorizer);
    }

    /// @notice Registers a new validator. It starts active. Open to anyone: the caller
    ///         becomes the validator's account and its Validator Authorizer.
    /// @dev `account` must equal `msg.sender`, so a validator cannot be registered on
    ///      someone else's behalf. One account per validator. The validator id is assigned
    ///      incrementally by the contract.
    /// @param name        Globally unique, already normalized (see {_validateValidatorName}).
    /// @param description Profile/summary, may be empty.
    /// @param account     Caller's address; Validator Authorizer and EIP-712 signer.
    /// @return validatorId The newly assigned incremental validator id.
    function registerValidator(
        string calldata name,
        string calldata description,
        address account
    ) external virtual returns (uint256 validatorId) {
        if (account != msg.sender) revert ValidatorAccountMustBeCaller();
        validatorId = _registerValidator(name, description, account);
    }

    /// @dev Internal validator registration. `account` becomes the validator account/authorizer.
    function _registerValidator(
        string calldata name,
        string calldata description,
        address account
    ) internal virtual returns (uint256 validatorId) {
        if (account == address(0)) revert InvalidValidatorAccount();
        if (bytes(description).length > MAX_VALIDATOR_DESCRIPTION_LENGTH) revert InvalidValidatorDescription();

        _validateValidatorName(name);

        bytes32 nameKey = keccak256(bytes(name));
        if (_validatorIdByName[nameKey] != 0) revert ValidatorNameTaken(name);
        if (_validatorIdByAccount[account] != 0) revert ValidatorAccountTaken(account);

        validatorId = ++_nextValidatorId;
        _validators[validatorId] = Validator({
            id: validatorId,
            name: name,
            description: description,
            account: account
        });
        _validatorActive[validatorId] = true;
        _validatorIdByName[nameKey] = validatorId;
        _validatorIdByAccount[account] = validatorId;
        _validatorIds.add(validatorId);

        emit ValidatorRegistered(validatorId, account, name);
        emit ValidatorActiveSet(validatorId, true);
    }

    /// @notice Permanently removes a validator from the registry.
    /// @dev Contract-level {Authorizer} only, and only when no Provider Code is linked
    ///      to it and no link request is pending. All lookups (id, name, account and
    ///      enumeration) are cleared. This is the Authorizer's only power.
    /// @param validatorId Validator to resolve.
    function resolveValidator(uint256 validatorId) external virtual onlyAuthorizer {
        Validator storage validator = _validators[validatorId];
        if (validator.id == 0) revert ValidatorNotFound(validatorId);

        uint256 count = validatorLinkCount[validatorId];
        if (count != 0) revert ValidatorStillLinked(validatorId, count);
        count = _validatorLinkRequests[validatorId].length;
        if (count != 0) revert ValidatorStillRequested(validatorId, count);

        delete _validatorIdByName[keccak256(bytes(validator.name))];
        delete _validatorIdByAccount[validator.account];
        _validatorIds.remove(validatorId);
        delete _validatorActive[validatorId];
        delete validatorLinkCount[validatorId];
        delete _validators[validatorId];

        emit ValidatorResolved(validatorId);
    }

    /// @notice Returns the ids of all registered validators.
    function listValidators() external view returns (uint256[] memory) {
        return _validatorIds.values();
    }


    // Per-code validator assignment

    /// @notice Requests that `validatorId` be linked to Provider Code `code`.
    /// @dev Caller must hold {ORG_CONTROLLER_ROLE} for the organisation that owns `code`.
    ///      This only creates a request: the validator's Validator Authorizer must call
    ///      {acceptValidatorLink} before the link takes effect. At most one request
    ///      may be pending for a given (`code`, `validatorId`) pair.
    /// @param code        Provider Code the validator would be linked to.
    /// @param validatorId Validator the request is addressed to.
    /// @param message     Short note for the Validator Authorizer, max 256 bytes.
    function requestValidatorLink(uint32 code, uint256 validatorId, string calldata message) external {
        _checkCodeRole(code, ClueOrgRoles.ORG_CONTROLLER_ROLE);
        _requestValidatorLink(code, validatorId, message);
    }

    /// @notice Accepts a pending link request and links `validatorId` to `code`.
    /// @dev Only the validator's Validator Authorizer (its registered `account`).
    ///      The accepted request is removed; other pending requests addressed to the
    ///      same validator are left untouched.
    function acceptValidatorLink(uint32 code, uint256 validatorId)
        external
        onlyValidatorAuthorizer(validatorId)
    {
        _removeValidatorLinkRequest(code, validatorId);
        _linkValidatorToCode(code, validatorId);
        emit ValidatorLinkRequestAccepted(code, validatorId);
    }

    /// @notice Unlinks `validatorId` from Provider Code `code`.
    /// @dev Caller must be the `account` of `validatorId` or hold {ORG_CONTROLLER_ROLE}
    ///      for the organisation that owns `code`. Either side may revoke a link,
    ///      without needing the other side's approval.
    function unlinkValidator(uint32 code, uint256 validatorId) external {
        uint256 orgId = _codeOrg(code);
        bool isOrgController =
            orgId != 0 && _clueCode.hasOrgRole(orgId, ClueOrgRoles.ORG_CONTROLLER_ROLE, msg.sender);
        if (!isOrgController && _validatorIdByAccount[msg.sender] != validatorId) {
            revert NotValidatorAccount(validatorId, msg.sender);
        }
        _unlinkValidatorFromCode(code, validatorId);
    }

    /// @notice Returns every pending link request message addressed to `validatorId`.
    /// @dev Public view, so the Validator Authorizer can read its messages. Results are
    ///      in insertion order. Accepted requests and requests passed to
    ///      {deleteValidatorLinkRequests} are removed from this list.
    function validatorLinkRequestsOf(uint256 validatorId) external view returns (ValidatorLinkRequest[] memory) {
        return _validatorLinkRequests[validatorId];
    }

    /// @notice Deletes the pending link requests of `validatorId` for the given codes.
    /// @dev Only that validator's Validator Authorizer, which is also the only account
    ///      that can accept its requests. Reverts if any code has no pending request,
    ///      so the batch is all-or-nothing.
    function deleteValidatorLinkRequests(uint256 validatorId, uint32[] calldata codes)
        external
        onlyValidatorAuthorizer(validatorId)
    {
        for (uint256 i = 0; i < codes.length; ++i) {
            _removeValidatorLinkRequest(codes[i], validatorId);
        }
    }

    /// @notice Returns true if `msg.sender` is a validator linked to `code`.
    function isValidatorLinked(uint32 code) external view returns (bool) {
        uint256 validatorId = _validatorIdByAccount[msg.sender];
        return validatorId != 0 && _codeValidators[code].contains(validatorId);
    }

    /// @notice Returns the ids of all validators linked to `code`.
    function linkedValidatorsOf(uint32 code) external view returns (uint256[] memory) {
        return _codeValidators[code].values();
    }

    /// @dev Records an organisation's request to link `validatorId` to `code`.
    function _requestValidatorLink(uint32 code, uint256 validatorId, string calldata message) internal {
        if (bytes(message).length > MAX_LINK_MESSAGE_LENGTH) revert InvalidLinkMessage();

        Validator storage validator = _validators[validatorId];
        if (validator.id == 0) revert ValidatorNotFound(validatorId);
        if (!_validatorActive[validatorId]) revert ValidatorNotActive(validatorId);
        if (_codeValidators[code].contains(validatorId)) revert CodeValidatorAlreadyAdded(code, validatorId);
        if (_validatorRequestIndex[validatorId][code] != 0) revert ValidatorLinkAlreadyRequested(code, validatorId);

        ValidatorLinkRequest[] storage requests = _validatorLinkRequests[validatorId];
        requests.push(ValidatorLinkRequest({code: code, validatorId: validatorId, message: message}));
        _validatorRequestIndex[validatorId][code] = requests.length; // 1-based: 0 means absent

        emit ValidatorLinkRequested(code, validatorId, message);
    }

    /// @dev Links an active validator to a Provider Code and increments its link count.
    function _linkValidatorToCode(uint32 code, uint256 validatorId) internal {
        Validator storage validator = _validators[validatorId];
        if (validator.id == 0) revert ValidatorNotFound(validatorId);
        if (!_validatorActive[validatorId]) revert ValidatorNotActive(validatorId);
        if (!_codeValidators[code].add(validatorId)) revert CodeValidatorAlreadyAdded(code, validatorId);

        validatorLinkCount[validatorId] += 1;
        emit CodeValidatorAdded(code, validatorId, validator.account);
    }

    /// @dev Unlinks a validator from a Provider Code and decrements its link count.
    function _unlinkValidatorFromCode(uint32 code, uint256 validatorId) internal {
        if (!_codeValidators[code].remove(validatorId)) revert CodeValidatorNotFound(code, validatorId);

        validatorLinkCount[validatorId] -= 1;
        emit CodeValidatorRemoved(code, validatorId, _validators[validatorId].account);
    }

    /// @dev Removes the pending request for (`code`, `validatorId`) using swap-and-pop,
    ///      keeping the 1-based index map in sync. Reverts if there is no such request.
    function _removeValidatorLinkRequest(uint32 code, uint256 validatorId) internal {
        uint256 position = _validatorRequestIndex[validatorId][code];
        if (position == 0) revert ValidatorLinkRequestNotFound(code, validatorId);

        ValidatorLinkRequest[] storage requests = _validatorLinkRequests[validatorId];
        uint256 index = position - 1;
        uint256 last = requests.length - 1;
        if (index != last) {
            uint32 movedCode = requests[last].code;
            requests[index] = requests[last];
            _validatorRequestIndex[validatorId][movedCode] = position;
        }
        requests.pop();
        delete _validatorRequestIndex[validatorId][code];

        emit ValidatorLinkRequestDeleted(code, validatorId);
    }

    /// @dev Clears validator links and pending link requests for a released code.
    function _cleanupCodeOnRelease(uint32 code, uint256 orgId) internal virtual override {
        uint256[] memory linked = _codeValidators[code].values();
        for (uint256 i = 0; i < linked.length; ++i) {
            _unlinkValidatorFromCode(code, linked[i]);
        }

        uint256[] memory validatorIds = _validatorIds.values();
        for (uint256 i = 0; i < validatorIds.length; ++i) {
            uint256 validatorId = validatorIds[i];
            if (_validatorRequestIndex[validatorId][code] != 0) {
                _removeValidatorLinkRequest(code, validatorId);
            }
        }

        super._cleanupCodeOnRelease(code, orgId);
    }

    /// @dev Name must be 3 to 32 bytes, start with a letter, and use only
    ///      lowercase a-z, 0-9, '.', '_' and '-'. The caller must send it already
    ///      normalized, so no case folding happens on-chain and `Alpha` and `alpha`
    ///      cannot both exist.
    function _validateValidatorName(string calldata name) internal pure {
        bytes calldata b = bytes(name);
        uint256 len = b.length;
        if (len < MIN_VALIDATOR_NAME_LENGTH || len > MAX_VALIDATOR_NAME_LENGTH) revert InvalidValidatorName();

        for (uint256 i = 0; i < len; ++i) {
            bytes1 c = b[i];
            bool isLower = c >= 0x61 && c <= 0x7a;                // a-z
            bool isDigit = c >= 0x30 && c <= 0x39;                // 0-9
            bool isSep   = c == 0x2e || c == 0x5f || c == 0x2d;   // . _ -
            if (!(isLower || isDigit || isSep)) revert InvalidValidatorName();
            if (i == 0 && !isLower) revert InvalidValidatorName();
        }
    }
}