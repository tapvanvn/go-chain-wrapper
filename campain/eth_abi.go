package campain

import (
	"encoding/hex"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/tapvanvn/go-jsonrpc-wrapper/system"
)

type EthereumABI struct {
	Abi abi.ABI
}

func NewEthereumABI(abiFileName string) (IABI, error) {
	// load contract ABI
	ethABI := &EthereumABI{}

	file, err := os.Open(system.RootPath + "/abi_file/" + abiFileName)

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
	/*
		// decode txInput method signature
		decodedSig, err := hex.DecodeString(txInput[2:10])
		if err != nil {
			log.Fatal(err)
		}

		// recover Method from signature and ABI
		method, err := abi.MethodById(decodedSig)
		if err != nil {
			log.Fatal(err)
		}

		// decode txInput Payload
		decodedData, err := hex.DecodeString(txInput[10:])
		if err != nil {
			log.Fatal(err)
		}

		// create strut that matches input names to unpack
		// for example my function takes 2 inputs, with names "Name1" and "Name2" and of type uint256 (solidity)
		type FunctionInputs struct {
			Name1 *big.Int // *big.Int for uint256 for example
			Name2 *big.Int
		}

		var data FunctionInputs

		// unpack method inputs
		err = method.Inputs.Unpack(&data, decodedData)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(data)
	*/
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
