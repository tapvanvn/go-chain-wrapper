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
		fmt.Println("not tool", task.ToolLabel())
	}
}

func (task *Task) ToolLabel() string {
	return task.tool
}

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

func (task *ClientTask) Process(tool interface{}) {

	if tool1, ok := tool.(ITool); ok {

		task.command.Do(tool1)

	} else {

		fmt.Println("not tool for client task", task.ToolLabel(), tool)
	}
}

func (task *ClientTask) ToolLabel() string {
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

func (task *ContractTask) Process(tool interface{}) {

	if tool1, ok := tool.(*ContractTool); ok {

		tool1.Process(task.call)

	} else {
		fmt.Println("not tool", task.ToolLabel())
	}
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

func (task *PubsubTask) Process(tool interface{}) {

	if tool1, ok := tool.(*ToolExportPubSub); ok {

		tool1.AddMessage(task.topic, task.message)

	} else {
		fmt.Println("not tool", task.ToolLabel())
	}
}

func (task *PubsubTask) ToolLabel() string {
	return task.tool
}
