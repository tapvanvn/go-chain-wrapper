package campain

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
)

type ContractTool struct {
	id       int
	contract *Contract
	campain  *Campain
	backend  *ethclient.Client
}

func NewContractTool(campain *Campain, contractName string, backendURL string) (*ContractTool, error) {

	__tool_id += 1
	tool := &ContractTool{id: __tool_id,
		campain: campain,
	}
	backend, err := ethclient.Dial(backendURL) //"https://bsc-dataseed1.binance.org")
	if err != nil {
		return nil, err
	}
	tool.backend = backend

	abiObj, ok := campain.abis[contractName]
	if !ok {
		return nil, errors.New("abi not load")
	}
	_ = abiObj
	/*ethAbiObj := abiObj.(*EthereumABI)
	address := common.HexToAddress(ethAbiObj.ContractAddress)

	contract, err := ethAbiObj.NewContract(address, tool.backend)
	if err != nil {
		return nil, err
	}
	tool.contract = contract
	*/
	return tool, nil
}

func (tool *ContractTool) Process(call *ContractCall) {

	var outs []interface{}
	var err error = nil

	fmt.Println("do contract:", call.FuncName)

	if call.Params == nil || len(call.Params) == 0 {

		err = tool.contract.ContractCaller.contract.Call(nil, &outs, call.FuncName)
	} else {
		err = tool.contract.ContractCaller.contract.Call(nil, &outs, call.FuncName, call.Params...)
	}

	if err != nil {

		fmt.Println("contract error", err)

	} else {

		results := [][]byte{}
		inputs := [][]byte{}
		for _, param := range call.Params {
			inputData, _ := json.Marshal(param)
			inputs = append(inputs, inputData)
		}
		for _, out := range outs {
			outData, _ := json.Marshal(out)
			results = append(results, outData)
		}
		call.Out = &results
		call.Input = &inputs

		if call.ReportName != "" && call.Topic != "" {

			tool.campain.Report(call.ReportName, call.Topic, call)
		}
	}
}
