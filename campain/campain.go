package campain

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/tapvanvn/go-jsonrpc-wrapper/entity"
	"github.com/tapvanvn/go-jsonrpc-wrapper/export"
	"github.com/tapvanvn/go-jsonrpc-wrapper/filter"
	"github.com/tapvanvn/go-jsonrpc-wrapper/repository"
	"github.com/tapvanvn/go-jsonrpc-wrapper/utility"
	goworker "github.com/tapvanvn/goworker/v2"
)

var __campmap map[string]*Campain = map[string]*Campain{}

func GetCampain(chainName string) *Campain {
	if camp, ok := __campmap[chainName]; ok {
		return camp
	}
	return nil
}

type Endpoint string
type Campain struct {
	mux sync.Mutex

	isRun              bool
	isAutoMine         bool
	chainName          string
	chnTransactions    chan entity.Transaction
	chnBlockNumber     chan uint64
	chnEvent           chan *ReportEvent
	endpoints          []Endpoint
	filters            map[filter.IFilter]*entity.Track
	lastBlockNumber    uint64
	miningBlockNumber  uint64
	abis               map[string]IABI
	contractAddress    map[string]string
	directContractTool map[string]*ContractTool
	quantityControl    map[string]int
	originEndpoint     map[string]Endpoint
	workerLabels       []string
}

func AddCampain(chain *entity.Chain) *Campain {

	camp := GetCampain(chain.Name)
	if camp != nil {
		return camp
	}
	camp = &Campain{
		isRun:              false,
		chainName:          chain.Name,
		isAutoMine:         chain.AutoMine,
		chnTransactions:    make(chan entity.Transaction),
		chnBlockNumber:     make(chan uint64),
		chnEvent:           make(chan *ReportEvent),
		filters:            make(map[filter.IFilter]*entity.Track),
		endpoints:          make([]Endpoint, 0),
		lastBlockNumber:    chain.MineFromBlock,
		miningBlockNumber:  0,
		abis:               map[string]IABI{},
		contractAddress:    map[string]string{},
		directContractTool: map[string]*ContractTool{},
		quantityControl:    map[string]int{},
		originEndpoint:     map[string]Endpoint{},
		workerLabels:       make([]string, 0),
	}
	__campmap[chain.Name] = camp

	goworker.AddToolWithControl(chain.Name, &ClientBlackSmith{
		Campain: camp,
	}, chain.NumWorker)

	camp.workerLabels = append(camp.workerLabels, chain.Name)

	for origin, endpoint := range chain.Endpoints {

		camp.AddEndpoint(origin, Endpoint(endpoint))
	}

	for _, contract := range chain.Contracts {
		err := camp.LoadContract(&contract)
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}

		camp.SetupContractWorker(&contract, chain.NumWorker)
	}

	for _, track := range chain.Tracking {

		camp.Tracking(track)
	}

	return camp
}

func (campain *Campain) AddEndpoint(origin string, endpoint Endpoint) {

	if len(campain.originEndpoint) == 0 {
		origin = "default"
	}
	campain.originEndpoint[origin] = endpoint
	campain.endpoints = append(campain.endpoints, endpoint)

}

func (campain *Campain) GetEndpoint(origin string) Endpoint {
	if endpoint, ok := campain.originEndpoint[origin]; ok {
		return endpoint
	}
	return ""
}

func (campain *Campain) EndpointQuantityReport(origin string, meta interface{}, failCount int) {
	campain.quantityControl[origin] = failCount

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
				//abiObj.Info()
				campain.abis[contract.Name] = abiObj
				campain.contractAddress[contract.Name] = contract.Address

			}
		} else if campain.chainName == "kai" {
			abiObj, err := NewKaiABI(contract.AbiName, contract.Address)
			if err != nil {

				return err
			} else {
				//abiObj.Info()
				campain.abis[contract.Name] = abiObj
				campain.contractAddress[contract.Name] = contract.Address
			}
		}

		contractTool, err := NewContractTool(campain, contract.Name, string(campain.endpoints[0]))
		if err == nil {
			campain.directContractTool[contract.Name] = contractTool
		} else {
			fmt.Println("", err)
		}
	}

	return nil
}

func (campain *Campain) SetupContractWorker(contract *entity.Contract, numWorker int) {

	labelContract := fmt.Sprintf("%s.%s", campain.chainName, contract.Name)
	labelTrans := fmt.Sprintf("%s.%s.trans", campain.chainName, contract.Name)

	goworker.AddToolWithControl(labelContract, &ContractBlackSmith{
		Campain:      campain,
		ContractName: contract.Name,
	}, numWorker)

	goworker.AddToolWithControl(labelTrans, &TransactionBlackSmith{
		Campain:      campain,
		ContractName: contract.Name,
	}, numWorker)

	campain.workerLabels = append(campain.workerLabels, labelTrans)
	campain.workerLabels = append(campain.workerLabels, labelContract)
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

		for _, sub := range report.Subjects {
			if sub == "transaction" {
				report.ReportTransaction = true
			} else if sub == "event" {
				report.ReportEvent = true
			}
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
		blockNumber := <-campain.chnBlockNumber
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
		trans := <-campain.chnTransactions

		//fmt.Println("transaction ", len(campain.filters), trans.BlockNumber)
		for filter, track := range campain.filters {
			if filter.Match(&trans) {
				event := map[string]interface{}{}

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
				fmt.Println("")
				for _, report := range track.Reports {
					if report.ReportTransaction {
						campain.report(report, event)
					}
				}

				if entity.AnyReportEvent(track.Reports) {
					//TODO: add switch by config
					task := TransactionTask{
						tool:        campain.chainName + "." + track.ContractName + ".trans",
						transaction: &trans,
						track:       track,
					}
					go goworker.AddTask(&task)
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
	}
}

func (campain *Campain) processEvent() {
	for {
		reportEvent := <-campain.chnEvent

		for _, report := range reportEvent.track.Reports {

			for _, event := range reportEvent.events {
				campain.report(report, event)
			}
		}
	}
}

func (campain *Campain) Run() {
	if campain.isRun {

		return
	}

	cacheLastBlock := repository.GetLastBlock(campain.chainName)
	//only fetch cache if we not set the init value
	if cacheLastBlock > campain.lastBlockNumber {
		campain.lastBlockNumber = cacheLastBlock
	}
	if campain.isAutoMine {

		fmt.Println("auto mine")
		utility.Schedule(campain.MineBlock, time.Second)

	} else {

		fmt.Println("not auto mine")
	}
	for _, toolLabel := range campain.workerLabels {

		for origin, _ := range campain.originEndpoint {

			if origin != "default" {

				goworker.AddOrigin(toolLabel, origin, nil, 100)
			}
		}
		goworker.SetQuantityReporter(toolLabel, campain.EndpointQuantityReport)
	}

	go campain.processBlockNumber()
	go campain.processTransaction()
	go campain.processEvent()

	campain.isRun = true
}

func (campain *Campain) GetDirectContractTool(contractName string) *ContractTool {
	if tool, ok := campain.directContractTool[contractName]; ok {
		return tool
	}
	return nil
}
