package htlc

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"os/exec"
	"testing"
	"time"
)

func init()  {

}

func listen(t *testing.T, cmd *exec.Cmd, runing *bool)  {
	if err:= cmd.Run();err != nil {
		*runing = false
		t.Fatal("unexpect ganache-cli exit: ", err)
	}
}

func TestDeployHtlc(t *testing.T)  {
	running := true
	pathOfGanache := "/usr/local/bin/ganache-cli"
	cmdGanache := exec.Command(
		pathOfGanache,
		"--account", "0xa5a1aca01671e2660f1ee47abfd7065d5d38f99fa4a53495f02df939cd5b86f6,111111111111111111111",
		"-p", "7545")

	defer func() {
		if running {
			cmdGanache.Process.Kill()
		}
	}()

	go listen(t, cmdGanache, &running)

	//waiting for rpc service
	time.Sleep(2*time.Second)

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

	//Deploy contract
	address, tx, _, err := DeployHtlc(auth, client)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("address hex = ", address.Hex())
	fmt.Println("address string =", address.String())
	fmt.Println("tx hash = ", tx.Hash().Hex())
}