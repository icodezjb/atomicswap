// Copyright 2019 icodezjb
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"context"

	"github.com/icodezjb/atomicswap/cmd"

	"github.com/spf13/cobra"
)

func init() {
	deployCmd.Flags().StringVar(
		&privateKey,
		"key",
		"",
		"the private key of the account without '0x' prefix. if specified, the keystore will no longer be used")
}

var privateKey string

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy the atomicswap contract",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return h.Config.ParseConfig(h.ConfigPath)
	},
	Run: func(_ *cobra.Command, args []string) {
		cmd.Must(h.Config.Connect(""))

		cmd.Must(h.Config.ValidateAddress(h.Config.Account))

		cmd.Must(h.Config.Unlock(privateKey))

		cmd.Must(h.DeployContract(context.Background()))
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}
