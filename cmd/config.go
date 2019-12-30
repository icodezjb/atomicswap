package cmd

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

//in uints
const gasLimit = uint64(3000000)

func init() {
	log.SetFlags(0)
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
		log.Fatalf("generate secret: %v", err)
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

func (c *Config) ParseConfig(cfgPath string) error {
	configFile, err := os.Open(cfgPath)
	defer configFile.Close() //nolint:staticcheck

	if err != nil {
		return errors.Wrap(err, "open config file")
	}

	configStr, err := ioutil.ReadAll(configFile)
	if err != nil {
		return errors.Wrap(err, "read config file")
	}

	if err := json.Unmarshal(configStr, c); err != nil {
		return errors.Wrapf(err, "parse config file (%s)", cfgPath)
	}

	return nil
}

func (c *Config) Connect(otherContract string) error {
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
		return errors.Wrapf(err, "connect to %v", c.Chain.URL)
	}

	c.client = client

	return nil
}

func (c *Config) Unlock(privateKey string) error {
	switch {
	case privateKey != "":
		key, err := crypto.HexToECDSA(privateKey)
		if err != nil {
			return errors.Wrapf(err, "parse private key (%v)", privateKey)
		}

		account := crypto.PubkeyToAddress(key.PublicKey).String()

		if strings.ToLower(c.Account) != strings.ToLower(account) {
			return errors.Errorf("mismatch private key (%s) and account (%s)", privateKey, c.Account)
		}
		c.key = key
	default:
		c.ks = keystore.NewKeyStore(c.KeyStore, keystore.StandardScryptN, keystore.StandardScryptP)
		fromAccount := accounts.Account{Address: common.HexToAddress(c.Account)}

		if c.ks.HasAddress(fromAccount.Address) {
			err := c.ks.Unlock(fromAccount, c.Password)
			if err != nil {
				return errors.Wrapf(err, "unlock %v keystore", c.Account)
			}
		} else {
			return errors.Errorf("not found %v in %v keystore (%v)", c.Account, c.KeyStore, c.ks.Accounts())
		}
	}

	return nil
}

func (c *Config) ValidateAddress(address string) error {
	if valid := regexp.MustCompile("^0x[0-9a-fA-F]{40}$").MatchString(address); !valid {
		return errors.Errorf("invalid address: %v", address)
	}
	return nil
}

func (c *Config) rotate(cfgPath string) error {
	replacement, err := os.OpenFile(cfgPath+".new", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	defer replacement.Close() //nolint:staticcheck

	if err != nil {
		return errors.Wrap(err, "open config file")
	}

	enc := json.NewEncoder(replacement)
	enc.SetIndent("", "    ")

	if err = enc.Encode(c); err != nil {
		return errors.Wrap(err, "encode config")
	}

	// Replace the live config with the newly generated one
	return os.Rename(cfgPath+".new", cfgPath)
}

func (c *Config) makeAuth(ctx context.Context, value int64) (*bind.TransactOpts, error) {
	fromAccount := accounts.Account{Address: common.HexToAddress(c.Account)}
	nonce, err := c.client.PendingNonceAt(ctx, fromAccount.Address)
	if err != nil {
		return nil, errors.Wrapf(err, "account=%v get nonce", c.Account)
	}

	gasPrice, err := c.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "account=%v get gasPrice", c.Account)
	}

	auth, err := bind.NewKeyStoreTransactor(c.ks, fromAccount)
	if err != nil {
		return nil, errors.Wrapf(err, "account=%v get keystore transactor", c.Account)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(value) //in wei
	auth.GasLimit = gasLimit

	gasPriceInt, _ := big.NewInt(0).SetString(gasPrice.String(), 10)
	auth.GasPrice = gasPriceInt

	return auth, nil
}

func (c *Config) promptConfirm(prefix string) {
	log.Printf("? Confirm to %v the contract on %v(chainID = %v)? [y/N]", prefix, c.Chain.Name, c.Chain.ID)

	if c.test {
		log.Println("Test chose: y")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	data, _, _ := reader.ReadLine()

	input := string(data)
	if len(input) > 0 && strings.ToLower(input[:1]) == "y" {
		log.Println("Your chose: y")
	} else {
		log.Println("Your chose: N")
		os.Exit(0)
	}
}

func Must(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
