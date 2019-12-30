// Copyright 2019 icodezjb
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package cmd

import (
	"context"
	"log"
	"math/big"
	"strings"

	htlc "github.com/icodezjb/atomicswap/contract"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

type Handler struct {
	ConfigPath string
	Config     *Config
}

func (h *Handler) estimateGas(ctx context.Context, auth *bind.TransactOpts, txType string, input []byte) error {
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
		return errors.Wrapf(err, "estimate gas (%v)", txType)
	}

	feeByWei := new(big.Int).Mul(new(big.Int).SetUint64(estimateGas), auth.GasPrice).String()

	balance, err := h.Config.client.BalanceAt(ctx, auth.From, nil)
	if err != nil {
		return errors.Wrapf(err, "account=%v get balance", auth.From.String())
	}

	log.Printf("from = %v, balance = %v", auth.From.String(), balance)
	log.Printf("%v Contract fee = gas(%v) * gasPrice(%v) = %v", txType, estimateGas, auth.GasPrice.String(), feeByWei)

	return nil
}

func (h *Handler) sendTx(ctx context.Context, auth *bind.TransactOpts, data []byte, contract *common.Address) (*types.Transaction, error) {
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
		return nil, errors.New("unexpected sendTx error")
	}

	if err != nil {
		return nil, errors.Wrapf(err, "account=%v sign tx ", h.Config.Account)
	}

	err = h.Config.client.SendTransaction(ctx, txSigned)
	if err != nil {
		return nil, errors.Wrapf(err, "account=%v send tx", h.Config.Account)
	}

	return txSigned, nil
}

func (h *Handler) DeployContract(ctx context.Context) error {
	auth, err := h.Config.makeAuth(ctx, 0)
	if err != nil {
		return err
	}

	log.Println("Deploy contract...")

	input := common.FromHex(htlc.HTLCBIN)

	//estimate deploy contract fee
	if err := h.estimateGas(ctx, auth, "Deploy", input); err != nil {
		return err
	}

	//deploy-contract prompt
	h.Config.promptConfirm("Deploy")

	//send tx
	txSigned, err := h.sendTx(ctx, auth, input, nil)
	if err != nil {
		return err
	}

	//update contract address
	h.Config.Contract = crypto.CreateAddress(auth.From, txSigned.Nonce()).String()

	log.Printf("contract address = %v", h.Config.Contract)
	log.Printf("transaction hash = %v", txSigned.Hash().String())

	//update config
	return h.Config.rotate(h.ConfigPath)
}

func (h *Handler) StatContract(ctx context.Context) error {
	//TODO
	return nil
}

func (h *Handler) NewContract(ctx context.Context, participant common.Address, amount int64, hashLock [32]byte, timeLock *big.Int) (*types.Transaction, error) {
	auth, err := h.Config.makeAuth(ctx, amount)
	if err != nil {
		return nil, errors.Wrapf(err, "make auth %v", h.Config.Account)
	}

	log.Println("Call NewContract ...")

	parsedABI, err := abi.JSON(strings.NewReader(htlc.HTLCABI))
	if err != nil {
		return nil, errors.Wrap(err, "parse HTLCABI")
	}

	input, err := parsedABI.Pack("newContract", participant, hashLock, timeLock)
	if err != nil {
		return nil, errors.Wrap(err, "pack newContract")
	}

	//estimate call contract fee
	if err = h.estimateGas(ctx, auth, "Call", input); err != nil {
		return nil, err
	}

	//call-contract prompt
	h.Config.promptConfirm("Call")

	//send tx
	contract := common.HexToAddress(h.Config.Contract)

	return h.sendTx(ctx, auth, input, &contract)
}

func (h *Handler) GetContractId(ctx context.Context, txID common.Hash) (*HtlcLogHTLCNew, error) {
	receipt, err := h.Config.client.TransactionReceipt(ctx, txID)
	if err != nil {
		return nil, errors.Wrapf(err, "get txid=%v receipt", txID.String())
	}

	var logHTLCEvent HtlcLogHTLCNew
	parsedABI, err := abi.JSON(strings.NewReader(htlc.HTLCABI))
	if err != nil {
		return nil, errors.Wrap(err, "parse HTLCABI")
	}

	if len(receipt.Logs) == 0 {
		return nil, errors.Errorf("len(receipt.Logs) == 0, receipt.Status = %v", receipt.Status)
	}

	if err := parsedABI.Unpack(&logHTLCEvent, "LogHTLCNew", receipt.Logs[0].Data); err != nil {
		return nil, errors.Wrap(err, "unpack log data for LogHTLCNew")
	}

	logHTLCEvent.ContractId = receipt.Logs[0].Topics[1]
	logHTLCEvent.Sender = common.HexToAddress(receipt.Logs[0].Topics[2].Hex())
	logHTLCEvent.Receiver = common.HexToAddress(receipt.Logs[0].Topics[3].Hex())

	return &logHTLCEvent, nil
}

func (h *Handler) AuditContract(ctx context.Context, result interface{}, contractId common.Hash) error {
	method := "getContract"

	parsedABI, err := abi.JSON(strings.NewReader(htlc.HTLCABI))
	if err != nil {
		return errors.Wrap(err, "parse HTLCABI")
	}

	input, err := parsedABI.Pack(method, contractId)
	if err != nil {
		return errors.Wrap(err, "pack getContract")
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
			return errors.Wrap(err, "call CodeAt")
		} else if len(code) == 0 {
			return errors.New("no contract code at given address")
		}
	}

	if err != nil {
		return errors.Wrap(err, "call CallContract")
	}

	if err = parsedABI.Unpack(result, method, output); err != nil {
		return errors.Wrap(err, "unpack result of contract call")
	}

	return nil
}

func (h *Handler) Redeem(ctx context.Context, contractId common.Hash, secret common.Hash) (*types.Transaction, error) {
	auth, err := h.Config.makeAuth(ctx, 0)
	if err != nil {
		return nil, errors.Wrapf(err, "make auth %v", h.Config.Account)
	}

	log.Println("Call Withdraw ...")

	parsedABI, err := abi.JSON(strings.NewReader(htlc.HTLCABI))
	if err != nil {
		return nil, errors.Wrap(err, "parse HTLCABI")
	}

	input, err := parsedABI.Pack("withdraw", contractId, secret)
	if err != nil {
		return nil, errors.Wrap(err, "pack withdraw")
	}

	//estimate call contract fee
	if err = h.estimateGas(ctx, auth, "Call", input); err != nil {
		return nil, err
	}

	//call-contract prompt
	h.Config.promptConfirm("Call")

	//send tx
	contract := common.HexToAddress(h.Config.Chain.Contract)

	return h.sendTx(ctx, auth, input, &contract)
}

func (h *Handler) Refund(ctx context.Context, contractId common.Hash) (*types.Transaction, error) {
	auth, err := h.Config.makeAuth(ctx, 0)
	if err != nil {
		return nil, errors.Wrapf(err, "make auth %v", h.Config.Account)
	}

	log.Println("Call Withdraw ...")

	parsedABI, err := abi.JSON(strings.NewReader(htlc.HTLCABI))
	if err != nil {
		return nil, errors.Wrap(err, "parse HTLCABI")
	}

	input, err := parsedABI.Pack("refund", contractId)
	if err != nil {
		return nil, errors.Wrap(err, "pack refund")
	}

	//estimate call contract fee
	if err = h.estimateGas(ctx, auth, "Call", input); err != nil {
		return nil, err
	}

	//call-contract prompt
	h.Config.promptConfirm("Call")

	//send tx
	contract := common.HexToAddress(h.Config.Chain.Contract)

	return h.sendTx(ctx, auth, input, &contract)
}
