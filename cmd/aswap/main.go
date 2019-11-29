package main

import (
	"fmt"
	"os"

	"github.com/icodezjb/atomicswap/cmd"

	"github.com/spf13/cobra"
)

const lock48Hour = 48 * 60 * 60

var (
	h cmd.Handle

	rootCmd = &cobra.Command{
		Use:   "aswap",
		Short: "atomic swap between two different blockchains which based on EVM",
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&h.ConfigPath,
		"config",
		"c",
		"config.json",
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
  "contract": "0xA63684e9aCfb86330Ce7e8001C6aAa3e23DD6Fe7"
}
`)

}

func main() {
	rootCmd.Version = cmd.VersionFunc()

	rootCmd.AddCommand(initiateCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
