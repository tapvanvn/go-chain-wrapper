package export

import (
	"github.com/tapvanvn/gopubsubengine"
	"github.com/tapvanvn/gopubsubengine/wspubsub"
)

type PubsubExporter struct {
	hub        gopubsubengine.Hub
	publisher  map[string]gopubsubengine.Publisher
	subscriber map[string]gopubsubengine.Subscriber
}

func NewPubsubExporter(endpointAddress string) (*PubsubExporter, error) {

	hub, err := wspubsub.NewWSPubSubHub(endpointAddress)

	if err != nil {
		return nil, err
	}
	return &PubsubExporter{hub: hub}, nil
}

func (ex *PubsubExporter) Export(topic string, message interface{}) {

	publisher, ok := ex.publisher[topic]
	if !ok {
		publisher, err := ex.hub.PublishOn(topic)
		if err != nil {
			return
		}
		ex.publisher[topic] = publisher
	}
	publisher.Publish(message)
}
