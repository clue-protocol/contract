# ClueCode

`ClueCode.sol` is the ERC-721 contract that represents ownership of a **Provider Code**, the numeric prefix of every Clue ID (for example `27` in `27-1234567890`).

ClueCode is intentionally small. It stores ownership and derives grade and canonical form from the token id. It stores no capacity, zones, links, partitions, versions or status. Those live in ClueRegistry, which also holds the Code Vault.

> Reference implementation, **not audited**.

## Token model

The ERC-721 `tokenId` is the Provider Code read as an integer. Grades are set by digit length, so the ranges never overlap and grade is read directly from the id.

| Grade | Enum | Digits | Range | Supply | Canonical width |
| --- | --- | --- | --- | --- | --- |
| S | 1 | 2 | 01–99 | 99 | 2 |
| A | 2 | 3 | 100–999 | 900 | 3 |
| B | 3 | 5 | 10000–99999 | 90,000 | 5 |

Total supply is **90,999**. Zero, `1000–9999` and anything above `99999` are invalid and can never be minted. Enum value `0` (`None`) marks an invalid code.

**Canonical form.** A code is written at the fixed width of its grade: `02` is valid, `2` is not. The Damm check digit is computed over the provider code plus the nine-digit serial, so a leading-zero ambiguity would give different results in different clients. `canonical()` is the single source of truth for Solidity, Go and TypeScript.

## Interface

| Function | Access | Purpose |
| --- | --- | --- |
| `gradeOf(code)` | public view | Grade, or reverts `InvalidProviderCode(code)` |
| `isValidCode(code)` | public view | True if the code is in a valid range |
| `canonical(code)` | public view | Zero-padded string, e.g. `"02"` |
| `mint(to, code)` | `MINTER_ROLE` | Mint one code |
| `mintBatch(to, codes)` | `MINTER_ROLE` | Mint several codes atomically |
| `totalMinted()` | public view | Codes minted so far |
| `MAX_SUPPLY()` | public view | `90999` (convenience for clients and tests) |
| `tokenURI(code)` | public view | On-chain base64 JSON with an SVG image |
| `beginDefaultAdminTransfer`, `acceptDefaultAdminTransfer` | admin / pending admin | Two-step, delayed admin change |
| `grantRole`, `revokeRole`, `renounceRole` | role admin / holder | OpenZeppelin role management |
| Standard ERC-721 | public | `ownerOf`, `balanceOf`, `approve`, `setApprovalForAll`, `transferFrom`, `safeTransferFrom` |

The only custom error is `InvalidProviderCode(uint256 code)`. Events are the inherited ERC-721 `Transfer` and `Approval`, plus the AccessControl events.

## Design decisions

- **Fixed, bounded supply.** Only valid codes can be minted and each token exists once, so supply cannot exceed 90,999 without a separate cap check.
- **Not upgradeable.** Roles are only the admin and the minter. Neither can transfer, burn or pause a token or change its metadata, so the Foundation holds no ultimate authority over code assignment.
- **Fully on-chain metadata.** `tokenURI` returns a base64 JSON data URI (name, grade, digit count) with an SVG of the canonical code on a grade-coloured background. Nothing is hosted and nothing can be changed. Status, links and capacity are read from ClueRegistry, never from the token.
- **No `ERC721Enumerable`.** It adds storage cost on every transfer. Indexers list holdings from `Transfer` events.
- **`_mint`, not `_safeMint`,** so a batch can go to a distributor contract.
- **No hooks into ClueRegistry.** The NFT stays usable on any marketplace and the registry can be redeployed without touching ownership.
- **Key separation.** Owning a code gives no signing authority. Profile Server keys are registered in ClueRegistry, so transferring the NFT cannot silently transfer provider keys.

## Access control

ClueCode uses OpenZeppelin's `AccessControlDefaultAdminRules` (OpenZeppelin 5.0.2).

| Role | Holder | Can do |
| --- | --- | --- |
| `DEFAULT_ADMIN_ROLE` | A multisig | Grant and revoke `MINTER_ROLE`. Changing the admin is two-step and takes effect after the configured delay |
| `MINTER_ROLE` | Genesis distributor | `mint` and `mintBatch` |

## Genesis and lock-down

1. Deploy with a multisig admin, an admin delay, and the distributor as minter (`scripts/deploy.ts`).
2. Mint all 99 + 900 + 90,000 codes in batches (`scripts/genesis.ts`).
3. Verify `totalMinted() == 90999` and spot-check `ownerOf` and `canonical`.
4. Distribute codes with ordinary ERC-721 transfers, or hold them for sale.
5. Renounce both roles:
   - The minter calls `renounceRole(MINTER_ROLE, minter)`.
   - The admin calls `beginDefaultAdminTransfer(address(0))`, waits out the delay, then calls `renounceRole(DEFAULT_ADMIN_ROLE, admin)`.
6. Confirm `hasRole` is false for both and that a further `mint` reverts.

Renouncing `MINTER_ROLE` alone is not enough, because the admin could grant it again. Note that `acceptDefaultAdminTransfer` cannot be used to hand the admin to the zero address; the delayed `renounceRole` above is the correct path.

## Interaction with ClueRegistry

ClueRegistry, through its Code Vault, is the only contract that treats ClueCode specially:

- **Activate:** the owner approves the registry; the registry checks `ownerOf` and pulls the token with `transferFrom`.
- **Root update:** the registry checks that it owns the code. This is the on-chain test for "held in the Code Vault".
- **Grade checks:** link and partition rules call `gradeOf`; Route codes must be Grade S or A.
- **Release:** the registry returns the token with `transferFrom`.

The registry must use `transferFrom` (not `safeTransferFrom`) on release so a depositor contract that cannot receive ERC-721 tokens does not lock the code, write state before the external transfer, and not rely on balance counts.

## Test vectors

These are shared with the Go backend and TypeScript clients and must give identical results everywhere.

| Input | Result |
| --- | --- |
| `2` | Valid, canonical `02`, Grade S |
| `"2"` as a string | Rejected: not canonical (off-chain clients) |
| `99` | Valid, `99`, Grade S |
| `100` / `999` | Valid, Grade A |
| `1000`–`9999` | Rejected: no four-digit grade |
| `10000` / `99999` | Valid, Grade B |
| `0` and `100000` | Rejected |

Tests live in [`../../test/ClueCode.ts`](../../test/ClueCode.ts).

## Build outputs

`make compile` writes the ABI to `out/abi/ClueCode.abi`, the bytecode to `out/bin/ClueCode.bin` (hex, no `0x`), and, if `abigen` is installed, Go bindings to `out/go/code/ClueCode.go` (package `code`, type `ClueCode`).

## Open items

- Initial distribution of the 90,999 tokens, including any reserved Grade S codes.
- Final SVG design and resulting bytecode size.
- Maximum batch size on the chosen L2.
- Royalties (ERC-2981). The default is none.
- What happens to Clue IDs under a code when the NFT leaves the vault (a ClueRegistry decision).
- Audit scope and timing.
