package tests

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"testing"
	"time"

	cluecode "github.com/clue-protocol/contract/out/go/code"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pat3icki/pkgt/envt"
)

// ======================================================================
// Deploy test (separate DEPLOY_* config)
// ======================================================================
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
		PanicErr: false,
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
	if supply, err := instance.MAXSUPPLY(nil); err != nil || supply.Cmp(big.NewInt(90_999)) != 0 {
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
