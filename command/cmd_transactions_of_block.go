package command

import (
	"fmt"

	"github.com/tapvanvn/go-jsonrpc-wrapper/entity"
)

type CmdTransactionsOfBlock struct {
	id           int
	BlockNumber  int
	Transactions []*entity.Transaction
}

func CreateCmdTransactionsOfBlock(blockNumber int) *CmdTransactionsOfBlock {
	return &CmdTransactionsOfBlock{
		BlockNumber: blockNumber,
	}
}

func (cmd *CmdTransactionsOfBlock) Init() {
	cmd.id = IssueID()
}

func (cmd *CmdTransactionsOfBlock) Debug() {
	fmt.Println(cmd.BlockNumber)
	fmt.Println("num transaction:", len(cmd.Transactions))
}

func (cmd *CmdTransactionsOfBlock) GetID() int {
	return cmd.id
}

func (cmd *CmdTransactionsOfBlock) GetCommand() string {
	str := "eth.getTransactionsByBlockNumber(\"latest\")"
	return str
}

func (cmd *CmdTransactionsOfBlock) GetResponseInterface() interface{} {
	return &cmd.Transactions
}
