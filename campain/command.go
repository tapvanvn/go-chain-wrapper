package campain

var __cmd_id int = 0

func IssueID() int {
	__cmd_id += 1
	return __cmd_id
}

type Command interface {
	Init()
	GetID() int
	GetCommand() string
	GetResponseInterface() interface{}
	Debug()
	Done(campain *Campain)
}

type ContractCall struct {
	FuncName string
	Params   []interface{}
	Out      *[]interface{}
}

func CreateContractCall(funcName string, params []interface{}) *ContractCall {
	return &ContractCall{
		FuncName: funcName,
		Params:   params,
		Out:      nil,
	}
}
