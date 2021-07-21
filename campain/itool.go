package campain

import "github.com/tapvanvn/go-jsonrpc-wrapper/entity"

type ITool interface {
	GetLatestBlockNumber() (uint64, error)
	GetCampain() *Campain
	GetBlockTransaction(blockNumber uint64) ([]*entity.Transaction, error)
	//GetTransactionReceipt(tx string)
}

type ITransactionTool interface {
	Parse(transaction *entity.Transaction, track *entity.Track)
}
