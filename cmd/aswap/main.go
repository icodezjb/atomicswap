package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aswap",
	Short: "atomic swap between two different block chains which based on EVM",
	Run: func(cmd *cobra.Command, args []string) {

	},
	//Args:cobra.OnlyValidArgs
}

var (
	Commit    = "unknown-commit"
	BuildTime = "unknown-buildtime"

	version = "0.1.0"
)

// VersionFunc holds the textual version string.
func VersionFunc() string {
	return fmt.Sprintf(": %s\ncommit: %s\nbuild time: %s\ngolang version: %s\n",
		version, Commit, BuildTime, runtime.Version()+" "+runtime.GOOS+"/"+runtime.GOARCH)
}

func main() {
	rootCmd.Version = VersionFunc()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
