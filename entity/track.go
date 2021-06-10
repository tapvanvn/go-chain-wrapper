package entity

type Track struct {
	Address  string   `json:"address"`
	Subjects []string `json:"subjects"`
	AbiName  string   `json:"abi,omitempty"`
}
