package campain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/tapvanvn/go-jsonrpc-wrapper/entity"
	"github.com/tapvanvn/go-kaiclient/kardia"
	"go.uber.org/zap"
)

type KaiTransactionTool struct {
	id      int
	ready   bool
	campain *Campain
	backend kardia.Node
	abi     *KaiABI
}

func NewKaiTransactionTool(campain *Campain, backendURL string, abi *KaiABI) (*KaiTransactionTool, error) {

	__tool_id += 1
	tool := &KaiTransactionTool{id: __tool_id,
		ready:   false,
		campain: campain,
		abi:     abi,
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
func (tool *KaiTransactionTool) Parse(transaction *entity.Transaction, track *entity.Track) {

	trans := transaction.OriginTransaction.(*kardia.Transaction)
	if trans != nil {
		if recept, err := tool.backend.GetTransactionReceipt(context.TODO(), transaction.Hash); err == nil {
			__count++
			events := []*entity.Event{}
			for _, log := range recept.Logs {
				for _, topic := range log.Topics {

					if event, err := tool.abi.Abi.EventByID(common.BytesToHash([]byte(topic))); err == nil {

						outs, err := event.Inputs.Unpack([]byte(log.Data))
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
										//value = outs[count].(common.Address).String()
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
		}
	}
	/*hash := common.HexToHash(transaction.Hash)
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
	}*/
}
