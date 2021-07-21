package entity

type Event struct {
	Name      string           `json:"name"`
	Arguments map[string]Param `json:"arguments"`
}
