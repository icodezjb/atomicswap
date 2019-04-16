// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package htlc

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// HtlcABI is the input ABI used to generate the binding from.
const HtlcABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_receiver\",\"type\":\"address\"},{\"name\":\"_hashlock\",\"type\":\"bytes32\"},{\"name\":\"_timelock\",\"type\":\"uint256\"}],\"name\":\"newContract\",\"outputs\":[{\"name\":\"contractId\",\"type\":\"bytes32\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_contractId\",\"type\":\"bytes32\"},{\"name\":\"_preimage\",\"type\":\"bytes32\"}],\"name\":\"withdraw\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_contractId\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_contractId\",\"type\":\"bytes32\"}],\"name\":\"getContract\",\"outputs\":[{\"name\":\"sender\",\"type\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"hashlock\",\"type\":\"bytes32\"},{\"name\":\"timelock\",\"type\":\"uint256\"},{\"name\":\"withdrawn\",\"type\":\"bool\"},{\"name\":\"refunded\",\"type\":\"bool\"},{\"name\":\"preimage\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"contractId\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"hashlock\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"timelock\",\"type\":\"uint256\"}],\"name\":\"LogHTLCNew\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"contractId\",\"type\":\"bytes32\"}],\"name\":\"LogHTLCWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"contractId\",\"type\":\"bytes32\"}],\"name\":\"LogHTLCRefund\",\"type\":\"event\"}]"

// HtlcBin is the compiled bytecode used for deploying new contracts.
const HtlcBin = `608060405234801561001057600080fd5b506110e1806100206000396000f3fe60806040526004361061003f5760003560e01c8063335ef5bd1461004457806363615149146100b05780637249fbb61461010d578063e16c7d9814610160575b600080fd5b61009a6004803603606081101561005a57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919080359060200190929190505050610240565b6040518082815260200191505060405180910390f35b3480156100bc57600080fd5b506100f3600480360360408110156100d357600080fd5b810190808035906020019092919080359060200190929190505050610658565b604051808215151515815260200191505060405180910390f35b34801561011957600080fd5b506101466004803603602081101561013057600080fd5b8101908080359060200190929190505050610ae5565b604051808215151515815260200191505060405180910390f35b34801561016c57600080fd5b506101996004803603602081101561018357600080fd5b8101908080359060200190929190505050610ebb565b604051808973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200187815260200186815260200185815260200184151515158152602001831515151581526020018281526020019850505050505050505060405180910390f35b60008034116102b7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f6d73672e76616c7565206d757374206265203e2030000000000000000000000081525060200191505060405180910390fd5b81428111610310576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602381526020018061103f6023913960400191505060405180910390fd5b60023386348787604051602001808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1660601b81526014018573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1660601b8152601401848152602001838152602001828152602001955050505050506040516020818303038152906040526040518082805190602001908083835b602083106103e357805182526020820191506020810190506020830392506103c0565b6001836020036101000a038019825116818451168082178552505050505050905001915050602060405180830381855afa158015610425573d6000803e3d6000fd5b5050506040513d602081101561043a57600080fd5b8101908080519060200190929190505050915061045682610fd0565b1561046057600080fd5b6040518061010001604052803373ffffffffffffffffffffffffffffffffffffffff1681526020018673ffffffffffffffffffffffffffffffffffffffff1681526020013481526020018581526020018481526020016000151581526020016000151581526020016000801b81525060008084815260200190815260200160002060008201518160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060208201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060408201518160020155606082015181600301556080820151816004015560a08201518160050160006101000a81548160ff02191690831515021790555060c08201518160050160016101000a81548160ff02191690831515021790555060e082015181600601559050508473ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16837f329a8316ed9c3b2299597538371c2944c5026574e803b1ec31d6113e1cd67bde34888860405180848152602001838152602001828152602001935050505060405180910390a4509392505050565b60008261066481610fd0565b6106d6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260198152602001807f636f6e7472616374496420646f6573206e6f742065786973740000000000000081525060200191505060405180910390fd5b8383600281604051602001808281526020019150506040516020818303038152906040526040518082805190602001908083835b6020831061072d578051825260208201915060208101905060208303925061070a565b6001836020036101000a038019825116818451168082178552505050505050905001915050602060405180830381855afa15801561076f573d6000803e3d6000fd5b5050506040513d602081101561078457600080fd5b8101908080519060200190929190505050600080848152602001908152602001600020600301541461081e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601c8152602001807f686173686c6f636b206861736820646f6573206e6f74206d617463680000000081525060200191505060405180910390fd5b853373ffffffffffffffffffffffffffffffffffffffff1660008083815260200190815260200160002060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16146108f5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f776974686472617761626c653a206e6f7420726563656976657200000000000081525060200191505060405180910390fd5b6000151560008083815260200190815260200160002060050160009054906101000a900460ff16151514610991576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601f8152602001807f776974686472617761626c653a20616c72656164792077697468647261776e0081525060200191505060405180910390fd5b4260008083815260200190815260200160002060040154116109fe576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260318152602001806110856031913960400191505060405180910390fd5b6000806000898152602001908152602001600020905086816006018190555060018160050160006101000a81548160ff0219169083151502179055508060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc82600201549081150290604051600060405180830381858888f19350505050158015610aa8573d6000803e3d6000fd5b50877fd6fd4c8e45bf0c70693141c7ce46451b6a6a28ac8386fca2ba914044e0e2391660405160405180910390a260019550505050505092915050565b600081610af181610fd0565b610b63576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260198152602001807f636f6e7472616374496420646f6573206e6f742065786973740000000000000081525060200191505060405180910390fd5b823373ffffffffffffffffffffffffffffffffffffffff1660008083815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614610c3a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f726566756e6461626c653a206e6f742073656e6465720000000000000000000081525060200191505060405180910390fd5b6000151560008083815260200190815260200160002060050160019054906101000a900460ff16151514610cd6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601c8152602001807f726566756e6461626c653a20616c726561647920726566756e6465640000000081525060200191505060405180910390fd5b6000151560008083815260200190815260200160002060050160009054906101000a900460ff16151514610d72576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601d8152602001807f726566756e6461626c653a20616c72656164792077697468647261776e00000081525060200191505060405180910390fd5b42600080838152602001908152602001600020600401541115610de0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260238152602001806110626023913960400191505060405180910390fd5b6000806000868152602001908152602001600020905060018160050160016101000a81548160ff0219169083151502179055508060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc82600201549081150290604051600060405180830381858888f19350505050158015610e81573d6000803e3d6000fd5b50847f989b3a845197c9aec15f8982bbb30b5da714050e662a7a287bb1a94c81e2e70e60405160405180910390a260019350505050919050565b60008060008060008060008060001515610ed48a610fd0565b15151415610f15576000806000806000806000808797508696508595508460001b94508393508060001b905097509750975097509750975097509750610fc5565b60008060008b815260200190815260200160002090508060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff168160010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff168260020154836003015484600401548560050160009054906101000a900460ff168660050160019054906101000a900460ff16876006015487975086965098509850985098509850985098509850505b919395975091939597565b60008073ffffffffffffffffffffffffffffffffffffffff1660008084815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161415905091905056fe74696d656c6f636b2074696d65206d75737420626520696e2074686520667574757265726566756e6461626c653a2074696d656c6f636b206e6f742079657420706173736564776974686472617761626c653a2074696d656c6f636b2074696d65206d75737420626520696e2074686520667574757265a165627a7a72305820c0c206dd0a8e7ec5f56c009b104c9e6bc4478e82b75b3d6217b792bc1bf73a860029`

// DeployHtlc deploys a new Ethereum contract, binding an instance of Htlc to it.
func DeployHtlc(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Htlc, error) {
	parsed, err := abi.JSON(strings.NewReader(HtlcABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(HtlcBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Htlc{HtlcCaller: HtlcCaller{contract: contract}, HtlcTransactor: HtlcTransactor{contract: contract}, HtlcFilterer: HtlcFilterer{contract: contract}}, nil
}

// Htlc is an auto generated Go binding around an Ethereum contract.
type Htlc struct {
	HtlcCaller     // Read-only binding to the contract
	HtlcTransactor // Write-only binding to the contract
	HtlcFilterer   // Log filterer for contract events
}

// HtlcCaller is an auto generated read-only Go binding around an Ethereum contract.
type HtlcCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HtlcTransactor is an auto generated write-only Go binding around an Ethereum contract.
type HtlcTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HtlcFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type HtlcFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HtlcSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type HtlcSession struct {
	Contract     *Htlc             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HtlcCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type HtlcCallerSession struct {
	Contract *HtlcCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// HtlcTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type HtlcTransactorSession struct {
	Contract     *HtlcTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HtlcRaw is an auto generated low-level Go binding around an Ethereum contract.
type HtlcRaw struct {
	Contract *Htlc // Generic contract binding to access the raw methods on
}

// HtlcCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type HtlcCallerRaw struct {
	Contract *HtlcCaller // Generic read-only contract binding to access the raw methods on
}

// HtlcTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type HtlcTransactorRaw struct {
	Contract *HtlcTransactor // Generic write-only contract binding to access the raw methods on
}

// NewHtlc creates a new instance of Htlc, bound to a specific deployed contract.
func NewHtlc(address common.Address, backend bind.ContractBackend) (*Htlc, error) {
	contract, err := bindHtlc(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Htlc{HtlcCaller: HtlcCaller{contract: contract}, HtlcTransactor: HtlcTransactor{contract: contract}, HtlcFilterer: HtlcFilterer{contract: contract}}, nil
}

// NewHtlcCaller creates a new read-only instance of Htlc, bound to a specific deployed contract.
func NewHtlcCaller(address common.Address, caller bind.ContractCaller) (*HtlcCaller, error) {
	contract, err := bindHtlc(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HtlcCaller{contract: contract}, nil
}

// NewHtlcTransactor creates a new write-only instance of Htlc, bound to a specific deployed contract.
func NewHtlcTransactor(address common.Address, transactor bind.ContractTransactor) (*HtlcTransactor, error) {
	contract, err := bindHtlc(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HtlcTransactor{contract: contract}, nil
}

// NewHtlcFilterer creates a new log filterer instance of Htlc, bound to a specific deployed contract.
func NewHtlcFilterer(address common.Address, filterer bind.ContractFilterer) (*HtlcFilterer, error) {
	contract, err := bindHtlc(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HtlcFilterer{contract: contract}, nil
}

// bindHtlc binds a generic wrapper to an already deployed contract.
func bindHtlc(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(HtlcABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Htlc *HtlcRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Htlc.Contract.HtlcCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Htlc *HtlcRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Htlc.Contract.HtlcTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Htlc *HtlcRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Htlc.Contract.HtlcTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Htlc *HtlcCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Htlc.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Htlc *HtlcTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Htlc.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Htlc *HtlcTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Htlc.Contract.contract.Transact(opts, method, params...)
}

// GetContract is a free data retrieval call binding the contract method 0xe16c7d98.
//
// Solidity: function getContract(bytes32 _contractId) constant returns(address sender, address receiver, uint256 amount, bytes32 hashlock, uint256 timelock, bool withdrawn, bool refunded, bytes32 preimage)
func (_Htlc *HtlcCaller) GetContract(opts *bind.CallOpts, _contractId [32]byte) (struct {
	Sender    common.Address
	Receiver  common.Address
	Amount    *big.Int
	Hashlock  [32]byte
	Timelock  *big.Int
	Withdrawn bool
	Refunded  bool
	Preimage  [32]byte
}, error) {
	ret := new(struct {
		Sender    common.Address
		Receiver  common.Address
		Amount    *big.Int
		Hashlock  [32]byte
		Timelock  *big.Int
		Withdrawn bool
		Refunded  bool
		Preimage  [32]byte
	})
	out := ret
	err := _Htlc.contract.Call(opts, out, "getContract", _contractId)
	return *ret, err
}

// GetContract is a free data retrieval call binding the contract method 0xe16c7d98.
//
// Solidity: function getContract(bytes32 _contractId) constant returns(address sender, address receiver, uint256 amount, bytes32 hashlock, uint256 timelock, bool withdrawn, bool refunded, bytes32 preimage)
func (_Htlc *HtlcSession) GetContract(_contractId [32]byte) (struct {
	Sender    common.Address
	Receiver  common.Address
	Amount    *big.Int
	Hashlock  [32]byte
	Timelock  *big.Int
	Withdrawn bool
	Refunded  bool
	Preimage  [32]byte
}, error) {
	return _Htlc.Contract.GetContract(&_Htlc.CallOpts, _contractId)
}

// GetContract is a free data retrieval call binding the contract method 0xe16c7d98.
//
// Solidity: function getContract(bytes32 _contractId) constant returns(address sender, address receiver, uint256 amount, bytes32 hashlock, uint256 timelock, bool withdrawn, bool refunded, bytes32 preimage)
func (_Htlc *HtlcCallerSession) GetContract(_contractId [32]byte) (struct {
	Sender    common.Address
	Receiver  common.Address
	Amount    *big.Int
	Hashlock  [32]byte
	Timelock  *big.Int
	Withdrawn bool
	Refunded  bool
	Preimage  [32]byte
}, error) {
	return _Htlc.Contract.GetContract(&_Htlc.CallOpts, _contractId)
}

// NewContract is a paid mutator transaction binding the contract method 0x335ef5bd.
//
// Solidity: function newContract(address _receiver, bytes32 _hashlock, uint256 _timelock) returns(bytes32 contractId)
func (_Htlc *HtlcTransactor) NewContract(opts *bind.TransactOpts, _receiver common.Address, _hashlock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _Htlc.contract.Transact(opts, "newContract", _receiver, _hashlock, _timelock)
}

// NewContract is a paid mutator transaction binding the contract method 0x335ef5bd.
//
// Solidity: function newContract(address _receiver, bytes32 _hashlock, uint256 _timelock) returns(bytes32 contractId)
func (_Htlc *HtlcSession) NewContract(_receiver common.Address, _hashlock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _Htlc.Contract.NewContract(&_Htlc.TransactOpts, _receiver, _hashlock, _timelock)
}

// NewContract is a paid mutator transaction binding the contract method 0x335ef5bd.
//
// Solidity: function newContract(address _receiver, bytes32 _hashlock, uint256 _timelock) returns(bytes32 contractId)
func (_Htlc *HtlcTransactorSession) NewContract(_receiver common.Address, _hashlock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _Htlc.Contract.NewContract(&_Htlc.TransactOpts, _receiver, _hashlock, _timelock)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 _contractId) returns(bool)
func (_Htlc *HtlcTransactor) Refund(opts *bind.TransactOpts, _contractId [32]byte) (*types.Transaction, error) {
	return _Htlc.contract.Transact(opts, "refund", _contractId)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 _contractId) returns(bool)
func (_Htlc *HtlcSession) Refund(_contractId [32]byte) (*types.Transaction, error) {
	return _Htlc.Contract.Refund(&_Htlc.TransactOpts, _contractId)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 _contractId) returns(bool)
func (_Htlc *HtlcTransactorSession) Refund(_contractId [32]byte) (*types.Transaction, error) {
	return _Htlc.Contract.Refund(&_Htlc.TransactOpts, _contractId)
}

// Withdraw is a paid mutator transaction binding the contract method 0x63615149.
//
// Solidity: function withdraw(bytes32 _contractId, bytes32 _preimage) returns(bool)
func (_Htlc *HtlcTransactor) Withdraw(opts *bind.TransactOpts, _contractId [32]byte, _preimage [32]byte) (*types.Transaction, error) {
	return _Htlc.contract.Transact(opts, "withdraw", _contractId, _preimage)
}

// Withdraw is a paid mutator transaction binding the contract method 0x63615149.
//
// Solidity: function withdraw(bytes32 _contractId, bytes32 _preimage) returns(bool)
func (_Htlc *HtlcSession) Withdraw(_contractId [32]byte, _preimage [32]byte) (*types.Transaction, error) {
	return _Htlc.Contract.Withdraw(&_Htlc.TransactOpts, _contractId, _preimage)
}

// Withdraw is a paid mutator transaction binding the contract method 0x63615149.
//
// Solidity: function withdraw(bytes32 _contractId, bytes32 _preimage) returns(bool)
func (_Htlc *HtlcTransactorSession) Withdraw(_contractId [32]byte, _preimage [32]byte) (*types.Transaction, error) {
	return _Htlc.Contract.Withdraw(&_Htlc.TransactOpts, _contractId, _preimage)
}

// HtlcLogHTLCNewIterator is returned from FilterLogHTLCNew and is used to iterate over the raw logs and unpacked data for LogHTLCNew events raised by the Htlc contract.
type HtlcLogHTLCNewIterator struct {
	Event *HtlcLogHTLCNew // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *HtlcLogHTLCNewIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HtlcLogHTLCNew)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(HtlcLogHTLCNew)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *HtlcLogHTLCNewIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HtlcLogHTLCNewIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HtlcLogHTLCNew represents a LogHTLCNew event raised by the Htlc contract.
type HtlcLogHTLCNew struct {
	ContractId [32]byte
	Sender     common.Address
	Receiver   common.Address
	Amount     *big.Int
	Hashlock   [32]byte
	Timelock   *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogHTLCNew is a free log retrieval operation binding the contract event 0x329a8316ed9c3b2299597538371c2944c5026574e803b1ec31d6113e1cd67bde.
//
// Solidity: event LogHTLCNew(bytes32 indexed contractId, address indexed sender, address indexed receiver, uint256 amount, bytes32 hashlock, uint256 timelock)
func (_Htlc *HtlcFilterer) FilterLogHTLCNew(opts *bind.FilterOpts, contractId [][32]byte, sender []common.Address, receiver []common.Address) (*HtlcLogHTLCNewIterator, error) {

	var contractIdRule []interface{}
	for _, contractIdItem := range contractId {
		contractIdRule = append(contractIdRule, contractIdItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _Htlc.contract.FilterLogs(opts, "LogHTLCNew", contractIdRule, senderRule, receiverRule)
	if err != nil {
		return nil, err
	}
	return &HtlcLogHTLCNewIterator{contract: _Htlc.contract, event: "LogHTLCNew", logs: logs, sub: sub}, nil
}

// WatchLogHTLCNew is a free log subscription operation binding the contract event 0x329a8316ed9c3b2299597538371c2944c5026574e803b1ec31d6113e1cd67bde.
//
// Solidity: event LogHTLCNew(bytes32 indexed contractId, address indexed sender, address indexed receiver, uint256 amount, bytes32 hashlock, uint256 timelock)
func (_Htlc *HtlcFilterer) WatchLogHTLCNew(opts *bind.WatchOpts, sink chan<- *HtlcLogHTLCNew, contractId [][32]byte, sender []common.Address, receiver []common.Address) (event.Subscription, error) {

	var contractIdRule []interface{}
	for _, contractIdItem := range contractId {
		contractIdRule = append(contractIdRule, contractIdItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _Htlc.contract.WatchLogs(opts, "LogHTLCNew", contractIdRule, senderRule, receiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HtlcLogHTLCNew)
				if err := _Htlc.contract.UnpackLog(event, "LogHTLCNew", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// HtlcLogHTLCRefundIterator is returned from FilterLogHTLCRefund and is used to iterate over the raw logs and unpacked data for LogHTLCRefund events raised by the Htlc contract.
type HtlcLogHTLCRefundIterator struct {
	Event *HtlcLogHTLCRefund // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *HtlcLogHTLCRefundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HtlcLogHTLCRefund)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(HtlcLogHTLCRefund)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *HtlcLogHTLCRefundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HtlcLogHTLCRefundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HtlcLogHTLCRefund represents a LogHTLCRefund event raised by the Htlc contract.
type HtlcLogHTLCRefund struct {
	ContractId [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogHTLCRefund is a free log retrieval operation binding the contract event 0x989b3a845197c9aec15f8982bbb30b5da714050e662a7a287bb1a94c81e2e70e.
//
// Solidity: event LogHTLCRefund(bytes32 indexed contractId)
func (_Htlc *HtlcFilterer) FilterLogHTLCRefund(opts *bind.FilterOpts, contractId [][32]byte) (*HtlcLogHTLCRefundIterator, error) {

	var contractIdRule []interface{}
	for _, contractIdItem := range contractId {
		contractIdRule = append(contractIdRule, contractIdItem)
	}

	logs, sub, err := _Htlc.contract.FilterLogs(opts, "LogHTLCRefund", contractIdRule)
	if err != nil {
		return nil, err
	}
	return &HtlcLogHTLCRefundIterator{contract: _Htlc.contract, event: "LogHTLCRefund", logs: logs, sub: sub}, nil
}

// WatchLogHTLCRefund is a free log subscription operation binding the contract event 0x989b3a845197c9aec15f8982bbb30b5da714050e662a7a287bb1a94c81e2e70e.
//
// Solidity: event LogHTLCRefund(bytes32 indexed contractId)
func (_Htlc *HtlcFilterer) WatchLogHTLCRefund(opts *bind.WatchOpts, sink chan<- *HtlcLogHTLCRefund, contractId [][32]byte) (event.Subscription, error) {

	var contractIdRule []interface{}
	for _, contractIdItem := range contractId {
		contractIdRule = append(contractIdRule, contractIdItem)
	}

	logs, sub, err := _Htlc.contract.WatchLogs(opts, "LogHTLCRefund", contractIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HtlcLogHTLCRefund)
				if err := _Htlc.contract.UnpackLog(event, "LogHTLCRefund", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// HtlcLogHTLCWithdrawIterator is returned from FilterLogHTLCWithdraw and is used to iterate over the raw logs and unpacked data for LogHTLCWithdraw events raised by the Htlc contract.
type HtlcLogHTLCWithdrawIterator struct {
	Event *HtlcLogHTLCWithdraw // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *HtlcLogHTLCWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HtlcLogHTLCWithdraw)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(HtlcLogHTLCWithdraw)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *HtlcLogHTLCWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HtlcLogHTLCWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HtlcLogHTLCWithdraw represents a LogHTLCWithdraw event raised by the Htlc contract.
type HtlcLogHTLCWithdraw struct {
	ContractId [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogHTLCWithdraw is a free log retrieval operation binding the contract event 0xd6fd4c8e45bf0c70693141c7ce46451b6a6a28ac8386fca2ba914044e0e23916.
//
// Solidity: event LogHTLCWithdraw(bytes32 indexed contractId)
func (_Htlc *HtlcFilterer) FilterLogHTLCWithdraw(opts *bind.FilterOpts, contractId [][32]byte) (*HtlcLogHTLCWithdrawIterator, error) {

	var contractIdRule []interface{}
	for _, contractIdItem := range contractId {
		contractIdRule = append(contractIdRule, contractIdItem)
	}

	logs, sub, err := _Htlc.contract.FilterLogs(opts, "LogHTLCWithdraw", contractIdRule)
	if err != nil {
		return nil, err
	}
	return &HtlcLogHTLCWithdrawIterator{contract: _Htlc.contract, event: "LogHTLCWithdraw", logs: logs, sub: sub}, nil
}

// WatchLogHTLCWithdraw is a free log subscription operation binding the contract event 0xd6fd4c8e45bf0c70693141c7ce46451b6a6a28ac8386fca2ba914044e0e23916.
//
// Solidity: event LogHTLCWithdraw(bytes32 indexed contractId)
func (_Htlc *HtlcFilterer) WatchLogHTLCWithdraw(opts *bind.WatchOpts, sink chan<- *HtlcLogHTLCWithdraw, contractId [][32]byte) (event.Subscription, error) {

	var contractIdRule []interface{}
	for _, contractIdItem := range contractId {
		contractIdRule = append(contractIdRule, contractIdItem)
	}

	logs, sub, err := _Htlc.contract.WatchLogs(opts, "LogHTLCWithdraw", contractIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HtlcLogHTLCWithdraw)
				if err := _Htlc.contract.UnpackLog(event, "LogHTLCWithdraw", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
