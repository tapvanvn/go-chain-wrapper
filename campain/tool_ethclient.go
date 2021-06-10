package campain

import "time"

type EthClientTool struct {
	id             int
	ready          bool
	commands       chan Command
	waitingCommand Command
	campain        *Campain
}

func NewEthClientTool(campain *Campain) (*EthClientTool, error) {

	__tool_id += 1
	tool := &EthClientTool{id: __tool_id,
		ready:          false,
		commands:       make(chan Command),
		waitingCommand: nil,
		campain:        campain,
	}

	go tool.process()

	return tool, nil
}

func (tool *EthClientTool) process() {

	for {

		if tool.waitingCommand != nil {
			time.Sleep(time.Microsecond * 20)
			continue
		}
		cmd, ok := <-tool.commands

		if !ok {
			break
		}

		tool.waitingCommand = cmd
		//process command
	}
}
