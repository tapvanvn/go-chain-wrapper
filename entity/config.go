package entity

type Chain struct {
	Name      string     `json:"name"`
	NumWorker int        `json:"num_worker"`
	Endpoints []string   `json:"endpoints"`
	Tracking  []Track    `json:"tracking"`
	Contracts []Contract `json:"contracts"`
}
type Config struct {
	Chains []Chain `json:"chains"`
}

func (config *Config) GetNumWorker() int {
	number := 0
	for _, chain := range config.Chains {
		number += chain.NumWorker
	}
	return number
}

var DefaultConfig Config = Config{Chains: []Chain{{Name: "bsc", NumWorker: 1, Endpoints: []string{"https://bsc-dataseed1.binance.org"}}}}
