package main

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func init() {
	participantCmd.Flags().StringVarP(
		&initiator,
		"initiator",
		"i",
		"",
		"initiator address")

	participantCmd.Flags().Int64VarP(
		&participateAmount,
		"amount",
		"a",
		-1,
		"amount of atomicswap asset")

	participantCmd.Flags().Int64VarP(
		&untilTime,
		"time",
		"t",
		0,
		"unix time of the initiator timelock")

	participantCmd.Flags().StringVar(
		&hash,
		"hash",
		"",
		"the hash of the initiator secret")

	_ = participantCmd.MarkFlagRequired("initiator")
	_ = participantCmd.MarkFlagRequired("amount")
	_ = participantCmd.MarkFlagRequired("time")
	_ = participantCmd.MarkFlagRequired("hash")
}

var (
	initiator         string
	participateAmount int64
	untilTime         int64
	hash              string
)

var participantCmd = &cobra.Command{
	Use:   "participant --initiator <initiator address> --amount <amount> --time <unix time> --hash <secret hash>",
	Short: "performed by the participant to create the second contract",
	PreRun: func(cmd *cobra.Command, args []string) {
		h.ParseConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		//check contract address
		h.ValidateAddress(h.Config.Contract)

		//check participant address
		h.ValidateAddress(initiator)

		//half of initiator timelock
		timeLock := new(big.Int).SetInt64(time.Now().Unix() + (untilTime-time.Now().Unix())/2)
		//TODO: check timelock

		secretHash := common.HexToHash(hash)
		//connect to chain
		h.Connect()

		h.NewContract(common.HexToAddress(initiator), participateAmount, secretHash, timeLock)
	},
}
