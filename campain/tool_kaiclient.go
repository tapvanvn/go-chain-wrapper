package campain

import (
	"context"
	"strconv"
	"strings"

	"github.com/tapvanvn/go-jsonrpc-wrapper/entity"
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

	__tool_id += 1
	tool := &KaiClientTool{id: __tool_id,
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
			BlockNumber:      strconv.FormatUint(blockNumber, 10),
			Gas:              strconv.FormatUint(trans.GasUsed, 10),
			GasPrice:         strconv.FormatUint(trans.GasPrice, 10),
			Hash:             trans.Hash,
			Input:            trans.InputData,
			From:             strings.ToLower(trans.From),
			To:               strings.ToLower(trans.To),
			TransactionIndex: strconv.FormatUint(uint64(trans.TransactionIndex), 10),
		}
		result = append(result, entityTrans)
	}

	return result, nil

}
