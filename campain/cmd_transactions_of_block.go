package campain

import (
	"fmt"
	"strconv"

	"github.com/tapvanvn/go-jsonrpc-wrapper/entity"
)

type CmdTransactionsOfBlock struct {
	id           int
	BlockNumber  uint64
	Transactions []*entity.Transaction
}

func CreateCmdTransactionsOfBlock(blockNumber uint64) *CmdTransactionsOfBlock {
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

func (cmd *CmdTransactionsOfBlock) Do(tool IClientTool) {
	transactions, err := tool.GetBlockTransaction(cmd.BlockNumber)
	if err != nil {
		return
	}
	cmd.Transactions = transactions
	cmd.Done(tool.GetCampain())
}

func (cmd *CmdTransactionsOfBlock) GetCommand(chain string) string {
	if chain == "kai" {
		if cmd.BlockNumber == 0 {
			return "kai.getTransactionsByBlockNumber(\"latest\")"
		}
		//num, err := strconv.ParseInt(hex_num, 16, 64)
		hexNum := strconv.FormatUint(cmd.BlockNumber, 16)
		return "kai.getTransactionsByBlockNumber(\"0x" + hexNum + "\")"
	}
	if cmd.BlockNumber == 0 {
		return "eth.getTransactionsByBlockNumber(\"latest\")"
	}
	//num, err := strconv.ParseInt(hex_num, 16, 64)
	hexNum := strconv.FormatUint(cmd.BlockNumber, 16)
	return "eth.getTransactionsByBlockNumber(\"0x" + hexNum + "\")"
}

func (cmd *CmdTransactionsOfBlock) GetResponseInterface() interface{} {
	return &cmd.Transactions
}

func (cmd *CmdTransactionsOfBlock) Done(campain *Campain) {
	fmt.Println("done fine transaction for block", cmd.BlockNumber)
	for _, trans := range cmd.Transactions {
		campain.chnTransactions <- *trans
	}
}
