package campain

import (
	"context"
	"strconv"

	"github.com/tapvanvn/go-chain-wrapper/entity"
	"github.com/tapvanvn/go-kaiclient/kardia"
	"go.uber.org/zap"
)

type KaiClientTool struct {
	id      int
	ready   bool
	campain *Campain
	backend kardia.Node
}

func NewKaiClientTool(campain *Campain, backendURL string) (*KaiClientTool, error) {

	tool := &KaiClientTool{
		id:      newToolID(),
		ready:   false,
		campain: campain,
	}
	lgr, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	node, err := kardia.NewNode(backendURL, lgr)

	if err != nil {
		return nil, err
	}
	tool.backend = node

	return tool, nil
}

func (tool *KaiClientTool) GetLatestBlockNumber() (uint64, error) {

	return tool.backend.LatestBlockNumber(context.TODO())
}
func (tool *KaiClientTool) GetCampain() *Campain {
	return tool.campain
}

func (tool *KaiClientTool) GetBlockTransaction(blockNumber uint64) ([]*entity.Transaction, error) {

	ctx := context.TODO()

	head, err := tool.backend.BlockHeaderByNumber(ctx, blockNumber)
	if err != nil {
		return nil, err
	}

	block, err := tool.backend.BlockByHash(ctx, head.Hash)

	result := []*entity.Transaction{}
	for _, trans := range block.Txs {

		entityTrans := &entity.Transaction{BlockHash: trans.BlockHash,
			BlockNumber:       strconv.FormatUint(blockNumber, 10),
			Gas:               strconv.FormatUint(trans.GasUsed, 10),
			GasPrice:          strconv.FormatUint(trans.GasPrice, 10),
			Hash:              trans.Hash,
			Input:             trans.InputData,
			From:              trans.From,
			To:                trans.To,
			TransactionIndex:  strconv.FormatUint(uint64(trans.TransactionIndex), 10),
			OriginTransaction: trans,
			Logs:              make([]*entity.Log, 0),
		}
		if recept, err := tool.backend.GetTransactionReceipt(context.TODO(), trans.Hash); err == nil {
			entityTrans.Success = recept.Status == 1
			for _, log := range recept.Logs {
				entityLog := &entity.Log{
					Topics: make([]string, 0),
					Data:   []byte(log.Data),
				}
				for _, topic := range log.Topics {
					entityLog.Topics = append(entityLog.Topics, topic)
				}
				entityTrans.Logs = append(entityTrans.Logs, entityLog)
			}
		}
		result = append(result, entityTrans)
	}

	return result, nil

}
