package filter

import "github.com/tapvanvn/go-chain-wrapper/entity"

type IFilter interface {
	Match(transaction *entity.Transaction) bool
}
