package campain

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"

	"github.com/kardiachain/go-kardia/lib/abi"
	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/types"

	"github.com/tapvanvn/go-jsonrpc-wrapper/system"
)

var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = common.Big1
	_ = types.BloomLookup
)

type KaiABI struct {
	Abi             abi.ABI
	ContractAddress string
}

func NewKaiABI(abiFileName string, address string) (IABI, error) {
	// load contract ABI
	kaiABI := &KaiABI{
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
	kaiABI.Abi = abiObj

	return kaiABI, nil
}

func (kaiAbi *KaiABI) Info() {
	fmt.Println("events:", len(kaiAbi.Abi.Events))
	for _, event := range kaiAbi.Abi.Events {
		fmt.Println("\t", event.Name, event.Inputs)
	}
	fmt.Println("methods:", len(kaiAbi.Abi.Methods))
	for _, method := range kaiAbi.Abi.Methods {
		fmt.Println("\t", method.Name, method.Inputs)
	}
}

func (kaiAbi *KaiABI) GetMethod(input string) (string, []interface{}, error) {

	if len(input) < 10 {
		return "", nil, errors.New("invalid input:" + input)
	}
	decodedSig, err := hex.DecodeString(input[2:10])
	if err != nil {

		return "", nil, err
	}

	method, err := kaiAbi.Abi.MethodById(decodedSig)
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

/*
// bindStore binds a generic wrapper to an already deployed contract.
func (ethAbi *KaiABI) bindContract(address common.Address,
	caller bind.ContractCaller,
	transactor bind.ContractTransactor,
	filterer bind.ContractFilterer) (*bind.BoundContract, error) {

	return kardia.NewBoundContract(address, ethAbi.Abi, caller, transactor, filterer), nil
}

// Store is an auto generated Go binding around an Ethereum contract.
type KaiContract struct {
	ContractCaller     // Read-only binding to the contract
	ContractTransactor // Write-only binding to the contract
	ContractFilterer   // Log filterer for contract events
}

// StoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type KaiContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type KaiContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KaiContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KaiContractSession struct {
	Contract     *Contract         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KaiContractCallerSession struct {
	Contract *ContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// StoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KaiContractTransactorSession struct {
	Contract     *ContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// StoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type KaiContractRaw struct {
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

// NewStore creates a new instance of Store, bound to a specific deployed contract.
func (ethAbi *KaiABI) NewContract(address common.Address, backend bind.ContractBackend) (*Contract, error) {
	contract, err := ethAbi.bindContract(address, backend, backend, backend)

	if err != nil {
		return nil, err
	}
	return &Contract{ContractCaller: ContractCaller{contract: contract},
		ContractTransactor: ContractTransactor{contract: contract},
		ContractFilterer:   ContractFilterer{contract: contract}}, nil
}

// NewStoreCaller creates a new read-only instance of Store, bound to a specific deployed contract.
func (ethAbi *KaiABI) NewContractCaller(address common.Address, caller bind.ContractCaller) (*ContractCaller, error) {
	contract, err := ethAbi.bindContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractCaller{contract: contract}, nil
}

// NewStoreTransactor creates a new write-only instance of Store, bound to a specific deployed contract.
func (ethAbi *KaiABI) NewStoreTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractTransactor, error) {
	contract, err := ethAbi.bindContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractTransactor{contract: contract}, nil
}

// NewStoreFilterer creates a new log filterer instance of Store, bound to a specific deployed contract.
func (ethAbi *KaiABI) NewStoreFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractFilterer, error) {
	contract, err := ethAbi.bindContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractFilterer{contract: contract}, nil
}

func (_Store *KaiContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {

	return _Store.Contract.ContractCaller.contract.Call(opts, result, method, params...)
}
/*
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
*/
