package main

import "github.com/spf13/cobra"

var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "stat the atomicswap contract",
	PreRun: func(cmd *cobra.Command, args []string) {
		h.Config.ParseConfig(h.ConfigPath)
	},
	Run: func(cmd *cobra.Command, args []string) {
		h.Config.Connect()
		h.StatContract()
	},
}
