package main

import (
	"context"
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
		h.Config.ParseConfig(h.ConfigPath)
	},
	Run: func(_ *cobra.Command, args []string) {
		h.Config.Connect(otherContract)

		h.Config.ValidateAddress(h.Config.Account)

		h.Config.ValidateAddress(h.Config.Chain.Contract)

		logger.Info("Call getContract ...")
		logger.Info("contract address: %v\n", h.Config.Chain.Contract)

		contractDetails := new(cmd.ContractDetails)

		h.AuditContract(context.Background(), contractDetails, "getContract", common.HexToHash(contractId))

		logger.Info("Sender     = %s", contractDetails.Sender.String())
		logger.Info("Receiver   = %s", contractDetails.Receiver.String())
		logger.Info("Amount     = %s (wei)", contractDetails.Amount)
		logger.Info("TimeLock   = %s (%s)", contractDetails.Timelock, time.Unix(contractDetails.Timelock.Int64(), 0))
		logger.Info("SecretHash = %s", hexutil.Encode(contractDetails.Hashlock[:]))
		logger.Info("Withdrawn  = %v", contractDetails.Withdrawn)
		logger.Info("Refunded   = %v", contractDetails.Refunded)
		logger.Info("Secret     = %s", hexutil.Encode(contractDetails.Preimage[:]))
	},
}
