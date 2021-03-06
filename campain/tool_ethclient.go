package campain

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tapvanvn/go-chain-wrapper/entity"
)

type EthClientTool struct {
	id      int
	ready   bool
	campain *Campain
	backend *ethclient.Client
}

func NewEthClientTool(campain *Campain, backendURL string) (*EthClientTool, error) {

	tool := &EthClientTool{
		id:      newToolID(),
		ready:   false,
		campain: campain,
	}
	backend, err := ethclient.Dial(backendURL)
	if err != nil {
		fmt.Println("new ethclient err", campain.chainName, err)
		return nil, err
	}
	tool.backend = backend
	fmt.Println("new ethclient tool", campain.chainName, backendURL)
	return tool, nil
}

func (tool *EthClientTool) GetLatestBlockNumber() (uint64, error) {
	return tool.backend.BlockNumber(context.TODO())
}

func (tool *EthClientTool) GetCampain() *Campain {
	return tool.campain
}

func (tool *EthClientTool) GetBlockTransaction(blockNumber uint64) ([]*entity.Transaction, error) {
	blockBigNum := big.NewInt(int64(blockNumber))
	ctx := context.TODO()
	block, err := tool.backend.BlockByNumber(ctx, blockBigNum)
	if err != nil {

		return nil, err
	}
	hash := block.Hash()

	result := []*entity.Transaction{}
	for i, trans := range block.Transactions() {
		if trans == nil {
			fmt.Println("transaction fail")
			continue
		}
		to := trans.To()

		toAddress := ""
		if to != nil {
			toAddress = to.Hex()
		}

		entityTrans := &entity.Transaction{BlockHash: trans.Hash().Hex(),
			BlockNumber:       strconv.FormatUint(blockNumber, 10),
			Gas:               strconv.FormatUint(trans.Gas(), 10),
			GasPrice:          trans.GasPrice().String(),
			Hash:              trans.Hash().Hex(),
			Input:             "0x" + hex.EncodeToString(trans.Data()),
			From:              "",
			To:                toAddress,
			TransactionIndex:  strconv.FormatUint(uint64(i), 10),
			Nonce:             fmt.Sprint(trans.Nonce()),
			OriginTransaction: trans,
			Logs:              make([]*entity.Log, 0),
		}

		if sender, err2 := tool.backend.TransactionSender(ctx, trans, hash, uint(i)); err2 == nil {

			entityTrans.From = sender.String()
		}

		if recept, err := tool.backend.TransactionReceipt(context.TODO(), hash); err == nil {

			for _, log := range recept.Logs {
				entityLog := &entity.Log{
					Topics: make([]string, 0),
					Data:   log.Data,
				}
				for _, topic := range log.Topics {
					entityLog.Topics = append(entityLog.Topics, topic.String())
				}
				entityTrans.Logs = append(entityTrans.Logs, entityLog)
			}
		}

		result = append(result, entityTrans)
	}
	return result, nil
}

func (tool *EthClientTool) GetTransactionReceipt(txHash []byte) {
	txHashParsed := common.BytesToHash(txHash)
	if recept, err := tool.backend.TransactionReceipt(context.TODO(), txHashParsed); err != nil {
		fmt.Println(recept, err)
	}
}
