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
	"log"

	"github.com/icodezjb/atomicswap/cmd"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func init() {
	refundCmd.Flags().StringVar(
		&contractId,
		"id",
		"",
		"the contractId of the atomicswap pair")

	refundCmd.Flags().StringVar(
		&privateKey,
		"key",
		"",
		"the private key of the account without '0x' prefix. if specified, the keystore will no longer be used")

	_ = refundCmd.MarkFlagRequired("id")
}

var refundCmd = &cobra.Command{
	Use:   "refund --id <contractId> [--key <private key>]",
	Short: "refund on the contract if there was no withdraw AND the time lock has expired",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return h.Config.ParseConfig(h.ConfigPath)
	},
	Run: func(_ *cobra.Command, args []string) {
		//check contract address
		cmd.Must(h.Config.ValidateAddress(h.Config.Contract))

		//connect to chain
		cmd.Must(h.Config.Connect(""))

		//Unlock account
		cmd.Must(h.Config.Unlock(privateKey))

		contractId := common.HexToHash(contractId)

		txSigned, err := h.Refund(context.Background(), contractId)
		cmd.Must(err)

		log.Printf("%v(%v) txid: %v", h.Config.Chain.Name, h.Config.Chain.ID, txSigned.Hash().String())
	},
}
