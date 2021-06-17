package campain

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/tapvanvn/go-jsonrpc-wrapper/entity"
	"github.com/tapvanvn/go-jsonrpc-wrapper/filter"
	"github.com/tapvanvn/go-jsonrpc-wrapper/utility"
	"github.com/tapvanvn/gopubsubengine"
	"github.com/tapvanvn/gopubsubengine/wspubsub"
	"github.com/tapvanvn/goworker"
)

type Campain struct {
	mux             sync.Mutex
	timeRange       time.Duration
	isRun           bool
	chainName       string
	ChnTransactions chan entity.Transaction
	ChnBlockNumber  chan uint64
	filters         map[filter.IFilter]*entity.Track
	lastBlockNumber uint64
	abis            map[string]IABI
	exportType      map[string]string
	pubsubHub       map[string]gopubsubengine.Hub
}

func NewCampain(chain string, timeRange time.Duration) *Campain {

	camp := &Campain{
		timeRange:       timeRange,
		isRun:           false,
		chainName:       chain,
		ChnTransactions: make(chan entity.Transaction),
		ChnBlockNumber:  make(chan uint64),
		filters:         make(map[filter.IFilter]*entity.Track),
		lastBlockNumber: 0,
		abis:            map[string]IABI{},
		exportType:      map[string]string{},
		pubsubHub:       make(map[string]gopubsubengine.Hub),
	}

	return camp
}

func (campain *Campain) AddExport(export *entity.Export) error {
	if export.Type == "wspubsub" {
		if _, ok := campain.pubsubHub[export.Name]; !ok {
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
			hub, err := wspubsub.NewWSPubSubHub(selectEndpoint)

			if err != nil {
				return err
			}

			campain.pubsubHub[export.Name] = hub
			campain.exportType[export.Name] = export.Type

			goworker.AddToolWithControl(campain.chainName+"."+export.Name,
				&ExportPubsubBlackSmith{
					Campain:    campain,
					ExportName: export.Name,
				},
				1)
		}
		return nil
	}
	return errors.New("export not supported")
}
func (campain *Campain) LoadContract(contract *entity.Contract) error {
	if contract.AbiName != "" {
		if campain.chainName == "bsc" {
			abiObj, err := NewEthereumABI(contract.AbiName, contract.Address)
			if err != nil {

				return err
			} else {
				abiObj.Info()
				campain.abis[contract.Name] = abiObj
			}
		} else if campain.chainName == "kai" {
			abiObj, err := NewKaiABI(contract.AbiName, contract.Address)
			if err != nil {

				return err
			} else {
				abiObj.Info()
				campain.abis[contract.Name] = abiObj
			}
		}
	}
	return nil
}

func (campain *Campain) Tracking(track entity.Track) error {

	campain.mux.Lock()
	for _, subject := range track.Subjects {

		if subject == "transaction.to" {

			filter := &filter.FilMatchTo{

				Address: track.Address,
			}
			campain.filters[filter] = &track
		}
	}
	for _, report := range track.Reports {
		exportType, ok := campain.exportType[report.Name]
		if !ok {
			return errors.New("export not loaded")
		}
		if exportType != "wspubsub" {
			return errors.New("export is not supported")
		}
	}

	campain.mux.Unlock()
	return nil
}
func (campain *Campain) Report(reportName string, topic string, message interface{}) {
	exportType, ok := campain.exportType[reportName]
	if !ok || exportType != "wspubsub" {
		return
	}
	toolName := campain.chainName + "." + reportName
	go goworker.AddTask(NewPubsubTask(toolName, topic, message))
}
func (campain *Campain) report(report *entity.Report, message interface{}) {
	exportType, ok := campain.exportType[report.Name]
	if !ok || exportType != "wspubsub" {
		return
	}
	toolName := campain.chainName + "." + report.Name
	go goworker.AddTask(NewPubsubTask(toolName, report.Topic, message))
}

func (campain *Campain) processBlockNumber() {
	for {
		blockNumber := <-campain.ChnBlockNumber
		if blockNumber <= campain.lastBlockNumber {
			continue
		}
		if campain.lastBlockNumber == 0 {

			campain.lastBlockNumber = blockNumber
		}
		for i := campain.lastBlockNumber + 1; i <= blockNumber; i++ {

			fmt.Println(campain.chainName, "block:", i)
			cmd := CreateCmdTransactionsOfBlock(i)
			cmd.Init()
			task := NewClientTask(campain.chainName, cmd)
			go goworker.AddTask(task)
		}
		campain.lastBlockNumber = blockNumber
	}
}
func (campain *Campain) processTransaction() {
	for {
		trans := <-campain.ChnTransactions
		campain.mux.Lock()
		//isFilted := true
		for filter, track := range campain.filters {
			if filter.Match(&trans) {
				event := map[string]interface{}{}
				//isFilted = false
				fmt.Println(campain.chainName, "found transaction:", trans.Hash)
				fmt.Println("\tfrom:", trans.From)
				fmt.Println("\tto:", trans.To)
				event["from"] = trans.From
				event["to"] = trans.To
				event["hash"] = trans.Hash
				if track.ContractName != "" {
					if abiObj, ok := campain.abis[track.ContractName]; ok {
						method, args, err := abiObj.GetMethod(trans.Input)
						event["method"] = method
						event["args"] = args
						if err == nil {
							fmt.Println("\tmethod:", method, args)
						}
					}
				}
				for _, report := range track.Reports {

					campain.report(&report, event)
				}
			}
		}

		campain.mux.Unlock()
	}
}

func (campain *Campain) run() {

	cmd := CreateCmdLatestBlockNumber()
	cmd.Init()
	task := NewClientTask(campain.chainName, cmd)
	goworker.AddTask(task)
}

func (campain *Campain) Run() {

	if campain.isRun {

		return
	}
	go campain.processBlockNumber()
	go campain.processTransaction()
	campain.isRun = true
	utility.Schedule(campain.run, campain.timeRange)
}
