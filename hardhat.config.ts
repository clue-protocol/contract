import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";

// Optional: set EVM_VERSION to match the target L2 (default "paris" is the safest).
const evmVersion = process.env.EVM_VERSION || "paris";

// RPC settings come from the environment so no URL or key is ever committed.
//   RPC_URL      JSON-RPC endpoint of the target EVM L2 (required for --network l2)
//   CHAIN_ID     optional; checked against the endpoint when set
//   PRIVATE_KEY  deployer key, 0x-prefixed (use a throwaway key on testnets)
const RPC_URL = process.env.RPC_URL;
const CHAIN_ID = process.env.CHAIN_ID ? Number(process.env.CHAIN_ID) : undefined;
const accounts = process.env.PRIVATE_KEY ? [process.env.PRIVATE_KEY] : [];

const config: HardhatUserConfig = {
  solidity: {
    version: "0.8.24",
    settings: {
      evmVersion,
      optimizer: { enabled: true, runs: 200 },
    },
  },
  paths: {
    sources: "./src",
    tests: "./test",
    cache: "./cache",
    artifacts: "./artifacts",
  },
  networks: {
    hardhat: {
      // Large block so genesis batches can be benchmarked and tested.
      blockGasLimit: 60_000_000,
      allowUnlimitedContractSize: false,
    },
    // Any EVM L2 (testnet or mainnet): npx hardhat run scripts/deploy.ts --network l2
    ...(RPC_URL ? { l2: { url: RPC_URL, chainId: CHAIN_ID, accounts } } : {}),
  },
  mocha: { timeout: 600000 },
};

export default config;