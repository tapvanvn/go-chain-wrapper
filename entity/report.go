package entity

type Report struct {
	Name              string   `json:"name"`
	Topic             string   `json:"topic"`
	Subjects          []string `json:"subjects,omitempty"`
	ReportTransaction bool     `json:"-"`
	ReportEvent       bool     `json:"-"`
}
