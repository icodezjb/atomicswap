package htlc

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/icodezjb/atomicswap/contract/helper"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/common"
)

const hourSeconds = 3600

type Ganache struct  {
	cmd *exec.Cmd
	running bool
}

var ganache = Ganache {
	cmd: exec.Command(
		"/usr/local/bin/ganache-cli",
		"--account", "0xa5a1aca01671e2660f1ee47abfd7065d5d38f99fa4a53495f02df939cd5b86f6,111111111111111111111",
		"-p", "7545"),
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

	//Deploy contract
	_, _, instance, err := DeployHtlc(senderAuth, client)
	if err != nil {
		t.Fatal(err)
	}

	senderAuth = makeAuth(t, senderKey, client, 10^18)

	tx, err := instance.NewContract(senderAuth, receiver, hashPair.Hash,big.NewInt(timeLock1Hour))
	if err != nil {
		t.Fatal(err)
	}

	t.Log(tx)
}