package campain

import (
	"fmt"
)

type Task struct {
	tool    string
	command Command
}

func NewTask(chain string, command Command) *Task {
	return &Task{
		tool:    chain,
		command: command,
	}
}

func (task *Task) Process(tool interface{}) {

	if tool1, ok := tool.(*Tool); ok {

		tool1.AddCommand(task.command)

	} else {
		fmt.Println("not tool")
	}
}

func (task *Task) ToolLabel() string {
	return task.tool
}

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

func (task *ContractTask) Process(tool interface{}) {

	if tool1, ok := tool.(*ContractTool); ok {

		tool1.AddCall(task.call)

	} else {
		fmt.Println("not tool")
	}
}

func (task *ContractTask) ToolLabel() string {
	return task.tool
}

type PubsubTask struct {
	message interface{}
	tool    string
}

func NewPubsubTask(tool string, message interface{}) *PubsubTask {
	return &PubsubTask{
		tool:    tool,
		message: message,
	}
}

func (task *PubsubTask) Process(tool interface{}) {

	if tool1, ok := tool.(*ToolExportPubSub); ok {

		tool1.AddMessage(task.message)

	} else {
		fmt.Println("not tool")
	}
}

func (task *PubsubTask) ToolLabel() string {
	return task.tool
}
