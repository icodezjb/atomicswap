package cmd

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

//in uints
const gasLimit = uint64(3000000)

type HashTimeLockContract interface {
	DeployContract(ctx context.Context)
	StatContract(ctx context.Context)
	NewContract(ctx context.Context, participant common.Address, amount int64, hashLock [32]byte, timeLock *big.Int) *types.Transaction
	GetContractId(ctx context.Context, txID common.Hash) HtlcLogHTLCNew
	AuditContract(ctx context.Context, result interface{}, method string, contractId common.Hash)
	Redeem(ctx context.Context, contractId common.Hash, secret common.Hash) *types.Transaction
	Refund(ctx context.Context, contractId common.Hash) *types.Transaction
}

type chain struct {
	ID       *big.Int
	Name     string
	URL      string
	Contract string
}

type Config struct {
	ChainID        *big.Int `json:"chainID"`
	ChainName      string   `json:"chainName"`
	URL            string   `json:"url"`
	OtherChainID   *big.Int `json:"otherChainID"`
	OtherChainName string   `json:"otherChainName"`
	OtherURL       string   `json:"otherURL"`
	Account        string   `json:"account"`
	Contract       string   `json:"contract"`
	KeyStore       string   `json:"keystoreDir"`
	Password       string   `json:"password"`
	Chain          *chain   `json:"-"`
	client         *ethclient.Client
	ks             *keystore.KeyStore
	key            *ecdsa.PrivateKey

	//only for test
	test bool
}

type SecretHashPair struct {
	Secret [32]byte
	Hash   [32]byte
}

func NewSecretHashPair() *SecretHashPair {
	s := new(SecretHashPair)
	_, err := rand.Read(s.Secret[:])
	if err != nil {
		logger.FatalError("Fatal to generate secret: %v", err)
	}
	s.Hash = sha256.Sum256(s.Secret[:])

	return s
}

type HtlcLogHTLCNew struct {
	ContractId [32]byte
	Sender     common.Address
	Receiver   common.Address
	Amount     *big.Int
	Hashlock   [32]byte
	Timelock   *big.Int
}

type ContractDetails struct {
	Sender    common.Address
	Receiver  common.Address
	Amount    *big.Int
	Hashlock  [32]byte
	Timelock  *big.Int
	Withdrawn bool
	Refunded  bool
	Preimage  [32]byte
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

func (c *Config) Connect(otherContract string) {
	c.Chain = &chain{
		ID:       c.ChainID,
		Name:     c.ChainName,
		URL:      c.URL,
		Contract: c.Contract,
	}

	if otherContract != "" {
		c.Chain = &chain{
			ID:       c.OtherChainID,
			Name:     c.OtherChainName,
			URL:      c.OtherURL,
			Contract: otherContract,
		}
	}

	client, err := ethclient.Dial(c.Chain.URL)
	if err != nil {
		logger.FatalError("Fatal to connect server: %v", err)
	}

	c.client = client
}

func (c *Config) Unlock(privateKey string) {
	switch {
	case privateKey != "":
		key, err := crypto.HexToECDSA(privateKey)
		if err != nil {
			logger.FatalError("Fatal to parse private key (%v): %v", privateKey, err)
		}

		account := crypto.PubkeyToAddress(key.PublicKey).String()

		if strings.ToLower(c.Account) != strings.ToLower(account) {
			logger.FatalError("Fatal to match the private key (%v) and account (%v)", privateKey, c.Account)
		}
		c.key = key
	default:
		c.ks = keystore.NewKeyStore(c.KeyStore, keystore.StandardScryptN, keystore.StandardScryptP)
		fromAccount := accounts.Account{Address: common.HexToAddress(c.Account)}

		if c.ks.HasAddress(fromAccount.Address) {
			err := c.ks.Unlock(fromAccount, c.Password)
			if err != nil {
				logger.FatalError("Fatal to unlock %v", c.Account)
			}
		} else {
			logger.FatalError("Fatal to find %v in %v keystore (%v)", c.Account, c.KeyStore, c.ks.Accounts())
		}
	}
}

func (c *Config) ValidateAddress(address string) {
	if valid := regexp.MustCompile("^0x[0-9a-fA-F]{40}$").MatchString(address); !valid {
		logger.FatalError("Fatal to validate address: %v", address)
	}
}

func (c *Config) rotate(cfgPath string) {
	replacement, err := os.OpenFile(cfgPath+".new", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	defer replacement.Close() //nolint:staticcheck

	if err != nil {
		logger.FatalError("Fatal to open file: %v", err)
	}

	enc := json.NewEncoder(replacement)
	enc.SetIndent("", "    ")

	if err = enc.Encode(c); err != nil {
		logger.FatalError("Fatal to encode config: %v", err)
	}

	// Replace the live config with the newly generated one
	if err = os.Rename(cfgPath+".new", cfgPath); err != nil {
		logger.FatalError("Fatal to replace config file: %v", err)
	}
}

func (c *Config) makeAuth(ctx context.Context, value int64) *bind.TransactOpts {
	fromAccount := accounts.Account{Address: common.HexToAddress(c.Account)}
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
	auth.Value = big.NewInt(value) //in wei
	auth.GasLimit = gasLimit

	gasPriceInt, _ := big.NewInt(0).SetString(gasPrice.String(), 10)
	auth.GasPrice = gasPriceInt

	return auth
}

func (c *Config) promptConfirm(prefix string) {
	logger.Info("? Confirm to %v the contract on %v(chainID = %v)? [y/N]", prefix, c.Chain.Name, c.Chain.ID)

	if c.test {
		logger.Info("Test chose: y")
		return
	}

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
