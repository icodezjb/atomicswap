package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/icodezjb/atomicswap/cmd"
	"github.com/icodezjb/atomicswap/logger"

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
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := h.Config.ParseConfig(h.ConfigPath); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		h.Config.Connect(otherContract)

		h.Config.ValidateAddress(h.Config.Chain.Contract)

		logHTLCEvent := h.GetContractId(context.Background(), common.HexToHash(txid))

		printEvent(logHTLCEvent)
	},
}

func printEvent(e cmd.HtlcLogHTLCNew) {
	logger.Info("ContractId = %s", hexutil.Encode(e.ContractId[:]))
	logger.Info("Sender     = %s", e.Sender.String())
	logger.Info("Receiver   = %s", e.Receiver.String())
	logger.Info("Amount     = %s", e.Amount)
	logger.Info("TimeLock   = %s (%s)", e.Timelock, time.Unix(e.Timelock.Int64(), 0).Format(time.RFC3339))
	logger.Info("SecretHash = %s", hexutil.Encode(e.Hashlock[:]))
}
