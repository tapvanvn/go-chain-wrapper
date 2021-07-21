package entity

import "github.com/tapvanvn/godashboard"

type Chain struct {
	Name          string     `json:"name"`
	NumWorker     int        `json:"num_worker"`
	Endpoints     []string   `json:"endpoints"`
	Tracking      []Track    `json:"tracking"`
	Contracts     []Contract `json:"contracts"`
	AutoMine      bool       `json:"auto_mine"`
	MineFromBlock uint64     `json:"mine_from_block"`
}
type Config struct {
	Dashboards []godashboard.Dashboard `json:"dashboards, omitempty"`
	Exports    []Export                `json:"exports,omitempty"`
	Chains     []Chain                 `json:"chains,omitempty"`
}

func (config *Config) GetNumWorker() int {
	number := 0
	for _, chain := range config.Chains {
		number += chain.NumWorker
		number += chain.NumWorker * len(chain.Contracts) * 2
	}
	number = number + 1
	return number
}

var DefaultConfig Config = Config{Chains: []Chain{{Name: "bsc", NumWorker: 1, Endpoints: []string{"https://bsc-dataseed1.binance.org"}}}}
