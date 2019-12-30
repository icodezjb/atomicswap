package main

import (
	"context"
	"fmt"

	"github.com/icodezjb/atomicswap/cmd"

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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return h.Config.ParseConfig(h.ConfigPath)
	},
	Run: func(_ *cobra.Command, args []string) {
		//check contract address
		cmd.Must(h.Config.ValidateAddress(h.Config.Contract))

		//connect to chain
		cmd.Must(h.Config.Connect(""))

		//Unlock account
		cmd.Must(h.Config.Unlock(privateKey))

		contractId := common.HexToHash(contractId)

		txSigned, err := h.Refund(context.Background(), contractId)
		cmd.Must(err)

		fmt.Printf("%v(%v) txid: %v\n", h.Config.Chain.Name, h.Config.Chain.ID, txSigned.Hash().String())
	},
}
