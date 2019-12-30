package main

import (
	"context"
	"fmt"
	"time"

	"github.com/icodezjb/atomicswap/cmd"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return h.Config.ParseConfig(h.ConfigPath)
	},
	Run: func(_ *cobra.Command, args []string) {
		cmd.Must(h.Config.Connect(otherContract))

		cmd.Must(h.Config.ValidateAddress(h.Config.Chain.Contract))

		fmt.Printf("%v(%v) txid: %v", h.Config.Chain.Name, h.Config.Chain.ID, txid)
		fmt.Printf("contract address: %v\n", h.Config.Chain.Contract)

		logHTLCEvent, err := h.GetContractId(context.Background(), common.HexToHash(txid))
		cmd.Must(err)

		printEvent(logHTLCEvent)
	},
}

func printEvent(e *cmd.HtlcLogHTLCNew) {
	fmt.Printf("ContractId = %s\n", hexutil.Encode(e.ContractId[:]))
	fmt.Printf("Sender     = %s\n", e.Sender.String())
	fmt.Printf("Receiver   = %s\n", e.Receiver.String())
	fmt.Printf("Amount     = %s\n", e.Amount)
	fmt.Printf("TimeLock   = %s (%s)\n", e.Timelock, time.Unix(e.Timelock.Int64(), 0).Format(time.RFC3339))
	fmt.Printf("SecretHash = %s\n", hexutil.Encode(e.Hashlock[:]))
}
