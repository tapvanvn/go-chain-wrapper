package dashboard

import (
	"errors"
	"net"
	"strings"
	"time"

	"github.com/tapvanvn/go-jsonrpc-wrapper/entity"
	"github.com/tapvanvn/gopubsubengine"
	"github.com/tapvanvn/gopubsubengine/wspubsub"
)

var __pubsubmap map[string]gopubsubengine.Hub = make(map[string]gopubsubengine.Hub)
var __reporter []DashboardReporter = make([]DashboardReporter, 0)

func AddDashboard(db *entity.Dashboard) error {
	if db.Type == "wspubsub" {
		if _, ok := __pubsubmap[db.ConnectionString]; !ok {

			endpoints := strings.Split(db.ConnectionString, ",")

			if len(endpoints) == 0 {

				return errors.New("connect string not found")
			}
			selectEndpoint := endpoints[0]

			timeout := time.Duration(1 * time.Second)
			for _, endpoint := range endpoints {
				_, err := net.DialTimeout("tcp", endpoint, timeout)
				if err == nil {

					selectEndpoint = endpoint
					break
				}
			}

			hub, err := wspubsub.NewWSPubSubHub(selectEndpoint)

			if err != nil {
				return err
			}

			__pubsubmap[db.ConnectionString] = hub

			reporter := NewPubsubDashboardReporter(hub)

			if reporter == nil {
				return errors.New("create dashboard exporter fail")
			}
			__reporter = append(__reporter, reporter)

		}
		return nil
	}
	return errors.New("dashboard type is not support")
}

func Report(signal *entity.Signal) {
	for _, reporter := range __reporter {
		reporter.Report(signal)
	}
}
