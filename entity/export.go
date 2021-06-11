package entity

type Export struct {
	Name             string `json:"name"`
	Type             string `json:"type"`
	ConnectionString string `json:"connection_string"`
}
