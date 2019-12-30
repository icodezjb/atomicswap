package main

import (
	"context"

	"github.com/icodezjb/atomicswap/cmd"

	"github.com/spf13/cobra"
)

var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "stat the atomicswap contract",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return h.Config.ParseConfig(h.ConfigPath)
	},
	Run: func(_ *cobra.Command, args []string) {
		cmd.Must(h.Config.Connect(""))

		cmd.Must(h.StatContract(context.Background()))
	},
}
