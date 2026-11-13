package tests

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	cluecode "github.com/clue-protocol/contract/out/go/code"
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
// Deploy test (separate DEPLOY_* config)
// ======================================================================

// Grade B code range used for random mint tests. Adjust to match the
// contract's valid code ranges (MAX_SUPPLY is 90_999).
const (
	minTestCode = 10000
	maxTestCode = 90998
	maxSupply   = 90_999
)

// randomCodes returns n unique random codes in [minTestCode, maxTestCode].
func randomCodes(n int) []*big.Int {
	seen := make(map[int]struct{}, n)
	codes := make([]*big.Int, 0, n)
	for len(codes) < n {
		c := gofakeit.Number(minTestCode, maxTestCode)
		if _, dup := seen[c]; dup {
			continue
		}
		seen[c] = struct{}{}
		codes = append(codes, big.NewInt(int64(c)))
	}
	return codes
}

// DeployConfig holds the parameters of TestDeploy only.
type DeployConfig struct {
	Env        map[string]string
	RPCURL     string
	ChainID    *big.Int
	Key        *ecdsa.PrivateKey
	Admin      common.Address
	Minter     common.Address
	AdminDelay *big.Int
	InitialURI string
	OutFile    string
	TxTimeout  time.Duration
	MintBatch  bool // when true, TestDeploy exercises MintBatch after deployment
}

// loadDeployConfig returns (nil, "", nil) when DEPLOY_RPC_URL is not set.
func loadDeployConfig() (*DeployConfig, string, error) {
	env, err := envt.Map(&envt.Config{
		EnvPath:  "tests/.env",
		PanicErr: false,
	})
	if err != nil {
		return nil, "", err
	}
	c := &DeployConfig{Env: env}

	rpc := c.Env["DEPLOY_RPC_URL"]
	if rpc == "" {
		return nil, "", nil
	}
	c.RPCURL = rpc
	c.AdminDelay = big.NewInt(86400)

	gofakeit.Seed(0) // 0 = random seed, so each run produces different values
	uuid := gofakeit.UUID()

	if s := c.Env["DEPLOY_INITIAL_URI"]; s != "" {
		c.InitialURI = s
	} else {
		c.InitialURI = fmt.Sprintf("%s/%s/", strings.TrimRight(gofakeit.URL(), "/"), uuid)
	}

	if s := c.Env["DEPLOY_OUT_FILE"]; s != "" {
		c.OutFile = s
	} else {
		c.OutFile = fmt.Sprintf("out_%s.json", strings.Split(uuid, "-")[0])
	}
	c.TxTimeout = 180 * time.Second

	id, ok := new(big.Int).SetString(c.Env["DEPLOY_CHAIN_ID"], 10)
	if !ok {
		return nil, "", fmt.Errorf("DEPLOY_CHAIN_ID is required and must be a decimal integer")
	}
	c.ChainID = id

	if c.Env["DEPLOY_PRIVATE_KEY"] == "" {
		return nil, "", fmt.Errorf("DEPLOY_PRIVATE_KEY is required")
	}
	k, err := parseKey(c.Env["DEPLOY_PRIVATE_KEY"])
	if err != nil {
		return nil, "", fmt.Errorf("DEPLOY_PRIVATE_KEY: %w", err)
	}
	c.Key = k

	deployer := addrOf(k)
	c.Admin, c.Minter = deployer, deployer
	for envKey, dst := range map[string]*common.Address{
		"DEPLOY_ADMIN":  &c.Admin,
		"DEPLOY_MINTER": &c.Minter,
	} {
		if s := c.Env[envKey]; s != "" {
			if !common.IsHexAddress(s) {
				return nil, "", fmt.Errorf("%s %q is not a valid address", envKey, s)
			}
			*dst = common.HexToAddress(s)
		}
	}
	if c.Admin == (common.Address{}) {
		return nil, "", fmt.Errorf("DEPLOY_ADMIN must not be the zero address")
	}
	if s := c.Env["DEPLOY_ADMIN_DELAY"]; s != "" {
		d, ok := new(big.Int).SetString(s, 10)
		if !ok || d.Sign() < 0 || d.BitLen() > 48 {
			return nil, "", fmt.Errorf("DEPLOY_ADMIN_DELAY %q must be an integer in [0, 2^48)", s)
		}
		c.AdminDelay = d
	}
	if s := c.Env["DEPLOY_TX_TIMEOUT"]; s != "" {
		n, err := strconv.Atoi(s)
		if err != nil || n <= 0 {
			return nil, "", fmt.Errorf("DEPLOY_TX_TIMEOUT %q must be a positive number of seconds", s)
		}
		c.TxTimeout = time.Duration(n) * time.Second
	}
	if s := c.Env["DEPLOY_MINT_BATCH"]; s != "" {
		b, err := strconv.ParseBool(s)
		if err != nil {
			return nil, "", fmt.Errorf("DEPLOY_MINT_BATCH %q must be a boolean (true/false)", s)
		}
		c.MintBatch = b
	}
	return c, uuid, nil
}

func TestDeploy(t *testing.T) {
	dc, uuid, err := loadDeployConfig()
	require.NoError(t, err, "deploy configuration error")
	if dc == nil {
		t.Skip("DEPLOY_RPC_URL not set; see tests/.env.example")
	}

	ctx, cancel := context.WithTimeout(context.Background(), dc.TxTimeout)
	defer cancel()

	cl, err := ethclient.DialContext(ctx, dc.RPCURL)
	require.NoError(t, err, "dial %s", dc.RPCURL)
	defer cl.Close()

	remote, err := cl.ChainID(ctx)
	require.NoError(t, err, "fetch chain id")
	require.Zero(t, remote.Cmp(dc.ChainID),
		"DEPLOY_CHAIN_ID=%s but the endpoint reports chain id %s", dc.ChainID, remote)

	deployer := addrOf(dc.Key)
	bal, err := cl.BalanceAt(ctx, deployer, nil)
	require.NoError(t, err, "balance")
	require.NotZero(t, bal.Sign(), "deployer %s has no funds on chain %s", deployer, dc.ChainID)

	t.Logf("deploying from %s on chain %s (admin %s, minter %s, delay %ss, uri %s)",
		deployer, dc.ChainID, dc.Admin, dc.Minter, dc.AdminDelay, dc.InitialURI)

	auth, err := bind.NewKeyedTransactorWithChainID(dc.Key, dc.ChainID)
	require.NoError(t, err)
	auth.Context = ctx

	address, tx, instance, err := cluecode.DeployClueCode(auth, cl, dc.Admin, dc.AdminDelay, dc.Minter, dc.InitialURI)
	require.NoError(t, err, "DeployClueCode")

	receipt, err := bind.WaitMined(ctx, cl, tx)
	require.NoError(t, err, "wait for deployment %s", tx.Hash())
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status, "deployment %s reverted", tx.Hash())
	require.Equal(t, address, receipt.ContractAddress, "receipt address != returned address")
	t.Logf("ClueCode deployed at %s (tx %s, block %s, gas %d)",
		address, tx.Hash(), receipt.BlockNumber, receipt.GasUsed)

	verifyDeployment(t, dc, instance)

	// --- Optionally exercise mintBatch immediately after deployment ---
	var codesToMint []*big.Int
	if dc.MintBatch {
		codesToMint = testMintBatch(t, ctx, cl, dc, instance, deployer)
	} else {
		t.Log("DEPLOY_MINT_BATCH not enabled; skipping mintBatch")
	}

	// Print the address in a form that can be pasted into the shell for the rest of the suite.
	fmt.Printf("\nexport CLUE_RPC_URL=%s\nexport CLUE_CHAIN_ID=%s\nexport CLUE_CONTRACT_ADDRESS=%s\n",
		dc.RPCURL, dc.ChainID, address)

	if dc.OutFile != "" {
		result := map[string]interface{}{
			"chainId":          dc.ChainID.String(),
			"rpc":              dc.RPCURL,
			"contract address": address.Hex(),
			"txHash":           tx.Hash().Hex(),
			"block":            receipt.BlockNumber.String(),
			"admin":            dc.Admin.Hex(),
			"minter":           dc.Minter.Hex(),
			"adminDelay":       dc.AdminDelay.String(),
			"initialURI":       dc.InitialURI,
			"mintBatch":        dc.MintBatch,
			"uuid":             uuid,
			// NOTE: this writes a private key to disk. Keep tests/deploy/ out of git.
			"admin key": hexutil.Encode(crypto.FromECDSA(dc.Key)),
		}
		if dc.MintBatch {
			result["mintedCodes"] = codesToMint
		}
		out, err := json.MarshalIndent(result, "", "  ")
		require.NoError(t, err, "marshal deploy output")
		require.NoError(t, envt.WriteFile(filepath.Join("tests", "deploy", dc.OutFile), out, 0o644),
			"write %s", dc.OutFile)
		t.Logf("wrote %s", dc.OutFile)
	}
}

// verifyDeployment checks that the constructor arguments took effect.
func verifyDeployment(t *testing.T, dc *DeployConfig, instance *cluecode.ClueCode) {
	t.Helper()

	name, err := instance.Name(nil)
	require.NoError(t, err)
	require.Equal(t, "Clue Provider Code", name)

	symbol, err := instance.Symbol(nil)
	require.NoError(t, err)
	require.Equal(t, "CPC", symbol)

	supply, err := instance.MAXSUPPLY(nil)
	require.NoError(t, err)
	require.Zero(t, supply.Cmp(big.NewInt(maxSupply)), "MAX_SUPPLY = %v", supply)

	minted, err := instance.TotalMinted(nil)
	require.NoError(t, err)
	require.Zero(t, minted.Sign(), "totalMinted = %v; want 0", minted)

	admin, err := instance.DefaultAdmin(nil)
	require.NoError(t, err)
	require.Equal(t, dc.Admin, admin)

	delay, err := instance.DefaultAdminDelay(nil)
	require.NoError(t, err)
	require.Zero(t, delay.Cmp(dc.AdminDelay), "defaultAdminDelay = %v; want %s", delay, dc.AdminDelay)

	minterRole, err := instance.MINTERROLE(nil)
	require.NoError(t, err)

	ok, err := instance.HasRole(nil, minterRole, dc.Minter)
	require.NoError(t, err)
	require.True(t, ok, "minter %s should hold MINTER_ROLE", dc.Minter)

	ok, err = instance.HasRole(nil, [32]byte{}, dc.Admin)
	require.NoError(t, err)
	require.True(t, ok, "admin %s should hold DEFAULT_ADMIN_ROLE", dc.Admin)

	if dc.Minter != dc.Admin {
		ok, err = instance.HasRole(nil, minterRole, dc.Admin)
		require.NoError(t, err)
		require.False(t, ok, "admin unexpectedly holds MINTER_ROLE")
	}
}

// testMintBatch mints a random batch to the deployer and verifies the result.
// It returns the codes that were minted.
func testMintBatch(
	t *testing.T,
	ctx context.Context,
	cl *ethclient.Client,
	dc *DeployConfig,
	instance *cluecode.ClueCode,
	deployer common.Address,
) []*big.Int {
	t.Helper()
	t.Log("exercising mintBatch post-deployment")

	// The deployer is also the minter by default in this test setup.
	mintAuth, err := bind.NewKeyedTransactorWithChainID(dc.Key, dc.ChainID)
	require.NoError(t, err, "create transactor for mintBatch")
	mintAuth.Context = ctx

	batchSize := gofakeit.Number(2, 5)
	codesToMint := randomCodes(batchSize)
	require.Len(t, codesToMint, batchSize)
	t.Logf("minting %d random codes: %v", batchSize, codesToMint)

	txMint, err := instance.MintBatch(mintAuth, deployer, codesToMint)
	require.NoError(t, err, "MintBatch")

	receiptMint, err := bind.WaitMined(ctx, cl, txMint)
	require.NoError(t, err, "wait for MintBatch tx %s", txMint.Hash())
	require.Equal(t, types.ReceiptStatusSuccessful, receiptMint.Status,
		"MintBatch tx %s reverted", txMint.Hash())

	totalMinted, err := instance.TotalMinted(nil)
	require.NoError(t, err, "TotalMinted")
	require.Zero(t, totalMinted.Cmp(big.NewInt(int64(batchSize))),
		"expected totalMinted=%d, got %v", batchSize, totalMinted)

	for _, code := range codesToMint {
		owner, err := instance.OwnerOf(nil, code)
		require.NoError(t, err, "OwnerOf(%v)", code)
		require.Equal(t, deployer, owner, "owner of code %v", code)
	}

	t.Logf("Successfully minted batch of %d codes to %s", batchSize, deployer)
	return codesToMint
}
