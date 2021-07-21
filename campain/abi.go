package campain

type IContract interface {
	Call(result *[]interface{}, method string, params ...interface{}) error
	ParseLog(result map[string]interface{}, event string, log interface{}) error
}
type IABI interface {
	Info()
	GetMethod(input string) (string, []interface{}, error)
	NewContract(address string, backendURL string) (IContract, error)
}
