// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {ClueAccounts} from "./ClueAccounts.sol";
import {IClueRegistry} from "./IClueRegistry.sol";

/**
 * @title ClueRegistry
 * @notice On-chain registry of confirmed Provider Code anchors. Inherits {ClueAccounts}
 *         so the leaf a client resolves against is checked against the same append-only
 *         record the registry registers into.
 * @dev Types and the external API are declared in {IClueRegistry}. The chain mapping is
 *      written by the confirmation logic (not assembled yet) through {_setChain}.
 *
 *      NOTE: reference implementation, not audited.
 */
contract ClueRegistry is IClueRegistry, ClueAccounts {
    // Storage

    /// @dev Provider Code => last confirmed chain record.
    mapping(uint32 code => ChainRegistry record) private _chains;

    // Views: resolve

    /// @inheritdoc IClueRegistry
    function resolve(uint32 code, bytes calldata leaf)
        external
        view
        override
        returns (ResolveResult memory)
    {
        ChainRegistry storage chain = _chains[code];
        if (chain.providerCode == 0) revert UnknownProviderCode(code);
        if (
            chain.status != CommitStatus.Success
                || chain.version == 0
                || chain.root == bytes32(0)
        ) revert NoConfirmedRoot(code);

        _requireRegistered(_leafHash(leaf));

        return ResolveResult({
            root: chain.root,
            version: uint64(chain.version),
            providerType: chain.providerType,
            organisationSigner: chain.organisationSigner
        });
    }

    // Internal: chain records

    /// @dev Stores the confirmed record for `code`. Only callable by registry logic that
    ///      inherits this contract; `providerCode` in the record must match the mapping key.
    function _setChain(uint32 code, ChainRegistry memory chain) internal {
        if (code == 0 || chain.providerCode != code) revert UnknownProviderCode(code);
        _chains[code] = chain;
    }

    /// @dev Accepts either a precomputed leaf hash or the bytes to hash into one. An empty
    ///      value fails closed; a zero hash can never be registered in {ClueAccounts}.
    function _leafHash(bytes calldata leaf) private pure returns (bytes32 h) {
        if (leaf.length == 0) revert EmptyLeaf();
        if (leaf.length == 32) {
            bytes calldata l = leaf;
            assembly {
                h := calldataload(l.offset)
            }
        } else {
            h = keccak256(leaf);
        }
    }
}
