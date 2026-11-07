// Usage: ADMIN=0x.. MINTER=0x.. ADMIN_DELAY=86400 npx hardhat run scripts/deploy.ts --network <net>
import { ethers } from "hardhat";

async function main() {
  const [deployer] = await ethers.getSigners();
  const admin = process.env.ADMIN || deployer.address;
  const minter = process.env.MINTER || deployer.address;
  const delay = Number(process.env.ADMIN_DELAY || 86400);

  const factory = await ethers.getContractFactory("ClueCode");
  const c = await factory.deploy(admin, delay, minter);
  await c.waitForDeployment();
  console.log("ClueCode:", await c.getAddress());
  console.log({ admin, minter, adminDelay: delay }); // record in the deployment manifest
}

main().catch((e) => {
  console.error(e);
  process.exitCode = 1;
});
