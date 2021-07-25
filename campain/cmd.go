package campain

var __cmd_id int = 0

func IssueID() int {
	__cmd_id += 1
	return __cmd_id
}

type Command interface {
	Init()
	GetID() int
	GetCommand(chain string) string
	Do(tool IClientTool) error
	GetResponseInterface() interface{}
	Debug()
	Done(campain *Campain)
}

type ContractCall struct {
	ReportName string        `json:"-"`
	Topic      string        `json:"-"`
	FuncName   string        `json:"func_name"`
	Params     []interface{} `json:"-"`
	Out        *[][]byte     `json:"output"`
	Input      *[][]byte     `json:"input"`
}

func CreateContractCall(funcName string, params []interface{}, reportName string, topic string) *ContractCall {
	return &ContractCall{
		ReportName: reportName,
		Topic:      topic,
		FuncName:   funcName,
		Params:     params,
		Out:        nil,
		Input:      nil,
	}
}
