package campain

import (
	"fmt"
)

type BlackSmith struct {
	Campain *Campain
}

//Make make tool
func (blacksmith *BlackSmith) Make() interface{} {

	tool, err := NewTool(blacksmith.Campain)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return tool
}
