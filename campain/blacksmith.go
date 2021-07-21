package campain

import (
	"fmt"

	"github.com/tapvanvn/go-jsonrpc-wrapper/export"
)

//MARK:JsonRpcBlackSmith
type JsonRpcBlackSmith struct {
	Campain     *Campain
	BackendURLS []string
}

//Make make tool
func (blacksmith *JsonRpcBlackSmith) Make() interface{} {

	//TODO: random backend
	tool, err := NewTool(blacksmith.Campain, blacksmith.BackendURLS[0])
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return tool
}

//MARK:Client
type ClientBlackSmith struct {
	Campain     *Campain
	BackendURLS []string
}

//Make make tool
func (blacksmith *ClientBlackSmith) Make() interface{} {

	//TODO: random backend
	if blacksmith.Campain.chainName == "bsc" {
		tool, err := NewEthClientTool(blacksmith.Campain, blacksmith.BackendURLS[0])
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return tool
	} else if blacksmith.Campain.chainName == "kai" {
		tool, err := NewKaiClientTool(blacksmith.Campain, blacksmith.BackendURLS[0])
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return tool
	}
	return nil
}

//MARK: Transaction
type TransactionBlackSmith struct {
	Campain      *Campain
	BackendURLS  []string
	ContractName string
}

//Make make tool
func (blacksmith *TransactionBlackSmith) Make() interface{} {
	if blacksmith.Campain.chainName == "bsc" {
		if rawAbi, ok := blacksmith.Campain.abis[blacksmith.ContractName]; ok {
			if abi := rawAbi.(*EthereumABI); abi != nil {

				tool, err := NewEthTransactionTool(blacksmith.Campain, blacksmith.BackendURLS[0], abi, blacksmith.ContractName)
				if err != nil {
					fmt.Println(err)
					return nil
				}
				return tool

			}
		}
	} else if blacksmith.Campain.chainName == "kai" {
		if rawAbi, ok := blacksmith.Campain.abis[blacksmith.ContractName]; ok {
			if abi := rawAbi.(*KaiABI); abi != nil {
				tool, err := NewKaiTransactionTool(blacksmith.Campain, blacksmith.BackendURLS[0], abi)
				if err != nil {
					fmt.Println(err)
					return nil
				}
				return tool
			}
		}
	}
	return nil
}

//MARK: contract call

type EthContractBlackSmith struct {
	Campain      *Campain
	ContractName string
	BackendURLS  []string
}

//Make make tool
func (blacksmith *EthContractBlackSmith) Make() interface{} {
	//TODO: random backend
	fmt.Println("create contract tool")
	tool, err := NewContractTool(blacksmith.Campain, blacksmith.ContractName, blacksmith.BackendURLS[0])
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return tool
}

type ExportBlackSmith struct {
	ExportName string
}

//Make make tool
func (blacksmith *ExportBlackSmith) Make() interface{} {

	ex := export.GetExport(blacksmith.ExportName)
	if ex != nil {

		tool, err := NewExportTool(ex)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return tool
	}
	return nil
}
