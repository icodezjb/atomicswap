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

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

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

func makeAuthandClient(t *testing.T)(*bind.TransactOpts, bind.ContractBackend)  {
	privateKey, err := crypto.HexToECDSA("a5a1aca01671e2660f1ee47abfd7065d5d38f99fa4a53495f02df939cd5b86f6")
	if err != nil {
		t.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	client, err := ethclient.Dial("http://127.0.0.1:7545")
	if err != nil {
		t.Fatal(err)
	}
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
	auth.Value = big.NewInt(0) //in wei
	auth.GasLimit = uint64(3000000) //in uints
	auth.GasPrice = gasPrice

	return auth, client
}

func TestDeployHtlc(t *testing.T)  {
	setup()
	defer cleanup()

	auth, client := makeAuthandClient(t)

	//Deploy contract
	_, _, _, err := DeployHtlc(auth, client)
	if err != nil {
		t.Fatal(err)
	}
}