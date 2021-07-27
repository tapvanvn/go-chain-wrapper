package campain

import "github.com/tapvanvn/go-chain-wrapper/entity"

type IContract interface {
	Call(result *[]interface{}, method string, params ...interface{}) error
	ParseLog(topic string, data []byte) (*entity.Event, error)
}
type IABI interface {
	Info()
	GetMethod(input string) (string, []interface{}, error)
	NewContract(address ContractAddress, backendURL Endpoint) (IContract, error)
}
