package worker

import (
	"fmt"

	"github.com/tapvanvn/go-jsonrpc-wrapper/command"
)

type Task struct {
	tool    string
	command command.Command
}

func NewTask(chain string, command command.Command) *Task {
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
