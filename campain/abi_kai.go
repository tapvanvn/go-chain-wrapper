package campain

import (
	"context"
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
	"go.uber.org/zap"

	"github.com/tapvanvn/go-jsonrpc-wrapper/system"
	"github.com/tapvanvn/go-kaiclient/kardia"
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

// bindStore binds a generic wrapper to an already deployed contract.
func (kaiABI *KaiABI) NewContract(address string, backendURL string) (IContract, error) {
	byteAddress := common.HexToAddress(address)
	lgr, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	node, err := kardia.NewNode(backendURL, lgr)

	if err != nil {
		return nil, err
	}
	contract := &KaiContract{
		node:     node,
		address:  address,
		Abi:      &kaiABI.Abi,
		contract: kardia.NewBoundContract(node, &kaiABI.Abi, byteAddress),
	}
	return contract, nil
}

type KaiContract struct {
	node     kardia.Node
	address  string
	Abi      *abi.ABI
	contract *kardia.BoundContract
}

func (contract *KaiContract) Call(result *[]interface{}, method string, params ...interface{}) error {

	payload, err := contract.contract.Abi.Pack(method, params...)
	if err != nil {
		fmt.Println("call error", err)
		return err
	}
	res, err := contract.node.KardiaCall(context.TODO(), kardia.ConstructCallArgs(contract.address, payload))

	resResult, err := contract.contract.Abi.Unpack(method, res)
	if result == nil {
		result = new([]interface{})
	}

	if err != nil {
		return err
	}

	*result = resResult

	return nil
}
