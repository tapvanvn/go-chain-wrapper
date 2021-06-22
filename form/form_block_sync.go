package form

type FormBlockSync struct {
	ReportName string `json:"report_name"`
	Topic      string `json:"topic"`
	ChainName  string `json:"chain"`
	BlockNum   Param  `json:"block_num"`
}

func (frm *FormBlockSync) IsValid() error {
	return nil
}
