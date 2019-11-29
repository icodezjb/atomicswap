package main

import (
	"fmt"
	"time"

	util "github.com/icodezjb/atomicswap/contract/helper"
	"github.com/icodezjb/atomicswap/logger"

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
		&amount,
		"amount",
		"a",
		-1,
		"amount of atomicswap asset")

	_ = initiateCmd.MarkFlagRequired("participant")
	_ = initiateCmd.MarkFlagRequired("amount")
}

var (
	participant string
	amount      int64
)

var initiateCmd = &cobra.Command{
	Use:   "initiate --participant <participant address> --amount <amount>",
	Short: "performed by the initiator to create the first contract",
	PreRun: func(cmd *cobra.Command, args []string) {
		h.ParseConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		//check contract address
		h.ValidateAddress(h.Config.Contract)

		//check participant address
		h.ValidateAddress(participant)

		timeLock := fmt.Sprintf("%x", time.Now().Unix()+lock48Hour)
		hashPair := util.NewSecretHashPair()
		hashLock := hexutil.Encode(hashPair.Hash[:])
		logger.Event("\nSecret = %s\nSecret Hash = %s", hexutil.Encode(hashPair.Secret[:]), hashLock)

		// connect to chain
		h.Connect()

		// without "0x" prefix
		h.NewContract(participant[2:], hashLock[2:], timeLock)
	},
}
