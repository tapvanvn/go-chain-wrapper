package campain

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/tapvanvn/go-jsonrpc-wrapper/entity"
	"github.com/tapvanvn/go-jsonrpc-wrapper/export"
	"github.com/tapvanvn/go-jsonrpc-wrapper/filter"
	"github.com/tapvanvn/go-jsonrpc-wrapper/repository"
	"github.com/tapvanvn/go-jsonrpc-wrapper/utility"
	"github.com/tapvanvn/godashboard"
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

	isRun              bool
	IsAutoMine         bool
	chainName          string
	ChnTransactions    chan entity.Transaction
	ChnBlockNumber     chan uint64
	Endpoints          []string
	filters            map[filter.IFilter]*entity.Track
	lastBlockNumber    uint64
	miningBlockNumber  uint64
	abis               map[string]IABI
	contractAddress    map[string]string
	exportType         map[string]string
	DirectContractTool map[string]*ContractTool
}

func AddCampain(chain *entity.Chain) *Campain {

	camp := GetCampain(chain.Name)
	if camp != nil {
		return camp
	}
	camp = &Campain{
		isRun:              false,
		chainName:          chain.Name,
		IsAutoMine:         chain.AutoMine,
		ChnTransactions:    make(chan entity.Transaction),
		ChnBlockNumber:     make(chan uint64),
		filters:            make(map[filter.IFilter]*entity.Track),
		Endpoints:          make([]string, 0),
		lastBlockNumber:    chain.MineFromBlock,
		miningBlockNumber:  0,
		abis:               map[string]IABI{},
		contractAddress:    map[string]string{},
		exportType:         map[string]string{},
		DirectContractTool: map[string]*ContractTool{},
	}
	__campmap[chain.Name] = camp

	return camp
}
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

//this function will be schedule if auto run config set to true
func (campain *Campain) MineBlock() {

	cmd := &CmdGetLatestBlockNumber{}
	cmd.Init()
	task := NewClientTask(campain.chainName, cmd)
	goworker.AddTask(task)
}

//LoadContract load a contract
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

		contractTool, err := NewContractTool(campain, contract.Name, campain.Endpoints[0])
		if err == nil {
			campain.DirectContractTool[contract.Name] = contractTool
		} else {
			fmt.Println("", err)
		}
	}
	return nil
}

//Tracking tracking transaction
func (campain *Campain) Tracking(track entity.Track) error {

	campain.mux.Lock()
	defer campain.mux.Unlock()
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

	return nil
}

func (campain *Campain) Report(reportName string, topic string, message interface{}) {
	exporter := export.GetExport(reportName)
	if exporter == nil {
		fmt.Println("report not found", reportName, topic)
		return
	}
	go exporter.Export(topic, message)

}

func (campain *Campain) report(report *entity.Report, message interface{}) {
	exporter := export.GetExport(report.Name)
	if exporter == nil {
		fmt.Println("report not found", report.Name, report.Topic)
		return
	}
	go exporter.Export(report.Topic, message)
}

func (campain *Campain) processBlockNumber() {
	for {
		blockNumber := <-campain.ChnBlockNumber
		if blockNumber <= campain.lastBlockNumber {
			continue
		}

		//we need to mine several block at a time
		for i := campain.lastBlockNumber + 1; i <= campain.lastBlockNumber+100; i++ {
			if i > blockNumber {
				break
			}
			if i < campain.miningBlockNumber {

				continue
			}
			campain.miningBlockNumber = i
			cmd := CreateCmdTransactionsOfBlock(i)
			cmd.Init()
			task := NewClientTask(campain.chainName, cmd)
			go goworker.AddTask(task)
		}
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

		if transBlock, err := strconv.ParseUint(trans.BlockNumber, 10, 64); err == nil {
			if transBlock > campain.lastBlockNumber {
				campain.lastBlockNumber = transBlock
			}
		} else {
			fmt.Println(err)
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
	//only fetch cache if we not set the init value
	if campain.lastBlockNumber == 0 {
		campain.lastBlockNumber = repository.GetLastBlock(campain.chainName)
	}
	if campain.IsAutoMine {
		fmt.Println("auto mine")
		utility.Schedule(campain.MineBlock, time.Second)
	} else {
		fmt.Println("not auto mine")
	}
	go campain.processBlockNumber()
	go campain.processTransaction()
	campain.isRun = true
}
