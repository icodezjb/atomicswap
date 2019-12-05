package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func init() {
	getContractIdCmd.Flags().StringVar(
		&txid,
		"txid",
		"",
		"the initiate txid")

	getContractIdCmd.Flags().StringVar(
		&otherContract,
		"other",
		"",
		"contract address")

	_ = getContractIdCmd.MarkFlagRequired("txid")
}

var txid string

var getContractIdCmd = &cobra.Command{
	Use:   "getcontractid --txid <initiator or participant txid> [--other <contract address>]",
	Short: "get the atomicswap contract id with the specified initiate txid",
	PreRun: func(cmd *cobra.Command, args []string) {
		h.Config.ParseConfig(h.ConfigPath)
	},
	Run: func(cmd *cobra.Command, args []string) {
		h.Config.Connect(otherContract)

		h.Config.ValidateAddress(h.Config.Chain.Contract)

		h.GetContractId(common.HexToHash(txid))
	},
}
