package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	htlc "github.com/icodezjb/atomicswap/contract"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("http://127.0.0.1:7545")
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("b80dbf638b9128e54f3222d2b6d3213d45521d49bb6317abdf34b219a55204b7")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce,err := client.PendingNonceAt(context.Background(),fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) //in wei
	auth.GasLimit = uint64(3000000) //in uints
	auth.GasPrice = gasPrice

	address, tx, instance, err := htlc.DeployHtlc(auth, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("address hex = ", address.Hex())
	fmt.Println("address string =", address.String())
	fmt.Println("tx hash = ", tx.Hash().Hex())

	_ = instance
}
