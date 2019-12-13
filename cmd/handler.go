package cmd

import (
	"context"
	"math/big"
	"strings"

	htlc "github.com/icodezjb/atomicswap/contract"
	"github.com/icodezjb/atomicswap/logger"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type Handler struct {
	ConfigPath string
	Config     *Config
}

func (h *Handler) estimateGas(ctx context.Context, auth *bind.TransactOpts, txType string, input []byte) {
	contract := func() *common.Address {
		switch {
		case txType == "Deploy":
			return nil
		case h.Config.Contract != "":
			tmpAddress := common.HexToAddress(h.Config.Chain.Contract)
			return &tmpAddress
		default:
			return nil
		}
	}()

	estimateGas, err := h.Config.client.EstimateGas(ctx, ethereum.CallMsg{
		From:     auth.From,
		To:       contract,
		Gas:      0,
		GasPrice: auth.GasPrice,
		Value:    auth.Value,
		Data:     input,
	})
	if err != nil {
		logger.Error("Fatal to estimate gas (%v): %v", txType, err)
		return
	}

	feeByWei := new(big.Int).Mul(new(big.Int).SetUint64(estimateGas), auth.GasPrice).String()

	balance, err := h.Config.client.BalanceAt(ctx, auth.From, nil)
	if err != nil {
		logger.Error("Fatal to get %v balance: %v", auth.From.String(), err)
		return
	}

	logger.Info("from = %v, balance = %v", auth.From.String(), balance)
	logger.Info("%v Contract fee = gas(%v) * gasPrice(%v) = %v", txType, estimateGas, auth.GasPrice.String(), feeByWei)

}

func (h *Handler) sendTx(ctx context.Context, auth *bind.TransactOpts, data []byte, contract *common.Address) *types.Transaction {
	var (
		rawTx    *types.Transaction
		txSigned *types.Transaction
		err      error
	)

	// Create and Sign the transaction, sign it and schedule it for execution
	if contract == nil {
		rawTx = types.NewContractCreation(auth.Nonce.Uint64(), auth.Value, auth.GasLimit, auth.GasPrice, data)
	} else {
		rawTx = types.NewTransaction(auth.Nonce.Uint64(), *contract, auth.Value, auth.GasLimit, auth.GasPrice, data)
	}

	switch {
	case h.Config.key != nil:
		txSigned, err = types.SignTx(rawTx, types.NewEIP155Signer(h.Config.Chain.ID), h.Config.key)
	case h.Config.ks != nil:
		txSigned, err = h.Config.ks.SignTxWithPassphrase(
			accounts.Account{Address: common.HexToAddress(h.Config.Account)},
			h.Config.Password,
			rawTx,
			h.Config.Chain.ID)
	default:
		logger.FatalError("Fatal to get sign function")
	}

	if err != nil {
		logger.FatalError("Fatal to sign tx %v: %v", h.Config.Account, err)
	}

	err = h.Config.client.SendTransaction(ctx, txSigned)
	if err != nil {
		logger.FatalError("Fatal to send tx: %v", err)
	}
	return txSigned
}

func (h *Handler) DeployContract(ctx context.Context) {
	auth := h.Config.makeAuth(ctx, 0)

	logger.Info("Deploy contract...")

	input := common.FromHex(htlc.HtlcBin)

	//estimate deploy contract fee
	h.estimateGas(ctx, auth, "Deploy", input)

	//deploy-contract prompt
	h.Config.promptConfirm("Deploy")

	//send tx
	txSigned := h.sendTx(ctx, auth, input, nil)

	//update contract address
	h.Config.Contract = crypto.CreateAddress(auth.From, txSigned.Nonce()).String()

	logger.Info("contract address = %v", h.Config.Contract)
	logger.Info("transaction hash = %v", txSigned.Hash().String())

	//update config
	h.Config.rotate(h.ConfigPath)
}

func (h *Handler) StatContract(ctx context.Context) {
	//TODO
}

func (h *Handler) NewContract(ctx context.Context, participant common.Address, amount int64, hashLock [32]byte, timeLock *big.Int) *types.Transaction {
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
	h.Config.promptConfirm("Call")

	//send tx
	contract := common.HexToAddress(h.Config.Contract)

	return h.sendTx(ctx, auth, input, &contract)
}

func (h *Handler) GetContractId(ctx context.Context, txID common.Hash) HtlcLogHTLCNew {
	logger.Info("%v(%v) txid: %v", h.Config.Chain.Name, h.Config.Chain.ID, txID.String())
	logger.Info("contract address: %v\n", h.Config.Chain.Contract)

	receipt, err := h.Config.client.TransactionReceipt(ctx, txID)
	if err != nil {
		logger.FatalError("Fatal to get tx %v receipt: %v", txID.String(), err)
	}

	var logHTLCEvent HtlcLogHTLCNew
	parsedABI, err := abi.JSON(strings.NewReader(htlc.HtlcABI))
	if err != nil {
		logger.FatalError("Fatal to parse HtlcABI: %v", err)
	}

	if len(receipt.Logs) == 0 {
		logger.FatalError("len(receipt.Logs) == 0, receipt.Status = %v", receipt.Status)
	}

	if err := parsedABI.Unpack(&logHTLCEvent, "LogHTLCNew", receipt.Logs[0].Data); err != nil {
		logger.FatalError("Fatal to unpack log data for LogHTLCNew: %v", err)
	}

	logHTLCEvent.ContractId = receipt.Logs[0].Topics[1]
	logHTLCEvent.Sender = common.HexToAddress(receipt.Logs[0].Topics[2].Hex())
	logHTLCEvent.Receiver = common.HexToAddress(receipt.Logs[0].Topics[3].Hex())

	return logHTLCEvent
}

func (h *Handler) AuditContract(ctx context.Context, result interface{}, method string, contractId common.Hash) {
	parsedABI, err := abi.JSON(strings.NewReader(htlc.HtlcABI))
	if err != nil {
		logger.FatalError("Fatal to parse HtlcABI: %v", err)
	}

	input, err := parsedABI.Pack(method, contractId)
	if err != nil {
		logger.FatalError("Fatal to pack %v: %v", method, err)
	}

	//Call
	from := common.HexToAddress(h.Config.Account)
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

func (h *Handler) Redeem(ctx context.Context, contractId common.Hash, secret common.Hash) *types.Transaction {
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
	h.Config.promptConfirm("Call")

	//send tx
	contract := common.HexToAddress(h.Config.Chain.Contract)

	return h.sendTx(ctx, auth, input, &contract)
}

func (h *Handler) Refund(ctx context.Context, contractId common.Hash) *types.Transaction {
	auth := h.Config.makeAuth(ctx, 0)

	logger.Info("Call Withdraw ...")

	parsedABI, err := abi.JSON(strings.NewReader(htlc.HtlcABI))
	if err != nil {
		logger.Error("Fatal to parse HtlcABI: %v", err)
		return nil
	}

	input, err := parsedABI.Pack("refund", contractId)
	if err != nil {
		logger.Error("Fatal to pack newContract: %v", err)
		return nil
	}

	//estimate call contract fee
	h.estimateGas(ctx, auth, "Call", input)

	//call-contract prompt
	h.Config.promptConfirm("Call")

	//send tx
	contract := common.HexToAddress(h.Config.Chain.Contract)

	return h.sendTx(ctx, auth, input, &contract)
}
