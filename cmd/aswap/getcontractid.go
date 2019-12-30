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
	"time"

	"github.com/icodezjb/atomicswap/cmd"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/spf13/cobra"
)

func init() {
	getContractIdCmd.Flags().StringVar(
		&txid,
		"txid",
		"",
		"the initiate txid")

	getContractIdCmd.Flags().StringVar(
		&otherContract,
		"other",
		"",
		"contract address")

	_ = getContractIdCmd.MarkFlagRequired("txid")
}

var txid string

var getContractIdCmd = &cobra.Command{
	Use:   "getcontractid --txid <initiator or participant txid> [--other <contract address>]",
	Short: "get the atomicswap contract id with the specified initiate txid",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return h.Config.ParseConfig(h.ConfigPath)
	},
	Run: func(_ *cobra.Command, args []string) {
		cmd.Must(h.Config.Connect(otherContract))

		cmd.Must(h.Config.ValidateAddress(h.Config.Chain.Contract))

		log.Printf("%s(%s) txid: %s", h.Config.Chain.Name, h.Config.Chain.ID, txid)
		log.Printf("contract address: %s", h.Config.Chain.Contract)

		logHTLCEvent, err := h.GetContractId(context.Background(), common.HexToHash(txid))
		cmd.Must(err)

		printEvent(logHTLCEvent)
	},
}

func printEvent(e *cmd.HtlcLogHTLCNew) {
	log.Printf("ContractId = %s", hexutil.Encode(e.ContractId[:]))
	log.Printf("Sender     = %s", e.Sender.String())
	log.Printf("Receiver   = %s", e.Receiver.String())
	log.Printf("Amount     = %s", e.Amount)
	log.Printf("TimeLock   = %s (%s)", e.Timelock, time.Unix(e.Timelock.Int64(), 0).Format(time.RFC3339))
	log.Printf("SecretHash = %s", hexutil.Encode(e.Hashlock[:]))
}
