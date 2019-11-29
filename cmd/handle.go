package cmd

import (
	"bufio"
	"context"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"os"
	"regexp"
	"strings"

	htlc "github.com/icodezjb/atomicswap/contract"
	util "github.com/icodezjb/atomicswap/contract/helper"
	"github.com/icodezjb/atomicswap/logger"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Config struct {
	ChainId   *big.Int `json:"chainId"`
	ChainName string   `json:"chainName"`
	Url       string   `json:"url"`
	From      string   `json:"from"`
	KeyStore  string   `json:"keystoreDir"`
	Password  string   `json:"password"`
	Contract  string   `json:"contract"`
}

type Handle struct {
	ConfigPath string
	Config     Config
	client     *ethclient.Client
	ks         *keystore.KeyStore
	account    accounts.Account
	htlc       *htlc.Htlc
}

func (h *Handle) ParseConfig() {
	configFile, err := os.Open(h.ConfigPath)
	defer configFile.Close() //nolint:staticcheck

	if err != nil {
		logger.FatalError("Fatal to open config file (%s): %v", h.ConfigPath, err)
	}

	configStr, err := ioutil.ReadAll(configFile)
	if err != nil {
		logger.FatalError("Fatal to read config file (%s): %v", h.ConfigPath, err)
	}

	if err := json.Unmarshal(configStr, &h.Config); err != nil {
		logger.FatalError("Fatal to parse config file (%s): %v", h.ConfigPath, err)
	}
}

func (h *Handle) Connect() {
	client, err := ethclient.Dial(h.Config.Url)
	if err != nil {
		logger.FatalError("Fatal to connect server: %v", err)
	}
	h.client = client
}

func (h *Handle) unlock() {
	h.ks = keystore.NewKeyStore(h.Config.KeyStore, keystore.StandardScryptN, keystore.StandardScryptP)
	h.account = accounts.Account{Address: common.HexToAddress(h.Config.From)}

	if h.ks.HasAddress(h.account.Address) {
		err := h.ks.Unlock(h.account, h.Config.Password)
		if err != nil {
			logger.FatalError("Fatal to unlock %v", h.Config.From)
		}
	} else {
		logger.FatalError("Fatal to find %v in %v keystore (%v)", h.Config.From, h.Config.KeyStore, h.ks.Accounts())
	}
}

func (h *Handle) makeAuth(ctx context.Context, value int64) *bind.TransactOpts {
	nonce, err := h.client.PendingNonceAt(ctx, h.account.Address)
	if err != nil {
		logger.FatalError("Fatal to get nonce: %v", err)
	}

	gasPrice, err := h.client.SuggestGasPrice(ctx)
	if err != nil {
		logger.FatalError("Fatal to get gasPrice: %v", err)
	}

	auth, err := bind.NewKeyStoreTransactor(h.ks, h.account)
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

func (h *Handle) promptConfirm(prefix string) {
	logger.Info("? Confirm to %v the contract on %v(chainID = %v)? [y/N]", prefix, h.Config.ChainName, h.Config.ChainId)

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

func (h *Handle) estimateGas(ctx context.Context, auth *bind.TransactOpts, txType string) {
	estimateGas, err := h.client.EstimateGas(ctx, ethereum.CallMsg{
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

	balance, err := h.client.BalanceAt(ctx, auth.From, nil)
	if err != nil {
		logger.FatalError("Fatal to get %v balance: %v", auth.From.String(), err)
	}

	logger.Info("from = %v, balance = %v", auth.From.String(), balance)
	logger.Info("%v fee = gas(%v) * gasPrice(%v) = %v", txType, estimateGas, auth.GasPrice.String(), feeByWei)

}

func (h *Handle) generate() {
	replacement, err := os.OpenFile("config-after-deployed.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	defer replacement.Close() //nolint:staticcheck

	if err != nil {
		logger.FatalError("Fatal to open file: %v", err)
	}

	enc := json.NewEncoder(replacement)
	enc.SetIndent("", "    ")

	if err = enc.Encode(h.Config); err != nil {
		logger.FatalError("Fatal to encode config: %v", err)
	}
}

func (h *Handle) DeployContract() {
	ctx := context.Background()

	//unlock account
	h.unlock()

	auth := h.makeAuth(ctx, 0)

	logger.Info("Deploy contract...")

	//estimate deploy contract fee
	h.estimateGas(ctx, auth, "contract")

	//deploy-contract prompt
	h.promptConfirm("deploy")

	address, tx, _, err := htlc.DeployHtlc(auth, h.client)
	if err != nil {
		logger.FatalError("Fatal to deploy contract: %v", err)
	}

	//update contract address
	h.Config.Contract = address.String()

	logger.Info("contract address = %v", address.String())
	logger.Info("transaction hash = %v", tx.Hash().String())

	//generate config-after-deployed.json file
	h.generate()
}

func (h *Handle) StatContract() {
	//TODO
}

func (h *Handle) ValidateAddress(address string) {
	valid := regexp.MustCompile("^0x[0-9a-fA-F]{40}$").MatchString(address)
	switch valid {
	case false:
		logger.FatalError("Fatal to validate address: %v", address)
	default:
	}
}

func (h *Handle) LoadContract(contractAddress string) {
	logger.Info("Initiate contract...")

	instance, err := htlc.NewHtlc(common.HexToAddress(contractAddress), h.client)
	if err != nil {
		logger.FatalError("Failed to load contract: %s", err)
	}

	h.htlc = instance
	logger.Info("contract at %s is loaded", contractAddress)
}

func (h *Handle) sendTx(ctx context.Context, auth *bind.TransactOpts, data []byte) common.Hash {
	from := new(accounts.Account)
	from.Address = common.HexToAddress(h.Config.From)

	tx := types.NewTransaction(auth.Nonce.Uint64(),
		common.HexToAddress(h.Config.Contract),
		auth.Value,
		auth.GasLimit,
		auth.GasPrice,
		data)

	txSigned, err := h.ks.SignTxWithPassphrase(*from, h.Config.Password, tx, h.Config.ChainId)
	if err != nil {
		logger.FatalError("Fatal to sign tx %v: %v", h.Config.From, err)
	}

	txHash := txSigned.Hash()
	err = h.client.SendTransaction(context.Background(), txSigned)
	if err != nil {
		logger.FatalError("Fatal to send tx: %v", err)
	}

	return txHash
}

func (h *Handle) NewContract(participant string, hashLock string, timeLock string) {
	ctx := context.Background()
	//unlock account
	h.unlock()

	auth := h.makeAuth(ctx, 0)

	logger.Info("Call NewContract ...")

	//estimate call contract fee
	h.estimateGas(ctx, auth, "contract")

	//call-contract prompt
	h.promptConfirm("call")

	//NewContract function signature + participant addr + hashlock + timelock
	dataStr := "335ef5bd" + util.PadTo32Bytes(participant) + util.PadTo32Bytes(hashLock) + util.PadTo32Bytes(timeLock)

	data, err := hex.DecodeString(dataStr)
	if err != nil {
		logger.FatalError("Fatal to make data call contract: %v", err)
	}

	txHash := h.sendTx(ctx, auth, data)

	logger.Info("initiate tx: %v", txHash.String())
}
