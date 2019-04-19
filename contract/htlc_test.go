package htlc

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/icodezjb/atomicswap/contract/helper"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const hourSeconds = 3600
const oneFinney = 1000000000000000
const port  = "7545"
const REQUIRE_FAILED_MSG = "VM Exception while processing transaction: revert"

type Ganache struct  {
	cmd *exec.Cmd
	running bool
}

var ganache = Ganache {
	cmd: exec.Command(
		"/usr/local/bin/ganache-cli",
		"--account", "0xa5a1aca01671e2660f1ee47abfd7065d5d38f99fa4a53495f02df939cd5b86f6,111111111111111111111",
		"-p", port),
	running:true,
}

func setup()  {
	go func() {
		if err:= ganache.cmd.Run();err != nil {
			ganache.running = false
			fmt.Println("unexpect ganache-cli exit: ", err)
			os.Exit(1)
		}
	}()

	//waiting for rpc service
	time.Sleep(2*time.Second)
}

func cleanup() {
	if ganache.running {
		_ = ganache.cmd.Process.Kill()
	}
}

func makeAuth(t *testing.T, private string, client *ethclient.Client, value int64)(*bind.TransactOpts)  {
	privateKey, err := crypto.HexToECDSA(private)
	if err != nil {
		t.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce,err := client.PendingNonceAt(context.Background(),fromAddress)
	if err != nil {
		t.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	//in wei
	auth.Value = big.NewInt(value)
	auth.GasLimit = uint64(3000000) //in uints
	auth.GasPrice = gasPrice

	return auth
}

func TestNewContract(t *testing.T) {
	setup()
	defer cleanup()

	var timeLock1Hour = time.Now().Unix() + hourSeconds
	var senderKey   = "a5a1aca01671e2660f1ee47abfd7065d5d38f99fa4a53495f02df939cd5b86f6"
	var receiver = common.HexToAddress("0x5eb1231fd20ac35dff0d4295959cbd11c9cdae40")
	var hashPair = htlc.NewSecretHashPair()

	client, err := ethclient.Dial("http://127.0.0.1:7545")
	if err != nil {
		t.Fatal(err)
	}

	senderAuth := makeAuth(t, senderKey, client, 0)

	_, _, instance, err := DeployHtlc(senderAuth, client)
	if err != nil {
		t.Fatal(err)
	}

	senderAuth = makeAuth(t, senderKey, client, oneFinney)

	tx, err := instance.NewContract(senderAuth, receiver, hashPair.Hash, big.NewInt(timeLock1Hour))
	if err != nil {
		t.Fatal(err)
	}

	receipt, err :=	client.TransactionReceipt(context.Background(),	tx.Hash())
	if err != nil	{
		t.Fatal(err)
	}

	contractABI, err := abi.JSON(strings.NewReader(string(HtlcABI)))
	if err != nil {
		t.Fatal("Fatal read contract abi", err)
	}

	var logHTLCEvent HtlcLogHTLCNew
	if err := contractABI.Unpack(&logHTLCEvent, "LogHTLCNew", receipt.Logs[0].Data); err != nil {
		t.Fatal("Fatal unpack log data for LogHTLCNew", err)
	}

	logHTLCEvent.ContractId = receipt.Logs[0].Topics[1]
	logHTLCEvent.Sender = common.HexToAddress(receipt.Logs[0].Topics[2].Hex())
	logHTLCEvent.Receiver = common.HexToAddress(receipt.Logs[0].Topics[3].Hex())

	//Check logHTLCEvent
	if match, _ := regexp.MatchString("^0x[0-9a-f]{64}$", hexutil.Encode(logHTLCEvent.ContractId[:])); match == false {
		t.Fatal("logHTLCEvent.ContractId should be Sha256Hash")
	}

	if senderAuth.From.Hex() != logHTLCEvent.Sender.Hex() {
		t.Fatal("logHTLCEvent.Sender should be the specified sender")
	}

	if receiver.Hex() != logHTLCEvent.Receiver.Hex() {
		t.Fatal("logHTLCEvent.Receiver should be the specified receiver")
	}

	if big.NewInt(oneFinney).Cmp(logHTLCEvent.Amount) != 0 {
		t.Fatal("logHTLCEvent.Amount should be equal oneFinney")
	}

	if hexutil.Encode(hashPair.Hash[:]) != hexutil.Encode(logHTLCEvent.Hashlock[:]) {
		t.Fatal("logHTLCEvent.Hashlock should be the specified hashlock")
	}

	if big.NewInt(timeLock1Hour).Cmp(logHTLCEvent.Timelock) != 0 {
		t.Fatal("logHTLCEvent.Timelock should be the specified timelock")
	}

	//Check contract details on-chain
	contractDetails, err := instance.GetContract(&bind.CallOpts{From:senderAuth.From}, logHTLCEvent.ContractId)
	if err != nil {
		t.Fatal("Fatal call GetContract")
	}

	if senderAuth.From.Hex() != contractDetails.Sender.Hex() {
		t.Fatal("contractDetails.Sender should be the specified sender")
	}

	if receiver.Hex() != contractDetails.Receiver.Hex() {
		t.Fatal("contractDetails.Receiver should be the specified receiver")
	}

	if big.NewInt(oneFinney).Cmp(contractDetails.Amount) != 0 {
		t.Fatal("contractDetails.Amount should be equal oneFinney")
	}

	if hexutil.Encode(hashPair.Hash[:]) != hexutil.Encode(contractDetails.Hashlock[:]) {
		t.Fatal("contractDetails.Hashlock should be the specified hashlock")
	}

	if big.NewInt(timeLock1Hour).Cmp(contractDetails.Timelock) != 0 {
		t.Fatal("contractDetails.Timelock should be the specified timelock")
	}

	if (hexutil.Encode(contractDetails.Preimage[:])) != "0x0000000000000000000000000000000000000000000000000000000000000000" {
		t.Fatal("contractDetails.Preimage string should be 0x0000000000000000000000000000000000000000000000000000000000000000")
	}

	if contractDetails.Refunded != false {
		t.Fatal("contractDetails.Refunded should be false")
	}

	if contractDetails.Withdrawn != false {
		t.Fatal("contractDetails.Refunded should be false")
	}

}

//newContract() should fail when no ETH sent
func TestNewContractWithNoETH(t *testing.T) {
	setup()
	defer cleanup()

	var timeLock1Hour = time.Now().Unix() + hourSeconds
	var senderKey   = "a5a1aca01671e2660f1ee47abfd7065d5d38f99fa4a53495f02df939cd5b86f6"
	var receiver = common.HexToAddress("0x5eb1231fd20ac35dff0d4295959cbd11c9cdae40")
	var hashPair = htlc.NewSecretHashPair()

	client, err := ethclient.Dial("http://127.0.0.1:7545")
	if err != nil {
		t.Fatal(err)
	}

	senderAuth := makeAuth(t, senderKey, client, 0)

	_, _, instance, err := DeployHtlc(senderAuth, client)
	if err != nil {
		t.Fatal(err)
	}

	senderAuth = makeAuth(t, senderKey, client, 0)

	_, err = instance.NewContract(senderAuth, receiver, hashPair.Hash, big.NewInt(timeLock1Hour))
	if !(err != nil && strings.HasPrefix(err.Error(),REQUIRE_FAILED_MSG)) {
		t.Fatal("expected failure due to 0 value transferred")
	}
}

//newContract() should fail with timelocks in the past
func TestNewContractWithPastTimelocks(t *testing.T)  {
	setup()
	defer cleanup()

	var timeLockPast = time.Now().Unix() - 1
	var senderKey   = "a5a1aca01671e2660f1ee47abfd7065d5d38f99fa4a53495f02df939cd5b86f6"
	var receiver = common.HexToAddress("0x5eb1231fd20ac35dff0d4295959cbd11c9cdae40")
	var hashPair = htlc.NewSecretHashPair()

	client, err := ethclient.Dial("http://127.0.0.1:7545")
	if err != nil {
		t.Fatal(err)
	}

	senderAuth := makeAuth(t, senderKey, client, 0)

	_, _, instance, err := DeployHtlc(senderAuth, client)
	if err != nil {
		t.Fatal(err)
	}

	senderAuth = makeAuth(t, senderKey, client, oneFinney)

	_, err = instance.NewContract(senderAuth, receiver, hashPair.Hash, big.NewInt(timeLockPast))
	if !(err != nil && strings.HasPrefix(err.Error(),REQUIRE_FAILED_MSG)) {
		t.Fatal("expected failure due past timelock")
	}
}