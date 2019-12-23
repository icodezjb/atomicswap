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
	auditContractCmd.Flags().StringVar(
		&contractId,
		"id",
		"",
		"the contractId of the atomicswap pair")

	auditContractCmd.Flags().StringVar(
		&otherContract,
		"other",
		"",
		"contract address")

	_ = auditContractCmd.MarkFlagRequired("id")
}

var auditContractCmd = &cobra.Command{
	Use:   "auditcontract --id <contractId> [--other <contract address>]",
	Short: "get the atomicswap pair details with the specified contractId",
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := h.Config.ParseConfig(h.ConfigPath); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	},
	Run: func(_ *cobra.Command, args []string) {
		h.Config.Connect(otherContract)

		h.Config.ValidateAddress(h.Config.Account)
		h.Config.ValidateAddress(h.Config.Chain.Contract)

		printBanner(h.Config.Chain.Contract)

		contractDetails := new(cmd.ContractDetails)

		h.AuditContract(context.Background(), contractDetails, "getContract", common.HexToHash(contractId))

		printContractDetails(contractDetails)
	},
}

func printBanner(contract string) {
	logger.Info("Call getContract ...")
	logger.Info("contract address: %v\n", contract)
}

func printContractDetails(d *cmd.ContractDetails) {
	logger.Info("Sender     = %s", d.Sender.String())
	logger.Info("Receiver   = %s", d.Receiver.String())
	logger.Info("Amount     = %s (wei)", d.Amount)
	logger.Info("TimeLock   = %s (%s)", d.Timelock, time.Unix(d.Timelock.Int64(), 0))
	logger.Info("SecretHash = %s", hexutil.Encode(d.Hashlock[:]))
	logger.Info("Withdrawn  = %v", d.Withdrawn)
	logger.Info("Refunded   = %v", d.Refunded)
	logger.Info("Secret     = %s", hexutil.Encode(d.Preimage[:]))
}
