package entity

/*
blockHash: "0xc1a1a1b756733afb89298b6a6b496e1be3514b3291b3237885a64e67b7b61d6d",
    blockNumber: "0x7c2a44",
    from: "0xbe807dddb074639cd9fa61b47676c064fc50d62c",
    gas: "0x7fffffffffffffff",
    gasPrice: "0x0",
    hash: "0x1d57b69cb81dc8a12818054233c6457f60a94fc8a623daf0696588a6dd0f1304",
    input: "0xf340fa01000000000000000000000000be807dddb074639cd9fa61b47676c064fc50d62c",
    nonce: "0x1f705",
    r: "0x106bcf346d14df2ed47c89f9f5b431356fb738f0e4389eec530dc73774bc2996",
    s: "0x14ea14ea7963235f3f57d418fac0f92b931b647dbee059bc3f1539d3a7b67e51",
    to: "0x0000000000000000000000000000000000001000",
    transactionIndex: "0x68",
    type: "0x0",
    v: "0x94",
    value: "0x10939a90f499bea"
*/
type Transaction struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	R                string `json:"r"`
	S                string `json:"s"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Type             string `json:"type"`
	V                string `json:"v"`
	Value            string `json:"value"`
}
