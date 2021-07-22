package export

import (
	"sync"

	"github.com/tapvanvn/gopubsubengine"
	"github.com/tapvanvn/gopubsubengine/wspubsub"
)

type PubsubExporter struct {
	pubMux     sync.Mutex
	subMux     sync.Mutex
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
	ex.pubMux.Lock()
	publisher, ok := ex.publisher[topic]
	ex.pubMux.Unlock()
	if !ok {
		publisher2, err := ex.hub.PublishOn(topic)
		if err != nil {
			return
		}
		ex.pubMux.Lock()
		ex.publisher[topic] = publisher2
		ex.pubMux.Unlock()
		publisher = publisher2
	}
	publisher.Publish(message)
}
