package main

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"runtime"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	htlc "github.com/icodezjb/atomicswap/contract"
	"github.com/icodezjb/atomicswap/logger"

	"github.com/spf13/cobra"
)

type Config struct {
	ChainId   uint   `json:"chainId"`
	ChainName string `json:"chainName"`
	Url       string `json:"url"`
	From      string `json:"from"`
	KeyStore  string `json:"keystore"`
	Password  string `json:"password"`
}

var (
	Commit    = "unknown-commit"
	BuildTime = "unknown-buildtime"

	version = "0.1.0"

	configPath string
	config     Config

	rootCmd = &cobra.Command{
		Use:   "aswap-deploy",
		Short: "deploy atomicswap smart contract",
	}
)

// VersionFunc holds the textual version string.
func VersionFunc() string {
	return fmt.Sprintf(": %s\ncommit: %s\nbuild time: %s\ngolang version: %s\n",
		version, Commit, BuildTime, runtime.Version()+" "+runtime.GOOS+"/"+runtime.GOARCH)
}

func init() {
	rootCmd.Flags().StringVar(
		&configPath,
		"config",
		"./config.json",
		"config file path",
	)
}

func main() {
	rootCmd.Version = VersionFunc()

	rootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		parseConfig(configPath)
	}

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		client, err := ethclient.Dial(config.Url)
		if err != nil {
			logger.FatalError("Fatal to connect server: %s", err)
		}

		deployContract(client)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func makeAuth(private string, client *ethclient.Client, value int64) *bind.TransactOpts {
	privateKey, err := crypto.HexToECDSA(private)
	if err != nil {
		logger.FatalError("Fatal to parse private key: %s", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		logger.FatalError("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	ctx := context.Background()
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		logger.FatalError("Fatal to get nonce: %s", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		logger.FatalError("Fatal to get gasPrice: %s", err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(value)  //in wei
	auth.GasLimit = uint64(3000000) //in uints

	gasPriceInt, _ := big.NewInt(0).SetString(gasPrice.String(), 10)
	auth.GasPrice = gasPriceInt

	return auth
}

func deployContract(client *ethclient.Client) {
	auth := makeAuth("b80dbf638b9128e54f3222d2b6d3213d45521d49bb6317abdf34b219a55204b7", client, 0)
	ctx := context.Background()

	//Deploy contract
	logger.Info("Deploy contract...")

	gas, err := client.EstimateGas(ctx, ethereum.CallMsg{
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

	feeByWei := new(big.Int).Mul(new(big.Int).SetUint64(gas), auth.GasPrice).String()

	balance, err := client.BalanceAt(ctx, auth.From, nil)
	if err != nil {
		logger.FatalError("Fatal to get %v balance: %s", auth.From.String(), err)
	}

	logger.Info("from = %v, balance = %v", auth.From.String(), balance)
	logger.Info("contract fee = gas(%v) * gasPrice(%v) = %v", gas, auth.GasPrice.String(), feeByWei)

	promptConfirm()

	address, tx, _, err := htlc.DeployHtlc(auth, client)
	if err != nil {
		logger.FatalError("Fatal to deploy contract: %s", err)
	}

	logger.Info("contract address = %s", address.Hex())
	logger.Info("transaction hash = %s", tx.Hash().Hex())
}

func promptConfirm() {
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

func parseConfig(path string) {
	configFile, err := os.Open(path)
	if err != nil {
		logger.FatalError("Fatal to open config file (%s): %v", path, err)
	}

	configStr, err := ioutil.ReadAll(configFile)
	if err != nil {
		logger.FatalError("Fatal to read config file (%s): %v", path, err)
	}

	if err := json.Unmarshal(configStr, &config); err != nil {
		logger.FatalError("Fatal to parse config file (%s): %v", path, err)
	}
}
