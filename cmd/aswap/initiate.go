package main

import (
	"github.com/spf13/cobra"
)

var initiateCmd = &cobra.Command{
	Use:   "initiate <participant address> <amount>",
	Short: "performed by the initiator to create the first contract",
	PreRun: func(cmd *cobra.Command, args []string) {
		h.ParseConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {

	},
}
