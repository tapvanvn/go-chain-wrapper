package campain

import "fmt"

type CmdGetLatestBlockNumber struct {
	id          int
	BlockNumber uint64
}

func CreateCmdLatestBlockNumber() *CmdGetLatestBlockNumber {
	return &CmdGetLatestBlockNumber{}
}

func (cmd *CmdGetLatestBlockNumber) Init() {
	cmd.id = IssueID()
}

func (cmd *CmdGetLatestBlockNumber) GetID() int {
	return cmd.id
}

func (cmd *CmdGetLatestBlockNumber) Do(tool IClientTool) {
	blockNumber, err := tool.GetLatestBlockNumber()
	if err != nil {
		//TODO: process error
		return
	}
	cmd.BlockNumber = blockNumber
	cmd.Done(tool.GetCampain())
}

func (cmd *CmdGetLatestBlockNumber) GetCommand(chain string) string {
	if chain == "kai" {
		return "kai.blockNumber"
	}
	return "eth.blockNumber"
}

func (cmd *CmdGetLatestBlockNumber) GetResponseInterface() interface{} {
	return &cmd.BlockNumber
}

func (cmd *CmdGetLatestBlockNumber) Debug() {
	fmt.Println(cmd.BlockNumber)
}

func (cmd *CmdGetLatestBlockNumber) Done(campain *Campain) {
	campain.chnBlockNumber <- cmd.BlockNumber
}
