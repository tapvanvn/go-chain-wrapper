package campain

import "github.com/tapvanvn/go-jsonrpc-wrapper/entity"

var __tool_id = 0

func newToolID() int {
	__tool_id++
	return __tool_id
}

type IClientTool interface {
	GetLatestBlockNumber() (uint64, error)
	GetCampain() *Campain
	GetBlockTransaction(blockNumber uint64) ([]*entity.Transaction, error)
}

type ITransactionTool interface {
	Parse(transaction *entity.Transaction, track *entity.Track)
}
