package entity

type Event struct {
	Name      string            `json:"name"`
	Arguments map[string]string `json:"arguments"`
}
