package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func init() {
	auditContractCmd.Flags().StringVar(
		&contractId,
		"id",
		"",
		"the contractId of the atomicswap pair")

	_ = auditContractCmd.MarkFlagRequired("id")
}

var auditContractCmd = &cobra.Command{
	Use:   "auditcontract --id <contractId>",
	Short: "get the atomicswap pair details with the specified contractId",
	PreRun: func(cmd *cobra.Command, args []string) {
		h.Config.ParseConfig(h.ConfigPath)
	},
	Run: func(cmd *cobra.Command, args []string) {
		h.Config.Connect()
		h.Config.ValidateAddress(h.Config.From)
		h.AuditContract(common.HexToAddress(h.Config.From), common.HexToHash(contractId))
	},
}
