# ============================================================
# ClueCode Smart Contracts — Makefile
# ============================================================
# Targets:
#   make init        – install npm dependencies, create out/ dirs
#   make compile     – compile Solidity and export out/abi + out/bin
#   make test        – run the Hardhat test suite
#   make bindings    – generate Go bindings with abigen (out/go/code)
#   make clean       – remove build artefacts
#   make clean-all   – remove build artefacts + node_modules
# ============================================================

.DEFAULT_GOAL := help

.PHONY: help init compile abi test bindings clean clean-abi clean-bin clean-bindings clean-all go

CONTRACT_DIR := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

# Writable cache dir for Hardhat (avoids EROFS errors on /home/.local in some envs)
HARDHAT_HOME := $(CONTRACT_DIR)/.hardhat-home

HARDHAT_ENV := HARDHAT_DISABLE_TELEMETRY_PROMPT=true \
               XDG_DATA_HOME=$(HARDHAT_HOME) \
               XDG_CONFIG_HOME=$(HARDHAT_HOME)

OUT_DIR      := $(CONTRACT_DIR)/out
CLUECODE_SRC := $(CONTRACT_DIR)/src/cluecode/ClueCode.sol

# abigen settings (Go package / type names for the generated binding)
GO_DIR  := $(OUT_DIR)/go/code
GO_PKG  := code
GO_TYPE := ClueCode

# ── Help ─────────────────────────────────────────────────────

help: ## Show this help
	@echo "ClueCode Smart Contracts"
	@echo ""
	@echo "Usage: make <target>"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# ── Dependency installation ──────────────────────────────────

init: ## Install npm dependencies and create out/ dirs
	@echo "→ Installing npm dependencies..."
	cd $(CONTRACT_DIR) && npm install 
	@echo "→ Setting up out/ directories..."
	mkdir -p $(OUT_DIR)/abi $(OUT_DIR)/bin
	@echo "✓ Init complete. Run 'make compile' to build."

# ── Build ────────────────────────────────────────────────────

compile: ## Compile Solidity and export ABI + bytecode
	@echo "→ Compiling Solidity contracts..."
	mkdir -p $(HARDHAT_HOME) $(OUT_DIR)/abi $(OUT_DIR)/bin
	cd $(CONTRACT_DIR) && $(HARDHAT_ENV) npx hardhat compile
	@$(MAKE) --no-print-directory abi
	@echo "✓ Compilation done. Artefacts in artifacts/, ABI in out/abi/, bytecode in out/bin/"

abi: ## Export ABI and bytecode from compiled artefacts
	@echo "→ Exporting ABI and bytecode..."
	cd $(CONTRACT_DIR) && npx ts-node scripts/export-abi-bin.ts

# ── Go bindings (abigen) ─────────────────────────────────────

go: bindings ## Alias for 'bindings'

bindings: compile ## Generate Go bindings with abigen (out/go/code)
	@command -v abigen >/dev/null 2>&1 || (echo "abigen not found: install go-ethereum's abigen"; exit 1)
	mkdir -p $(GO_DIR)
	abigen --abi $(OUT_DIR)/abi/$(GO_TYPE).abi \
	       --bin $(OUT_DIR)/bin/$(GO_TYPE).bin \
	       --pkg $(GO_PKG) --type $(GO_TYPE) \
	       --out $(GO_DIR)/$(GO_TYPE).go
	@echo "✓ Bindings written to out/go/code/$(GO_TYPE).go"

# ── Test ─────────────────────────────────────────────────────

test: ## Run the Hardhat test suite
	@echo "→ Running Hardhat tests..."
	mkdir -p $(HARDHAT_HOME)
	cd $(CONTRACT_DIR) && $(HARDHAT_ENV) npx hardhat test

# ── Clean ────────────────────────────────────────────────────

clean: ## Remove build artefacts (artifacts/, cache/, out/)
	@echo "→ Removing build artefacts..."
	mkdir -p $(HARDHAT_HOME)
	cd $(CONTRACT_DIR) && $(HARDHAT_ENV) npx hardhat clean
	rm -rf $(HARDHAT_HOME) $(OUT_DIR)
	@echo "✓ Clean done."

clean-abi: ## Remove exported ABI (out/abi)
	rm -rf $(OUT_DIR)/abi

clean-bin: ## Remove exported bytecode (out/bin)
	rm -rf $(OUT_DIR)/bin

clean-bindings: ## Remove generated Go bindings (out/go)
	rm -rf $(OUT_DIR)/go

clean-all: clean ## Remove build artefacts + node_modules
	@echo "→ Removing node_modules..."
	rm -rf $(CONTRACT_DIR)/node_modules
	@echo "✓ Full clean done."