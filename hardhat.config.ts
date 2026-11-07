import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";

// Optional: set EVM_VERSION to match the target L2 (default "paris" is the safest).
const evmVersion = process.env.EVM_VERSION || "paris";

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
  },
  mocha: { timeout: 600000 },
};

export default config;
