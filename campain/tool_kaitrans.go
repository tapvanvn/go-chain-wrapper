package campain

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tapvanvn/go-jsonrpc-wrapper/entity"
)

type KaiTransactionTool struct {
	id      int
	ready   bool
	campain *Campain
	backend *ethclient.Client
	abi     *KaiABI
}

func NewKaiTransactionTool(campain *Campain, backendURL string, abi *KaiABI) (*KaiTransactionTool, error) {

	__tool_id += 1
	tool := &KaiTransactionTool{id: __tool_id,
		ready:   false,
		campain: campain,
		abi:     abi,
	}
	backend, err := ethclient.Dial(backendURL)
	if err != nil {
		return nil, err
	}
	tool.backend = backend

	return tool, nil
}
func (tool *KaiTransactionTool) Parse(transaction *entity.Transaction, track *entity.Track) {

}
