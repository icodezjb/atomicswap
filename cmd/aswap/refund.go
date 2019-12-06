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

	refundCmd.Flags().StringVar(
		&privateKey,
		"key",
		"",
		"the private key of the account without '0x' prefix. if specified, the keystore will no longer be used")

	_ = refundCmd.MarkFlagRequired("id")
}

var refundCmd = &cobra.Command{
	Use:   "refund --id <contractId> [--key <private key>]",
	Short: "refund on the contract if there was no withdraw AND the time lock has expired",
	PreRun: func(cmd *cobra.Command, args []string) {
		h.Config.ParseConfig(h.ConfigPath)
	},
	Run: func(cmd *cobra.Command, args []string) {
		//check contract address
		h.Config.ValidateAddress(h.Config.Contract)

		//connect to chain
		h.Config.Connect("")

		//Unlock account
		h.Config.Unlock(privateKey)

		contractId := common.HexToHash(contractId)

		h.Refund(contractId)
	},
}
