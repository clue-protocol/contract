// Package tests holds the Go tests for the ClueCode bindings in out/go/code.
//
// Every test except TestDeploy uses the global config (CLUE_* env vars) and runs against an
// already deployed contract. TestDeploy has its own DEPLOY_* parameters and deploys a fresh
// contract. See tests/.env.example.
package tests

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	cluecode "github.com/clue-protocol/contract/out/go/code"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pat3icki/pkgt/envt"
)

// ======================================================================
// Global configuration (all tests except TestDeploy)
// ======================================================================
// Global test configuration.
//
// Every test except TestDeploy uses this shared config. It points at an
// already deployed ClueCode contract and is read from the environment:
//
//	CLUE_RPC_URL             JSON-RPC endpoint                      (required; tests skip when unset)
//	CLUE_CHAIN_ID            chain id; checked against the endpoint (optional, derived when unset)
//	CLUE_PRIVATE_KEY         default signer, hex                    (required for state-changing tests)
//	CLUE_CONTRACT_ADDRESS    deployed ClueCode address              (required)
//	CLUE_MINTER_PRIVATE_KEY  key holding MINTER_ROLE                (optional, defaults to CLUE_PRIVATE_KEY)
//	CLUE_ADMIN_ADDRESS       expected default admin                 (optional, enables an extra check)
//	CLUE_TX_TIMEOUT          seconds to wait for a transaction      (optional, default 120)
//
// State-changing tests (mint, transfer) skip themselves when the signer is not a minter.
// Config is the global configuration shared by all tests except TestDeploy.
type Config struct {
	RPCURL    string
	ChainID   *big.Int
	Contract  common.Address
	Key       *ecdsa.PrivateKey // default signer
	MinterKey *ecdsa.PrivateKey // signer used for mint calls
	Admin     *common.Address   // optional expected default admin
	TxTimeout time.Duration
}

var (
	cfg    *Config
	cfgErr error
	client *ethclient.Client
	cc     *cluecode.ClueCode
)

func TestMain(m *testing.M) {
	cfg, cfgErr = loadConfig()
	if cfg != nil && cfgErr == nil {
		cfgErr = connect()
	}
	rc := m.Run()
	if client != nil {
		client.Close()
	}
	os.Exit(rc)
}

// loadConfig returns (nil, nil) when CLUE_RPC_URL is not set, so the suite skips cleanly.
func loadConfig() (*Config, error) {
	env, err := envt.Map(&envt.Config{
		EnvPath:  "tests/.env",
		PanicErr: true,
	})
	if err != nil {
		return nil, err
	}

	rpc := env["CLUE_RPC_URL"]
	if rpc == "" {
		return nil, nil
	}
	c := &Config{RPCURL: rpc, TxTimeout: 120 * time.Second}

	if s := env["CLUE_CHAIN_ID"]; s != "" {
		id, ok := new(big.Int).SetString(s, 10)
		if !ok {
			return nil, fmt.Errorf("CLUE_CHAIN_ID %q is not a decimal integer", s)
		}
		c.ChainID = id
	}

	addr := env["CLUE_CONTRACT_ADDRESS"]
	if !common.IsHexAddress(addr) {
		return nil, fmt.Errorf("CLUE_CONTRACT_ADDRESS %q is not a valid address", addr)
	}
	c.Contract = common.HexToAddress(addr)

	if s := env["CLUE_PRIVATE_KEY"]; s != "" {
		k, err := parseKey(s)
		if err != nil {
			return nil, fmt.Errorf("CLUE_PRIVATE_KEY: %w", err)
		}
		c.Key = k
		c.MinterKey = k
	}

	if s := env["CLUE_MINTER_PRIVATE_KEY"]; s != "" {
		k, err := parseKey(s)
		if err != nil {
			return nil, fmt.Errorf("CLUE_MINTER_PRIVATE_KEY: %w", err)
		}
		c.MinterKey = k
	}

	if s := env["CLUE_ADMIN_ADDRESS"]; s != "" {
		if !common.IsHexAddress(s) {
			return nil, fmt.Errorf("CLUE_ADMIN_ADDRESS %q is not a valid address", s)
		}
		a := common.HexToAddress(s)
		c.Admin = &a
	}

	if s := env["CLUE_TX_TIMEOUT"]; s != "" {
		n, err := strconv.Atoi(s)
		if err != nil || n <= 0 {
			return nil, fmt.Errorf("CLUE_TX_TIMEOUT %q must be a positive number of seconds", s)
		}
		c.TxTimeout = time.Duration(n) * time.Second
	}

	return c, nil
}

func parseKey(s string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(strings.TrimPrefix(strings.TrimSpace(s), "0x"))
}

// connect dials the RPC, verifies the chain id and that code exists at the address, and binds ClueCode.
func connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cl, err := ethclient.DialContext(ctx, cfg.RPCURL)
	if err != nil {
		return fmt.Errorf("dial %s: %w", cfg.RPCURL, err)
	}
	client = cl
	remote, err := cl.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("fetch chain id: %w", err)
	}
	if cfg.ChainID == nil {
		cfg.ChainID = remote
	} else if cfg.ChainID.Cmp(remote) != 0 {
		return fmt.Errorf("CLUE_CHAIN_ID=%s but the endpoint reports chain id %s", cfg.ChainID, remote)
	}
	bytecode, err := cl.CodeAt(ctx, cfg.Contract, nil)
	if err != nil {
		return fmt.Errorf("fetch code at %s: %w", cfg.Contract, err)
	}
	if len(bytecode) == 0 {
		return fmt.Errorf("no contract code at %s on chain %s", cfg.Contract, cfg.ChainID)
	}
	cc, err = cluecode.NewClueCode(cfg.Contract, cl)
	return err
}

// requireConfig skips the test when the suite is not configured, and fails it when the config is broken.
func requireConfig(t *testing.T) {
	t.Helper()
	if cfgErr != nil {
		t.Fatalf("test configuration error: %v", cfgErr)
	}
	if cfg == nil {
		t.Skip("CLUE_RPC_URL not set; see tests/.env.example")
	}
}

// requireSigner additionally needs CLUE_PRIVATE_KEY.
func requireSigner(t *testing.T) {
	t.Helper()
	requireConfig(t)
	if cfg.Key == nil {
		t.Skip("CLUE_PRIVATE_KEY not set")
	}
}

// requireMinter additionally needs a signer that currently holds MINTER_ROLE.
func requireMinter(t *testing.T) {
	t.Helper()
	requireSigner(t)
	role, err := cc.MINTERROLE(nil)
	if err != nil {
		t.Fatalf("MINTER_ROLE: %v", err)
	}
	has, err := cc.HasRole(nil, role, crypto.PubkeyToAddress(cfg.MinterKey.PublicKey))
	if err != nil {
		t.Fatalf("hasRole: %v", err)
	}
	if !has {
		t.Skip("signer does not hold MINTER_ROLE (renounced or wrong key)")
	}
}

func addrOf(k *ecdsa.PrivateKey) common.Address { return crypto.PubkeyToAddress(k.PublicKey) }

func txOpts(t *testing.T, k *ecdsa.PrivateKey) *bind.TransactOpts {
	t.Helper()
	auth, err := bind.NewKeyedTransactorWithChainID(k, cfg.ChainID)
	if err != nil {
		t.Fatalf("transactor: %v", err)
	}
	return auth
}

// mine waits for the transaction and fails unless it succeeded.
func mine(t *testing.T, tx *types.Transaction, err error) *types.Receipt {
	t.Helper()
	if err != nil {
		t.Fatalf("send transaction: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), cfg.TxTimeout)
	defer cancel()
	r, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		t.Fatalf("wait for %s: %v", tx.Hash(), err)
	}
	if r.Status != types.ReceiptStatusSuccessful {
		t.Fatalf("transaction %s reverted", tx.Hash())
	}
	return r
}

func randomAddress(t *testing.T) common.Address {
	t.Helper()
	k, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	return addrOf(k)
}

func bigInt(n uint64) *big.Int { return new(big.Int).SetUint64(n) }

// ---- revert assertions ----
func revertData(err error) []byte {
	var de interface{ ErrorData() interface{} }
	if errors.As(err, &de) {
		if s, ok := de.ErrorData().(string); ok {
			if b, e := hexutil.Decode(s); e == nil {
				return b
			}
		}
	}
	return nil
}

// isRevert reports whether err is a revert with the named custom error from the ClueCode ABI.
// When the node returns no revert data, any error is accepted.
func isRevert(err error, name string) bool {
	if err == nil {
		return false
	}
	parsed, perr := cluecode.ClueCodeMetaData.GetAbi()
	if perr != nil {
		return false
	}
	e, ok := parsed.Errors[name]
	if !ok {
		return false
	}
	data := revertData(err)
	if data == nil {
		return true
	}
	return len(data) >= 4 && bytes.Equal(data[:4], e.ID[:4])
}

func requireRevert(t *testing.T, err error, name string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected revert %s, call succeeded", name)
	}
	if !isRevert(err, name) {
		t.Fatalf("expected revert %s, got: %v", name, err)
	}
}

// ======================================================================
// Deploy test (separate DEPLOY_* config)
// ======================================================================
// TestDeploy is deliberately separate from the global config above: it does not
// use an existing contract address. It deploys a fresh ClueCode and takes its own parameters
// from DEPLOY_* environment variables:
//
//	DEPLOY_RPC_URL        JSON-RPC endpoint            (required; the test skips when unset)
//	DEPLOY_CHAIN_ID       chain id, checked vs the RPC (required)
//	DEPLOY_PRIVATE_KEY    deployer key, hex            (required; must hold gas funds)
//	DEPLOY_ADMIN          default admin, e.g. a multisig   (optional, default: deployer)
//	DEPLOY_MINTER         genesis distributor              (optional, default: deployer)
//	DEPLOY_ADMIN_DELAY    admin-change delay in seconds    (optional, default 86400)
//	DEPLOY_OUT_FILE       write {chainId, address, ...} JSON here (optional)
//	DEPLOY_TX_TIMEOUT     seconds to wait for the deploy tx (optional, default 180)
//
// Point CLUE_CONTRACT_ADDRESS at the printed address to run the rest of the suite against it.
// DeployConfig holds the parameters of TestDeploy only.
type DeployConfig struct {
	RPCURL     string
	ChainID    *big.Int
	Key        *ecdsa.PrivateKey
	Admin      common.Address
	Minter     common.Address
	AdminDelay *big.Int
	OutFile    string
	TxTimeout  time.Duration
}

// loadDeployConfig returns (nil, nil) when DEPLOY_RPC_URL is not set.
func loadDeployConfig() (*DeployConfig, error) {
	env, err := envt.Map(&envt.Config{
		EnvPath:  "tests/.env",
		PanicErr: true,
	})
	if err != nil {
		return nil, err
	}

	rpc := env["DEPLOY_RPC_URL"]
	if rpc == "" {
		return nil, nil
	}
	c := &DeployConfig{
		RPCURL:     rpc,
		AdminDelay: big.NewInt(86400),
		OutFile:    env["DEPLOY_OUT_FILE"],
		TxTimeout:  180 * time.Second,
	}

	id, ok := new(big.Int).SetString(env["DEPLOY_CHAIN_ID"], 10)
	if !ok {
		return nil, fmt.Errorf("DEPLOY_CHAIN_ID is required and must be a decimal integer")
	}
	c.ChainID = id

	if env["DEPLOY_PRIVATE_KEY"] == "" {
		return nil, fmt.Errorf("DEPLOY_PRIVATE_KEY is required")
	}
	k, err := parseKey(env["DEPLOY_PRIVATE_KEY"])
	if err != nil {
		return nil, fmt.Errorf("DEPLOY_PRIVATE_KEY: %w", err)
	}
	c.Key = k

	deployer := addrOf(k)
	c.Admin, c.Minter = deployer, deployer

	for envKey, dst := range map[string]*common.Address{
		"DEPLOY_ADMIN":  &c.Admin,
		"DEPLOY_MINTER": &c.Minter,
	} {
		if s := env[envKey]; s != "" {
			if !common.IsHexAddress(s) {
				return nil, fmt.Errorf("%s %q is not a valid address", envKey, s)
			}
			*dst = common.HexToAddress(s)
		}
	}

	if c.Admin == (common.Address{}) {
		return nil, fmt.Errorf("DEPLOY_ADMIN must not be the zero address")
	}

	if s := env["DEPLOY_ADMIN_DELAY"]; s != "" {
		d, ok := new(big.Int).SetString(s, 10)
		if !ok || d.Sign() < 0 || d.BitLen() > 48 {
			return nil, fmt.Errorf("DEPLOY_ADMIN_DELAY %q must be an integer in [0, 2^48)", s)
		}
		c.AdminDelay = d
	}

	if s := env["DEPLOY_TX_TIMEOUT"]; s != "" {
		n, err := strconv.Atoi(s)
		if err != nil || n <= 0 {
			return nil, fmt.Errorf("DEPLOY_TX_TIMEOUT %q must be a positive number of seconds", s)
		}
		c.TxTimeout = time.Duration(n) * time.Second
	}

	return c, nil
}

func TestDeploy(t *testing.T) {
	dc, err := loadDeployConfig()
	if err != nil {
		t.Fatalf("deploy configuration error: %v", err)
	}
	if dc == nil {
		t.Skip("DEPLOY_RPC_URL not set; see tests/.env.example")
	}
	ctx, cancel := context.WithTimeout(context.Background(), dc.TxTimeout)
	defer cancel()
	cl, err := ethclient.DialContext(ctx, dc.RPCURL)
	if err != nil {
		t.Fatalf("dial %s: %v", dc.RPCURL, err)
	}
	defer cl.Close()
	remote, err := cl.ChainID(ctx)
	if err != nil {
		t.Fatalf("fetch chain id: %v", err)
	}
	if remote.Cmp(dc.ChainID) != 0 {
		t.Fatalf("DEPLOY_CHAIN_ID=%s but the endpoint reports chain id %s", dc.ChainID, remote)
	}
	deployer := addrOf(dc.Key)
	bal, err := cl.BalanceAt(ctx, deployer, nil)
	if err != nil {
		t.Fatalf("balance: %v", err)
	}
	if bal.Sign() == 0 {
		t.Fatalf("deployer %s has no funds on chain %s", deployer, dc.ChainID)
	}
	t.Logf("deploying from %s on chain %s (admin %s, minter %s, delay %ss)",
		deployer, dc.ChainID, dc.Admin, dc.Minter, dc.AdminDelay)
	auth, err := bind.NewKeyedTransactorWithChainID(dc.Key, dc.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	auth.Context = ctx
	address, tx, instance, err := cluecode.DeployClueCode(auth, cl, dc.Admin, dc.AdminDelay, dc.Minter)
	if err != nil {
		t.Fatalf("DeployClueCode: %v", err)
	}
	receipt, err := bind.WaitMined(ctx, cl, tx)
	if err != nil {
		t.Fatalf("wait for deployment %s: %v", tx.Hash(), err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		t.Fatalf("deployment %s reverted", tx.Hash())
	}
	if receipt.ContractAddress != address {
		t.Fatalf("receipt address %s != returned address %s", receipt.ContractAddress, address)
	}
	t.Logf("ClueCode deployed at %s (tx %s, block %s, gas %d)",
		address, tx.Hash(), receipt.BlockNumber, receipt.GasUsed)
	// --- verify the constructor arguments took effect ---
	if name, err := instance.Name(nil); err != nil || name != "Clue Provider Code" {
		t.Errorf("name = %q, err %v", name, err)
	}
	if sym, err := instance.Symbol(nil); err != nil || sym != "CPC" {
		t.Errorf("symbol = %q, err %v", sym, err)
	}
	if supply, err := instance.MAXSUPPLY(nil); err != nil || supply.Cmp(big.NewInt(totalSupply)) != 0 {
		t.Errorf("MAX_SUPPLY = %v, err %v", supply, err)
	}
	if n, err := instance.TotalMinted(nil); err != nil || n.Sign() != 0 {
		t.Errorf("totalMinted = %v, err %v; want 0", n, err)
	}
	if a, err := instance.DefaultAdmin(nil); err != nil || a != dc.Admin {
		t.Errorf("defaultAdmin = %s, err %v; want %s", a, err, dc.Admin)
	}
	if d, err := instance.DefaultAdminDelay(nil); err != nil || d.Cmp(dc.AdminDelay) != 0 {
		t.Errorf("defaultAdminDelay = %v, err %v; want %s", d, err, dc.AdminDelay)
	}
	minterRole, err := instance.MINTERROLE(nil)
	if err != nil {
		t.Fatal(err)
	}
	if ok, err := instance.HasRole(nil, minterRole, dc.Minter); err != nil || !ok {
		t.Errorf("hasRole(MINTER_ROLE, %s) = %v, err %v; want true", dc.Minter, ok, err)
	}
	if ok, err := instance.HasRole(nil, [32]byte{}, dc.Admin); err != nil || !ok {
		t.Errorf("hasRole(DEFAULT_ADMIN_ROLE, %s) = %v, err %v; want true", dc.Admin, ok, err)
	}
	if dc.Minter != dc.Admin {
		if ok, err := instance.HasRole(nil, minterRole, dc.Admin); err != nil || ok {
			t.Errorf("admin unexpectedly holds MINTER_ROLE (%v, err %v)", ok, err)
		}
	}
	// Print the address in a form that can be pasted into the shell for the rest of the suite.
	fmt.Printf("\nexport CLUE_RPC_URL=%s\nexport CLUE_CHAIN_ID=%s\nexport CLUE_CONTRACT_ADDRESS=%s\n\n",
		dc.RPCURL, dc.ChainID, address)
	if dc.OutFile != "" {
		out, _ := json.MarshalIndent(map[string]interface{}{
			"chainId":    dc.ChainID.String(),
			"rpc":        dc.RPCURL,
			"address":    address.Hex(),
			"txHash":     tx.Hash().Hex(),
			"block":      receipt.BlockNumber.String(),
			"admin":      dc.Admin.Hex(),
			"minter":     dc.Minter.Hex(),
			"adminDelay": dc.AdminDelay.String(),
		}, "", "  ")
		if err := os.WriteFile(dc.OutFile, out, 0o644); err != nil {
			t.Fatalf("write %s: %v", dc.OutFile, err)
		}
		t.Logf("wrote %s", dc.OutFile)
	}
}

// ======================================================================
// Reference implementation and Section 8 test vectors
// ======================================================================
// Reference implementation of the Provider Code rules, used to cross-check the on-chain
// functions. Keep identical to the vectors in the ClueCode implementation note (Section 8).
const (
	gradeNone   uint8 = 0
	gradeS      uint8 = 1
	gradeA      uint8 = 2
	gradeB      uint8 = 3
	totalSupply       = 90_999
)

func refGrade(c uint64) uint8 {
	switch {
	case c >= 1 && c <= 99:
		return gradeS
	case c >= 100 && c <= 999:
		return gradeA
	case c >= 10000 && c <= 99999:
		return gradeB
	}
	return gradeNone
}

func refWidth(g uint8) int {
	switch g {
	case gradeS:
		return 2
	case gradeA:
		return 3
	}
	return 5
}

func refCanonical(c uint64) string {
	return fmt.Sprintf("%0*d", refWidth(refGrade(c)), c)
}

type validVector struct {
	code      uint64
	grade     uint8
	canonical string
}

var validVectors = []validVector{
	{1, gradeS, "01"},
	{2, gradeS, "02"},
	{27, gradeS, "27"},
	{99, gradeS, "99"},
	{100, gradeA, "100"},
	{999, gradeA, "999"},
	{10000, gradeB, "10000"},
	{54321, gradeB, "54321"},
	{99999, gradeB, "99999"},
}

// 0, the four-digit gap and anything above 99999 are never valid.
var invalidVectors = []uint64{0, 1000, 1001, 5000, 9998, 9999, 100000, 100001, 1 << 40}

// ======================================================================
// Read-only: metadata, grade, canonical form, parity
// ======================================================================
func TestTokenMetadata(t *testing.T) {
	requireConfig(t)
	name, err := cc.Name(nil)
	if err != nil || name != "Clue Provider Code" {
		t.Fatalf("name = %q, err %v", name, err)
	}
	symbol, err := cc.Symbol(nil)
	if err != nil || symbol != "CPC" {
		t.Fatalf("symbol = %q, err %v", symbol, err)
	}
	maxSupply, err := cc.MAXSUPPLY(nil)
	if err != nil || maxSupply.Cmp(big.NewInt(totalSupply)) != 0 {
		t.Fatalf("MAX_SUPPLY = %v, err %v", maxSupply, err)
	}
	minted, err := cc.TotalMinted(nil)
	if err != nil {
		t.Fatal(err)
	}
	if minted.Sign() < 0 || minted.Cmp(maxSupply) > 0 {
		t.Fatalf("totalMinted %v outside [0, %v]", minted, maxSupply)
	}
}

func TestIsValidCode(t *testing.T) {
	requireConfig(t)
	for _, v := range validVectors {
		ok, err := cc.IsValidCode(nil, bigInt(v.code))
		if err != nil || !ok {
			t.Errorf("isValidCode(%d) = %v, err %v; want true", v.code, ok, err)
		}
	}
	for _, c := range invalidVectors {
		ok, err := cc.IsValidCode(nil, bigInt(c))
		if err != nil || ok {
			t.Errorf("isValidCode(%d) = %v, err %v; want false", c, ok, err)
		}
	}
	huge := new(big.Int).Lsh(big.NewInt(1), 255)
	if ok, err := cc.IsValidCode(nil, huge); err != nil || ok {
		t.Errorf("isValidCode(2^255) = %v, err %v; want false", ok, err)
	}
}

func TestGradeOf(t *testing.T) {
	requireConfig(t)
	for _, v := range validVectors {
		g, err := cc.GradeOf(nil, bigInt(v.code))
		if err != nil {
			t.Errorf("gradeOf(%d): %v", v.code, err)
			continue
		}
		if g != v.grade {
			t.Errorf("gradeOf(%d) = %d, want %d", v.code, g, v.grade)
		}
	}
}

func TestGradeOfInvalidReverts(t *testing.T) {
	requireConfig(t)
	for _, c := range invalidVectors {
		_, err := cc.GradeOf(nil, bigInt(c))
		if err == nil {
			t.Errorf("gradeOf(%d) succeeded, want InvalidProviderCode", c)
			continue
		}
		if !isRevert(err, "InvalidProviderCode") {
			t.Errorf("gradeOf(%d): wrong error: %v", c, err)
		}
	}
}

func TestCanonical(t *testing.T) {
	requireConfig(t)
	for _, v := range validVectors {
		s, err := cc.Canonical(nil, bigInt(v.code))
		if err != nil {
			t.Errorf("canonical(%d): %v", v.code, err)
			continue
		}
		if s != v.canonical {
			t.Errorf("canonical(%d) = %q, want %q", v.code, s, v.canonical)
		}
	}
}

func TestCanonicalInvalidReverts(t *testing.T) {
	requireConfig(t)
	for _, c := range invalidVectors {
		_, err := cc.Canonical(nil, bigInt(c))
		if err == nil {
			t.Errorf("canonical(%d) succeeded, want InvalidProviderCode", c)
			continue
		}
		if !isRevert(err, "InvalidProviderCode") {
			t.Errorf("canonical(%d): wrong error: %v", c, err)
		}
	}
}

// TestParitySweep compares the contract with the Go reference over the full Grade S and A
// ranges, the boundaries of the invalid gap, and a stride through Grade B.
// Skipped with -short because it makes one RPC call per code.
func TestParitySweep(t *testing.T) {
	requireConfig(t)
	if testing.Short() {
		t.Skip("skipped in -short mode")
	}
	check := func(c uint64) {
		t.Helper()
		ok, err := cc.IsValidCode(nil, bigInt(c))
		if err != nil {
			t.Fatalf("isValidCode(%d): %v", c, err)
		}
		g := refGrade(c)
		if ok != (g != gradeNone) {
			t.Fatalf("isValidCode(%d) = %v, reference grade %d", c, ok, g)
		}
		if g == gradeNone {
			return
		}
		onchainGrade, err := cc.GradeOf(nil, bigInt(c))
		if err != nil || onchainGrade != g {
			t.Fatalf("gradeOf(%d) = %d, err %v; want %d", c, onchainGrade, err, g)
		}
		s, err := cc.Canonical(nil, bigInt(c))
		if err != nil || s != refCanonical(c) {
			t.Fatalf("canonical(%d) = %q, err %v; want %q", c, s, err, refCanonical(c))
		}
	}
	for c := uint64(1); c <= 999; c++ {
		check(c)
	}
	for _, c := range []uint64{1000, 1001, 9998, 9999, 10000, 10001, 99998, 99999, 100000, 100001} {
		check(c)
	}
	for c := uint64(10000); c <= 99999; c += 997 {
		check(c)
	}
}

// ======================================================================
// Roles, access control and ERC-165
// ======================================================================
func TestRoleConstants(t *testing.T) {
	requireConfig(t)
	minter, err := cc.MINTERROLE(nil)
	if err != nil {
		t.Fatal(err)
	}
	if want := crypto.Keccak256Hash([]byte("MINTER_ROLE")); common.Hash(minter) != want {
		t.Fatalf("MINTER_ROLE = %x, want keccak256(\"MINTER_ROLE\") = %s", minter, want)
	}
	admin, err := cc.DEFAULTADMINROLE(nil)
	if err != nil {
		t.Fatal(err)
	}
	if admin != [32]byte{} {
		t.Fatalf("DEFAULT_ADMIN_ROLE = %x, want zero", admin)
	}
}

func TestDefaultAdmin(t *testing.T) {
	requireConfig(t)
	got, err := cc.DefaultAdmin(nil)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Admin != nil && got != *cfg.Admin {
		t.Fatalf("defaultAdmin = %s, want %s", got, *cfg.Admin)
	}
	// After the admin renounces (genesis lock-down) defaultAdmin is the zero address.
	if got == (common.Address{}) {
		t.Log("default admin renounced: role management is permanently frozen")
		return
	}
	delay, err := cc.DefaultAdminDelay(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("defaultAdmin %s, delay %s seconds", got, delay)
	hasRole, err := cc.HasRole(nil, [32]byte{}, got)
	if err != nil || !hasRole {
		t.Fatalf("hasRole(DEFAULT_ADMIN_ROLE, %s) = %v, err %v", got, hasRole, err)
	}
}

func TestSupportsInterface(t *testing.T) {
	requireConfig(t)
	cases := []struct {
		name string
		id   [4]byte
		want bool
	}{
		{"ERC-165", [4]byte{0x01, 0xff, 0xc9, 0xa7}, true},
		{"ERC-721", [4]byte{0x80, 0xac, 0x58, 0xcd}, true},
		{"ERC-721 Metadata", [4]byte{0x5b, 0x5e, 0x13, 0x9f}, true},
		{"AccessControl", [4]byte{0x79, 0x65, 0xdb, 0x0b}, true},
		{"ERC-721 Enumerable (deliberately absent)", [4]byte{0x78, 0x0e, 0x9d, 0x63}, false},
		{"invalid 0xffffffff", [4]byte{0xff, 0xff, 0xff, 0xff}, false},
	}
	for _, c := range cases {
		got, err := cc.SupportsInterface(nil, c.id)
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
			continue
		}
		if got != c.want {
			t.Errorf("supportsInterface(%s) = %v, want %v", c.name, got, c.want)
		}
	}
}

// callAs performs an eth_call of a ClueCode method from an arbitrary (unfunded) sender.
func callAs(t *testing.T, from common.Address, method string, args ...interface{}) error {
	t.Helper()
	parsed, err := cluecode.ClueCodeMetaData.GetAbi()
	if err != nil {
		t.Fatal(err)
	}
	data, err := parsed.Pack(method, args...)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err = client.CallContract(ctx, ethereum.CallMsg{From: from, To: &cfg.Contract, Data: data}, nil)
	return err
}

// Anyone without MINTER_ROLE must be rejected. Uses eth_call from a random address,
// so it needs no funds and changes no state.
func TestMintRequiresMinterRole(t *testing.T) {
	requireConfig(t)
	stranger := randomAddress(t)
	requireRevert(t, callAs(t, stranger, "mint", stranger, bigInt(27)), "AccessControlUnauthorizedAccount")
	requireRevert(t, callAs(t, stranger, "mintBatch", stranger, []*big.Int{bigInt(27)}), "AccessControlUnauthorizedAccount")
}

func TestGrantRoleRequiresAdmin(t *testing.T) {
	requireConfig(t)
	stranger := randomAddress(t)
	minter, err := cc.MINTERROLE(nil)
	if err != nil {
		t.Fatal(err)
	}
	requireRevert(t, callAs(t, stranger, "grantRole", minter, stranger), "AccessControlUnauthorizedAccount")
}

func TestDefaultAdminRulesEnforced(t *testing.T) {
	requireConfig(t)
	stranger := randomAddress(t)
	// With AccessControlDefaultAdminRules the admin role can never be granted or revoked directly.
	requireRevert(t, callAs(t, stranger, "grantRole", [32]byte{}, stranger), "AccessControlEnforcedDefaultAdminRules")
	requireRevert(t, callAs(t, stranger, "revokeRole", [32]byte{}, stranger), "AccessControlEnforcedDefaultAdminRules")
	requireRevert(t, callAs(t, stranger, "beginDefaultAdminTransfer", stranger), "AccessControlUnauthorizedAccount")
}

func TestTokenURIUnmintedReverts(t *testing.T) {
	requireConfig(t)
	// Valid code, but pick one very unlikely to be minted on a fresh genesis-in-progress chain.
	for _, c := range []uint64{99999, 99998, 99997, 99996, 99995} {
		if _, err := cc.OwnerOf(nil, bigInt(c)); err == nil {
			continue // minted: try the next one
		}
		_, err := cc.TokenURI(nil, bigInt(c))
		requireRevert(t, err, "ERC721NonexistentToken")
		return
	}
	t.Skip("codes 99995-99999 are all minted; nothing unminted to test")
}

// ======================================================================
// State-changing: mint, batch, transfer
// ======================================================================
// State-changing tests. They need CLUE_PRIVATE_KEY (and, if different, CLUE_MINTER_PRIVATE_KEY)
// to hold MINTER_ROLE and enough gas funds, and skip themselves otherwise. They mint a few
// codes from the top of the Grade B range (99999 downwards), skipping any already minted.
// freeCodes returns n valid, currently unminted codes, searching downward from 99999.
func freeCodes(t *testing.T, n int) []uint64 {
	t.Helper()
	var out []uint64
	for c := uint64(99999); c >= 10000 && len(out) < n; c-- {
		_, err := cc.OwnerOf(nil, bigInt(c))
		if err == nil {
			continue // already minted
		}
		if !isRevert(err, "ERC721NonexistentToken") {
			t.Fatalf("ownerOf(%d): unexpected error: %v", c, err)
		}
		out = append(out, c)
	}
	if len(out) < n {
		t.Skipf("fewer than %d unminted Grade B codes available", n)
	}
	return out
}

type tokenMeta struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Attributes  []struct {
		TraitType   string      `json:"trait_type"`
		DisplayType string      `json:"display_type"`
		Value       interface{} `json:"value"`
	} `json:"attributes"`
}

func decodeDataURI(t *testing.T, uri, prefix string) []byte {
	t.Helper()
	if !strings.HasPrefix(uri, prefix) {
		t.Fatalf("URI does not start with %q: %.60s", prefix, uri)
	}
	b, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(uri, prefix))
	if err != nil {
		t.Fatalf("base64 decode: %v", err)
	}
	return b
}

func TestMintTransferAndMetadata(t *testing.T) {
	requireMinter(t)
	owner := addrOf(cfg.Key)
	minterAuth := txOpts(t, cfg.MinterKey)
	code := freeCodes(t, 1)[0]
	id := bigInt(code)
	before, err := cc.TotalMinted(nil)
	if err != nil {
		t.Fatal(err)
	}
	// --- mint ---
	tx, err := cc.Mint(minterAuth, owner, id)
	mine(t, tx, err)
	got, err := cc.OwnerOf(nil, id)
	if err != nil || got != owner {
		t.Fatalf("ownerOf(%d) = %s, err %v; want %s", code, got, err, owner)
	}
	after, err := cc.TotalMinted(nil)
	if err != nil || new(big.Int).Sub(after, before).Int64() != 1 {
		t.Fatalf("totalMinted %v -> %v, want +1 (err %v)", before, after, err)
	}
	// --- tokenURI: on-chain JSON + SVG ---
	uri, err := cc.TokenURI(nil, id)
	if err != nil {
		t.Fatal(err)
	}
	var meta tokenMeta
	if err := json.Unmarshal(decodeDataURI(t, uri, "data:application/json;base64,"), &meta); err != nil {
		t.Fatalf("tokenURI JSON: %v", err)
	}
	canon := refCanonical(code)
	if meta.Name != "Clue Provider Code "+canon {
		t.Errorf("name = %q", meta.Name)
	}
	if len(meta.Attributes) != 2 || meta.Attributes[0].Value != "B" || meta.Attributes[1].Value != float64(5) {
		t.Errorf("attributes = %+v, want Grade B / 5 digits", meta.Attributes)
	}
	svg := string(decodeDataURI(t, meta.Image, "data:image/svg+xml;base64,"))
	if !strings.Contains(svg, ">"+canon+"</text>") || !strings.Contains(svg, "Grade B</text>") {
		t.Errorf("svg does not show code %s / Grade B: %s", canon, svg)
	}
	// --- rejected mints ---
	_, err = cc.Mint(minterAuth, owner, id)
	requireRevert(t, err, "ERC721InvalidSender") // duplicate
	_, err = cc.Mint(minterAuth, owner, big.NewInt(1000))
	requireRevert(t, err, "InvalidProviderCode") // four-digit gap
	_, err = cc.Mint(minterAuth, owner, big.NewInt(0))
	requireRevert(t, err, "InvalidProviderCode")
	_, err = cc.Mint(minterAuth, common.Address{}, bigInt(code-1))
	requireRevert(t, err, "ERC721InvalidReceiver")
	// --- a role holder cannot move someone else's token ---
	if addrOf(cfg.MinterKey) != owner {
		_, err = cc.TransferFrom(minterAuth, owner, addrOf(cfg.MinterKey), id)
		requireRevert(t, err, "ERC721InsufficientApproval")
	}
	// --- approve, then transfer (approval is cleared) ---
	ownerAuth := txOpts(t, cfg.Key)
	operator := randomAddress(t)
	tx, err = cc.Approve(ownerAuth, operator, id)
	mine(t, tx, err)
	if a, err := cc.GetApproved(nil, id); err != nil || a != operator {
		t.Fatalf("getApproved = %s, err %v; want %s", a, err, operator)
	}
	recipient := randomAddress(t)
	tx, err = cc.TransferFrom(ownerAuth, owner, recipient, id)
	mine(t, tx, err)
	if o, err := cc.OwnerOf(nil, id); err != nil || o != recipient {
		t.Fatalf("after transfer ownerOf = %s, err %v; want %s", o, err, recipient)
	}
	if a, err := cc.GetApproved(nil, id); err != nil || a != (common.Address{}) {
		t.Fatalf("approval not cleared after transfer: %s, err %v", a, err)
	}
	if b, err := cc.BalanceOf(nil, recipient); err != nil || b.Int64() != 1 {
		t.Fatalf("balanceOf(recipient) = %v, err %v", b, err)
	}
}

func TestMintBatch(t *testing.T) {
	requireMinter(t)
	owner := addrOf(cfg.Key)
	auth := txOpts(t, cfg.MinterKey)
	codes := freeCodes(t, 3)
	// A batch containing an invalid code is all-or-nothing: nothing is minted.
	bad := []*big.Int{bigInt(codes[0]), bigInt(codes[1]), big.NewInt(1000)}
	_, err := cc.MintBatch(auth, owner, bad)
	requireRevert(t, err, "InvalidProviderCode")
	if _, err := cc.OwnerOf(nil, bigInt(codes[0])); !isRevert(err, "ERC721NonexistentToken") {
		t.Fatalf("code %d was minted by a failed batch (err %v)", codes[0], err)
	}
	// A valid batch mints every code and advances totalMinted by the batch size.
	before, err := cc.TotalMinted(nil)
	if err != nil {
		t.Fatal(err)
	}
	good := []*big.Int{bigInt(codes[0]), bigInt(codes[1]), bigInt(codes[2])}
	tx, err := cc.MintBatch(auth, owner, good)
	mine(t, tx, err)
	for _, c := range codes {
		if o, err := cc.OwnerOf(nil, bigInt(c)); err != nil || o != owner {
			t.Errorf("ownerOf(%d) = %s, err %v; want %s", c, o, err, owner)
		}
	}
	after, err := cc.TotalMinted(nil)
	if err != nil || new(big.Int).Sub(after, before).Int64() != 3 {
		t.Fatalf("totalMinted %v -> %v, want +3 (err %v)", before, after, err)
	}
}
