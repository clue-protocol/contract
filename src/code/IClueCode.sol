// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {IERC721} from "@openzeppelin/contracts/token/ERC721/IERC721.sol";

/**
 * @title IClueCode - Provider Code NFT
 * @notice Interface of the ClueCode ERC-721 (Provider Code ownership).
 *         ClueCode stores ownership, grade, canonical form and lockable metadata.
 *         Organisation/vault state is declared in {IClueProviderCode}; capacity,
 *         zones and partitions still live in ClueRegistry. `supportsInterface` is
 *         inherited from IERC721; role administration is inherited from
 *         OpenZeppelin AccessControl in the implementation.
 */
interface IClueCode is IERC721 {
    enum Grade {
        None,
        S,
        A,
        B
    }

    struct Trait {
        string traitType;
        string value;
    }

    struct CodeMetadata {
        string name;
        string description;
        string image;
        string externalUrl;
        Trait[] traits;
    }

    // Events

    event MetadataUpdate(uint256 _tokenId); // ERC-4906
    event BatchMetadataUpdate(uint256 _fromTokenId, uint256 _toTokenId); // ERC-4906
    event MetadataSet(uint256 indexed code);
    event CodeMetadataFrozen(uint256 indexed code);
    event BaseURIUpdated(string newBaseURI);
    event MetadataFrozen();

    // Errors

    error CodeAlreadyMinted(uint256 code);
    error InvalidProviderCode(uint256 code);
    error MetadataIsFrozen();
    error CodeMetadataIsFrozen(uint256 code);
    error LengthMismatch();

    // Constants and state views

    function MINTER_ROLE() external view returns (bytes32);
    function METADATA_ROLE() external view returns (bytes32);
    function MAX_SUPPLY() external view returns (uint256);
    function totalMinted() external view returns (uint256);
    function baseURI() external view returns (string memory);
    function metadataFrozen() external view returns (bool);
    function codeFrozen(uint256 code) external view returns (bool);

    // Grade and canonical form

    function isValidCode(uint256 code) external pure returns (bool);
    function gradeOf(uint256 code) external pure returns (Grade);
    function canonical(uint256 code) external pure returns (string memory);

    /// @notice True when the Provider Code token has been minted.
    function isMinted(uint256 code) external view returns (bool);

    // Minting (MINTER_ROLE)

    function mint(address to, uint256 code) external;
    function mintBatch(address to, uint256[] calldata codes) external;
    function mintWithMetadata(address to, uint256 code, CodeMetadata calldata md) external;

    // Metadata (METADATA_ROLE)

    function setMetadata(uint256 code, CodeMetadata calldata md) external;
    function setMetadataBatch(uint256[] calldata codes, CodeMetadata[] calldata mds) external;
    function clearMetadata(uint256 code) external;
    function setBaseURI(string calldata newBaseURI) external;
    function freezeCode(uint256 code) external;
    function freezeMetadata() external;

    // Metadata reads

    function metadata(uint256 code) external view returns (string memory);
    function metadataOf(uint256 code) external view returns (CodeMetadata memory);
    function hasMetadata(uint256 code) external view returns (bool);
}
