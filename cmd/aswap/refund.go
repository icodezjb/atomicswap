package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func init() {
	refundCmd.Flags().StringVar(
		&contractId,
		"id",
		"",
		"the contractId of the atomicswap pair")

	_ = auditContractCmd.MarkFlagRequired("id")
}

var refundCmd = &cobra.Command{
	Use:   "refund --id <contractId>",
	Short: "refund from the atomicswap pair",
	PreRun: func(cmd *cobra.Command, args []string) {
		h.ParseConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		//check contract address
		h.ValidateAddress(h.Config.Contract)

		//connect to chain
		h.Connect()

		//Unlock account
		h.Unlock()

		contractId := common.HexToHash(contractId)

		h.Refund(contractId)
	},
}
