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
	redeemCmd.Flags().StringVar(
		&contractId,
		"id",
		"",
		"the contractId of the atomicswap pair")

	redeemCmd.Flags().StringVar(
		&secret,
		"secret",
		"",
		"the secret of hashlock")

	redeemCmd.Flags().StringVar(
		&otherContract,
		"other",
		"",
		"contract address")

	redeemCmd.Flags().StringVar(
		&privateKey,
		"key",
		"",
		"the private key of the account without '0x' prefix. if specified, the keystore will no longer be used")

	_ = redeemCmd.MarkFlagRequired("id")
	_ = redeemCmd.MarkFlagRequired("secret")
	_ = redeemCmd.MarkFlagRequired("other")
}

var secret string

var redeemCmd = &cobra.Command{
	Use:   "redeem --id <contractId> --secret <secret> --other <contract address> [--key <private key>]",
	Short: "redeem once they know secret which is the preimage of the hashlock AND the time lock has no expired ",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return h.Config.ParseConfig(h.ConfigPath)
	},
	Run: func(_ *cobra.Command, args []string) {
		cmd.Must(h.Config.Connect(otherContract))

		cmd.Must(h.Config.ValidateAddress(h.Config.Chain.Contract))

		cmd.Must(h.Config.Unlock(privateKey))

		contractId := common.HexToHash(contractId)
		secret := common.HexToHash(secret)

		txSigned, err := h.Redeem(context.Background(), contractId, secret)
		cmd.Must(err)

		log.Printf("%v(%v) txid: %v", h.Config.Chain.Name, h.Config.Chain.ID, txSigned.Hash().String())
	},
}
