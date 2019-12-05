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

	redeemCmd.Flags().StringVar(
		&otherContract,
		"other",
		"",
		"contract address")

	_ = redeemCmd.MarkFlagRequired("id")
	_ = redeemCmd.MarkFlagRequired("secret")
	_ = redeemCmd.MarkFlagRequired("other")
}

var secret string

var redeemCmd = &cobra.Command{
	Use:   "redeem --id <contractId> --secret <secret> --other <contract address>",
	Short: "redeem on the other contract by the secret",
	PreRun: func(cmd *cobra.Command, args []string) {
		h.Config.ParseConfig(h.ConfigPath)
	},
	Run: func(cmd *cobra.Command, args []string) {
		h.Config.Connect(otherContract)

		h.Config.ValidateAddress(h.Config.Chain.Contract)

		h.Config.Unlock()

		contractId := common.HexToHash(contractId)
		secret := common.HexToHash(secret)

		h.Redeem(contractId, secret)
	},
}
