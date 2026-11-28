// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title ClueAuthorizer
 * @notice Reusable set of contract-level authorizers, managed by an OpenZeppelin {Ownable}
 *         owner. Following {ClueValidators}, authorizers are contract-wide accounts, not
 *         per-code validator authorizers.
 * @dev The deployer is the initial owner and the initial authorizer. Only the owner can add
 *      or revoke authorizers; an authorizer cannot transfer the role to someone else.
 *      Inherit this contract to get `isAuthorizer`, `onlyAuthorizer`, {addAuthorizer} and
 *      {revokeAuthorizer}.
 */
abstract contract ClueAuthorizer is Ownable {
    // Storage

    /// @dev Set of contract-level authorizers. Only addresses marked here may perform gated
    ///      actions. Managed exclusively by the contract owner via {addAuthorizer} and
    ///      {revokeAuthorizer}.
    mapping(address => bool) public isAuthorizer;

    // Errors

    /// @dev The address is already an authorizer.
    error AlreadyAuthorizer(address account);

    /// @dev The address is not an authorizer.
    error NotAnAuthorizer(address account);

    /// @dev Caller is not an authorizer.
    error NotAuthorizer(address caller);

    // Events

    /// @notice An address was added to the authorizer set.
    event AuthorizerAdded(address indexed authorizer);

    /// @notice An address was removed from the authorizer set.
    event AuthorizerRevoked(address indexed authorizer);

    // Constructor

    /// @dev Sets the deploying account as the initial owner and the initial authorizer.
    constructor() Ownable(msg.sender) {
        isAuthorizer[msg.sender] = true;
        emit AuthorizerAdded(msg.sender);
    }

    // Modifiers

    /// @dev Contract-level authorizer: gated actions only.
    modifier onlyAuthorizer() {
        if (!isAuthorizer[msg.sender]) revert NotAuthorizer(msg.sender);
        _;
    }

    // Owner: authorizer set management

    /// @notice Adds `account` to the authorizer set.
    /// @dev Owner only. Reverts if `account` is already an authorizer.
    function addAuthorizer(address account) external onlyOwner {
        if (isAuthorizer[account]) revert AlreadyAuthorizer(account);
        isAuthorizer[account] = true;
        emit AuthorizerAdded(account);
    }

    /// @notice Removes `account` from the authorizer set.
    /// @dev Owner only. Reverts if `account` is not an authorizer.
    function revokeAuthorizer(address account) external onlyOwner {
        if (!isAuthorizer[account]) revert NotAnAuthorizer(account);
        isAuthorizer[account] = false;
        emit AuthorizerRevoked(account);
    }
}
