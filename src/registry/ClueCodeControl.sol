// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/**
 * @title ClueCodeControl
 * @notice Per-code, role-based access control for a Provider Code in the registry.
 * @dev Follows the ideology of OpenZeppelin's AccessControl, scoped to an individual
 *      code: roles are identified by `bytes32` and administered per `(code, role)`.
 *      Control of a code is represented by {CONTROLLER_ROLE}. Vault/custody and status
 *      concerns live in ClueCodeVault, which accesses this control layer for authorization.
 */
abstract contract ClueCodeControl {

    // Types

    struct RoleData {
        mapping(address account => bool) hasRole;
        bytes32 adminRole;
    }

    // Storage

    /// @dev code => role => role data
    mapping(uint32 => mapping(bytes32 => RoleData)) private _roles;

    /// @dev code => current controller (singular holder of CONTROLLER_ROLE), or zero.
    mapping(uint32 => address) private _controller;

    /// @dev Role admin that cannot be held, used to lock roles from public administration.
    bytes32 internal constant NO_ROLE_ADMIN = bytes32(type(uint256).max);

    bytes32 public constant DEFAULT_ADMIN_ROLE = 0x00;
    bytes32 public constant CONTROLLER_ROLE = keccak256("CONTROLLER_ROLE");

    // Errors

    error NotVaulted(uint32 code);
    error AlreadyVaulted(uint32 code);
    error NotDepositor(uint32 code, address caller);
    error CodeAccessControlUnauthorizedAccount(uint32 code, address account, bytes32 neededRole);
    error CodeAccessControlBadConfirmation();

    // Events

    event CodeRoleAdminChanged(uint32 indexed code, bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 newAdminRole);
    event CodeRoleGranted(uint32 indexed code, bytes32 indexed role, address indexed account, address sender);
    event CodeRoleRevoked(uint32 indexed code, bytes32 indexed role, address indexed account, address sender);

    // Role views (OZ-style, per code)

    /// @notice Returns true if `account` has been granted `role` for `code`.
    function hasRole(uint32 code, bytes32 role, address account) internal view virtual returns (bool) {
        return _roles[code][role].hasRole[account];
    }

    /// @notice Returns the admin role that controls `role` for `code`.
    /// @dev The controller role is locked from public administration and reports {NO_ROLE_ADMIN}.
    function getRoleAdmin(uint32 code, bytes32 role) internal view virtual returns (bytes32) {
        if (role == CONTROLLER_ROLE) return NO_ROLE_ADMIN;
        bytes32 admin = _roles[code][role].adminRole;
        // DEFAULT_ADMIN_ROLE is the default admin for all other roles when unset.
        return admin;
    }

    // Role administration (OZ-style, per code)

    /// @notice Grants `role` for `code` to `account`.
    /// @dev Caller must hold the role's admin role for `code`. {CONTROLLER_ROLE} cannot
    ///      be granted publicly; it is managed by the vault lifecycle.
    function grantRole(uint32 code, bytes32 role, address account) internal virtual {
        _checkRole(code, getRoleAdmin(code, role));
        _grantRole(code, role, account);
    }

    /// @notice Revokes `role` for `code` from `account`.
    /// @dev Caller must hold the role's admin role for `code`. {CONTROLLER_ROLE} cannot
    ///      be revoked publicly; it is managed by the vault lifecycle.
    function revokeRole(uint32 code, bytes32 role, address account) internal virtual {
        _checkRole(code, getRoleAdmin(code, role));
        _revokeRole(code, role, account);
    }

    /// @notice Revokes `role` for `code` from the calling account.
    /// @dev {CONTROLLER_ROLE} cannot be renounced publicly; use the vault release flow.
    function renounceRole(uint32 code, bytes32 role, address callerConfirmation) internal virtual {
        if (callerConfirmation != msg.sender) revert CodeAccessControlBadConfirmation();
        // Disallow public renounce of controller role to preserve vault invariants.
        _checkRole(code, getRoleAdmin(code, role));
        _revokeRole(code, role, callerConfirmation);
    }

    // Controller / depositor views (vault compatibility)

    /// @notice Returns the controller (depositor) of `code`, or zero if none.
    function controllerOf(uint32 code) internal view returns (address) {
        return _controller[code];
    }

    /// @notice Returns the depositor of `code`, or zero if none.
    /// @dev Alias of {controllerOf} kept for compatibility.
    function depositorOf(uint32 code) internal view returns (address) {
        return _controller[code];
    }

    /// @notice Returns true if `code` is under control (vaulted).
    function isControlled(uint32 code) internal view returns (bool) {
        return _controller[code] != address(0);
    }

    // Internal: role checks

    /// @dev Reverts if `msg.sender` is missing `role` for `code`.
    function _checkRole(uint32 code, bytes32 role) internal view virtual {
        _checkRole(code, role, msg.sender);
    }

    /// @dev Reverts if `account` is missing `role` for `code`.
    function _checkRole(uint32 code, bytes32 role, address account) internal view virtual {
        if (!hasRole(code, role, account)) {
            revert CodeAccessControlUnauthorizedAccount(code, account, role);
        }
    }

    /// @dev Sets `adminRole` as `role`'s admin role for `code`.
    function _setRoleAdmin(uint32 code, bytes32 role, bytes32 adminRole) internal virtual {
        bytes32 previousAdminRole = getRoleAdmin(code, role);
        _roles[code][role].adminRole = adminRole;
        emit CodeRoleAdminChanged(code, role, previousAdminRole, adminRole);
    }

    /// @dev Grants `role` for `code` to `account`. Internal, no access restriction.
    function _grantRole(uint32 code, bytes32 role, address account) internal virtual returns (bool) {
        if (!hasRole(code, role, account)) {
            _roles[code][role].hasRole[account] = true;
            if (role == CONTROLLER_ROLE) {
                _controller[code] = account;
            }
            emit CodeRoleGranted(code, role, account, msg.sender);
            return true;
        }
        return false;
    }

    /// @dev Revokes `role` for `code` from `account`. Internal, no access restriction.
    function _revokeRole(uint32 code, bytes32 role, address account) internal virtual returns (bool) {
        if (hasRole(code, role, account)) {
            _roles[code][role].hasRole[account] = false;
            if (role == CONTROLLER_ROLE && _controller[code] == account) {
                _controller[code] = address(0);
            }
            emit CodeRoleRevoked(code, role, account, msg.sender);
            return true;
        }
        return false;
    }

    // Internal: controller helpers (vault compatibility)

    /// @dev Returns the controller of `code`, or zero if none.
    function _controllerOf(uint32 code) internal view returns (address) {
        return _controller[code];
    }

    /// @dev Reverts if `code` is not under control.
    function _requireControlled(uint32 code) internal view returns (address controller) {
        controller = _controller[code];
        if (controller == address(0)) revert NotVaulted(code);
    }

    /// @dev Reverts if `code` is already under control.
    function _requireNotControlled(uint32 code) internal view {
        if (_controller[code] != address(0)) revert AlreadyVaulted(code);
    }

    /// @dev Reverts unless the caller is the controller of `code`.
    ///      Reverts with NotVaulted if the code is not under control.
    function _requireController(uint32 code) internal view returns (address controller) {
        controller = _controller[code];
        if (controller == address(0)) revert NotVaulted(code);
        if (controller != msg.sender) revert NotDepositor(code, msg.sender);
    }

    /// @dev Assigns control of `code` to `controller` by granting {CONTROLLER_ROLE}
    ///      and {DEFAULT_ADMIN_ROLE}. Reverts if already controlled or if `controller` is zero.
    function _assignController(uint32 code, address controller) internal {
        if (controller == address(0)) revert NotVaulted(code);
        if (_controller[code] != address(0)) revert AlreadyVaulted(code);
        _grantRole(code, CONTROLLER_ROLE, controller);
        _grantRole(code, DEFAULT_ADMIN_ROLE, controller);
    }

    /// @dev Clears control of `code` by revoking {CONTROLLER_ROLE} and {DEFAULT_ADMIN_ROLE}.
    ///      Reverts if not controlled. Returns the previous controller.
    function _clearController(uint32 code) internal returns (address controller) {
        controller = _controller[code];
        if (controller == address(0)) revert NotVaulted(code);
        _revokeRole(code, CONTROLLER_ROLE, controller);
        _revokeRole(code, DEFAULT_ADMIN_ROLE, controller);
        // Reset admin configuration for the code.
        _roles[code][DEFAULT_ADMIN_ROLE].adminRole = bytes32(0);
    }
}
