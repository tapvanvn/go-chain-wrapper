package filter

import "github.com/tapvanvn/go-jsonrpc-wrapper/entity"

type IFilter interface {
	Match(transaction *entity.Transaction) bool
}
