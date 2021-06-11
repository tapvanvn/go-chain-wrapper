package campain

import (
	"fmt"
)

//MARK:JsonRpcBlackSmith
type JsonRpcBlackSmith struct {
	Campain *Campain
}

//Make make tool
func (blacksmith *JsonRpcBlackSmith) Make() interface{} {

	tool, err := NewTool(blacksmith.Campain)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return tool
}

type EthContractBlackSmith struct {
	Campain      *Campain
	ContractName string
}

//Make make tool
func (blacksmith *EthContractBlackSmith) Make() interface{} {

	tool, err := NewContractTool(blacksmith.Campain, blacksmith.ContractName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return tool
}

type ExportPubsubBlackSmith struct {
	Campain    *Campain
	Topic      string
	ExportName string
}

//Make make tool
func (blacksmith *ExportPubsubBlackSmith) Make() interface{} {

	if hub, ok := blacksmith.Campain.pubsubHub[blacksmith.ExportName]; ok {
		tool, err := NewExportPubSubTool(hub, blacksmith.Topic)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return tool
	}
	return nil
}
