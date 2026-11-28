# ClueCode

`ClueCode.sol` is the ERC-721 contract that represents ownership of a **Provider Code**, the numeric prefix of every Clue ID (for example `27` in `27-1234567890`).

The deployed contract combines three layers:

- **Code** — ERC-721 ownership, grade, canonical form and lockable metadata, implemented in `ClueCode.sol`.
- **Organisation / vault** — organisations, org roles, Ed25519 signers and code custody, implemented in the inherited `ClueCodeOrganisation.sol`.
- **Constrict** — a small ERC-165-based trust layer (`../zlibs/constrict/ClueConstrict.sol`) that lets the admin grant a sibling contract the right to set code status.

Capacity, zones, links, partitions and versions are **not** stored here; they live in ClueRegistry.

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

`IClueCode.sol` declares the NFT surface; `IClueProviderCode.sol` extends it with the organisation/vault surface and is the interface used for ERC-165 and constrict role checks. `IClueConstrict.sol` declares the constrict errors and events once for both sides.

| Function | Access | Purpose |
| --- | --- | --- |
| `gradeOf(code)` | public view | Grade, or reverts `InvalidProviderCode(code)` |
| `isValidCode(code)` | public view | True if the code is in a valid range |
| `canonical(code)` | public view | Zero-padded string, e.g. `"02"` |
| `isMinted(code)` | public view | True once the token exists |
| `mint`, `mintBatch`, `mintWithMetadata` | `MINTER_ROLE` | Mint codes, optionally with metadata |
| `totalMinted()`, `MAX_SUPPLY()` | public view | Supply counters |
| `setMetadata`, `setMetadataBatch`, `clearMetadata` | `METADATA_ROLE` | Change stored metadata |
| `setBaseURI` | `METADATA_ROLE` | Change the fallback base URI |
| `freezeCode`, `freezeMetadata` | `METADATA_ROLE` | Permanently lock metadata |
| `metadata`, `metadataOf`, `hasMetadata`, `tokenURI` | public view | Read metadata |
| `registry()` | public view | Address allowed to drive org/vault state |
| `setRegistry` | admin | Set the trusted registry |
| `codeOrg`, `codeStatus`, `totalOrganisation` | public view | Vault and code status views |
| `organisationOf`, `organisationSignersOf`, `organisationSigner`, `hasOrgRole` | public view | Organisation views |
| `createOrganisationFor`, `addOrganisationSignerFor`, `grantOrgRole`, `revokeOrgRole`, `setOrgRoleAdmin`, `setOrganisationVerificationFlag` | registry | Organisation administration on behalf of a user |
| `vaultCode`, `releaseVaultCode`, `setCodeStatus` | registry | Code custody and lifecycle |
| `grantConstrictRole` | admin | Grant a sibling contract a capability by interface id |
| `setProviderCodeStatus` | constrict role | Set status from a granted sibling contract |
| `beginDefaultAdminTransfer`, `acceptDefaultAdminTransfer`, `grantRole`, `revokeRole`, `renounceRole` | admin / role admin | OpenZeppelin access control |
| Standard ERC-721 | public | `ownerOf`, `balanceOf`, `approve`, `setApprovalForAll`, `transferFrom`, `safeTransferFrom` |

Errors and events for the metadata surface are declared in `IClueCode`, for the org/vault surface in `IClueProviderCode`, and for constrict in `IClueConstrict`. `supportsInterface` also reports ERC-4906 (`0x49064906`) and `type(IClueProviderCode).interfaceId`.

## Metadata model

- Each code may have stored metadata: name, description, image, external URL and extra traits.
- `tokenURI` prefers the stored JSON when present. A code without stored metadata falls back to `baseURI + canonical(code)` when a base URI is set, otherwise to generated on-chain JSON.
- Grade and Digits traits are always derived from the token id and cannot be overridden.
- `freezeCode(code)` locks one code; `freezeMetadata()` locks every code and the base URI. Both are permanent.
- `setBaseURI` emits `BatchMetadataUpdate(1, MAX_SUPPLY)` so indexers can refresh existing tokens.

## Access control

| Role / capability | Holder | Can do |
| --- | --- | --- |
| `DEFAULT_ADMIN_ROLE` | multisig | Grant/revoke roles, set the registry, grant constrict roles. Two-step delayed admin transfer |
| `MINTER_ROLE` | genesis distributor | `mint`, `mintBatch`, `mintWithMetadata` |
| `METADATA_ROLE` | metadata operator | Set/clear metadata, set base URI, freeze metadata |
| `registry` address | ClueRegistry | Create organisations, manage org roles/signers, vault/release codes and set status |
| Constrict role for `IClueProviderCode` | granted sibling | `setProviderCodeStatus` |

## Organisation / vault model

- `createOrganisationFor` creates an organisation and gives the initiator the default admin, controller and updater roles.
- Ed25519 signers are append-only. `organisationSigner` returns the first signer; `organisationSignersOf` returns all of them.
- `verificationFlag` mirrors the backend `Organisation.VerificationFlag` and is set by an org admin through the registry.
- Role admins default to the org admin role. An org admin may delegate a role's admin, or permanently lock it with `ClueOrgRoles.NO_ROLE_ADMIN`.
- `CodeStatus` moves between `Open` (not vaulted), `Stateless`, `Active` and `Archived`. `Open` is set only by the release path.
- `registry` is the only address allowed to call organisation/vault mutations; it forwards the original caller explicitly for role checks.

## Genesis and lock-down

1. Deploy with a multisig admin, an admin delay, the distributor as minter and an initial base URI.
2. Mint all 99 + 900 + 90,000 codes in batches.
3. Verify `totalMinted() == 90999` and spot-check `ownerOf` and `canonical`.
4. Distribute codes with ordinary ERC-721 transfers, or hold them for sale.
5. Set the registry with `setRegistry(registryAddress)`.
6. Renounce roles when done:
   - minter: `renounceRole(MINTER_ROLE, minter)`;
   - metadata operator: `renounceRole(METADATA_ROLE, account)`;
   - admin: `beginDefaultAdminTransfer(address(0))`, wait out the delay, then `renounceRole(DEFAULT_ADMIN_ROLE, admin)`.
7. Confirm no role holder can mint, change metadata or change the registry.

Renouncing `MINTER_ROLE` alone is not enough, because the admin could grant it again. `acceptDefaultAdminTransfer` cannot hand the admin role to the zero address; the delayed `renounceRole` above is the correct path.

## Interaction with ClueRegistry

The registry is the only contract that treats a Provider Code specially:

- **Vault:** the registry checks the owner and transfers the token into this contract. No separate ERC-721 approval is needed because the registry forwards the owning caller.
- **Status:** the registry drives `Stateless → Active → Archived` and back through `setCodeStatus`; a granted constrict sibling can also set status.
- **Release:** the registry returns the token with an internal transfer, and only while the code is `Stateless`.
- **Resolve:** the registry reads `codeStatus` and the organisation signer to build the on-chain anchor; ownership alone gives no signing authority.

Internal `_transfer` is used for vaulting and release, so no receiver callback runs during registry-driven custody changes.

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

Hardhat tests are still to be added under `../../test/`.

## Build outputs

`make compile` writes the ABI to `out/abi/ClueCode.abi`, the bytecode to `out/bin/ClueCode.bin` (hex, no `0x`), and, if `abigen` is installed, Go bindings to `out/go/code/ClueCode.go` (package `code`, type `ClueCode`).

## Open items

- Hardhat tests for grade/canonical, metadata freeze, org roles/signers, vault/release and ERC-165.
- Registry-side wiring: `setRegistry` must be called after deployment, and the registry must satisfy the constrict interface if it needs to set code status.
- Final metadata image strategy; deployed bytecode is close to the EIP-170 limit.
- Maximum batch size on the chosen L2.
- Royalties (ERC-2981). The default is none.
- Audit scope and timing.
