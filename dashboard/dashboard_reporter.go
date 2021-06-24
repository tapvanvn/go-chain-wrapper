package dashboard

import (
	"github.com/tapvanvn/go-jsonrpc-wrapper/entity"
	"github.com/tapvanvn/gopubsubengine"
)

type DashboardReporter interface {
	Report(signal *entity.Signal)
}

func NewPubsubDashboardReporter(hub gopubsubengine.Hub) *PubsubDashboardReporter {
	publisher, err := hub.PublishOn("dashboard")
	if err != nil {
		return nil
	}
	return &PubsubDashboardReporter{
		publisher: publisher,
	}
}

type PubsubDashboardReporter struct {
	publisher gopubsubengine.Publisher
}

func (dbr *PubsubDashboardReporter) Report(signal *entity.Signal) {
	dbr.publisher.Publish(signal)
}
