package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func init() {
	redeemCmd.Flags().StringVar(
		&contractId,
		"id",
		"",
		"the contractId of the atomicswap pair")

	redeemCmd.Flags().StringVar(
		&secret,
		"secret",
		"",
		"the secret of hashlock")

	_ = auditContractCmd.MarkFlagRequired("id")
	_ = auditContractCmd.MarkFlagRequired("secret")
}

var secret string

var redeemCmd = &cobra.Command{
	Use:   "redeem --id <contractId> --secret <secret>",
	Short: "redeem from the atomicswap pair by the secret",
	PreRun: func(cmd *cobra.Command, args []string) {
		h.Config.ParseConfig(h.ConfigPath)
	},
	Run: func(cmd *cobra.Command, args []string) {
		//check contract address
		h.Config.ValidateAddress(h.Config.Contract)

		//connect to chain
		h.Config.Connect()

		//Unlock account
		h.Config.Unlock()

		contractId := common.HexToHash(contractId)
		secret := common.HexToHash(secret)

		h.Redeem(contractId, secret)
	},
}
