package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func init() {
	getContractIdCmd.Flags().StringVar(
		&initiateTxid,
		"txid",
		"",
		"the initiate txid")
}

var initiateTxid string

var getContractIdCmd = &cobra.Command{
	Use:   "getcontractid --id <initiate txid>",
	Short: "get the atomicswap contract id with the specified initiate txid",
	PreRun: func(cmd *cobra.Command, args []string) {
		h.ParseConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		h.Connect()
		h.GetContractId(common.HexToHash(initiateTxid))
	},
}
