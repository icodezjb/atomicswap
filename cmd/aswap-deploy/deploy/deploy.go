package deploy

import (
	"bufio"
	"context"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"os"
	"strings"

	htlc "github.com/icodezjb/atomicswap/contract"
	"github.com/icodezjb/atomicswap/logger"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Config struct {
	ChainId   uint   `json:"chainId"`
	ChainName string `json:"chainName"`
	Url       string `json:"url"`
	From      string `json:"from"`
	KeyStore  string `json:"keystoreDir"`
	Password  string `json:"password"`
}

type DeployHandle struct {
	ConfigPath string
	Config     Config
	Client     *ethclient.Client
	Ks         *keystore.KeyStore
	Account    accounts.Account
}

func (d *DeployHandle) ParseConfig() {
	configFile, err := os.Open(d.ConfigPath)
	if err != nil {
		logger.FatalError("Fatal to open config file (%s): %v", d.ConfigPath, err)
	}

	configStr, err := ioutil.ReadAll(configFile)
	if err != nil {
		logger.FatalError("Fatal to read config file (%s): %v", d.ConfigPath, err)
	}

	if err := json.Unmarshal(configStr, &d.Config); err != nil {
		logger.FatalError("Fatal to parse config file (%s): %v", d.ConfigPath, err)
	}
}

func (d *DeployHandle) Connect() {
	client, err := ethclient.Dial(d.Config.Url)
	if err != nil {
		logger.FatalError("Fatal to connect server: %s", err)
	}
	d.Client = client
}

func (d *DeployHandle) Unlock() {
	d.Ks = keystore.NewKeyStore(d.Config.KeyStore, keystore.StandardScryptN, keystore.StandardScryptP)
	d.Account = accounts.Account{Address: common.HexToAddress(d.Config.From)}

	if d.Ks.HasAddress(d.Account.Address) {
		err := d.Ks.Unlock(d.Account, d.Config.Password)
		if err != nil {
			logger.FatalError("Fatal to unlock %v", d.Config.From)
		}
	} else {
		logger.FatalError("Fatal to find %v in %v keystore (%v)", d.Config.From, d.Config.KeyStore, d.Ks.Accounts())
	}
}

func (d *DeployHandle) MakeAuth(ctx context.Context, value int64) *bind.TransactOpts {
	nonce, err := d.Client.PendingNonceAt(ctx, d.Account.Address)
	if err != nil {
		logger.FatalError("Fatal to get nonce: %s", err)
	}

	gasPrice, err := d.Client.SuggestGasPrice(ctx)
	if err != nil {
		logger.FatalError("Fatal to get gasPrice: %s", err)
	}

	auth, err := bind.NewKeyStoreTransactor(d.Ks, d.Account)
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

func (d *DeployHandle) PromptConfirm() {
	logger.Info("? Confirm to deploy the contract ? [y/N]")

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

func (d *DeployHandle) EstimateGas(ctx context.Context, auth *bind.TransactOpts) {
	estimateGas, err := d.Client.EstimateGas(ctx, ethereum.CallMsg{
		From:     auth.From,
		To:       nil,
		Gas:      0,
		GasPrice: auth.GasPrice,
		Value:    auth.Value,
		Data:     common.FromHex(htlc.HtlcBin),
	})
	if err != nil {
		logger.FatalError("Fatal to estimate gas: %v", err)
	}

	feeByWei := new(big.Int).Mul(new(big.Int).SetUint64(estimateGas), auth.GasPrice).String()

	balance, err := d.Client.BalanceAt(ctx, auth.From, nil)
	if err != nil {
		logger.FatalError("Fatal to get %v balance: %s", auth.From.String(), err)
	}

	logger.Info("from = %v, balance = %v", auth.From.String(), balance)
	logger.Info("contract fee = gas(%v) * gasPrice(%v) = %v", estimateGas, auth.GasPrice.String(), feeByWei)

}

func (d *DeployHandle) DeployContract() {
	ctx := context.Background()

	//unlock account
	d.Unlock()

	auth := d.MakeAuth(ctx, 0)

	//Deploy contract
	logger.Info("Deploy contract...")

	d.EstimateGas(ctx, auth)

	d.PromptConfirm()

	address, tx, _, err := htlc.DeployHtlc(auth, d.Client)
	if err != nil {
		logger.FatalError("Fatal to deploy contract: %s", err)
	}

	logger.Info("contract address = %s", address.Hex())
	logger.Info("transaction hash = %s", tx.Hash().Hex())
}
