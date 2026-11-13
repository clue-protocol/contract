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
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	cluecode "github.com/clue-protocol/contract/out/go/code"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pat3icki/pkgt/envt"
	"github.com/stretchr/testify/require"
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
	gofakeit.Seed(0) // 0 = random seed, so each run picks different codes

	cfg, cfgErr = loadConfigJson()
	if cfg != nil && cfgErr == nil {
		cfgErr = connect()
	}
	rc := m.Run()
	if client != nil {
		client.Close()
	}
	os.Exit(rc)
}

func findProjectRoot() (string, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "", errors.New("envt: cannot get caller info to determine project root")
	}

	dir := filepath.Dir(filename)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		if _, err := os.Stat(filepath.Join(dir, "go.sum")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("envt: go.mod/go.sum not found; are you inside a Go module?")
		}
		dir = parent
	}
}

// loadConfigJson builds the Config from the newest deploy output written by
// TestDeploy (tests/deploy/*.json). If no such file exists it falls back to
// the CLUE_* variables in tests/.env.
func loadConfigJson() (*Config, error) {
	root, err := findProjectRoot()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(root, "tests", "deploy")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return loadConfig()
		}
		return nil, err
	}

	// Pick the most recently modified .json file (ReadDir already returns
	// bare names, so join with dir exactly once).
	var newest string
	var newestMod time.Time
	for _, e := range entries {
		if e.IsDir() || !strings.EqualFold(filepath.Ext(e.Name()), ".json") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		if newest == "" || info.ModTime().After(newestMod) {
			newest, newestMod = filepath.Join(dir, e.Name()), info.ModTime()
		}
	}
	if newest == "" {
		return loadConfig()
	}

	data, err := os.ReadFile(newest)
	if err != nil {
		return nil, err
	}
	// Values are not all strings (mintBatch is a bool, mintedCodes an array).
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse %s: %w", newest, err)
	}
	str := func(keys ...string) string {
		for _, k := range keys {
			if v, ok := m[k].(string); ok && v != "" {
				return v
			}
		}
		return ""
	}

	chainID, ok := new(big.Int).SetString(str("chainId", "chainid"), 10)
	if !ok {
		return nil, fmt.Errorf("%s: missing or invalid chainId", newest)
	}
	addr := str("contract address", "contractAddress")
	if !common.IsHexAddress(addr) {
		return nil, fmt.Errorf("%s: invalid contract address %q", newest, addr)
	}
	key, err := parseKey(str("admin key", "adminKey"))
	if err != nil {
		return nil, fmt.Errorf("%s: admin key: %w", newest, err)
	}
	c := &Config{
		RPCURL:    str("rpc"),
		ChainID:   chainID,
		Contract:  common.HexToAddress(addr),
		Key:       key,
		MinterKey: key, // TestDeploy makes the deployer the minter by default
		TxTimeout: 120 * time.Second,
	}
	if a := str("admin"); common.IsHexAddress(a) {
		ad := common.HexToAddress(a)
		c.Admin = &ad
	}
	return c, nil
}

func loadConfig() (*Config, error) {
	env, err := envt.Map(&envt.Config{
		EnvPath: "tests/.env",
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

// ======================================================================
// Test helpers
// ======================================================================
func requireConfig(t *testing.T) {
	t.Helper()
	require.NoError(t, cfgErr, "test configuration error")
	if cfg == nil {
		t.Skip("CLUE_RPC_URL not set; see tests/.env.example")
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
	require.NoError(t, err, "MINTER_ROLE")
	has, err := cc.HasRole(nil, role, crypto.PubkeyToAddress(cfg.MinterKey.PublicKey))
	require.NoError(t, err, "hasRole")
	if !has {
		t.Skip("signer does not hold MINTER_ROLE (renounced or wrong key)")
	}
}

func addrOf(k *ecdsa.PrivateKey) common.Address { return crypto.PubkeyToAddress(k.PublicKey) }

func txOpts(t *testing.T, k *ecdsa.PrivateKey) *bind.TransactOpts {
	t.Helper()
	auth, err := bind.NewKeyedTransactorWithChainID(k, cfg.ChainID)
	require.NoError(t, err, "transactor")
	return auth
}

func mine(t *testing.T, tx *types.Transaction, err error) *types.Receipt {
	t.Helper()
	require.NoError(t, err, "send transaction")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.TxTimeout)
	defer cancel()
	r, err := bind.WaitMined(ctx, client, tx)
	require.NoError(t, err, "wait for %s", tx.Hash())
	require.Equal(t, types.ReceiptStatusSuccessful, r.Status, "transaction %s reverted", tx.Hash())
	return r
}

func randomAddress(t *testing.T) common.Address {
	t.Helper()
	k, err := crypto.GenerateKey()
	require.NoError(t, err)
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
	require.Error(t, err, "expected revert %s, call succeeded", name)
	require.True(t, isRevert(err, name), "expected revert %s, got: %v", name, err)
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

	gradeBMin = 10000
	gradeBMax = 99999
)

func refGrade(c uint64) uint8 {
	switch {
	case c >= 1 && c <= 99:
		return gradeS
	case c >= 100 && c <= 999:
		return gradeA
	case c >= gradeBMin && c <= gradeBMax:
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
	require.NoError(t, err)
	require.Equal(t, "Clue Provider Code", name)

	symbol, err := cc.Symbol(nil)
	require.NoError(t, err)
	require.Equal(t, "CPC", symbol)

	maxSupply, err := cc.MAXSUPPLY(nil)
	require.NoError(t, err)
	require.Zero(t, maxSupply.Cmp(big.NewInt(totalSupply)), "MAX_SUPPLY = %v", maxSupply)

	minted, err := cc.TotalMinted(nil)
	require.NoError(t, err)
	require.False(t, minted.Sign() < 0 || minted.Cmp(maxSupply) > 0,
		"totalMinted %v outside [0, %v]", minted, maxSupply)
}

func TestIsValidCode(t *testing.T) {
	requireConfig(t)

	for _, v := range validVectors {
		t.Run(fmt.Sprintf("valid_%d", v.code), func(t *testing.T) {
			ok, err := cc.IsValidCode(nil, bigInt(v.code))
			require.NoError(t, err)
			require.True(t, ok)
		})
	}
	for _, c := range invalidVectors {
		t.Run(fmt.Sprintf("invalid_%d", c), func(t *testing.T) {
			ok, err := cc.IsValidCode(nil, bigInt(c))
			require.NoError(t, err)
			require.False(t, ok)
		})
	}
	t.Run("invalid_2^255", func(t *testing.T) {
		huge := new(big.Int).Lsh(big.NewInt(1), 255)
		ok, err := cc.IsValidCode(nil, huge)
		require.NoError(t, err)
		require.False(t, ok)
	})
}

func TestGradeOf(t *testing.T) {
	requireConfig(t)
	for _, v := range validVectors {
		t.Run(fmt.Sprintf("code_%d", v.code), func(t *testing.T) {
			g, err := cc.GradeOf(nil, bigInt(v.code))
			require.NoError(t, err)
			require.Equal(t, v.grade, g)
		})
	}
}

func TestGradeOfInvalidReverts(t *testing.T) {
	requireConfig(t)
	for _, c := range invalidVectors {
		t.Run(fmt.Sprintf("code_%d", c), func(t *testing.T) {
			_, err := cc.GradeOf(nil, bigInt(c))
			requireRevert(t, err, "InvalidProviderCode")
		})
	}
}

func TestCanonical(t *testing.T) {
	requireConfig(t)
	for _, v := range validVectors {
		t.Run(fmt.Sprintf("code_%d", v.code), func(t *testing.T) {
			s, err := cc.Canonical(nil, bigInt(v.code))
			require.NoError(t, err)
			require.Equal(t, v.canonical, s)
		})
	}
}

func TestCanonicalInvalidReverts(t *testing.T) {
	requireConfig(t)
	for _, c := range invalidVectors {
		t.Run(fmt.Sprintf("code_%d", c), func(t *testing.T) {
			_, err := cc.Canonical(nil, bigInt(c))
			requireRevert(t, err, "InvalidProviderCode")
		})
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
		require.NoError(t, err, "isValidCode(%d)", c)
		g := refGrade(c)
		require.Equal(t, g != gradeNone, ok, "isValidCode(%d), reference grade %d", c, g)
		if g == gradeNone {
			return
		}
		onchainGrade, err := cc.GradeOf(nil, bigInt(c))
		require.NoError(t, err, "gradeOf(%d)", c)
		require.Equal(t, g, onchainGrade, "gradeOf(%d)", c)

		s, err := cc.Canonical(nil, bigInt(c))
		require.NoError(t, err, "canonical(%d)", c)
		require.Equal(t, refCanonical(c), s, "canonical(%d)", c)
	}
	for c := uint64(1); c <= 999; c++ {
		check(c)
	}
	for _, c := range []uint64{1000, 1001, 9998, 9999, 10000, 10001, 99998, 99999, 100000, 100001} {
		check(c)
	}
	for c := uint64(gradeBMin); c <= gradeBMax; c += 997 {
		check(c)
	}
}

// ======================================================================
// Roles, access control and ERC-165
// ======================================================================
func TestRoleConstants(t *testing.T) {
	requireConfig(t)

	minter, err := cc.MINTERROLE(nil)
	require.NoError(t, err)
	want := crypto.Keccak256Hash([]byte("MINTER_ROLE"))
	require.Equal(t, want, common.Hash(minter), `MINTER_ROLE != keccak256("MINTER_ROLE")`)

	admin, err := cc.DEFAULTADMINROLE(nil)
	require.NoError(t, err)
	require.Equal(t, [32]byte{}, admin, "DEFAULT_ADMIN_ROLE should be zero")
}

func TestDefaultAdmin(t *testing.T) {
	requireConfig(t)

	got, err := cc.DefaultAdmin(nil)
	require.NoError(t, err)
	if cfg.Admin != nil {
		require.Equal(t, *cfg.Admin, got, "defaultAdmin")
	}
	if got == (common.Address{}) {
		t.Log("default admin renounced: role management is permanently frozen")
		return
	}
	delay, err := cc.DefaultAdminDelay(nil)
	require.NoError(t, err)
	t.Logf("defaultAdmin %s, delay %s seconds", got, delay)

	hasRole, err := cc.HasRole(nil, [32]byte{}, got)
	require.NoError(t, err)
	require.True(t, hasRole, "hasRole(DEFAULT_ADMIN_ROLE, %s)", got)
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
		t.Run(c.name, func(t *testing.T) {
			got, err := cc.SupportsInterface(nil, c.id)
			require.NoError(t, err)
			require.Equal(t, c.want, got)
		})
	}
}

func callAs(t *testing.T, from common.Address, method string, args ...interface{}) error {
	t.Helper()
	parsed, err := cluecode.ClueCodeMetaData.GetAbi()
	require.NoError(t, err)
	data, err := parsed.Pack(method, args...)
	require.NoError(t, err)

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
	require.NoError(t, err)
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
	code := freeCodes(t, 1)[0]
	_, err := cc.TokenURI(nil, bigInt(code))
	requireRevert(t, err, "ERC721NonexistentToken")
}

// ======================================================================
// State-changing: mint, batch, transfer
// ======================================================================

// freeCodes returns n distinct, random, unminted Grade B codes. Random picks
// (rather than scanning from 99999) keep repeated or concurrent runs from
// colliding on the same codes.
func freeCodes(t *testing.T, n int) []uint64 {
	t.Helper()
	seen := make(map[uint64]struct{}, n)
	out := make([]uint64, 0, n)
	for attempts := 0; len(out) < n && attempts < 50*n+200; attempts++ {
		c := uint64(gofakeit.Number(gradeBMin, gradeBMax))
		if _, dup := seen[c]; dup {
			continue
		}
		seen[c] = struct{}{}

		_, err := cc.OwnerOf(nil, bigInt(c))
		if err == nil {
			continue // already minted
		}
		require.True(t, isRevert(err, "ERC721NonexistentToken"),
			"ownerOf(%d): unexpected error: %v", c, err)
		out = append(out, c)
	}
	if len(out) < n {
		t.Skipf("fewer than %d unminted Grade B codes found", n)
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
	require.True(t, strings.HasPrefix(uri, prefix), "URI does not start with %q: %.60s", prefix, uri)
	b, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(uri, prefix))
	require.NoError(t, err, "base64 decode")
	return b
}

func TestMintTransferAndMetadata(t *testing.T) {
	requireMinter(t)
	owner := addrOf(cfg.Key)
	minterAuth := txOpts(t, cfg.MinterKey)
	code := freeCodes(t, 1)[0]
	id := bigInt(code)

	before, err := cc.TotalMinted(nil)
	require.NoError(t, err)

	tx, err := cc.Mint(minterAuth, owner, id)
	mine(t, tx, err)

	got, err := cc.OwnerOf(nil, id)
	require.NoError(t, err, "ownerOf(%d)", code)
	require.Equal(t, owner, got, "ownerOf(%d)", code)

	after, err := cc.TotalMinted(nil)
	require.NoError(t, err)
	require.Equal(t, int64(1), new(big.Int).Sub(after, before).Int64(),
		"totalMinted %v -> %v, want +1", before, after)

	// --- tokenURI ---
	uri, err := cc.TokenURI(nil, id)
	require.NoError(t, err)
	canon := refCanonical(code)
	base, err := cc.BaseURI(nil)
	require.NoError(t, err)
	if base != "" {
		// No stored metadata + non-empty baseURI => tokenURI = baseURI + canonical code.
		require.Equal(t, base+canon, uri, "tokenURI should be baseURI+canonical")
	} else {
		// Empty baseURI => default on-chain JSON.
		var meta tokenMeta
		require.NoError(t, json.Unmarshal(decodeDataURI(t, uri, "data:application/json;base64,"), &meta),
			"tokenURI JSON")
		require.Equal(t, "Clue Provider Code "+canon, meta.Name)
		require.Len(t, meta.Attributes, 2, "attributes = %+v, want Grade B / 5 digits", meta.Attributes)
		require.Equal(t, "B", meta.Attributes[0].Value)
		require.Equal(t, float64(5), meta.Attributes[1].Value)
		require.Empty(t, meta.Image, "no stored metadata set")
	}

	// --- invalid mints ---
	_, err = cc.Mint(minterAuth, owner, id)
	requireRevert(t, err, "ERC721InvalidSender")
	_, err = cc.Mint(minterAuth, owner, big.NewInt(1000))
	requireRevert(t, err, "InvalidProviderCode")
	_, err = cc.Mint(minterAuth, owner, big.NewInt(0))
	requireRevert(t, err, "InvalidProviderCode")
	_, err = cc.Mint(minterAuth, common.Address{}, bigInt(code-1))
	requireRevert(t, err, "ERC721InvalidReceiver")

	// A minter that is not the owner has no transfer rights.
	if addrOf(cfg.MinterKey) != owner {
		_, err = cc.TransferFrom(minterAuth, owner, addrOf(cfg.MinterKey), id)
		requireRevert(t, err, "ERC721InsufficientApproval")
	}

	// --- approve and transfer ---
	ownerAuth := txOpts(t, cfg.Key)
	operator := randomAddress(t)
	tx, err = cc.Approve(ownerAuth, operator, id)
	mine(t, tx, err)
	approved, err := cc.GetApproved(nil, id)
	require.NoError(t, err)
	require.Equal(t, operator, approved, "getApproved")

	recipient := randomAddress(t)
	tx, err = cc.TransferFrom(ownerAuth, owner, recipient, id)
	mine(t, tx, err)

	newOwner, err := cc.OwnerOf(nil, id)
	require.NoError(t, err)
	require.Equal(t, recipient, newOwner, "ownerOf after transfer")

	approved, err = cc.GetApproved(nil, id)
	require.NoError(t, err)
	require.Equal(t, common.Address{}, approved, "approval not cleared after transfer")

	bal, err := cc.BalanceOf(nil, recipient)
	require.NoError(t, err)
	require.Equal(t, int64(1), bal.Int64(), "balanceOf(recipient)")
}

func TestMintBatch(t *testing.T) {
	requireMinter(t)
	owner := addrOf(cfg.Key)
	auth := txOpts(t, cfg.MinterKey)
	codes := freeCodes(t, 3)

	// A batch containing one invalid code must revert atomically.
	bad := []*big.Int{bigInt(codes[0]), bigInt(codes[1]), big.NewInt(1000)}
	_, err := cc.MintBatch(auth, owner, bad)
	requireRevert(t, err, "InvalidProviderCode")

	_, err = cc.OwnerOf(nil, bigInt(codes[0]))
	require.True(t, isRevert(err, "ERC721NonexistentToken"),
		"code %d was minted by a failed batch (err %v)", codes[0], err)

	before, err := cc.TotalMinted(nil)
	require.NoError(t, err)

	good := []*big.Int{bigInt(codes[0]), bigInt(codes[1]), bigInt(codes[2])}
	tx, err := cc.MintBatch(auth, owner, good)
	mine(t, tx, err)

	for _, c := range codes {
		o, err := cc.OwnerOf(nil, bigInt(c))
		require.NoError(t, err, "ownerOf(%d)", c)
		require.Equal(t, owner, o, "ownerOf(%d)", c)
	}

	after, err := cc.TotalMinted(nil)
	require.NoError(t, err)
	require.Equal(t, int64(3), new(big.Int).Sub(after, before).Int64(),
		"totalMinted %v -> %v, want +3", before, after)
}
