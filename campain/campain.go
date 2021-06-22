package campain

import (
	"errors"
	"fmt"
	"sync"

	"github.com/tapvanvn/go-jsonrpc-wrapper/entity"
	"github.com/tapvanvn/go-jsonrpc-wrapper/filter"
	"github.com/tapvanvn/goworker"
)

var __campmap map[string]*Campain = map[string]*Campain{}

func GetCampain(chainName string) *Campain {
	if camp, ok := __campmap[chainName]; ok {
		return camp
	}
	return nil
}

type Campain struct {
	mux sync.Mutex
	//timeRange       time.Duration
	isRun           bool
	chainName       string
	ChnTransactions chan entity.Transaction
	ChnBlockNumber  chan uint64
	filters         map[filter.IFilter]*entity.Track
	lastBlockNumber uint64
	abis            map[string]IABI
	contractAddress map[string]string
	exportType      map[string]string
	//pubsubHub       map[string]gopubsubengine.Hub
}

func AddCampain(chain string) *Campain {

	camp := GetCampain(chain)
	if camp != nil {
		return camp
	}
	camp = &Campain{
		//timeRange:       timeRange,
		isRun:           false,
		chainName:       chain,
		ChnTransactions: make(chan entity.Transaction),
		ChnBlockNumber:  make(chan uint64),
		filters:         make(map[filter.IFilter]*entity.Track),
		lastBlockNumber: 0,
		abis:            map[string]IABI{},
		contractAddress: map[string]string{},
		exportType:      map[string]string{},
		//pubsubHub:       make(map[string]gopubsubengine.Hub),
	}
	__campmap[chain] = camp

	return camp
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
				campain.contractAddress[contract.Name] = contract.Address
			}
		} else if campain.chainName == "kai" {
			abiObj, err := NewKaiABI(contract.AbiName, contract.Address)
			if err != nil {

				return err
			} else {
				abiObj.Info()
				campain.abis[contract.Name] = abiObj
				campain.contractAddress[contract.Name] = contract.Address
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
	toolName := reportName
	go goworker.AddTask(NewPubsubTask(toolName, topic, message))
}

func (campain *Campain) report(report *entity.Report, message interface{}) {
	exportType, ok := campain.exportType[report.Name]
	if !ok || exportType != "wspubsub" {
		return
	}
	toolName := report.Name
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
	//#utility.Schedule(campain.run, campain.timeRange)
}
