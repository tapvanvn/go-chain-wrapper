package filter

import "github.com/tapvanvn/go-jsonrpc-wrapper/entity"

type FilMatchTo struct {
	Address string
}

func (filter *FilMatchTo) Match(transaction *entity.Transaction) bool {

	return transaction.To == filter.Address
}
