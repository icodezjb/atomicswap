package main

import (
	"context"

	"github.com/icodezjb/atomicswap/cmd"

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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return h.Config.ParseConfig(h.ConfigPath)
	},
	Run: func(_ *cobra.Command, args []string) {
		cmd.Must(h.Config.Connect(""))

		cmd.Must(h.Config.ValidateAddress(h.Config.Account))

		cmd.Must(h.Config.Unlock(privateKey))

		cmd.Must(h.DeployContract(context.Background()))
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}
