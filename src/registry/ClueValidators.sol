// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {ClueCodeControl} from "./ClueCodeControl.sol";

/**
 * @title ClueValidators
 * @notice Registry of validators and per-code validator assignments.
 * @dev Global validators are registered by name/id with their details. Each Provider Code
 *      maintains its own set of validator assignments that reference global validators.
 *      Global registration is gated by `Authorizer`, initially the deployer and
 *      transferable by the current Authorizer.
 */
abstract contract ClueValidators is ClueCodeControl {
    using EnumerableSet for EnumerableSet.Bytes32Set;

    // Constants

    uint256 internal constant MIN_VALIDATOR_NAME_LENGTH = 3;
    uint256 internal constant MAX_VALIDATOR_NAME_LENGTH = 32;
    uint256 internal constant MAX_VALIDATOR_DESCRIPTION_LENGTH = 256;

    // Types

    struct Validator {
        bytes32 id;           // Stable validator id (equals its key in _validators, never zero)
        string  name;         // Registered name (globally unique)
        string  description;  // Profile/summary
        address account;      // Address that signs EIP-712 approvals and is checked on-chain
    }

    // Storage

    /// @dev Address allowed to register validators. Set to the deployer at construction,
    ///      then transferable by the current Authorizer.
    address public Authorizer;

    /// @dev Global validators by stable id. Invariant: `_validators[k].id == k`.
    ///      Existence: `id != bytes32(0)`.
    mapping(bytes32 => Validator) internal _validators;

    /// @dev Globally enabled flag by validator id.
    mapping(bytes32 => bool) internal _validatorActive;

    /// @dev keccak256(normalized name) => validator id. Enforces global name uniqueness.
    mapping(bytes32 => bytes32) internal _validatorIdByName;

    /// @dev Validator account => validator id. One account maps to one validator,
    ///      and is the reverse lookup used when recovering EIP-712 signers.
    mapping(address => bytes32) internal _validatorIdByAccount;

    /// @dev All registered validator ids, for enumeration.
    EnumerableSet.Bytes32Set internal _validatorIds;

    /// @dev Provider Code => set of assigned validator ids.
    mapping(uint32 => EnumerableSet.Bytes32Set) internal _codeValidators;

    /// @dev Validator id => number of Provider Codes it is linked to.
    ///      Incremented on link, decremented on unlink; must be zero to resolve a validator.
    mapping(bytes32 => uint256) public validatorLinkCount;

    // Errors

    error InvalidValidator();
    error InvalidValidatorId();
    error InvalidValidatorName();
    error InvalidValidatorDescription();
    error InvalidValidatorAccount();
    error InvalidValidatorAuthorizer();
    error NotValidatorAuthorizer(address caller);
    error ValidatorAlreadyExists(bytes32 validatorId);
    error ValidatorNotFound(bytes32 validatorId);
    error ValidatorNameTaken(string name);
    error ValidatorAccountTaken(address account);
    error ValidatorNotActive(bytes32 validatorId);
    error ValidatorStillLinked(bytes32 validatorId, uint256 count);
    error CodeValidatorAlreadyAdded(uint32 code, bytes32 validatorId);
    error CodeValidatorNotFound(uint32 code, bytes32 validatorId);
    error NotValidatorAccount(bytes32 validatorId, address caller);

    // Events

    event ValidatorAuthorizerSet(address indexed previousAuthorizer, address indexed newAuthorizer);
    event ValidatorRegistered(bytes32 indexed validatorId, address indexed account, string name);
    event ValidatorResolved(bytes32 indexed validatorId);
    event ValidatorUpdated(bytes32 indexed validatorId, address indexed account, string name);
    event ValidatorActiveSet(bytes32 indexed validatorId, bool active);
    event CodeValidatorAdded(uint32 indexed code, bytes32 indexed validatorId, address indexed account);
    event CodeValidatorRemoved(uint32 indexed code, bytes32 indexed validatorId, address indexed account);

    // Constructor

    /// @dev Sets the deploying account as the initial Authorizer.
    constructor() {
        Authorizer = msg.sender;
    }

    // Modifiers

    modifier onlyValidatorAuthorizer() {
        if (msg.sender != Authorizer) revert NotValidatorAuthorizer(msg.sender);
        _;
    }

    // Global validator administration

    /// @notice Transfers the Authorizer role to `newAuthorizer`.
    /// @dev Current Authorizer only. The zero address is rejected.
    function transferAuthorizer(address newAuthorizer) internal onlyValidatorAuthorizer {
        if (newAuthorizer == address(0)) revert InvalidValidatorAuthorizer();
        address previous = Authorizer;
        Authorizer = newAuthorizer;
        emit ValidatorAuthorizerSet(previous, newAuthorizer);
    }

    /// @notice Registers a new global validator. It starts active.
    /// @param validatorId   Stable id, must be non-zero and unused.
    /// @param name          Globally unique, already normalized (see {_validateValidatorName}).
    /// @param description   Profile/summary, may be empty.
    /// @param account       Address that signs EIP-712 approvals. One account per validator.
    function registerValidator(
        bytes32 validatorId,
        string calldata name,
        string calldata description,
        address account
    ) external virtual onlyValidatorAuthorizer {
        if (validatorId == bytes32(0)) revert InvalidValidatorId();
        if (account == address(0)) revert InvalidValidatorAccount();
        if (_validators[validatorId].id != bytes32(0)) revert ValidatorAlreadyExists(validatorId);
        if (bytes(description).length > MAX_VALIDATOR_DESCRIPTION_LENGTH) revert InvalidValidatorDescription();

        _validateValidatorName(name);

        bytes32 nameKey = keccak256(bytes(name));
        if (_validatorIdByName[nameKey] != bytes32(0)) revert ValidatorNameTaken(name);
        if (_validatorIdByAccount[account] != bytes32(0)) revert ValidatorAccountTaken(account);

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
    /// @dev Authorizer only, and only when no Provider Code is linked to it. All
    ///      lookups (id, name, account and enumeration) are cleared.
    /// @param validatorId Validator to resolve.
    function resolveValidator(bytes32 validatorId) external virtual onlyValidatorAuthorizer {
        Validator storage validator = _validators[validatorId];
        if (validator.id == bytes32(0)) revert ValidatorNotFound(validatorId);

        uint256 count = validatorLinkCount[validatorId];
        if (count != 0) revert ValidatorStillLinked(validatorId, count);

        delete _validatorIdByName[keccak256(bytes(validator.name))];
        delete _validatorIdByAccount[validator.account];
        _validatorIds.remove(validatorId);
        delete _validatorActive[validatorId];
        delete validatorLinkCount[validatorId];
        delete _validators[validatorId];

        emit ValidatorResolved(validatorId);
    }

    /// @notice Returns the ids of all registered validators.
    function listValidators() external view returns (bytes32[] memory) {
        return _validatorIds.values();
    }

    // Per-code validator assignment

    /// @notice Links `validatorId` to Provider Code `code`.
    /// @dev Caller must hold {CONTROLLER_ROLE} for `code`.
    function linkValidator(uint32 code, bytes32 validatorId) external {
        _checkRole(code, CONTROLLER_ROLE);
        _linkValidatorToCode(code, validatorId);
    }

    /// @notice Unlinks `validatorId` from Provider Code `code`.
    /// @dev Caller must hold {CONTROLLER_ROLE} for `code`.
    function unlinkValidator(uint32 code, bytes32 validatorId) external {
        _checkRole(code, CONTROLLER_ROLE);
        _unlinkValidatorFromCode(code, validatorId);
    }

    /// @notice Returns true if `msg.sender` is a validator linked to `code`.
    function isValidatorLinked(uint32 code) external view returns (bool) {
        bytes32 validatorId = _validatorIdByAccount[msg.sender];
        return validatorId != bytes32(0) && _codeValidators[code].contains(validatorId);
    }

    /// @notice Returns the ids of all validators linked to `code`.
    function linkedValidatorsOf(uint32 code) external view returns (bytes32[] memory) {
        return _codeValidators[code].values();
    }

    /// @dev Links an active validator to a Provider Code and increments its link count.
    function _linkValidatorToCode(uint32 code, bytes32 validatorId) internal {
        Validator storage validator = _validators[validatorId];
        if (validator.id == bytes32(0)) revert ValidatorNotFound(validatorId);
        if (!_validatorActive[validatorId]) revert ValidatorNotActive(validatorId);
        if (!_codeValidators[code].add(validatorId)) revert CodeValidatorAlreadyAdded(code, validatorId);

        validatorLinkCount[validatorId] += 1;
        emit CodeValidatorAdded(code, validatorId, validator.account);
    }

    /// @dev Unlinks a validator from a Provider Code and decrements its link count.
    function _unlinkValidatorFromCode(uint32 code, bytes32 validatorId) internal {
        if (!_codeValidators[code].remove(validatorId)) revert CodeValidatorNotFound(code, validatorId);

        validatorLinkCount[validatorId] -= 1;
        emit CodeValidatorRemoved(code, validatorId, _validators[validatorId].account);
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