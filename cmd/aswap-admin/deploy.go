package main

import "github.com/spf13/cobra"

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy the atomicswap contract",
	PreRun: func(cmd *cobra.Command, args []string) {
		h.ParseConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		h.Connect()
		h.ValidateAddress(h.Config.From)
		h.DeployContract()
	},
}
