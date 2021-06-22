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
	//fmt.Println("make tool", blacksmith.Campain.chainName)
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
