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
	return &PubsubExporter{hub: hub,
		publisher:  make(map[string]gopubsubengine.Publisher),
		subscriber: map[string]gopubsubengine.Subscriber{},
	}, nil
}

func (ex *PubsubExporter) Export(topic string, message interface{}) {

	publisher, ok := ex.publisher[topic]
	if !ok {
		publisher2, err := ex.hub.PublishOn(topic)
		if err != nil {
			return
		}
		ex.publisher[topic] = publisher2
		publisher = publisher2
	}
	publisher.Publish(message)
}
