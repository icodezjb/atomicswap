package main

import "github.com/spf13/cobra"

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy the atomicswap contract",
	PreRun: func(cmd *cobra.Command, args []string) {
		h.Config.ParseConfig(h.ConfigPath)
	},
	Run: func(cmd *cobra.Command, args []string) {
		h.Config.Connect()
		h.Config.ValidateAddress(h.Config.From)
		h.DeployContract()
	},
}
