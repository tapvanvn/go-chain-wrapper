package campain

import (
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ContractTool struct {
	id             int
	ready          bool
	contract       *Contract
	commands       chan *ContractCall
	waitingCommand *ContractCall
	campain        *Campain
	backend        *ethclient.Client
}

func NewContractTool(campain *Campain, contractName string) (*ContractTool, error) {

	__tool_id += 1
	tool := &ContractTool{id: __tool_id,
		ready:          false,
		commands:       make(chan *ContractCall),
		waitingCommand: nil,
		campain:        campain,
	}
	backend, err := ethclient.Dial("https://bsc-dataseed1.binance.org")
	if err != nil {
		return nil, err
	}
	tool.backend = backend

	abiObj, ok := campain.abis[contractName]
	if !ok {
		return nil, errors.New("abi not load")
	}
	ethAbiObj := abiObj.(*EthereumABI)
	address := common.HexToAddress(ethAbiObj.ContractAddress)

	contract, err := ethAbiObj.NewContract(address, tool.backend)
	if err != nil {
		return nil, err
	}
	tool.contract = contract

	go tool.process()

	return tool, nil
}
func (tool *ContractTool) AddCall(call *ContractCall) {
	tool.commands <- call
}

func (tool *ContractTool) process() {

	for {

		if tool.waitingCommand != nil {
			time.Sleep(time.Microsecond * 20)
			continue
		}
		cmd, ok := <-tool.commands

		if !ok {
			break
		}

		tool.waitingCommand = cmd
		var out []interface{}
		var err error = nil

		if cmd.Params == nil || len(cmd.Params) == 0 {

			err = tool.contract.ContractCaller.contract.Call(nil, &out, cmd.FuncName)
		} else {
			err = tool.contract.ContractCaller.contract.Call(nil, &out, cmd.FuncName, cmd.Params...)
		}

		if err != nil {

			fmt.Println("contract error", err)
		} else {
			fmt.Println("contract call:", out)
			cmd.Out = &out
		}
		//error
		tool.waitingCommand = nil

		//process command
	}
}
