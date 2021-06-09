package worker

import "fmt"

type BSCBlacksmith struct {
}

//Make make tool
func (blacksmith *BSCBlacksmith) Make() interface{} {

	fmt.Println("make tool")
	tool, err := NewTool("bsc")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return tool
}
