package campain

import (
	"fmt"

	"github.com/tapvanvn/gopubsubengine"
)

var __publishers map[string]gopubsubengine.Publisher = make(map[string]gopubsubengine.Publisher)

type ToolExportPubSub struct {
	id  int
	hub gopubsubengine.Hub
}

func NewExportPubSubTool(hub gopubsubengine.Hub) (*ToolExportPubSub, error) {

	__tool_id += 1
	tool := &ToolExportPubSub{id: __tool_id, hub: hub}

	return tool, nil
}

func (tool *ToolExportPubSub) AddMessage(topic string, msg interface{}) {

	if publisher, ok := __publishers[topic]; ok {
		go publisher.Publish(msg)
		return
	}
	publisher, err := tool.hub.PublishOn(topic)
	if err != nil {
		fmt.Println(err)
		return
	}
	__publishers[topic] = publisher
	go publisher.Publish(msg)
}
