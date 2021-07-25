package export

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/tapvanvn/go-chain-wrapper/entity"
	"github.com/tapvanvn/gopubsubengine"
)

type Exporter interface {
	Export(topic string, message interface{})
}

func GetExport(name string) Exporter {
	if exp, ok := __exportmap[name]; ok {
		return exp
	}
	return nil
}

func GetPubSubHub(name string) gopubsubengine.Hub {
	if pubsubhub, ok := __pubsubmap[name]; ok {
		return pubsubhub
	}
	return nil
}

var __pubsubmap map[string]gopubsubengine.Hub = make(map[string]gopubsubengine.Hub)
var __exporttype map[string]string = map[string]string{}
var __exportmap map[string]Exporter = map[string]Exporter{}

func AddExport(export *entity.Export) error {

	if export.Type == "wspubsub" {

		if _, ok := __pubsubmap[export.Name]; !ok {

			endpoints := strings.Split(export.ConnectionString, ",")

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
			fmt.Println(selectEndpoint)
			ex, err := NewPubsubExporter(selectEndpoint)

			if err != nil {
				return err
			}

			__pubsubmap[export.Name] = ex.hub
			__exporttype[export.Name] = export.Type
			__exportmap[export.Name] = ex

		} else {
			panic("cannot found pubsub")
		}
		return nil
	}
	return errors.New("export not supported")
}
