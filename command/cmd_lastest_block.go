package command

import "fmt"

type CmdGetLatestBlockNumber struct {
	id          int
	BlockNumber int
}

func (cmd *CmdGetLatestBlockNumber) Init() {
	cmd.id = IssueID()
}

func (cmd *CmdGetLatestBlockNumber) GetID() int {
	return cmd.id
}

func (cmd *CmdGetLatestBlockNumber) GetCommand() string {
	return "eth.blockNumber"
}

func (cmd *CmdGetLatestBlockNumber) GetResponseInterface() interface{} {
	return &cmd.BlockNumber
}

func (cmd *CmdGetLatestBlockNumber) Debug() {
	fmt.Println(cmd.BlockNumber)
}
