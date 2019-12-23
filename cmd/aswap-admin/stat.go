package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "stat the atomicswap contract",
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := h.Config.ParseConfig(h.ConfigPath); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		h.Config.Connect("")
		h.StatContract(context.Background())
	},
}
