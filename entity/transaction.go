package entity

type Log struct {
	Topics []string `json:"topic"`
	Data   []byte   `json:"data"`
}

type Transaction struct {
	BlockHash         string      `json:"blockHash"`
	BlockNumber       string      `json:"blockNumber"`
	Gas               string      `json:"gas"`
	GasPrice          string      `json:"gasPrice"`
	Hash              string      `json:"hash"`
	Input             string      `json:"input"`
	Nonce             string      `json:"nonce"`
	R                 string      `json:"r"`
	S                 string      `json:"s"`
	From              string      `json:"from"`
	To                string      `json:"to"`
	TransactionIndex  string      `json:"transactionIndex"`
	Type              string      `json:"type"`
	V                 string      `json:"v"`
	Value             string      `json:"value"`
	OriginTransaction interface{} `json:"-"`
	Success           bool        `json:"Success"`
	Logs              []*Log      `json:"Logs"`
}
