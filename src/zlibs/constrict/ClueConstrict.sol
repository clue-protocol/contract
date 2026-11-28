// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ERC165Checker} from "@openzeppelin/contracts/utils/introspection/ERC165Checker.sol";
import {IClueConstrict} from "./IClueConstrict.sol";

/// @title ClueConstrict
/// @notice Lightweight trust layer for Clue contracts split across several
///         deployments to stay under the 24 KB (EIP-170) limit.
///         The "type" of a sibling is its ERC-165 interface id. An address
///         can only be granted a type if it answers `supportsInterface(id)`.
/// @dev Two internal functions and one modifier. Errors and events are declared
///      once in {IClueConstrict} so clients can decode them from the interface.
abstract contract ClueConstrict is IClueConstrict {
    using ERC165Checker for address;

    // ───────────── Storage ─────────────
    mapping(bytes4 => bool) private _constrictInterfaces;
    mapping(bytes4 => mapping(address => bool)) private _constrictRoles;

    // ───────────── Modifier ─────────────
    /// @notice Reverts unless `account` was granted the constrict role of `interfaceId`.
    /// @dev Usage: `hasConstrictRole(type(ICodeVault).interfaceId, msg.sender)`
    modifier hasConstrictRole(bytes4 interfaceId, address account) {
        if (!_constrictRoles[interfaceId][account]) {
            revert ConstrictMissingRole(interfaceId, account);
        }
        _;
    }

    // ───────────── Internal functions ─────────────

    /// @notice Registers an interface as a valid constrict type.
    /// @param name        Human-readable name (event only), e.g. "ICodeVault".
    /// @param interfaceId ERC-165 id, e.g. `type(ICodeVault).interfaceId`.
    function registerContractInterface(
        string memory name,
        bytes4 interfaceId
    ) internal {
        if (bytes(name).length == 0) revert ConstrictEmptyName();
        // 0x00000000 is meaningless; 0xffffffff is reserved invalid by ERC-165
        if (interfaceId == bytes4(0) || interfaceId == 0xffffffff) {
            revert ConstrictInvalidInterfaceId(interfaceId);
        }
        if (_constrictInterfaces[interfaceId]) {
            revert ConstrictInterfaceAlreadyRegistered(interfaceId);
        }

        _constrictInterfaces[interfaceId] = true;
        emit ContractInterfaceRegistered(interfaceId, name);
    }

    /// @notice Grants `account` the role of a registered interface type,
    ///         provided it reports support for that interface via ERC-165.
    function grantConstrict(bytes4 interfaceId, address account) internal {
        if (account == address(0)) revert ConstrictZeroAddress();
        if (!_constrictInterfaces[interfaceId]) {
            revert ConstrictInterfaceNotRegistered(interfaceId);
        }
        if (!account.supportsInterface(interfaceId)) {
            revert ConstrictInterfaceNotSupported(interfaceId, account);
        }

        _constrictRoles[interfaceId][account] = true;
        emit ConstrictGranted(interfaceId, account);
    }
}