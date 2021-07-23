package campain

import (
	"github.com/tapvanvn/go-jsonrpc-wrapper/export"
)

type ToolExport struct {
	id int
	ex export.Exporter
}

func NewExportTool(ex export.Exporter) (*ToolExport, error) {

	tool := &ToolExport{
		id: newToolID(),
		ex: ex}

	return tool, nil
}

func (tool *ToolExport) Export(topic string, msg interface{}) {

	go tool.ex.Export(topic, msg)
}
