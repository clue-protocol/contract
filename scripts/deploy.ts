// Usage:
//   RPC_URL=https://<l2-rpc> PRIVATE_KEY=0x.. [CHAIN_ID=..] \
//   ADMIN=0x.. MINTER=0x.. ADMIN_DELAY=86400 \
//   npx hardhat run scripts/deploy.ts --network l2
// The `l2` network is defined in hardhat.config.ts from RPC_URL / CHAIN_ID / PRIVATE_KEY.
import { ethers, network } from "hardhat";

async function main() {
  const signers = await ethers.getSigners();
  if (signers.length === 0) throw new Error("No deployer account: set PRIVATE_KEY");
  const [deployer] = signers;
  const chain = await ethers.provider.getNetwork();
  console.log(`Network: ${network.name} (chainId ${chain.chainId}), deployer ${deployer.address}`);
  const admin = process.env.ADMIN || deployer.address;
  const minter = process.env.MINTER || deployer.address;
  const delay = Number(process.env.ADMIN_DELAY || 86400);

  const factory = await ethers.getContractFactory("ClueCode");
  const c = await factory.deploy(admin, delay, minter);
  await c.waitForDeployment();
  console.log("ClueCode:", await c.getAddress());
  console.log({
    chainId: Number(chain.chainId),
    rpc: process.env.RPC_URL ?? "(hardhat in-process)",
    admin,
    minter,
    adminDelay: delay,
  }); // record in the deployment manifest
}

main().catch((e) => {
  console.error(e);
  process.exitCode = 1;
});