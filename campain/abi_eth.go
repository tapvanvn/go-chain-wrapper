package campain

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	"github.com/tapvanvn/go-jsonrpc-wrapper/system"
)

var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

type EthereumABI struct {
	Abi             abi.ABI
	ContractAddress string
}

func NewEthereumABI(abiFileName string, address string) (IABI, error) {
	// load contract ABI
	ethABI := &EthereumABI{
		ContractAddress: address,
	}
	filepath := system.RootPath + "/abi_file/" + abiFileName
	file, err := os.Open(filepath)
	fmt.Println("load abi:", filepath)
	if err != nil {

		return nil, err
	}

	defer file.Close()
	bytes, _ := ioutil.ReadAll(file)
	abiObj, err := abi.JSON(strings.NewReader(string(bytes)))
	if err != nil {
		return nil, err
	}
	ethABI.Abi = abiObj

	return ethABI, nil
}

func (ethAbi *EthereumABI) Info() {
	fmt.Println("events:", len(ethAbi.Abi.Events))
	for _, event := range ethAbi.Abi.Events {
		fmt.Println("\t", event.Name, event.Inputs)
	}
	fmt.Println("methods:", len(ethAbi.Abi.Methods))
	for _, method := range ethAbi.Abi.Methods {
		fmt.Println("\t", method.Name, method.Inputs)
	}
}

func (ethAbi *EthereumABI) GetMethod(input string) (string, []interface{}, error) {

	if len(input) < 10 {
		return "", nil, errors.New("invalid input:" + input)
	}
	decodedSig, err := hex.DecodeString(input[2:10])
	if err != nil {

		return "", nil, err
	}

	method, err := ethAbi.Abi.MethodById(decodedSig)
	if err != nil {
		return "", nil, err
	}

	decodedData, err := hex.DecodeString(input[10:])
	if err != nil {
		return "", nil, err
	}
	args, err := method.Inputs.Unpack(decodedData)
	if err != nil {
		return "", nil, err
	}

	return method.Name, args, nil

}
func (ethAbi *EthereumABI) NewContract(address string, backendURL string) (IContract, error) {
	backend, err := ethclient.Dial(backendURL) //"https://bsc-dataseed1.binance.org")
	if err != nil {
		return nil, err
	}
	contract, err := ethAbi.bindContract(common.HexToAddress(address), backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EthContract{contract: contract}, nil
}

type EthContract struct {
	contract *bind.BoundContract
}

func (contract *EthContract) Call(result *[]interface{}, method string, params ...interface{}) error {

	return contract.contract.Call(nil, result, method, params...)
}

// bindStore binds a generic wrapper to an already deployed contract.
func (ethAbi *EthereumABI) bindContract(address common.Address,
	caller bind.ContractCaller,
	transactor bind.ContractTransactor,
	filterer bind.ContractFilterer) (*bind.BoundContract, error) {

	return bind.NewBoundContract(address, ethAbi.Abi, caller, transactor, filterer), nil
}

// Store is an auto generated Go binding around an Ethereum contract.
type Contract struct {
	ContractCaller     // Read-only binding to the contract
	ContractTransactor // Write-only binding to the contract
	ContractFilterer   // Log filterer for contract events
}

// StoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractSession struct {
	Contract     *Contract         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractCallerSession struct {
	Contract *ContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// StoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractTransactorSession struct {
	Contract     *ContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// StoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractRaw struct {
	Contract *Contract // Generic contract binding to access the raw methods on
}

// StoreCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractCallerRaw struct {
	Contract *ContractCaller // Generic read-only contract binding to access the raw methods on
}

// StoreTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractTransactorRaw struct {
	Contract *ContractTransactor // Generic write-only contract binding to access the raw methods on
}

/*
// NewStore creates a new instance of Store, bound to a specific deployed contract.
func (ethAbi *EthereumABI) NewContract(address common.Address, backend bind.ContractBackend) (*Contract, error) {
	contract, err := ethAbi.bindContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contract{ContractCaller: ContractCaller{contract: contract},
		ContractTransactor: ContractTransactor{contract: contract},
		ContractFilterer:   ContractFilterer{contract: contract}}, nil
}
*/

// NewStoreCaller creates a new read-only instance of Store, bound to a specific deployed contract.
func (ethAbi *EthereumABI) NewContractCaller(address common.Address, caller bind.ContractCaller) (*ContractCaller, error) {
	contract, err := ethAbi.bindContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractCaller{contract: contract}, nil
}

// NewStoreTransactor creates a new write-only instance of Store, bound to a specific deployed contract.
func (ethAbi *EthereumABI) NewStoreTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractTransactor, error) {
	contract, err := ethAbi.bindContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractTransactor{contract: contract}, nil
}

// NewStoreFilterer creates a new log filterer instance of Store, bound to a specific deployed contract.
func (ethAbi *EthereumABI) NewStoreFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractFilterer, error) {
	contract, err := ethAbi.bindContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractFilterer{contract: contract}, nil
}

func (_Store *ContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {

	return _Store.Contract.ContractCaller.contract.Call(opts, result, method, params...)
}

func (_Store *ContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Store.Contract.ContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Store *ContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Store.Contract.ContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Store *ContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Store.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Store *ContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Store.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Store *ContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Store.Contract.contract.Transact(opts, method, params...)
}
