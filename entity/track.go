package entity

type Track struct {
	Address      string   `json:"address"`
	Subjects     []string `json:"subjects"`
	ContractName string   `json:"contract,omitempty"`
	Reports      []Report `json:"reports,omitempty"`
}
