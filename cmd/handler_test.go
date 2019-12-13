package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/icodezjb/atomicswap/logger"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	. "github.com/smartystreets/goconvey/convey"
)

func testDeployFunc(ctx context.Context, cfg string) *Handler {
	h := new(Handler)
	h.Config = new(Config)
	h.Config.test = true

	h.ConfigPath = cfg
	h.Config.ParseConfig(h.ConfigPath)

	h.Config.Connect("")
	h.Config.ValidateAddress(h.Config.Account)
	h.Config.Unlock("")
	h.DeployContract(ctx)

	return h
}

func testHashPair() struct {
	Secret      [32]byte
	Hash        [32]byte
	InputSecret string
	InputHash   string
} {
	hashPair := struct {
		Secret      [32]byte
		Hash        [32]byte
		InputSecret string
		InputHash   string
	}{
		InputSecret: "0xe8a1c89faa4d21a522912e5e4eed39c744c8892224c7365972408fa3f26698eb",
		InputHash:   "0x3a9fb66bfb804cc8694c442bf2b18e7a32e4eca03b79c3700d354b4927930105",
	}

	decode, _ := hexutil.Decode("0xe8a1c89faa4d21a522912e5e4eed39c744c8892224c7365972408fa3f26698eb")
	copy(hashPair.Secret[:], decode)
	decode, _ = hexutil.Decode("0x3a9fb66bfb804cc8694c442bf2b18e7a32e4eca03b79c3700d354b4927930105")
	copy(hashPair.Hash[:], decode)

	return hashPair
}

func testGetBalance(t *testing.T, ctx context.Context, client *ethclient.Client, account common.Address) *big.Int {
	balance, err := client.BalanceAt(ctx, account, nil)
	if err != nil {
		t.Fatal("Failed to get balance: ", err)
	}
	return balance
}

func testTransfer(ctx context.Context, h *Handler, account common.Address, value int64) *big.Int {
	auth := h.Config.makeAuth(ctx, value)

	rawTx := types.NewTransaction(auth.Nonce.Uint64(), account, auth.Value, auth.GasLimit, auth.GasPrice, nil)
	txSigned, err := h.Config.ks.SignTxWithPassphrase(
		accounts.Account{Address: common.HexToAddress(h.Config.Account)},
		h.Config.Password,
		rawTx,
		h.Config.Chain.ID)

	if err != nil {
		logger.FatalError("Fatal to sign tx %v: %v", h.Config.Account, err)
	}

	estimateGas, err := h.Config.client.EstimateGas(ctx, ethereum.CallMsg{
		From:     auth.From,
		To:       &account,
		Gas:      0,
		GasPrice: auth.GasPrice,
		Value:    auth.Value,
		Data:     nil,
	})
	if err != nil {
		logger.FatalError("Fatal to estimate gas (transfer): %v", err)
	}
	feeByWei := new(big.Int).Mul(new(big.Int).SetUint64(estimateGas), auth.GasPrice)

	err = h.Config.client.SendTransaction(ctx, txSigned)
	if err != nil {
		logger.FatalError("Fatal to send tx: %v", err)
	}

	return feeByWei
}

func TestHandlerAll_MainFlow(t *testing.T) {
	fmt.Println(`0xae6e5fee5161cede9bc4d89effbbf9944867127d      <===========atomicswap===========>       0x75a8f951632c2e550906f31b53b7923f45be5157
                chain1                                                                                    chain2         
                node1                                                                                     node2
              initiator                                                                                 participant
            100 coin1(wei)                                                                            10000 coin2(wei)`)

	const (
		paySome = 7949900000000000
	)

	var (
		h1 *Handler
		h2 *Handler

		//by wei
		initiatorAmount   int64 = 100
		participantAmount int64 = 10000

		ctx      = context.Background()
		hashPair = testHashPair()
		//set expire time 100 seconds
		timeLockOnChain1 = new(big.Int).SetInt64(time.Now().Unix() + 100)
		timeLockOnChain2 = new(big.Int)

		initiatorLockOnChain1Tx   *types.Transaction
		participantLockOnChain2Tx *types.Transaction
		ContractIDOnChain1        [32]byte
		ContractIDOnChain2        [32]byte
	)

	Convey("Test atomicswap between chain1 and chain2", t, func() {

		Convey("[1] node1 deploy contract address on chain1 should be 0x12D51a18385542d53acC27011aD27E57115b8e0b", func() {
			var b bytes.Buffer
			logger.SetOutput(&b)
			defer logger.SetOutput(os.Stderr)

			h1 = testDeployFunc(ctx, node1Config)

			var expect = "\x1b[92minfo:\x1b[0m Deploy contract...\n" +
				"\x1b[92minfo:\x1b[0m from = 0xAE6e5feE5161CedE9bc4D89eFFBbF9944867127d, balance = 1000000000000000000\n" +
				"\x1b[92minfo:\x1b[0m Deploy Contract fee = gas(1183431) * gasPrice(1000000000) = 1183431000000000\n" +
				"\x1b[92minfo:\x1b[0m ? Confirm to Deploy the contract on node1(chainID = 110)? [y/N]\n" +
				"\x1b[92minfo:\x1b[0m Test chose: y\n" +
				"\x1b[92minfo:\x1b[0m contract address = 0x12D51a18385542d53acC27011aD27E57115b8e0b\n" +
				"\x1b[92minfo:\x1b[0m transaction hash = 0x55733a60d225adfed8c6687a19de28f2ca5812b08570a5a9029fd57d4ccaee42\n"

			So(b.String(), ShouldEqual, expect)

			//reload config
			cfg := new(Config)
			cfg.ParseConfig(node1Config)
			So(cfg.Contract, ShouldEqual, "0x12D51a18385542d53acC27011aD27E57115b8e0b")

			//check initiator balance on chain1
			initiatorBalance := testGetBalance(t, ctx, h1.Config.client, common.HexToAddress(h1.Config.Account))
			expectBalance, _ := new(big.Int).SetString("1000000000000000000", 10)
			So(initiatorBalance.Cmp(expectBalance), ShouldEqual, 0)
		})

		Convey("[2] node2 deploy contract address on chain2 should be 0x071C14E8f6379c4f1d727fDf833024AE9C73C574", func() {
			var b bytes.Buffer
			logger.SetOutput(&b)
			defer logger.SetOutput(os.Stderr)

			h2 = testDeployFunc(ctx, node2Config)

			var expect = "\x1b[92minfo:\x1b[0m Deploy contract...\n" +
				"\x1b[92minfo:\x1b[0m from = 0x75A8f951632C2e550906f31b53b7923F45Be5157, balance = 1000000000000000000\n" +
				"\x1b[92minfo:\x1b[0m Deploy Contract fee = gas(1183431) * gasPrice(1000000000) = 1183431000000000\n" +
				"\x1b[92minfo:\x1b[0m ? Confirm to Deploy the contract on node2(chainID = 111)? [y/N]\n" +
				"\x1b[92minfo:\x1b[0m Test chose: y\n" +
				"\x1b[92minfo:\x1b[0m contract address = 0x071C14E8f6379c4f1d727fDf833024AE9C73C574\n" +
				"\x1b[92minfo:\x1b[0m transaction hash = 0x8a4596851e01df0dfa7f2053399c8bca7e5a0a965651acd5f401625e973efc7f\n"

			So(b.String(), ShouldEqual, expect)

			//reload config
			cfg := new(Config)
			cfg.ParseConfig(node2Config)
			So(cfg.Contract, ShouldEqual, "0x071C14E8f6379c4f1d727fDf833024AE9C73C574")

			//check participant balance on chain2
			participantBalance := testGetBalance(t, ctx, h2.Config.client, common.HexToAddress(h2.Config.Account))
			expectBalance, _ := new(big.Int).SetString("1000000000000000000", 10)
			So(participantBalance.Cmp(expectBalance), ShouldEqual, 0)
		})

		Convey("[3] initiator lock 100 coin1 on chain1 0x12D51a18385542d53acC27011aD27E57115b8e0b \n"+
			"    with hashlock=0x3a9fb66bfb804cc8694c442bf2b18e7a32e4eca03b79c3700d354b4927930105 and timeLockOnChain1", func() {
			var b bytes.Buffer
			logger.SetOutput(&b)
			defer logger.SetOutput(os.Stderr)

			var expect = "\x1b[92minfo:\x1b[0m Call NewContract ...\n" +
				"\x1b[92minfo:\x1b[0m from = 0xAE6e5feE5161CedE9bc4D89eFFBbF9944867127d, balance = 1000000000000000000\n" +
				"\x1b[92minfo:\x1b[0m Call Contract fee = gas(147040) * gasPrice(1000000000) = 147040000000000\n" +
				"\x1b[92minfo:\x1b[0m ? Confirm to Call the contract on node1(chainID = 110)? [y/N]\n" +
				"\x1b[92minfo:\x1b[0m Test chose: y\n"

			initiatorLockOnChain1Tx = h1.NewContract(ctx, common.HexToAddress(h2.Config.Account), initiatorAmount, hashPair.Hash, timeLockOnChain1)

			So(b.String(), ShouldEqual, expect)

			time.Sleep(2 * time.Second)
			//check initiator balance on chain1
			initiatorBalance := testGetBalance(t, ctx, h1.Config.client, common.HexToAddress(h1.Config.Account))
			expectBalance, _ := new(big.Int).SetString("999999999999999900", 10)
			So(initiatorBalance.Cmp(expectBalance), ShouldEqual, 0)
		})

		Convey("[4] initiator query the lock tx on chain1 and get the ContractIDOnChain1", func() {
			logger.SetOutput(ioutil.Discard)
			defer logger.SetOutput(os.Stderr)

			logHTLCEvent := h1.GetContractId(ctx, initiatorLockOnChain1Tx.Hash())

			So(h1.Config.Account, ShouldEqual, strings.ToLower(logHTLCEvent.Sender.String()))
			So(h2.Config.Account, ShouldEqual, strings.ToLower(logHTLCEvent.Receiver.String()))
			So(initiatorAmount, ShouldEqual, logHTLCEvent.Amount.Int64())
			So(timeLockOnChain1.Int64(), ShouldEqual, logHTLCEvent.Timelock.Int64())
			So(hashPair.Hash, ShouldEqual, logHTLCEvent.Hashlock)

			copy(ContractIDOnChain1[:], logHTLCEvent.ContractId[:])
		})

		Convey("[5] participant audit the contract on chain1 by the ContractIDOnChain1", func() {
			//here used h1 client for test
			//refer: the '--other' flag of auditcontract cmd
			tmpClient := h2.Config.client
			tmpChain := h2.Config.Chain
			h2.Config.client = h1.Config.client
			h2.Config.Chain = h1.Config.Chain
			defer func() {
				h2.Config.client = tmpClient
				h2.Config.Chain = tmpChain
			}()

			contractDetails := new(ContractDetails)
			h2.AuditContract(ctx, contractDetails, "getContract", ContractIDOnChain1)

			So(h1.Config.Account, ShouldEqual, strings.ToLower(contractDetails.Sender.String()))
			So(h2.Config.Account, ShouldEqual, strings.ToLower(contractDetails.Receiver.String()))
			So(initiatorAmount, ShouldEqual, contractDetails.Amount.Int64())
			So(timeLockOnChain1.Int64(), ShouldEqual, contractDetails.Timelock.Int64())
			So(hashPair.Hash, ShouldEqual, contractDetails.Hashlock)
			So(contractDetails.Withdrawn, ShouldBeFalse)
			So(contractDetails.Refunded, ShouldBeFalse)
			So("0x0000000000000000000000000000000000000000000000000000000000000000", ShouldEqual, hexutil.Encode(contractDetails.Preimage[:]))

		})

		Convey("[6] participant lock 10000 coin2 on chain2 0x071C14E8f6379c4f1d727fDf833024AE9C73C574 \n"+
			"    with hashlock=0x3a9fb66bfb804cc8694c442bf2b18e7a32e4eca03b79c3700d354b4927930105 and timeLockOnChain2", func() {
			var b bytes.Buffer
			logger.SetOutput(&b)
			defer logger.SetOutput(os.Stderr)

			var expect = "\x1b[92minfo:\x1b[0m Call NewContract ...\n" +
				"\x1b[92minfo:\x1b[0m from = 0x75A8f951632C2e550906f31b53b7923F45Be5157, balance = 1000000000000000000\n" +
				"\x1b[92minfo:\x1b[0m Call Contract fee = gas(147040) * gasPrice(1000000000) = 147040000000000\n" +
				"\x1b[92minfo:\x1b[0m ? Confirm to Call the contract on node2(chainID = 111)? [y/N]\n" +
				"\x1b[92minfo:\x1b[0m Test chose: y\n"

			timeLockOnChain2 = new(big.Int).SetInt64(time.Now().Unix() + (timeLockOnChain1.Int64()-time.Now().Unix())/2)
			participantLockOnChain2Tx = h2.NewContract(ctx, common.HexToAddress(h1.Config.Account), participantAmount, hashPair.Hash, timeLockOnChain2)

			So(b.String(), ShouldEqual, expect)

			time.Sleep(2 * time.Second)
			//check participant balance on chain2
			participantBalance := testGetBalance(t, ctx, h2.Config.client, common.HexToAddress(h2.Config.Account))
			expectBalance, _ := new(big.Int).SetString("999999999999990000", 10)
			So(participantBalance.Cmp(expectBalance), ShouldEqual, 0)
		})

		Convey("[7] participant query the lock tx on chain1 and get the ContractIDOnChain2", func() {
			logger.SetOutput(ioutil.Discard)
			defer logger.SetOutput(os.Stderr)

			logHTLCEvent := h2.GetContractId(ctx, participantLockOnChain2Tx.Hash())

			So(h2.Config.Account, ShouldEqual, strings.ToLower(logHTLCEvent.Sender.String()))
			So(h1.Config.Account, ShouldEqual, strings.ToLower(logHTLCEvent.Receiver.String()))
			So(participantAmount, ShouldEqual, logHTLCEvent.Amount.Int64())
			So(timeLockOnChain2.Int64(), ShouldEqual, logHTLCEvent.Timelock.Int64())
			So(hashPair.Hash, ShouldEqual, logHTLCEvent.Hashlock)

			copy(ContractIDOnChain2[:], logHTLCEvent.ContractId[:])
		})

		Convey("[8] initiator audit the contract on chain2 by the ContractIDOnChain2", func() {
			//here used h2 client for test
			//refer: the '--other' flag of auditcontract cmd
			tmpClient := h1.Config.client
			tmpChain := h1.Config.Chain
			h1.Config.client = h2.Config.client
			h1.Config.Chain = h2.Config.Chain
			defer func() {
				h1.Config.client = tmpClient
				h1.Config.Chain = tmpChain
			}()

			contractDetails := new(ContractDetails)
			h1.AuditContract(ctx, contractDetails, "getContract", ContractIDOnChain2)

			So(h2.Config.Account, ShouldEqual, strings.ToLower(contractDetails.Sender.String()))
			So(h1.Config.Account, ShouldEqual, strings.ToLower(contractDetails.Receiver.String()))
			So(participantAmount, ShouldEqual, contractDetails.Amount.Int64())
			So(timeLockOnChain2.Int64(), ShouldEqual, contractDetails.Timelock.Int64())
			So(hashPair.Hash, ShouldEqual, contractDetails.Hashlock)
			So(contractDetails.Withdrawn, ShouldBeFalse)
			So(contractDetails.Refunded, ShouldBeFalse)
			So("0x0000000000000000000000000000000000000000000000000000000000000000", ShouldEqual, hexutil.Encode(contractDetails.Preimage[:]))

		})

		Convey("[only for this test] check balance, if not enough, pay some to the other for redeem fee", func() {
			Convey("initiator balance on chain1 should be 999999999999999900", func() {
				initiatorBalance := testGetBalance(t, ctx, h1.Config.client, common.HexToAddress(h1.Config.Account))
				expectBalance, _ := new(big.Int).SetString("999999999999999900", 10)
				So(initiatorBalance.Cmp(expectBalance), ShouldEqual, 0)
			})

			Convey("participant balance on chain1 should be 0", func() {
				participantBalance := testGetBalance(t, ctx, h1.Config.client, common.HexToAddress(h2.Config.Account))
				expectBalance := new(big.Int)
				So(participantBalance.Cmp(expectBalance), ShouldEqual, 0)
			})

			Convey("initiator balance on chain2 should be 0", func() {
				initiatorBalance := testGetBalance(t, ctx, h2.Config.client, common.HexToAddress(h1.Config.Account))
				expectBalance := new(big.Int)
				So(initiatorBalance.Cmp(expectBalance), ShouldEqual, 0)
			})

			Convey("participant balance on chain2 should be 999999999999990000", func() {
				participantBalance := testGetBalance(t, ctx, h2.Config.client, common.HexToAddress(h2.Config.Account))
				expectBalance, _ := new(big.Int).SetString("999999999999990000", 10)
				So(participantBalance.Cmp(expectBalance), ShouldEqual, 0)
			})

			Convey("initiator transfer paySome=7949900000000000 wei to participant on chain1", func() {
				wastage := testTransfer(ctx, h1, common.HexToAddress(h2.Config.Account), paySome)
				time.Sleep(2 * time.Second)

				initiatorBalance := testGetBalance(t, ctx, h1.Config.client, common.HexToAddress(h1.Config.Account))
				participantBalance := testGetBalance(t, ctx, h1.Config.client, common.HexToAddress(h2.Config.Account))

				expectWastage, _ := new(big.Int).SetString("21000000000000", 10)
				expectBalance1, _ := new(big.Int).SetString("992050099999999900", 10)
				expectBalance2 := new(big.Int).SetInt64(paySome)

				So(wastage.String(), ShouldEqual, expectWastage.String())
				So(initiatorBalance.String(), ShouldEqual, expectBalance1.String())
				So(participantBalance.String(), ShouldEqual, expectBalance2.String())
			})

			Convey("participant transfer paySome=7949900000000000 wei to initiator on chain2", func() {
				wastage := testTransfer(ctx, h2, common.HexToAddress(h1.Config.Account), paySome)
				time.Sleep(2 * time.Second)

				initiatorBalance := testGetBalance(t, ctx, h2.Config.client, common.HexToAddress(h1.Config.Account))
				participantBalance := testGetBalance(t, ctx, h2.Config.client, common.HexToAddress(h2.Config.Account))

				expectWastage, _ := new(big.Int).SetString("21000000000000", 10)
				expectBalance1, _ := new(big.Int).SetString("992050099999990000", 10)
				expectBalance2 := new(big.Int).SetInt64(paySome)

				So(wastage.String(), ShouldEqual, expectWastage.String())
				So(participantBalance.String(), ShouldEqual, expectBalance1.String())
				So(initiatorBalance.String(), ShouldEqual, expectBalance2.String())
			})

		})

		Convey("[9] initiator redeem on chain2 0x071C14E8f6379c4f1d727fDf833024AE9C73C574 with the preimage or secret of hashlock", func() {
			logger.SetOutput(ioutil.Discard)
			defer logger.SetOutput(os.Stderr)

			//here used h2 client for test
			//refer: the '--other' flag of auditcontract cmd
			tmpClient := h1.Config.client
			tmpChain := h1.Config.Chain
			h1.Config.client = h2.Config.client
			h1.Config.Chain = h2.Config.Chain
			defer func() {
				h1.Config.client = tmpClient
				h1.Config.Chain = tmpChain
			}()

			h1.Redeem(ctx, ContractIDOnChain2, hashPair.Secret)

			Convey("after redeem, initiator balance should be somePay-redeemFee+participantAmount ~ 7871412000010000", func() {
				time.Sleep(2 * time.Second)
				initiatorBalanceAfterSwap := testGetBalance(t, ctx, h2.Config.client, common.HexToAddress(h1.Config.Account))
				//expectBalance=somePay-redeemFee+participantAmount
				expectBalance := fmt.Sprintf("%v", participantAmount)

				So(initiatorBalanceAfterSwap.String(), ShouldEndWith, expectBalance)
			})
		})

		Convey("[10] participant audit the contract on chain2 0x071C14E8f6379c4f1d727fDf833024AE9C73C574 to get the preimage or secret of hashlock", func() {
			contractDetails := new(ContractDetails)
			time.Sleep(2 * time.Second)
			h2.AuditContract(ctx, contractDetails, "getContract", ContractIDOnChain2)

			So(h2.Config.Account, ShouldEqual, strings.ToLower(contractDetails.Sender.String()))
			So(h1.Config.Account, ShouldEqual, strings.ToLower(contractDetails.Receiver.String()))
			So(participantAmount, ShouldEqual, contractDetails.Amount.Int64())
			So(timeLockOnChain2.Int64(), ShouldEqual, contractDetails.Timelock.Int64())
			So(hashPair.Hash, ShouldEqual, contractDetails.Hashlock)
			So(contractDetails.Withdrawn, ShouldBeTrue)
			So(contractDetails.Refunded, ShouldBeFalse)
			So(hashPair.InputSecret, ShouldEqual, hexutil.Encode(contractDetails.Preimage[:]))
		})

		Convey("[11] participant redeem on chain1 0x12D51a18385542d53acC27011aD27E57115b8e0b with the preimage or secret of hashlock", func() {
			logger.SetOutput(ioutil.Discard)
			defer logger.SetOutput(os.Stderr)
			//here used h1 client for test
			//refer: the '--other' flag of auditcontract cmd
			tmpClient := h2.Config.client
			tmpChain := h2.Config.Chain
			h2.Config.client = h1.Config.client
			h2.Config.Chain = h1.Config.Chain
			defer func() {
				h2.Config.client = tmpClient
				h2.Config.Chain = tmpChain
			}()

			h2.Redeem(ctx, ContractIDOnChain1, hashPair.Secret)

			Convey("after redeem, participant balance should be expectBalance=somePay-redeemFee+initiatorAmount ~ 7871412000000100", func() {
				time.Sleep(2 * time.Second)
				participantBalanceAfterSwap := testGetBalance(t, ctx, h1.Config.client, common.HexToAddress(h2.Config.Account))
				//expectBalance=somePay-redeemFee+initiatorAmount
				expectBalance := fmt.Sprintf("%v", initiatorAmount)

				So(participantBalanceAfterSwap.String(), ShouldEndWith, expectBalance)
			})
		})

		Convey("[12] initiator audit the contract on chain1 0x12D51a18385542d53acC27011aD27E57115b8e0b to complete atomic swap", func() {
			contractDetails := new(ContractDetails)
			time.Sleep(2 * time.Second)
			h1.AuditContract(ctx, contractDetails, "getContract", ContractIDOnChain1)

			So(h1.Config.Account, ShouldEqual, strings.ToLower(contractDetails.Sender.String()))
			So(h2.Config.Account, ShouldEqual, strings.ToLower(contractDetails.Receiver.String()))
			So(initiatorAmount, ShouldEqual, contractDetails.Amount.Int64())
			So(timeLockOnChain1.Int64(), ShouldEqual, contractDetails.Timelock.Int64())
			So(hashPair.Hash, ShouldEqual, contractDetails.Hashlock)
			So(contractDetails.Withdrawn, ShouldBeTrue)
			So(contractDetails.Refunded, ShouldBeFalse)
			So(hashPair.InputSecret, ShouldEqual, hexutil.Encode(contractDetails.Preimage[:]))
		})

		Convey("[13] initiator refund on chain1 0x12D51a18385542d53acC27011aD27E57115b8e0b should be fail", func() {
			var b bytes.Buffer
			logger.SetOutput(&b)
			defer logger.SetOutput(os.Stderr)
			h1.Refund(ctx, ContractIDOnChain1)

			So(b.String(), ShouldContainSubstring, "Fatal to estimate gas (Call): gas required exceeds allowance")
		})

		Convey("[14] participant refund on chain2 0x071C14E8f6379c4f1d727fDf833024AE9C73C574 should be fail", func() {
			var b bytes.Buffer
			logger.SetOutput(&b)
			defer logger.SetOutput(os.Stderr)
			h2.Refund(ctx, ContractIDOnChain2)

			So(b.String(), ShouldContainSubstring, "Fatal to estimate gas (Call): gas required exceeds allowance")

		})
	})
}
