package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/icodezjb/atomicswap/cmd/aswap-admin/deploy"

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

// versionFunc holds the textual version string.
func versionFunc() string {
	return fmt.Sprintf(": %s\ncommit: %s\nbuild time: %s\ngolang version: %s\n",
		version, Commit, BuildTime, runtime.Version()+" "+runtime.GOOS+"/"+runtime.GOARCH)
}

func init() {
	rootCmd.Flags().StringVarP(
		&dHandle.ConfigPath,
		"config",
		"c",
		"./config.json",
		"config file path",
	)

	rootCmd.SetUsageTemplate(`Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}

config-example.json 
{
  "chainId": 1,
  "chainName": "eth",
  "url": "http://127.0.0.1:7545",
  "from": "0xffd79941b7085805f48ded97298694c6bb950e2c",
  "keystoreDir": "./",
  "password": "password",
  "contract": ""
}
`)

}

func main() {
	rootCmd.Version = versionFunc()

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
