package form

type Param struct {
	Type  string `json:"type"`
	Value []byte `json:"value"`
}
type FormCallContract struct {
	ReportName string   `json:"report_name"`
	Topic      string   `json:"topic"`
	ChainName  string   `json:"chain"`
	Name       string   `json:"name"`
	Params     []*Param `json:"params"`
	OutType    []string `json:"out_types"`
}

func (frm *FormCallContract) IsValid() error {
	return nil
}
