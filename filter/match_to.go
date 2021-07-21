package filter

import (
	"strings"

	"github.com/tapvanvn/go-jsonrpc-wrapper/entity"
)

type FilMatchTo struct {
	Address string
}

func (filter *FilMatchTo) Match(transaction *entity.Transaction) bool {

	return strings.ToLower(transaction.To) == strings.ToLower(filter.Address)
}
