package campain

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tapvanvn/go-jsonrpc-wrapper/entity"
)

type EthClientTool struct {
	id      int
	ready   bool
	campain *Campain
	backend *ethclient.Client
}

func NewEthClientTool(campain *Campain, backendURL string) (*EthClientTool, error) {

	__tool_id += 1
	tool := &EthClientTool{id: __tool_id,
		ready:   false,
		campain: campain,
	}
	backend, err := ethclient.Dial(backendURL) //"https://bsc-dataseed1.binance.org")
	if err != nil {
		return nil, err
	}
	tool.backend = backend

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
	transCount, err := tool.backend.TransactionCount(ctx, hash)
	if err != nil {
		return nil, err
	}
	result := []*entity.Transaction{}
	for i := uint(0); i < transCount; i++ {
		trans, err := tool.backend.TransactionInBlock(ctx, hash, i)
		if err != nil {
			return nil, err
		}
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
