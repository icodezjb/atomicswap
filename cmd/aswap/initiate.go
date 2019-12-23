package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/icodezjb/atomicswap/cmd"
	"github.com/icodezjb/atomicswap/logger"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/spf13/cobra"
)

func init() {
	initiateCmd.Flags().StringVarP(
		&participant,
		"participant",
		"p",
		"",
		"participant address")

	initiateCmd.Flags().Int64VarP(
		&initiateAmount,
		"amount",
		"a",
		-1,
		"amount of atomicswap asset")

	initiateCmd.Flags().StringVar(
		&privateKey,
		"key",
		"",
		"the private key of the account without '0x' prefix. if specified, the keystore will no longer be used")

	_ = initiateCmd.MarkFlagRequired("participant")
	_ = initiateCmd.MarkFlagRequired("amount")
}

var (
	participant    string
	initiateAmount int64
)

var initiateCmd = &cobra.Command{
	Use:   "initiate --participant <participant address> --amount <amount> [--key <private key>]",
	Short: "performed by the initiator to create the first contract",
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := h.Config.ParseConfig(h.ConfigPath); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	},
	Run: func(init *cobra.Command, args []string) {
		//check contract address
		h.Config.ValidateAddress(h.Config.Contract)

		//check participant address
		h.Config.ValidateAddress(participant)

		timeLock := new(big.Int).SetInt64(time.Now().Unix() + lock48Hour)
		hashPair := cmd.NewSecretHashPair()
		logger.Event("\nSecret = %s\nSecret Hash = %s", hexutil.Encode(hashPair.Secret[:]), hexutil.Encode(hashPair.Hash[:]))

		//connect to chain
		h.Config.Connect("")

		//Unlock account
		h.Config.Unlock(privateKey)

		txSigned := h.NewContract(context.Background(), common.HexToAddress(participant), initiateAmount, hashPair.Hash, timeLock)
		logger.Info("%v(%v) txid: %v\n", h.Config.Chain.Name, h.Config.Chain.ID, txSigned.Hash().String())
	},
}
