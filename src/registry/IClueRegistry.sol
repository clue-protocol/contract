// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/**
 * @title IClueRegistry
 * @notice Interface and shared types of the Clue Protocol on-chain registry.
 * @dev {ClueRegistry} also inherits {ClueAccounts}, so the deployed ABI additionally
 *      contains the account registration and validation functions declared there. They
 *      are deliberately not redeclared here: the interface describes the registry's own
 *      surface, while ClueAccounts remains the single source of truth for its own API.
 */
interface IClueRegistry {
    // Types

    /// @notice Confirmation state of a Provider Code's accepted commit.
    enum CommitStatus { None, Success }

    /// @notice Route provenance for an atomically updated Static code.
    struct InfectedBy {
        uint32 routeCode;    // zero when not infected
        uint32 routeVersion; // zero when not infected
    }

    /// @notice On-chain registry anchor for one Provider Code.
    struct ChainRegistry {
        uint32 providerCode;
        string providerName;
        uint8 providerType;          // ProviderType: Static or Route

        // --- roots ---
        bytes32 root;                // combined RootHash accepted on-chain
        bytes32 prevRoot;            // previous accepted root (chain continuity)
        bytes32 totalAccountHash;    // account hash
        bytes32 organisationSigner;  // Ed25519 public key of the owning organisation

        // --- versioning ---
        uint32 version;              // current confirmed version
        uint32 prevVersion;
        CommitStatus status;         // None / Success

        // --- commit window ---
        uint64 commitSequence;       // first Count in the version
        uint64 commitSequenceEnd;    // last Count in the version
        uint32 updatedSerialCount;   // serials touched by this version
        uint64 serialCount;          // live occupancy (CREATE +1, DELETE -1)

        // --- route / static atomicity ---
        InfectedBy infectedBy;       // Route code + Route version, zero if none

        // --- validator approval ---
        uint256[] validatorsID;      // validator ids that signed
        uint32 approvalCount;
        uint32 required;             // threshold at confirmation time

        // --- anchoring evidence ---
        uint64 blockNumber;          // uint64(block.number)
        uint64 timeStamp;            // uint48(block.timestamp)
    }

    /// @notice Compact commit record for a Provider Code.
    struct ChainRecord {
        bytes root;
        bytes version;
        uint32 providerCode;
        uint64 commitSequence;
        uint64 serialCount;
    }

    /// @notice Trusted on-chain anchor returned by {resolve}.
    struct ResolveResult {
        bytes32 root;               // current confirmed root
        uint64 version;             // current confirmed version
        uint8 providerType;         // 0 - Static, 1 - Route
        bytes32 organisationSigner; // organisation Ed25519 public key
    }

    // Errors

    /// @dev No chain record is stored for `code`.
    error UnknownProviderCode(uint32 code);

    /// @dev The chain has no confirmed root yet.
    error NoConfirmedRoot(uint32 code);

    /// @dev No leaf was supplied to {resolve}.
    error EmptyLeaf();

    // Functions

    /**
     * @notice Returns the trusted on-chain anchor for `code`, for off-chain resolve
     *         comparison. The leaf must be registered in the inherited {ClueAccounts}
     *         record; account authenticity is verified client-side.
     * @dev View, no inclusion proof and no Ed25519 verification. `leaf` is either the
     *      32-byte account leaf hash, or its preimage, which is hashed with keccak256.
     * @param code Provider Code to resolve.
     * @param leaf Registered leaf hash, or the bytes to hash into one.
     */
    function resolve(uint32 code, bytes calldata leaf)
        external
        view
        returns (ResolveResult memory);
}
