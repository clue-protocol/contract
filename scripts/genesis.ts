// Mints all 90,999 codes to the distributor, verifies, then optionally renounces MINTER_ROLE.
// Usage: CLUECODE=0x.. BATCH=250 [DISTRIBUTOR=0x..] [RENOUNCE_MINTER=1] \
//        npx hardhat run scripts/genesis.ts --network <net>
// Signer must hold MINTER_ROLE. Benchmark BATCH against the target L2 gas limit first.
// The admin renounce (beginDefaultAdminTransfer(0) -> wait delay -> renounceRole) is a
// separate step for the admin multisig; it is not performed here.
import { ethers } from "hardhat";

function* allCodes(): Generator<number> {
  for (let c = 1; c <= 99; c++) yield c;
  for (let c = 100; c <= 999; c++) yield c;
  for (let c = 10000; c <= 99999; c++) yield c;
}

async function main() {
  const address = process.env.CLUECODE;
  if (!address) throw new Error("CLUECODE env var is required");

  const [minter] = await ethers.getSigners();
  const to = process.env.DISTRIBUTOR || minter.address;
  const batch = Number(process.env.BATCH || 250);
  const c = await ethers.getContractAt("ClueCode", address, minter);

  let buf: number[] = [];
  const flush = async () => {
    if (!buf.length) return;
    await (await c.mintBatch(to, buf)).wait();
    console.log(`minted ${buf[0]}..${buf[buf.length - 1]} (${buf.length})`);
    buf = [];
  };
  for (const code of allCodes()) {
    buf.push(code);
    if (buf.length === batch) await flush();
  }
  await flush();

  const total = await c.totalMinted();
  if (total !== 90999n) throw new Error(`totalMinted=${total}, expected 90999`);
  const spot: Array<[number, string]> = [[2, "02"], [99, "99"], [100, "100"], [10000, "10000"]];
  for (const [code, canon] of spot) {
    if ((await c.canonical(code)) !== canon || (await c.ownerOf(code)) !== to) {
      throw new Error(`spot check failed for ${code}`);
    }
  }
  console.log("verified: 90,999 minted");

  if (process.env.RENOUNCE_MINTER === "1") {
    await (await c.renounceRole(await c.MINTER_ROLE(), minter.address)).wait();
    console.log("MINTER_ROLE renounced");
  }
}

main().catch((e) => {
  console.error(e);
  process.exitCode = 1;
});
