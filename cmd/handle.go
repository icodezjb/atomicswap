package cmd

import (
	"context"
	"math/big"
	"strings"
	"time"

	htlc "github.com/icodezjb/atomicswap/contract"
	"github.com/icodezjb/atomicswap/logger"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type Handle struct {
	ConfigPath string
	Config     *Config
}

func (h *Handle) estimateGas(ctx context.Context, auth *bind.TransactOpts, txType string, input []byte) {
	var (
		contract *common.Address
	)
	if h.Config.Contract != "" {
		tmpAddress := common.HexToAddress(h.Config.Contract)
		contract = &tmpAddress
	}

	estimateGas, err := h.Config.client.EstimateGas(ctx, ethereum.CallMsg{
		From:     auth.From,
		To:       contract,
		Gas:      0,
		GasPrice: auth.GasPrice,
		Value:    auth.Value,
		Data:     input,
	})
	if err != nil {
		logger.FatalError("Fatal to estimate gas: %v", err)
	}

	feeByWei := new(big.Int).Mul(new(big.Int).SetUint64(estimateGas), auth.GasPrice).String()

	balance, err := h.Config.client.BalanceAt(ctx, auth.From, nil)
	if err != nil {
		logger.FatalError("Fatal to get %v balance: %v", auth.From.String(), err)
	}

	logger.Info("from = %v, balance = %v", auth.From.String(), balance)
	logger.Info("%v Contract fee = gas(%v) * gasPrice(%v) = %v", txType, estimateGas, auth.GasPrice.String(), feeByWei)

}

func (h *Handle) sendTx(ctx context.Context, auth *bind.TransactOpts, data []byte, contract *common.Address) *types.Transaction {
	var rawTx *types.Transaction

	// Create and Sign the transaction, sign it and schedule it for execution
	if contract == nil {
		rawTx = types.NewContractCreation(auth.Nonce.Uint64(), auth.Value, auth.GasLimit, auth.GasPrice, data)
	} else {
		rawTx = types.NewTransaction(auth.Nonce.Uint64(), *contract, auth.Value, auth.GasLimit, auth.GasPrice, data)
	}
	txSigned, err := h.Config.ks.SignTxWithPassphrase(
		accounts.Account{Address: common.HexToAddress(h.Config.Account)},
		h.Config.Password,
		rawTx,
		h.Config.Chain.ID)
	if err != nil {
		logger.FatalError("Fatal to sign tx %v: %v", h.Config.Account, err)
	}

	from := new(accounts.Account)
	from.Address = common.HexToAddress(h.Config.Account)

	err = h.Config.client.SendTransaction(ctx, txSigned)
	if err != nil {
		logger.FatalError("Fatal to send tx: %v", err)
	}
	return txSigned
}

func (h *Handle) DeployContract() {
	ctx := context.Background()

	auth := h.Config.makeAuth(ctx, 0)

	logger.Info("Deploy contract...")

	input := common.FromHex(htlc.HtlcBin)

	//estimate deploy contract fee
	h.estimateGas(ctx, auth, "Deploy", input)

	//deploy-contract prompt
	h.Config.promptConfirm("deploy")

	//send tx
	txSigned := h.sendTx(ctx, auth, input, nil)

	//update contract address
	h.Config.Contract = crypto.CreateAddress(auth.From, txSigned.Nonce()).String()

	logger.Info("contract address = %v", h.Config.Contract)
	logger.Info("transaction hash = %v", txSigned.Hash().String())

	//update config
	h.Config.rotate(h.ConfigPath)
}

func (h *Handle) StatContract() {
	//TODO
}

func (h *Handle) NewContract(participant common.Address, amount int64, hashLock [32]byte, timeLock *big.Int) {
	ctx := context.Background()

	auth := h.Config.makeAuth(ctx, amount)

	logger.Info("Call NewContract ...")

	parsedABI, err := abi.JSON(strings.NewReader(htlc.HtlcABI))
	if err != nil {
		logger.FatalError("Fatal to parse HtlcABI: %v", err)
	}

	input, err := parsedABI.Pack("newContract", participant, hashLock, timeLock)
	if err != nil {
		logger.FatalError("Fatal to pack newContract: %v", err)
	}

	//estimate call contract fee
	h.estimateGas(ctx, auth, "Call", input)

	//call-contract prompt
	h.Config.promptConfirm("call")

	//send tx
	contract := common.HexToAddress(h.Config.Contract)
	txSigned := h.sendTx(ctx, auth, input, &contract)

	logger.Info("%v(%v) txid: %v\n", h.Config.Chain.Name, h.Config.Chain.ID, txSigned.Hash().String())
}

func (h *Handle) GetContractId(txID common.Hash) {
	logger.Info("%v(%v) txid: %v", h.Config.Chain.Name, h.Config.Chain.ID, txID.String())
	logger.Info("contract address: %v\n", h.Config.Chain.Contract)

	receipt, err := h.Config.client.TransactionReceipt(context.Background(), txID)
	if err != nil {
		logger.FatalError("Failed to get tx %v receipt: %v", txID.String(), err)
	}

	var logHTLCEvent htlc.HtlcLogHTLCNew
	parsedABI, err := abi.JSON(strings.NewReader(htlc.HtlcABI))
	if err != nil {
		logger.FatalError("Fatal to parse HtlcABI: %v", err)
	}
	if err := parsedABI.Unpack(&logHTLCEvent, "LogHTLCNew", receipt.Logs[0].Data); err != nil {
		logger.FatalError("Failed to unpack log data for LogHTLCNew: %v", err)
	}

	if receipt.Logs == nil {
		logger.FatalError("initiateTx receipt.Logs == nil, receipt.Status = %v", receipt.Status)
	}

	logHTLCEvent.ContractId = receipt.Logs[0].Topics[1]
	logHTLCEvent.Sender = common.HexToAddress(receipt.Logs[0].Topics[2].Hex())
	logHTLCEvent.Receiver = common.HexToAddress(receipt.Logs[0].Topics[3].Hex())

	logger.Info("ContractId = %s", hexutil.Encode(logHTLCEvent.ContractId[:]))
	logger.Info("Sender     = %s", logHTLCEvent.Sender.String())
	logger.Info("Receiver   = %s", logHTLCEvent.Receiver.String())
	logger.Info("Amount     = %s", logHTLCEvent.Amount)
	logger.Info("TimeLock   = %s (%s)", logHTLCEvent.Timelock, time.Unix(logHTLCEvent.Timelock.Int64(), 0).Format(time.RFC3339))
	logger.Info("SecretHash = %s", hexutil.Encode(logHTLCEvent.Hashlock[:]))
}

func (h *Handle) AuditContract(from common.Address, contractId common.Hash) {
	ctx := context.Background()

	logger.Info("Call getContract ...")
	logger.Info("contract address: %v\n", h.Config.Chain.Contract)

	contractDetails := new(struct {
		Sender    common.Address
		Receiver  common.Address
		Amount    *big.Int
		Hashlock  [32]byte
		Timelock  *big.Int
		Withdrawn bool
		Refunded  bool
		Preimage  [32]byte
	})

	h.auditContract(ctx, contractDetails, "getContract", from, contractId)

	logger.Info("Sender     = %s", contractDetails.Sender.String())
	logger.Info("Receiver   = %s", contractDetails.Receiver.String())
	logger.Info("Amount     = %s (wei)", contractDetails.Amount)
	logger.Info("TimeLock   = %s (%s)", contractDetails.Timelock, time.Unix(contractDetails.Timelock.Int64(), 0))
	logger.Info("SecretHash = %s", hexutil.Encode(contractDetails.Hashlock[:]))
	logger.Info("Withdrawn  = %v", contractDetails.Withdrawn)
	logger.Info("Refunded   = %v", contractDetails.Refunded)
	logger.Info("Secret     = %s", hexutil.Encode(contractDetails.Preimage[:]))
}

func (h *Handle) auditContract(ctx context.Context, result interface{}, method string, from common.Address, contractId common.Hash) {
	parsedABI, err := abi.JSON(strings.NewReader(htlc.HtlcABI))
	if err != nil {
		logger.FatalError("Fatal to parse HtlcABI: %v", err)
	}

	input, err := parsedABI.Pack(method, contractId)
	if err != nil {
		logger.FatalError("Fatal to pack %v: %v", method, err)
	}

	//Call
	contract := common.HexToAddress(h.Config.Chain.Contract)
	msg := ethereum.CallMsg{From: from, To: &contract, Data: input}
	opts := bind.CallOpts{From: from}
	var output []byte

	output, err = h.Config.client.CallContract(ctx, msg, opts.BlockNumber)
	if err == nil && len(output) == 0 {
		// Make sure we have a contract to operate on, and bail out otherwise.
		if code, err := h.Config.client.CodeAt(ctx, contract, opts.BlockNumber); err != nil {
			logger.FatalError("Fatal to call CodeAt: %v", err)
		} else if len(code) == 0 {
			logger.FatalError("no contract code at given address")
		}
	}

	if err != nil {
		logger.FatalError("Fatal to call CallContract: %v", err)
	}

	if err = parsedABI.Unpack(result, method, output); err != nil {
		logger.FatalError("Fatal to unpack result of contract call: %v", err)
	}
}

func (h *Handle) Redeem(contractId common.Hash, secret common.Hash) {
	ctx := context.Background()

	auth := h.Config.makeAuth(ctx, 0)

	logger.Info("Call Withdraw ...")

	parsedABI, err := abi.JSON(strings.NewReader(htlc.HtlcABI))
	if err != nil {
		logger.FatalError("Fatal to parse HtlcABI: %v", err)
	}

	input, err := parsedABI.Pack("withdraw", contractId, secret)
	if err != nil {
		logger.FatalError("Fatal to pack newContract: %v", err)
	}

	//estimate call contract fee
	h.estimateGas(ctx, auth, "Call", input)

	//call-contract prompt
	h.Config.promptConfirm("call")

	//send tx
	contract := common.HexToAddress(h.Config.Chain.Contract)
	txSigned := h.sendTx(ctx, auth, input, &contract)

	logger.Info("%v(%v) txid: %v\n", h.Config.Chain.Name, h.Config.Chain.ID, txSigned.Hash().String())
}

func (h *Handle) Refund(contractId common.Hash) {
	ctx := context.Background()

	auth := h.Config.makeAuth(ctx, 0)

	logger.Info("Call Withdraw ...")

	parsedABI, err := abi.JSON(strings.NewReader(htlc.HtlcABI))
	if err != nil {
		logger.FatalError("Fatal to parse HtlcABI: %v", err)
	}

	input, err := parsedABI.Pack("refund", contractId)
	if err != nil {
		logger.FatalError("Fatal to pack newContract: %v", err)
	}

	//send tx
	contract := common.HexToAddress(h.Config.Chain.Contract)
	txSigned := h.sendTx(ctx, auth, input, &contract)

	logger.Info("%v(%v) txid: %v\n", h.Config.Chain.Name, h.Config.Chain.ID, txSigned.Hash().String())
}
