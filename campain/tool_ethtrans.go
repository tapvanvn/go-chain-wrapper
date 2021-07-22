package campain

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tapvanvn/go-jsonrpc-wrapper/entity"
)

type EthTransactionTool struct {
	id       int
	ready    bool
	campain  *Campain
	backend  *ethclient.Client
	contract IContract
	abi      *EthereumABI
}

func NewEthTransactionTool(campain *Campain, backendURL string, abi *EthereumABI, contractName string) (*EthTransactionTool, error) {

	__tool_id += 1
	tool := &EthTransactionTool{id: __tool_id,
		ready:   false,
		campain: campain,
		abi:     abi,
	}
	backend, err := ethclient.Dial(backendURL)
	if err != nil {
		return nil, err
	}
	tool.backend = backend

	contractAddress, ok := campain.contractAddress[contractName]
	if !ok {
		return nil, errors.New("contract address not found")
	}
	contract, err := abi.NewContract(contractAddress, backendURL)
	if err != nil {
		return nil, err
	}
	tool.contract = contract

	return tool, nil
}

var __count = 0

func (tool *EthTransactionTool) Parse(transaction *entity.Transaction, track *entity.Track) {

	hash := common.HexToHash(transaction.Hash)
	trans := transaction.OriginTransaction.(*types.Transaction)
	if trans != nil {
		if recept, err := tool.backend.TransactionReceipt(context.TODO(), hash); err == nil {
			__count++
			events := []*entity.Event{}
			for _, log := range recept.Logs {

				for _, topic := range log.Topics {

					if event, err := tool.abi.Abi.EventByID(topic); err == nil {

						outs, err := event.Inputs.Unpack(log.Data)
						if err == nil {

							count := 0
							if len(outs) > 0 {
								evt := &entity.Event{
									Name:      event.Name,
									Arguments: make(map[string]string),
								}
								for _, args := range event.Inputs {
									argType := args.Type.String()
									value := ""
									if argType == "uint256" {

										tryBig := outs[count].(*big.Int)
										value = tryBig.String()

									} else if argType == "address" {
										value = outs[count].(common.Address).String()
									} else {
										value = "unsupported"
									}
									if err != nil {
										break
									}
									evt.Arguments[args.Name] = fmt.Sprintf("%s.%s", args.Type.String(), value)
									count++
								}

								events = append(events, evt)
							}
						}
					}
				}
			}
			if len(events) > 0 {

				report := &ReportEvent{
					track:  track,
					events: events,
				}
				tool.campain.ChnEvent <- report
			}
		}
	}
}
