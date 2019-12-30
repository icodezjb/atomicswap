package main

import (
	"context"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/icodezjb/atomicswap/cmd"
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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return h.Config.ParseConfig(h.ConfigPath)
	},
	Run: func(_ *cobra.Command, args []string) {
		cmd.Must(h.Config.Connect(otherContract))

		cmd.Must(h.Config.ValidateAddress(h.Config.Account))
		cmd.Must(h.Config.ValidateAddress(h.Config.Chain.Contract))

		log.Print("Call getContract ...")
		log.Printf("contract address: %s", h.Config.Chain.Contract)

		contractDetails := new(cmd.ContractDetails)

		cmd.Must(h.AuditContract(context.Background(), contractDetails, common.HexToHash(contractId)))

		printContractDetails(contractDetails)
	},
}

func printContractDetails(d *cmd.ContractDetails) {
	log.Printf("Sender     = %s", d.Sender.String())
	log.Printf("Receiver   = %s", d.Receiver.String())
	log.Printf("Amount     = %s (wei)", d.Amount)
	log.Printf("TimeLock   = %s (%s)", d.Timelock, time.Unix(d.Timelock.Int64(), 0))
	log.Printf("SecretHash = %s", hexutil.Encode(d.Hashlock[:]))
	log.Printf("Withdrawn  = %t", d.Withdrawn)
	log.Printf("Refunded   = %t", d.Refunded)
	log.Printf("Secret     = %s", hexutil.Encode(d.Preimage[:]))
}
