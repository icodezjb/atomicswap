package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"os"
	"regexp"
	"strings"

	"github.com/icodezjb/atomicswap/logger"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type HashTimeLockContract interface {
	DeployContract()
	StatContract()
	NewContract(participant common.Address, amount int64, hashLock [32]byte, timeLock *big.Int)
	GetContractId(initiateTx common.Hash)
	AuditContract(from common.Address, contractId common.Hash)
	Redeem(contractId common.Hash, secret common.Hash)
	Refund(contractId common.Hash)
}

type Config struct {
	ChainId   *big.Int `json:"chainId"`
	ChainName string   `json:"chainName"`
	Url       string   `json:"url"`
	From      string   `json:"from"`
	KeyStore  string   `json:"keystoreDir"`
	Password  string   `json:"password"`
	Contract  string   `json:"contract"`
	client    *ethclient.Client
	ks        *keystore.KeyStore
}

func (c *Config) ParseConfig(cfgPath string) {
	configFile, err := os.Open(cfgPath)
	defer configFile.Close() //nolint:staticcheck

	if err != nil {
		logger.FatalError("Fatal to open config file (%s): %v", cfgPath, err)
	}

	configStr, err := ioutil.ReadAll(configFile)
	if err != nil {
		logger.FatalError("Fatal to read config file (%s): %v", cfgPath, err)
	}

	if err := json.Unmarshal(configStr, c); err != nil {
		logger.FatalError("Fatal to parse config file (%s): %v", cfgPath, err)
	}
}

func (c *Config) Connect() {
	client, err := ethclient.Dial(c.Url)
	if err != nil {
		logger.FatalError("Fatal to connect server: %v", err)
	}
	c.client = client
}

func (c *Config) Unlock() {
	c.ks = keystore.NewKeyStore(c.KeyStore, keystore.StandardScryptN, keystore.StandardScryptP)
	fromAccount := accounts.Account{Address: common.HexToAddress(c.From)}

	if c.ks.HasAddress(fromAccount.Address) {
		err := c.ks.Unlock(fromAccount, c.Password)
		if err != nil {
			logger.FatalError("Fatal to unlock %v", c.From)
		}
	} else {
		logger.FatalError("Fatal to find %v in %v keystore (%v)", c.From, c.KeyStore, c.ks.Accounts())
	}
}

func (c *Config) ValidateAddress(address string) {
	valid := regexp.MustCompile("^0x[0-9a-fA-F]{40}$").MatchString(address)
	switch valid {
	case false:
		logger.FatalError("Fatal to validate address: %v", address)
	default:
	}
}

func (c *Config) update() {
	replacement, err := os.OpenFile("config-after-deployed.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	defer replacement.Close() //nolint:staticcheck

	if err != nil {
		logger.FatalError("Fatal to open file: %v", err)
	}

	enc := json.NewEncoder(replacement)
	enc.SetIndent("", "    ")

	if err = enc.Encode(c); err != nil {
		logger.FatalError("Fatal to encode config: %v", err)
	}
}

func (c *Config) makeAuth(ctx context.Context, value int64) *bind.TransactOpts {
	fromAccount := accounts.Account{Address: common.HexToAddress(c.From)}
	nonce, err := c.client.PendingNonceAt(ctx, fromAccount.Address)
	if err != nil {
		logger.FatalError("Fatal to get nonce: %v", err)
	}

	gasPrice, err := c.client.SuggestGasPrice(ctx)
	if err != nil {
		logger.FatalError("Fatal to get gasPrice: %v", err)
	}

	auth, err := bind.NewKeyStoreTransactor(c.ks, fromAccount)
	if err != nil {
		logger.FatalError("Fatal to make new keystore transactor: %v", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(value)  //in wei
	auth.GasLimit = uint64(3000000) //in uints

	gasPriceInt, _ := big.NewInt(0).SetString(gasPrice.String(), 10)
	auth.GasPrice = gasPriceInt

	return auth
}

func (c *Config) promptConfirm(prefix string) {
	logger.Info("? Confirm to %v the contract on %v(chainID = %v)? [y/N]", prefix, c.ChainName, c.ChainId)

	reader := bufio.NewReader(os.Stdin)
	data, _, _ := reader.ReadLine()

	input := string(data)
	if len(input) > 0 && strings.ToLower(input[:1]) == "y" {
		logger.Info("Your chose: y")
	} else {
		logger.Info("Your chose: N")
		os.Exit(0)
	}
}
