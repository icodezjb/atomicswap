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
	"math/big"
	"time"

	"github.com/icodezjb/atomicswap/cmd"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/spf13/cobra"
)

func init() {
	initiateCmd.Flags().StringVarP(
		&participant,
		"participant",
		"p",
		"",
		"participant address")

	initiateCmd.Flags().Int64VarP(
		&initiateAmount,
		"amount",
		"a",
		-1,
		"amount of atomicswap asset")

	initiateCmd.Flags().StringVar(
		&privateKey,
		"key",
		"",
		"the private key of the account without '0x' prefix. if specified, the keystore will no longer be used")

	_ = initiateCmd.MarkFlagRequired("participant")
	_ = initiateCmd.MarkFlagRequired("amount")
}

var (
	participant    string
	initiateAmount int64
)

var initiateCmd = &cobra.Command{
	Use:   "initiate --participant <participant address> --amount <amount> [--key <private key>]",
	Short: "performed by the initiator to create the first contract",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return h.Config.ParseConfig(h.ConfigPath)
	},
	Run: func(_ *cobra.Command, args []string) {
		//check contract address
		cmd.Must(h.Config.ValidateAddress(h.Config.Contract))

		//check participant address
		cmd.Must(h.Config.ValidateAddress(participant))

		timeLock := new(big.Int).SetInt64(time.Now().Unix() + lock48Hour)
		hashPair := cmd.NewSecretHashPair()
		log.Printf("\nSecret = %s\nSecret Hash = %s",
			hexutil.Encode(hashPair.Secret[:]), hexutil.Encode(hashPair.Hash[:]))

		//connect to chain
		cmd.Must(h.Config.Connect(""))

		//Unlock account
		cmd.Must(h.Config.Unlock(privateKey))

		txSigned, err := h.NewContract(context.Background(), common.HexToAddress(participant), initiateAmount, hashPair.Hash, timeLock)
		cmd.Must(err)

		log.Printf("%s(%s) txid: %s", h.Config.Chain.Name, h.Config.Chain.ID, txSigned.Hash().String())
	},
}
