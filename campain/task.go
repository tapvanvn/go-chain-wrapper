package campain

import (
	"fmt"

	"github.com/tapvanvn/go-jsonrpc-wrapper/entity"
	goworker "github.com/tapvanvn/goworker/v2"
)

//MARK: ClientTask
type ClientTask struct {
	tool    string
	command Command
}

func NewClientTask(tool string, command Command) *ClientTask {

	return &ClientTask{

		tool:    tool,
		command: command,
	}
}

func (task *ClientTask) Process(tool interface{}) goworker.ToolQuantity {

	if tool1, ok := tool.(IClientTool); ok {

		task.command.Do(tool1)

	} else {

		fmt.Println("not tool for client task", task.ToolLabel(), tool)
	}

	return goworker.ToolQuantityGood
}

func (task *ClientTask) ToolLabel() string {
	return task.tool
}

//MARK: Transaction Task
type TransactionTask struct {
	tool        string
	transaction *entity.Transaction
	track       *entity.Track
}

func NewTransactionTask(tool string, command Command) *ClientTask {
	return &ClientTask{
		tool:    tool,
		command: command,
	}
}

func (task *TransactionTask) Process(tool interface{}) goworker.ToolQuantity {

	if tool1, ok := tool.(ITransactionTool); ok {

		tool1.Parse(task.transaction, task.track)

	} else {

		fmt.Println("not tool for client task", task.ToolLabel(), tool)
	}
	return goworker.ToolQuantityGood
}

func (task *TransactionTask) ToolLabel() string {
	return task.tool
}

//MARK: Contract task
type ContractTask struct {
	tool string
	call *ContractCall
}

func NewContractTask(chain string, call *ContractCall) *ContractTask {
	return &ContractTask{
		tool: chain,
		call: call,
	}
}

func (task *ContractTask) Process(tool interface{}) goworker.ToolQuantity {

	if tool1, ok := tool.(*ContractTool); ok {

		tool1.Process(task.call)

	} else {
		fmt.Println("not tool2", task.ToolLabel())
	}
	return goworker.ToolQuantityGood
}

func (task *ContractTask) ToolLabel() string {
	return task.tool
}

type PubsubTask struct {
	message interface{}
	tool    string
	topic   string
}

func NewPubsubTask(tool string, topic string, message interface{}) *PubsubTask {
	return &PubsubTask{
		tool:    tool,
		message: message,
		topic:   topic,
	}
}

func (task *PubsubTask) Process(tool interface{}) goworker.ToolQuantity {

	if tool1, ok := tool.(*ToolExport); ok {

		tool1.Export(task.topic, task.message)

	} else {
		fmt.Println("not tool3", task.ToolLabel())
	}
	return goworker.ToolQuantityGood
}

func (task *PubsubTask) ToolLabel() string {
	return task.tool
}
