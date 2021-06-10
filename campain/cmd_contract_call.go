package campain

type CmdContractCall struct {
	FuncName string
	Params   []interface{}
	Out      []interface{}
}

func CreateCmdContractCall(funcName string, params []interface{}) *CmdContractCall {
	return &CmdContractCall{
		FuncName: funcName,
		Params:   params,
		Out:      make([]interface{}, 0),
	}
}

/*
var out []interface{}
	err := _Store.contract.Call(opts, &out, "tokensOfOwner", _owner)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err
*/
