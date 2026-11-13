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

// deployFile is the JSON written by TestDeploy (DEPLOY_OUT_FILE). It is the
// only source of the chain / contract values; nothing falls back to .env.
const deployFile = "deploy/clue_code.json"

// deployRecord mirrors the JSON written by TestDeploy.
type deployRecord struct {
	ChainID    string `json:"chainId"`
	RPC        string `json:"rpc"`
	Address    string `json:"address"`
	TxHash     string `json:"txHash"`
	Block      string `json:"block"`
	Admin      string `json:"admin"`
	Minter     string `json:"minter"`
	AdminDelay string `json:"adminDelay"`
}

// loadConfig returns (nil, nil) when the deploy JSON does not exist, so the
// suite skips. RPC URL, chain id, contract address and admin come from the
// JSON only. Private keys are secrets and are never written to the JSON, so
// they (and the tx timeout) are still read from tests/.env.
func loadConfig() (*Config, error) {
	raw, err := os.ReadFile(deployFile)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", deployFile, err)
	}
	var d deployRecord
	if err := json.Unmarshal(raw, &d); err != nil {
		return nil, fmt.Errorf("parse %s: %w", deployFile, err)
	}
	if d.RPC == "" {
		return nil, fmt.Errorf("%s: \"rpc\" is missing", deployFile)
	}
	c := &Config{RPCURL: d.RPC, TxTimeout: 120 * time.Second}
	id, ok := new(big.Int).SetString(d.ChainID, 10)
	if !ok {
		return nil, fmt.Errorf("%s: chainId %q is not a decimal integer", deployFile, d.ChainID)
	}
	c.ChainID = id
	if !common.IsHexAddress(d.Address) {
		return nil, fmt.Errorf("%s: address %q is not a valid address", deployFile, d.Address)
	}
	c.Contract = common.HexToAddress(d.Address)
	if d.Admin != "" {
		if !common.IsHexAddress(d.Admin) {
			return nil, fmt.Errorf("%s: admin %q is not a valid address", deployFile, d.Admin)
		}
		a := common.HexToAddress(d.Admin)
		c.Admin = &a
	}

	// Secrets and tuning only: tests/.env
	env, err := envt.Map(&envt.Config{
		EnvPath: "tests/.env",
	})
	if err != nil {
		return nil, err
	}
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
		return fmt.Errorf("chainId=%s in the deploy JSON but the endpoint reports chain id %s", cfg.ChainID, remote)
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

func requireConfig(t *testing.T) {
	t.Helper()
	if cfgErr != nil {
		t.Fatalf("test configuration error: %v", cfgErr)
	}
	if cfg == nil {
		t.Skip("deploy JSON " + deployFile + " not found; run TestDeploy with DEPLOY_OUT_FILE set")
	}
}

func requireSigner(t *testing.T) {
	t.Helper()
	requireConfig(t)
	if cfg.Key == nil {
		t.Skip("CLUE_PRIVATE_KEY not set")
	}
}

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
// Reference implementation and Section 8 test vectors
// ======================================================================
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
	requireRevert(t, callAs(t, stranger, "grantRole", [32]byte{}, stranger), "AccessControlEnforcedDefaultAdminRules")
	requireRevert(t, callAs(t, stranger, "revokeRole", [32]byte{}, stranger), "AccessControlEnforcedDefaultAdminRules")
	requireRevert(t, callAs(t, stranger, "beginDefaultAdminTransfer", stranger), "AccessControlUnauthorizedAccount")
}

func TestTokenURIUnmintedReverts(t *testing.T) {
	requireConfig(t)
	for _, c := range []uint64{99999, 99998, 99997, 99996, 99995} {
		if _, err := cc.OwnerOf(nil, bigInt(c)); err == nil {
			continue
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
func freeCodes(t *testing.T, n int) []uint64 {
	t.Helper()
	var out []uint64
	for c := uint64(99999); c >= 10000 && len(out) < n; c-- {
		_, err := cc.OwnerOf(nil, bigInt(c))
		if err == nil {
			continue
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
	_, err = cc.Mint(minterAuth, owner, id)
	requireRevert(t, err, "ERC721InvalidSender")
	_, err = cc.Mint(minterAuth, owner, big.NewInt(1000))
	requireRevert(t, err, "InvalidProviderCode")
	_, err = cc.Mint(minterAuth, owner, big.NewInt(0))
	requireRevert(t, err, "InvalidProviderCode")
	_, err = cc.Mint(minterAuth, common.Address{}, bigInt(code-1))
	requireRevert(t, err, "ERC721InvalidReceiver")
	if addrOf(cfg.MinterKey) != owner {
		_, err = cc.TransferFrom(minterAuth, owner, addrOf(cfg.MinterKey), id)
		requireRevert(t, err, "ERC721InsufficientApproval")
	}
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
	bad := []*big.Int{bigInt(codes[0]), bigInt(codes[1]), big.NewInt(1000)}
	_, err := cc.MintBatch(auth, owner, bad)
	requireRevert(t, err, "InvalidProviderCode")
	if _, err := cc.OwnerOf(nil, bigInt(codes[0])); !isRevert(err, "ERC721NonexistentToken") {
		t.Fatalf("code %d was minted by a failed batch (err %v)", codes[0], err)
	}
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
