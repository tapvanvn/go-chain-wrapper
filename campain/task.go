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
