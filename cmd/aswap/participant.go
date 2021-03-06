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
	"github.com/spf13/cobra"
)

func init() {
	participantCmd.Flags().StringVarP(
		&initiator,
		"initiator",
		"i",
		"",
		"initiator address")

	participantCmd.Flags().Int64VarP(
		&participateAmount,
		"amount",
		"a",
		-1,
		"amount of atomicswap asset")

	participantCmd.Flags().Int64VarP(
		&untilTime,
		"time",
		"t",
		0,
		"unix time of the initiator timelock")

	participantCmd.Flags().StringVar(
		&hash,
		"hash",
		"",
		"the hash of the initiator secret")

	participantCmd.Flags().StringVar(
		&privateKey,
		"key",
		"",
		"the private key of the account without '0x' prefix. if specified, the keystore will no longer be used")

	_ = participantCmd.MarkFlagRequired("initiator")
	_ = participantCmd.MarkFlagRequired("amount")
	_ = participantCmd.MarkFlagRequired("time")
	_ = participantCmd.MarkFlagRequired("hash")
}

var (
	initiator         string
	participateAmount int64
	untilTime         int64
	hash              string
)

var participantCmd = &cobra.Command{
	Use:   "participant --initiator <initiator address> --amount <amount> --time <unix time> --hash <secret hash> [--key <private key>]",
	Short: "performed by the participant to create the second contract",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return h.Config.ParseConfig(h.ConfigPath)
	},
	Run: func(_ *cobra.Command, args []string) {
		//check contract address
		cmd.Must(h.Config.ValidateAddress(h.Config.Contract))

		//check participant address
		cmd.Must(h.Config.ValidateAddress(initiator))

		//half of initiator timelock
		timeLock := new(big.Int).SetInt64(time.Now().Unix() + (untilTime-time.Now().Unix())/2)
		//TODO: check timelock

		secretHash := common.HexToHash(hash)
		//connect to chain
		cmd.Must(h.Config.Connect(""))

		//Unlock account
		cmd.Must(h.Config.Unlock(privateKey))

		txSigned, err := h.NewContract(context.Background(), common.HexToAddress(initiator), participateAmount, secretHash, timeLock)
		cmd.Must(err)

		log.Printf("%v(%v) txid: %v", h.Config.Chain.Name, h.Config.Chain.ID, txSigned.Hash().String())
	},
}
