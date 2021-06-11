package campain

import (
	"github.com/tapvanvn/gopubsubengine"
)

type ToolExportPubSub struct {
	id        int
	topic     string
	publisher gopubsubengine.Publisher
}

func NewExportPubSubTool(hub gopubsubengine.Hub, topic string) (*ToolExportPubSub, error) {

	__tool_id += 1
	tool := &ToolExportPubSub{id: __tool_id,
		topic: topic,
	}
	publisher, err := hub.PublishOn(topic)
	if err != nil {
		return nil, err
	}
	tool.publisher = publisher

	return tool, nil
}

func (tool *ToolExportPubSub) AddMessage(msg interface{}) {

	go tool.publisher.Publish(msg)
}
