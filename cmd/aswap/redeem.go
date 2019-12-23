package main

import (
	"context"
	"fmt"
	"os"

	"github.com/icodezjb/atomicswap/logger"

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

	redeemCmd.Flags().StringVar(
		&privateKey,
		"key",
		"",
		"the private key of the account without '0x' prefix. if specified, the keystore will no longer be used")

	_ = redeemCmd.MarkFlagRequired("id")
	_ = redeemCmd.MarkFlagRequired("secret")
	_ = redeemCmd.MarkFlagRequired("other")
}

var secret string

var redeemCmd = &cobra.Command{
	Use:   "redeem --id <contractId> --secret <secret> --other <contract address> [--key <private key>]",
	Short: "redeem once they know secret which is the preimage of the hashlock AND the time lock has no expired ",
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := h.Config.ParseConfig(h.ConfigPath); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		h.Config.Connect(otherContract)

		h.Config.ValidateAddress(h.Config.Chain.Contract)

		h.Config.Unlock(privateKey)

		contractId := common.HexToHash(contractId)
		secret := common.HexToHash(secret)

		txSigned := h.Redeem(context.Background(), contractId, secret)
		logger.Info("%v(%v) txid: %v\n", h.Config.Chain.Name, h.Config.Chain.ID, txSigned.Hash().String())
	},
}
