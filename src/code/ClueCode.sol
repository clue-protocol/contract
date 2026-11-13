// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {ERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import {AccessControlDefaultAdminRules} from "@openzeppelin/contracts/access/extensions/AccessControlDefaultAdminRules.sol";
import {Strings} from "@openzeppelin/contracts/utils/Strings.sol";
import {Base64} from "@openzeppelin/contracts/utils/Base64.sol";

/**
 * @title ClueCode
 * @notice ERC-721 representing ownership of a Provider Code of the Clue Protocol.
 *
 * Metadata model
 *  - Each code may have metadata stored in this contract (name, description,
 *    image, external URL, extra traits). tokenURI builds the JSON from it.
 *  - Grade and Digits traits are always derived automatically.
 *  - A code with no stored metadata falls back to baseURI + canonical code
 *    (if baseURI is set), otherwise to a default on-chain JSON.
 *
 * NOTE: reference implementation, not audited.
 */
contract ClueCode is ERC721, AccessControlDefaultAdminRules {
    enum Grade { None, S, A, B }

    struct Trait {
        string traitType;
        string value;
    }

    struct CodeMetadata {
        string name;         // empty => "Clue Provider Code <canonical>"
        string description;  // empty => default description
        string image;        // e.g. "ipfs://<cid>/02.png"; empty => omitted
        string externalUrl;  // empty => omitted
        Trait[] traits;      // extra traits, added after Grade and Digits
    }

    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
    bytes32 public constant METADATA_ROLE = keccak256("METADATA_ROLE");
    uint256 public constant MAX_SUPPLY = 90_999;

    uint256 public totalMinted;

    /// @notice Fallback base, used only for codes with no stored metadata.
    string public baseURI;

    /// @notice Once true, no metadata or baseURI can ever change again.
    bool public metadataFrozen;

    mapping(uint256 => CodeMetadata) private _metadata;
    mapping(uint256 => bool) private _hasMetadata;
    mapping(uint256 => bool) public codeFrozen;

    /// @dev ERC-4906: tells marketplaces to refresh a token.
    event MetadataUpdate(uint256 _tokenId);
    event MetadataSet(uint256 indexed code);
    event CodeMetadataFrozen(uint256 indexed code);
    event BaseURIUpdated(string newBaseURI);
    event MetadataFrozen();

    error InvalidProviderCode(uint256 code);
    error MetadataIsFrozen();
    error CodeMetadataIsFrozen(uint256 code);
    error LengthMismatch();

    constructor(address admin, uint48 adminDelay, address minter, string memory initialURI)
        ERC721("Clue Provider Code", "CPC")
        AccessControlDefaultAdminRules(adminDelay, admin)
    {
        _grantRole(MINTER_ROLE, minter);
        _grantRole(METADATA_ROLE, admin);
        // _grantRole(MINTER_ROLE, admin);
        baseURI = initialURI;
    }

    // ------------------------------------------------------------------
    // Grade and canonical form
    // ------------------------------------------------------------------

    function _grade(uint256 c) private pure returns (Grade) {
        if (c >= 1 && c <= 99) return Grade.S;
        if (c >= 100 && c <= 999) return Grade.A;
        if (c >= 10000 && c <= 99999) return Grade.B;
        return Grade.None;
    }

    function isValidCode(uint256 code) public pure returns (bool) {
        return _grade(code) != Grade.None;
    }

    function gradeOf(uint256 code) public pure returns (Grade g) {
        g = _grade(code);
        if (g == Grade.None) revert InvalidProviderCode(code);
    }

    function canonical(uint256 code) public pure returns (string memory) {
        Grade g = gradeOf(code);
        uint256 width = g == Grade.S ? 2 : g == Grade.A ? 3 : 5;
        bytes memory digits = bytes(Strings.toString(code));
        bytes memory out = new bytes(width);
        uint256 pad = width - digits.length;
        for (uint256 i; i < pad; ++i) out[i] = "0";
        for (uint256 i; i < digits.length; ++i) out[pad + i] = digits[i];
        return string(out);
    }

    // ------------------------------------------------------------------
    // Minting (MINTER_ROLE only)
    // ------------------------------------------------------------------

    function mint(address to, uint256 code) external onlyRole(MINTER_ROLE) {
        _mintCode(to, code);
    }

    function mintBatch(address to, uint256[] calldata codes) external onlyRole(MINTER_ROLE) {
        for (uint256 i; i < codes.length; ++i) {
            _mintCode(to, codes[i]);
        }
    }

    /// @notice Mint a code and store its metadata in the same transaction.
    function mintWithMetadata(address to, uint256 code, CodeMetadata calldata md)
        external
        onlyRole(MINTER_ROLE)
    {
        _mintCode(to, code);
        _store(code, md);
    }

    function _mintCode(address to, uint256 code) private {
        gradeOf(code);
        _mint(to, code);
        unchecked {
            ++totalMinted;
        }
    }

    // ------------------------------------------------------------------
    // Metadata storage (METADATA_ROLE only)
    // ------------------------------------------------------------------

    function setMetadata(uint256 code, CodeMetadata calldata md)
        external
        onlyRole(METADATA_ROLE)
    {
        _store(code, md);
    }

    function setMetadataBatch(uint256[] calldata codes, CodeMetadata[] calldata mds)
        external
        onlyRole(METADATA_ROLE)
    {
        if (codes.length != mds.length) revert LengthMismatch();
        for (uint256 i; i < codes.length; ++i) {
            _store(codes[i], mds[i]);
        }
    }

    /// @notice Remove stored metadata so the code falls back to baseURI.
    function clearMetadata(uint256 code) external onlyRole(METADATA_ROLE) {
        _checkWritable(code);
        delete _metadata[code];
        _hasMetadata[code] = false;
        emit MetadataSet(code);
        emit MetadataUpdate(code);
    }

    function _checkWritable(uint256 code) private view {
        gradeOf(code);
        if (metadataFrozen) revert MetadataIsFrozen();
        if (codeFrozen[code]) revert CodeMetadataIsFrozen(code);
    }

    function _store(uint256 code, CodeMetadata calldata md) private {
        _checkWritable(code);
        CodeMetadata storage s = _metadata[code];
        s.name = md.name;
        s.description = md.description;
        s.image = md.image;
        s.externalUrl = md.externalUrl;
        delete s.traits;
        for (uint256 i; i < md.traits.length; ++i) {
            s.traits.push(Trait(md.traits[i].traitType, md.traits[i].value));
        }
        _hasMetadata[code] = true;
        emit MetadataSet(code);
        emit MetadataUpdate(code);
    }

    function setBaseURI(string calldata newBaseURI) external onlyRole(METADATA_ROLE) {
        if (metadataFrozen) revert MetadataIsFrozen();
        baseURI = newBaseURI;
        emit BaseURIUpdated(newBaseURI);
    }

    /// @notice Permanently lock one code's metadata.
    function freezeCode(uint256 code) external onlyRole(METADATA_ROLE) {
        gradeOf(code);
        codeFrozen[code] = true;
        emit CodeMetadataFrozen(code);
    }

    /// @notice Permanently lock all metadata and the base URI.
    function freezeMetadata() external onlyRole(METADATA_ROLE) {
        metadataFrozen = true;
        emit MetadataFrozen();
    }

    // ------------------------------------------------------------------
    // Metadata reads
    // ------------------------------------------------------------------

    /// @notice The metadata JSON for `code`, exactly as the contract builds it.
    function metadata(uint256 code) public view returns (string memory) {
        gradeOf(code);
        return _json(code);
    }

    /// @notice The raw stored fields for `code`.
    function metadataOf(uint256 code) external view returns (CodeMetadata memory) {
        gradeOf(code);
        return _metadata[code];
    }

    function hasMetadata(uint256 code) external view returns (bool) {
        return _hasMetadata[code];
    }

    function tokenURI(uint256 code) public view override returns (string memory) {
        _requireOwned(code);
        if (!_hasMetadata[code] && bytes(baseURI).length != 0) {
            return string.concat(baseURI, canonical(code));
        }
        return string.concat(
            "data:application/json;base64,",
            Base64.encode(bytes(_json(code)))
        );
    }

    // ------------------------------------------------------------------
    // JSON builder
    // ------------------------------------------------------------------

    function _json(uint256 code) private view returns (string memory) {
        Grade g = gradeOf(code);
        string memory c = canonical(code);
        CodeMetadata storage m = _metadata[code];

        string memory name = bytes(m.name).length != 0
            ? m.name
            : string.concat("Clue Provider Code ", c);
        string memory desc = bytes(m.description).length != 0
            ? m.description
            : "A Provider Code of the Clue Protocol.";

        string memory head = string.concat(
            '{"name":"', _esc(name), '","description":"', _esc(desc), '"'
        );

        string memory mid = "";
        if (bytes(m.image).length != 0) {
            mid = string.concat(',"image":"', _esc(m.image), '"');
        }
        if (bytes(m.externalUrl).length != 0) {
            mid = string.concat(mid, ',"external_url":"', _esc(m.externalUrl), '"');
        }

        return string.concat(head, mid, ',"attributes":', _attributes(c, g, m), "}");
    }

    function _attributes(string memory c, Grade g, CodeMetadata storage m)
        private
        view
        returns (string memory out)
    {
        out = string.concat(
            '[{"trait_type":"Grade","value":"', _gradeName(g), '"},',
            '{"trait_type":"Digits","display_type":"number","value":',
            Strings.toString(bytes(c).length),
            "}"
        );
        uint256 n = m.traits.length;
        for (uint256 i; i < n; ++i) {
            out = string.concat(
                out,
                ',{"trait_type":"', _esc(m.traits[i].traitType),
                '","value":"', _esc(m.traits[i].value), '"}'
            );
        }
        out = string.concat(out, "]");
    }

    function _gradeName(Grade g) private pure returns (string memory) {
        return g == Grade.S ? "S" : g == Grade.A ? "A" : "B";
    }

    /// @dev JSON string escaping: quote, backslash and control characters.
    function _esc(string memory s) private pure returns (string memory) {
        bytes memory b = bytes(s);
        bytes memory out = new bytes(b.length * 6);
        uint256 j;
        bytes16 hexChars = "0123456789abcdef";
        for (uint256 i; i < b.length; ++i) {
            bytes1 ch = b[i];
            if (ch == 0x22 || ch == 0x5c) {
                out[j++] = 0x5c;
                out[j++] = ch;
            } else if (uint8(ch) < 0x20) {
                out[j++] = 0x5c;
                out[j++] = "u";
                out[j++] = "0";
                out[j++] = "0";
                out[j++] = hexChars[uint8(ch) >> 4];
                out[j++] = hexChars[uint8(ch) & 0x0f];
            } else {
                out[j++] = ch;
            }
        }
        assembly {
            mstore(out, j)
        }
        return string(out);
    }

    // ------------------------------------------------------------------
    // ERC-165 (includes ERC-4906)
    // ------------------------------------------------------------------

    function supportsInterface(bytes4 id)
        public
        view
        override(ERC721, AccessControlDefaultAdminRules)
        returns (bool)
    {
        return id == bytes4(0x49064906) || super.supportsInterface(id);
    }
}