package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/icodezjb/atomicswap/cmd/aswap-deploy/deploy"

	"github.com/spf13/cobra"
)

var (
	Commit    = "unknown-commit"
	BuildTime = "unknown-buildtime"

	version = "0.2.0"

	dHandle deploy.Handle

	rootCmd = &cobra.Command{
		Use:   "aswap-deploy",
		Short: "deploy atomicswap smart contract",
	}
)

// VersionFunc holds the textual version string.
func VersionFunc() string {
	return fmt.Sprintf(": %s\ncommit: %s\nbuild time: %s\ngolang version: %s\n",
		version, Commit, BuildTime, runtime.Version()+" "+runtime.GOOS+"/"+runtime.GOARCH)
}

func init() {
	rootCmd.Flags().StringVar(
		&dHandle.ConfigPath,
		"config",
		"./config.json",
		"config file path",
	)
}

func main() {
	rootCmd.Version = VersionFunc()

	rootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		dHandle.ParseConfig()
	}

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		dHandle.Connect()
		dHandle.DeployContract()
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
