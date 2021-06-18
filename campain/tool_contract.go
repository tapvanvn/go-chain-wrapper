package campain

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ContractTool struct {
	id         int
	contract   IContract
	campain    *Campain
	backendURL string
}

func NewContractTool(campain *Campain, contractName string, backendURL string) (*ContractTool, error) {

	__tool_id += 1
	tool := &ContractTool{id: __tool_id,
		campain: campain,
	}

	abiObj, ok := campain.abis[contractName]
	if !ok {
		return nil, errors.New("abi not load")
	}
	contractAddress, ok := campain.contractAddress[contractName]
	if !ok {
		return nil, errors.New("contract address not found")
	}
	contract, err := abiObj.NewContract(contractAddress, backendURL)
	if err != nil {
		return nil, err
	}
	tool.contract = contract

	return tool, nil
}

func (tool *ContractTool) Process(call *ContractCall) {

	var outs []interface{}
	var err error = nil

	fmt.Println("do contract:", call.FuncName)

	if call.Params == nil || len(call.Params) == 0 {

		err = tool.contract.Call(&outs, call.FuncName)
	} else {
		err = tool.contract.Call(&outs, call.FuncName, call.Params...)
	}
	fmt.Println("result", outs)
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

		fmt.Println("call", call)

		if call.ReportName != "" && call.Topic != "" {

			tool.campain.Report(call.ReportName, call.Topic, call)
		}
	}
}
