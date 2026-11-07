# Clue Protocol: Contracts

Smart contracts for the **Clue Protocol** (Crypto Identification Secure Protocol, CISP): an identity layer for digital finance anchored on an EVM Layer 2.

The protocol uses two contracts. This repository currently contains the first.

| Contract | Role | Status |
| --- | --- | --- |
| **ClueCode** | ERC-721 representing ownership of a Provider Code (the numeric prefix of a Clue ID, e.g. `27` in `27-1234567890`) | Implemented in [`src/code`](src/code) |
| **ClueRegistry** | Code Vault, versioning, validator-approved Merkle roots, capacity, links and partitions | Not yet in this repo |

ClueCode answers one question: *who owns this code?* Everything else (is it live, what version, which root) belongs to ClueRegistry.

> The contracts are a reference implementation. They have **not been audited** and must not be deployed to mainnet before an independent audit.

## Repository layout

```
.
├── src/code/            Solidity sources (ClueCode)
├── test/                Hardhat tests (TypeScript)
├── scripts/             deploy, genesis mint, ABI/bin export (TypeScript)
├── out/
│   ├── abi/             ClueCode.abi
│   ├── bin/             ClueCode.bin
│   └── go/code/         abigen Go bindings (package `code`)
├── hardhat.config.ts
├── Makefile
├── package.json / tsconfig.json
└── go.mod               Go module for the generated bindings
```

`artifacts/`, `cache/`, `typechain-types/`, `node_modules/` and `.hardhat-home/` are build output.

## Requirements

- Node.js 18+ and npm
- GNU Make
- Optional, for Go bindings: [`abigen`](https://geth.ethereum.org/docs/tools/abigen) from go-ethereum

## Quick start

```sh
make init       # npm install, create out/ directories
make compile    # compile, export out/abi + out/bin, and out/go/code if abigen is installed
make test       # run the Hardhat test suite
```

| Target | What it does |
| --- | --- |
| `make init` | Installs npm dependencies and creates `out/abi` and `out/bin` |
| `make compile` | Compiles Solidity, generates typechain types, exports ABI and bytecode, and exports Go bindings when `abigen` is on `PATH` (skips quietly otherwise) |
| `make abi` | Re-exports `out/abi/ClueCode.abi` and `out/bin/ClueCode.bin` from existing artifacts |
| `make export-go` | Runs `abigen` into `out/go/code/ClueCode.go` (fails if `abigen` is missing). Aliases: `make go`, `make bindings` |
| `make test` | Runs the Hardhat tests |
| `make clean` | Removes Hardhat artifacts and `out/` |
| `make clean-abi` / `clean-bin` / `clean-bindings` | Removes `out/abi`, `out/bin`, or `out/go` |
| `make clean-all` | `clean` plus `node_modules` |

To also mint all 90,999 codes in a test (slow): `FULL_GENESIS=1 make test`.

Set `EVM_VERSION` (default `paris`) to match the target L2, for example `EVM_VERSION=cancun make compile`.

## Using the Go bindings

The generated package lives in `out/go/code`. Assuming the module path in `go.mod` is `github.com/clue-protocol/contract`:

```go
import "github.com/clue-protocol/contract/out/go/code"

c, err := code.NewClueCode(address, client) // client: *ethclient.Client
grade, err := c.GradeOf(nil, big.NewInt(27))
canon, err := c.Canonical(nil, big.NewInt(2)) // "02"
```

## Deployment

```sh
# 1. Deploy (ADMIN is a multisig; ADMIN_DELAY is in seconds)
ADMIN=0x… MINTER=0x… ADMIN_DELAY=86400 \
  npx hardhat run scripts/deploy.ts --network <network>

# 2. Genesis: mint all 90,999 codes, verify, optionally renounce MINTER_ROLE
CLUECODE=0x… BATCH=250 DISTRIBUTOR=0x… RENOUNCE_MINTER=1 \
  npx hardhat run scripts/genesis.ts --network <network>
```

Benchmark `BATCH` against the target L2's block gas limit before mainnet. After genesis the admin multisig finishes the lock-down (see [`src/code/README.md`](src/code/README.md#genesis-and-lock-down)). Add network entries and keys to `hardhat.config.ts` locally; never commit private keys.

## Design principles

- **Ownership is an on-chain asset.** A Provider Code is a token, not a permission granted by the Foundation.
- **Minimal footprint.** Contracts hold ownership and roots, never identity records or KYC/KYB evidence.
- **No silent key transfer.** Moving a token never moves provider signing keys.
- **One canonical form.** Every code has exactly one decimal representation shared by Solidity, Go and TypeScript.
- **Immutable.** No proxies; no role can move, burn or pause a token.

## License

MIT. See [LICENSE](LICENSE).
