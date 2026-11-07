// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {ERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import {AccessControlDefaultAdminRules} from "@openzeppelin/contracts/access/extensions/AccessControlDefaultAdminRules.sol";
import {Strings} from "@openzeppelin/contracts/utils/Strings.sol";
import {Base64} from "@openzeppelin/contracts/utils/Base64.sol";

/**
 * @title ClueCode
 * @notice ERC-721 representing ownership of a Provider Code of the Clue Protocol
 *         (the numeric prefix of a Clue ID, e.g. 27 in 27-1234567890).
 *
 * Token model
 *  - tokenId == Provider Code read as an integer.
 *  - Grade S: 01-99 (2 digits), Grade A: 100-999 (3), Grade B: 10000-99999 (5).
 *  - Total supply is 90,999. Zero, 1000-9999 and >99999 can never be minted.
 *
 * ClueCode stores ownership only. Grade and canonical form are derived from the
 * token id. Capacity, zones, links, partitions, versions and status live in
 * ClueRegistry. The contract is not upgradeable and no role can transfer, burn,
 * pause or re-describe a token.
 *
 * NOTE: reference implementation, not audited.
 */
contract ClueCode is ERC721, AccessControlDefaultAdminRules {
    /// @dev Matches the pack enumeration: S = 1, A = 2, B = 3, 0 = invalid.
    enum Grade { None, S, A, B }

    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
    uint256 public constant MAX_SUPPLY = 90_999;

    /// @notice Number of codes minted so far.
    uint256 public totalMinted;

    error InvalidProviderCode(uint256 code);

    /**
     * @param admin      Multisig that manages roles.
     * @param adminDelay Seconds before a default-admin change takes effect.
     * @param minter     Genesis distributor that receives MINTER_ROLE.
     */
    constructor(address admin, uint48 adminDelay, address minter)
        ERC721("Clue Provider Code", "CPC")
        AccessControlDefaultAdminRules(adminDelay, admin)
    {
        _grantRole(MINTER_ROLE, minter);
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

    /// @notice True if `code` lies within a valid range.
    function isValidCode(uint256 code) public pure returns (bool) {
        return _grade(code) != Grade.None;
    }

    /// @notice Grade of `code`; reverts with InvalidProviderCode for an invalid code.
    function gradeOf(uint256 code) public pure returns (Grade g) {
        g = _grade(code);
        if (g == Grade.None) revert InvalidProviderCode(code);
    }

    /// @notice Zero-padded canonical string, e.g. 2 => "02", 100 => "100".
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

    function _mintCode(address to, uint256 code) private {
        gradeOf(code); // reverts for any value that is not a valid code
        _mint(to, code); // reverts if the code already exists (or `to` is zero)
        unchecked {
            ++totalMinted;
        }
    }

    // ------------------------------------------------------------------
    // Metadata: fully on-chain, static
    // ------------------------------------------------------------------

    function tokenURI(uint256 code) public view override returns (string memory) {
        _requireOwned(code);
        Grade g = gradeOf(code);
        string memory c = canonical(code);

        string memory image = string.concat(
            "data:image/svg+xml;base64,",
            Base64.encode(bytes(_svg(c, g)))
        );

        string memory head = string.concat(
            '{"name":"Clue Provider Code ',
            c,
            '","description":"A Provider Code of the Clue Protocol.",'
        );
        string memory tail = string.concat(
            '"image":"',
            image,
            '","attributes":',
            _attributes(c, g),
            "}"
        );

        return string.concat(
            "data:application/json;base64,",
            Base64.encode(bytes(string.concat(head, tail)))
        );
    }

    // Built in separate steps to stay clear of the compiler's stack limit.
    function _svg(string memory c, Grade g) private pure returns (string memory) {
        string memory p1 = string.concat(
            '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 400 400">',
            '<rect width="400" height="400" fill="',
            _gradeColor(g),
            '"/>'
        );
        string memory p2 = string.concat(
            '<text x="200" y="215" font-size="96" text-anchor="middle" fill="#fff" font-family="monospace">',
            c,
            "</text>"
        );
        string memory p3 = string.concat(
            '<text x="200" y="300" font-size="28" text-anchor="middle" fill="#fff" font-family="monospace">Grade ',
            _gradeName(g),
            "</text></svg>"
        );
        return string.concat(p1, p2, p3);
    }

    function _attributes(string memory c, Grade g) private pure returns (string memory) {
        string memory a1 = string.concat('[{"trait_type":"Grade","value":"', _gradeName(g), '"},');
        string memory a2 = string.concat(
            '{"trait_type":"Digits","display_type":"number","value":',
            Strings.toString(bytes(c).length),
            "}]"
        );
        return string.concat(a1, a2);
    }

    function _gradeName(Grade g) private pure returns (string memory) {
        return g == Grade.S ? "S" : g == Grade.A ? "A" : "B";
    }

    function _gradeColor(Grade g) private pure returns (string memory) {
        return g == Grade.S ? "#1d1d1f" : g == Grade.A ? "#0b5fff" : "#0f766e";
    }

    // ------------------------------------------------------------------
    // ERC-165: both parents implement it
    // ------------------------------------------------------------------

    function supportsInterface(bytes4 id)
        public
        view
        override(ERC721, AccessControlDefaultAdminRules)
        returns (bool)
    {
        return super.supportsInterface(id);
    }
}
