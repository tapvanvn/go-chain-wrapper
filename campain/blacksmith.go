package campain

import (
	"fmt"

	"github.com/tapvanvn/go-jsonrpc-wrapper/export"
)

//MARK:Client
type ClientBlackSmith struct {
	Campain *Campain
}

//Make make tool
func (blacksmith *ClientBlackSmith) Make(origin string, meta interface{}) interface{} {

	endpoint := string(blacksmith.Campain.GetEndpoint(origin))
	if len(endpoint) == 0 {
		return nil
	}
	//TODO: random backend
	if blacksmith.Campain.chainName == "bsc" {
		tool, err := NewEthClientTool(blacksmith.Campain, endpoint)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return tool
	} else if blacksmith.Campain.chainName == "kai" {
		tool, err := NewKaiClientTool(blacksmith.Campain, endpoint)
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
	ContractName string
}

//Make make tool
func (blacksmith *TransactionBlackSmith) Make(origin string, meta interface{}) interface{} {
	endpoint := string(blacksmith.Campain.GetEndpoint(origin))
	if len(endpoint) == 0 {
		return nil
	}
	if blacksmith.Campain.chainName == "bsc" {
		if rawAbi, ok := blacksmith.Campain.abis[blacksmith.ContractName]; ok {
			if abi := rawAbi.(*EthereumABI); abi != nil {

				tool, err := NewEthTransactionTool(blacksmith.Campain, endpoint, abi, blacksmith.ContractName)
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
				tool, err := NewKaiTransactionTool(blacksmith.Campain, endpoint, abi)
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
type ContractBlackSmith struct {
	Campain      *Campain
	ContractName string
}

//Make make tool
func (blacksmith *ContractBlackSmith) Make(origin string, meta interface{}) interface{} {
	endpoint := string(blacksmith.Campain.GetEndpoint(origin))
	if len(endpoint) == 0 {
		return nil
	}
	//TODO: random backend
	fmt.Println("create contract tool")
	tool, err := NewContractTool(blacksmith.Campain, blacksmith.ContractName, endpoint)
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
func (blacksmith *ExportBlackSmith) Make(origin string, meta interface{}) interface{} {

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
