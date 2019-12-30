package main

import (
	"context"
	"fmt"
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

		fmt.Println("Call getContract ...")
		fmt.Printf("contract address: %v\n", h.Config.Chain.Contract)

		contractDetails := new(cmd.ContractDetails)

		cmd.Must(h.AuditContract(context.Background(), contractDetails, common.HexToHash(contractId)))

		printContractDetails(contractDetails)
	},
}

func printContractDetails(d *cmd.ContractDetails) {
	fmt.Printf("Sender     = %s\n", d.Sender.String())
	fmt.Printf("Receiver   = %s\n", d.Receiver.String())
	fmt.Printf("Amount     = %s (wei)\n", d.Amount)
	fmt.Printf("TimeLock   = %s (%s)\n", d.Timelock, time.Unix(d.Timelock.Int64(), 0))
	fmt.Printf("SecretHash = %s\n", hexutil.Encode(d.Hashlock[:]))
	fmt.Printf("Withdrawn  = %v\n", d.Withdrawn)
	fmt.Printf("Refunded   = %v\n", d.Refunded)
	fmt.Printf("Secret     = %s\n", hexutil.Encode(d.Preimage[:]))
}
