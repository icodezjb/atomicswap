package main

import (
	"github.com/icodezjb/atomicswap/cmd"
	"log"

	"github.com/spf13/cobra"
)

const lock48Hour = 48 * 60 * 60

var (
	h cmd.Handler

	rootCmd = &cobra.Command{
		Use:   "aswap",
		Short: "atomic swap between two different blockchains which based on EVM",
	}

	//auditContractCmd, redeemCmd
	contractId string
	//auditContractCmd, redeemCmd, getContractIdCmd
	otherContract string
	//initiateCmd, participantCmd, redeemCmd, refundCmd
	privateKey string
)

func init() {
	h.Config = new(cmd.Config)

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
    "chainID": 110,
    "chainName": "ETH1",
    "url": "http://127.0.0.1:8545",
    "otherChainID": 111,
    "otherChainName": "ETH2",
    "otherURL": "http://127.0.0.1:7545",
    "account": "0xffd79941b7085805f48ded97298694c6bb950e2c",
    "keystoreDir": "/absolute/path/",
    "password": "password"
}
`)

}

func main() {
	rootCmd.Version = cmd.VersionFunc()

	rootCmd.AddCommand(initiateCmd)
	rootCmd.AddCommand(participantCmd)
	rootCmd.AddCommand(getContractIdCmd)
	rootCmd.AddCommand(auditContractCmd)
	rootCmd.AddCommand(redeemCmd)
	rootCmd.AddCommand(refundCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
