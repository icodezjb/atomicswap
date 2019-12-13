package main

import (
	"context"

	"github.com/spf13/cobra"
)

func init() {
	deployCmd.Flags().StringVar(
		&privateKey,
		"key",
		"",
		"the private key of the account without '0x' prefix. if specified, the keystore will no longer be used")
}

var privateKey string

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy the atomicswap contract",
	PreRun: func(cmd *cobra.Command, args []string) {
		h.Config.ParseConfig(h.ConfigPath)
	},
	Run: func(cmd *cobra.Command, args []string) {
		h.Config.Connect("")
		h.Config.ValidateAddress(h.Config.Account)
		h.Config.Unlock(privateKey)
		h.DeployContract(context.Background())
	},
}
