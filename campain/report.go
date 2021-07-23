package campain

import (
	"fmt"

	"github.com/tapvanvn/go-jsonrpc-wrapper/entity"
	"github.com/tapvanvn/go-jsonrpc-wrapper/repository"
	"github.com/tapvanvn/godashboard"
)

type ReportEvent struct {
	track  *entity.Track
	events []*entity.Event
}

//ReportLive get report meta for dashboard
func ReportLive() map[string]godashboard.Param {
	report := map[string]godashboard.Param{}
	for _, camp := range __campmap {
		report[fmt.Sprintf("%s_lastest", camp.chainName)] = godashboard.Param{
			Type:  "uint64",
			Value: []byte(fmt.Sprintf("%d", camp.lastBlockNumber)),
		}
		repository.PutLastBlock(camp.chainName, camp.lastBlockNumber)
	}
	return report
}
