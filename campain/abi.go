package campain

type IContract interface {
	Call(result *[]interface{}, method string, params ...interface{}) error
}
type IABI interface {
	Info()
	GetMethod(input string) (string, []interface{}, error)
	NewContract(address ContractAddress, backendURL Endpoint) (IContract, error)
}
