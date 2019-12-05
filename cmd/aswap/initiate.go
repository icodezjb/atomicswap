package main

import (
	"math/big"
	"time"

	util "github.com/icodezjb/atomicswap/contract/helper"
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

	_ = initiateCmd.MarkFlagRequired("participant")
	_ = initiateCmd.MarkFlagRequired("amount")
}

var (
	participant    string
	initiateAmount int64
)

var initiateCmd = &cobra.Command{
	Use:   "initiate --participant <participant address> --amount <amount>",
	Short: "performed by the initiator to create the first contract",
	PreRun: func(cmd *cobra.Command, args []string) {
		h.Config.ParseConfig(h.ConfigPath)
	},
	Run: func(cmd *cobra.Command, args []string) {
		//check contract address
		h.Config.ValidateAddress(h.Config.Contract)

		//check participant address
		h.Config.ValidateAddress(participant)

		timeLock := new(big.Int).SetInt64(time.Now().Unix() + lock48Hour)
		hashPair := util.NewSecretHashPair()
		logger.Event("\nSecret = %s\nSecret Hash = %s", hexutil.Encode(hashPair.Secret[:]), hexutil.Encode(hashPair.Hash[:]))

		//connect to chain
		h.Config.Connect("")

		//Unlock account
		h.Config.Unlock()

		h.NewContract(common.HexToAddress(participant), initiateAmount, hashPair.Hash, timeLock)
	},
}
