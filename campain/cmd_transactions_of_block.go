package campain

import (
	"fmt"
	"strconv"

	"github.com/tapvanvn/go-jsonrpc-wrapper/entity"
)

type CmdTransactionsOfBlock struct {
	id           int
	BlockNumber  int64
	Transactions []*entity.Transaction
}

func CreateCmdTransactionsOfBlock(blockNumber int64) *CmdTransactionsOfBlock {
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
	if cmd.BlockNumber == -1 {
		return "eth.getTransactionsByBlockNumber(\"latest\")"
	}
	//num, err := strconv.ParseInt(hex_num, 16, 64)
	hexNum := strconv.FormatInt(cmd.BlockNumber, 16)
	return "eth.getTransactionsByBlockNumber(\"0x" + hexNum + "\")"
}

func (cmd *CmdTransactionsOfBlock) GetResponseInterface() interface{} {
	return &cmd.Transactions
}

func (cmd *CmdTransactionsOfBlock) Done(campain *Campain) {
	for _, trans := range cmd.Transactions {
		campain.ChnTransactions <- *trans
	}
}
