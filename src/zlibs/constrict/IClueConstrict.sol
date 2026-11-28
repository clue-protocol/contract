// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title IClueConstrict
 * @notice ABI surface of the {ClueConstrict} trust layer: registered interface types
 *         and the accounts granted a role for each type. The role operations themselves
 *         are internal to the inheriting contract.
 */
interface IClueConstrict {
    // Errors

    error ConstrictEmptyName();
    error ConstrictInvalidInterfaceId(bytes4 interfaceId);
    error ConstrictInterfaceAlreadyRegistered(bytes4 interfaceId);
    error ConstrictInterfaceNotRegistered(bytes4 interfaceId);
    error ConstrictZeroAddress();
    error ConstrictInterfaceNotSupported(bytes4 interfaceId, address account);
    error ConstrictMissingRole(bytes4 interfaceId, address account);

    // Events

    event ContractInterfaceRegistered(bytes4 indexed interfaceId, string name);
    event ConstrictGranted(bytes4 indexed interfaceId, address indexed account);
}
