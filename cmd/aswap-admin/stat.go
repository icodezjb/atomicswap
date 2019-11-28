package main

import "github.com/spf13/cobra"

var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "stat the atomicswap contract",
	PreRun: func(cmd *cobra.Command, args []string) {
		h.ParseConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		h.Connect()
		h.StatContract()
	},
}
