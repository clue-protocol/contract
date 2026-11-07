#!/usr/bin/env ts-node
// ============================================================
// Export ClueCode ABI + bytecode from Hardhat artifacts
// ============================================================
import * as fs from "fs";
import * as path from "path";

const ROOT = path.resolve(__dirname, "..");
const OUT_DIR = path.join(ROOT, "out");
const CONTRACT_NAME = "ClueCode";

// Optional override (used by Makefile: CLUECODE_ARTIFACT=...)
const ARTIFACT_OVERRIDE = process.env.CLUECODE_ARTIFACT;

// ── Artifact discovery ──────────────────────────────────────
function findArtifact(): string {
  if (ARTIFACT_OVERRIDE) {
    if (!fs.existsSync(ARTIFACT_OVERRIDE)) {
      throw new Error(`CLUECODE_ARTIFACT points to a non-existent file: ${ARTIFACT_OVERRIDE}`);
    }
    return ARTIFACT_OVERRIDE;
  }

  // Common candidates depending on paths.sources config
  const candidates = [
    path.join(ROOT, "artifacts/cluecode/ClueCode.sol", `${CONTRACT_NAME}.json`),          // sources = ./src
    path.join(ROOT, "artifacts/src/cluecode/ClueCode.sol", `${CONTRACT_NAME}.json`),      // sources = .
    path.join(ROOT, "artifacts/contracts/src/cluecode/ClueCode.sol", `${CONTRACT_NAME}.json`), // sources = ./contracts
  ];
  for (const c of candidates) {
    if (fs.existsSync(c)) return c;
  }

  // Last resort: recursive search
  const found = recursiveFind(path.join(ROOT, "artifacts"), `${CONTRACT_NAME}.json`);
  if (found) return found;

  throw new Error(
    `Artifact not found for ${CONTRACT_NAME}. ` +
    `Run 'make compile' first. Searched:\n  ` + candidates.join("\n  ")
  );
}

function recursiveFind(dir: string, name: string): string | null {
  if (!fs.existsSync(dir)) return null;
  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    const full = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      const hit = recursiveFind(full, name);
      if (hit) return hit;
    } else if (entry.name === name) {
      return full;
    }
  }
  return null;
}

// ── Main ────────────────────────────────────────────────────
function main() {
  const artifactPath = findArtifact();
  console.log(`→ Using artifact: ${path.relative(ROOT, artifactPath)}`);

  const artifact = JSON.parse(fs.readFileSync(artifactPath, "utf8"));

  fs.mkdirSync(path.join(OUT_DIR, "abi"), { recursive: true });
  fs.mkdirSync(path.join(OUT_DIR, "bin"), { recursive: true });

  const abiPath = path.join(OUT_DIR, "abi", `${CONTRACT_NAME}.abi`);
  const binPath = path.join(OUT_DIR, "bin", `${CONTRACT_NAME}.bin`);

  fs.writeFileSync(abiPath, JSON.stringify(artifact.abi, null, 2));

  // abigen expects raw hex bytecode (no 0x prefix)
  const bytecode: string = artifact.bytecode ?? "";
  fs.writeFileSync(binPath, bytecode.replace(/^0x/, ""));

  console.log(`✓ ABI  → ${path.relative(ROOT, abiPath)}`);
  console.log(`✓ BIN  → ${path.relative(ROOT, binPath)}`);
}

main();