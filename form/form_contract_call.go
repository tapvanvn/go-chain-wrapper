package form

type FormCallContract struct {
	ReportName string   `json:"report_name"` //report name
	Topic      string   `json:"topic"`       //report on topic
	ChainName  string   `json:"chain"`       //chain.contractName
	Name       string   `json:"name"`        //function name on contract
	Params     []*Param `json:"params"`      //input param for function call
	OutType    []string `json:"out_types,omitempty"`
}

func (frm *FormCallContract) IsValid() error {
	return nil
}
