package campain

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/tapvanvn/go-chain-wrapper/entity"
	"github.com/tapvanvn/go-chain-wrapper/export"
	"github.com/tapvanvn/go-chain-wrapper/filter"
	"github.com/tapvanvn/go-chain-wrapper/repository"
	"github.com/tapvanvn/go-chain-wrapper/utility"
	goworker "github.com/tapvanvn/goworker/v2"
)

var __campmap map[string]*Campain = map[string]*Campain{}

func GetCampain(chainName string) *Campain {
	if camp, ok := __campmap[chainName]; ok {
		return camp
	}
	return nil
}

var windowSize = 10

type Endpoint string
type ContractAddress string
type ContractName string

type Window struct {
	mux         sync.Mutex
	blockStatus []bool
}

func (w *Window) finishBlock(index int) bool {
	if index < 0 || index > windowSize {
		panic("Block Index invalid")
	}
	w.mux.Lock()
	defer w.mux.Unlock()
	w.blockStatus[index] = true

	for _, status := range w.blockStatus {
		if !status {
			return false
		}
	}
	return true
}

func newWindowSync() *Window {
	return &Window{
		blockStatus: make([]bool, windowSize),
	}
}

type Campain struct {
	mux sync.Mutex

	isRun              bool
	isAutoMine         bool
	chainName          string
	chnTransactions    chan []*entity.Transaction
	chnBlockNumber     chan uint64
	chnEvent           chan *ReportEvent
	endpoints          []Endpoint
	filters            map[filter.IFilter]*entity.Track
	lastBlockNumber    uint64
	miningBlockNumber  uint64
	abis               map[ContractName]IABI
	contractAddress    map[ContractName]ContractAddress
	directContractTool map[ContractName]*ContractTool
	quantityControl    map[string]int
	originEndpoint     map[string]Endpoint
	toolLabels         []string
	//block sync controll
	syncWindow map[uint64]*Window
	lowWindow  uint64 //lowest waiting syncing window
	highWindow uint64 //highest waiting syncing window
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
		chnTransactions:    make(chan []*entity.Transaction),
		chnBlockNumber:     make(chan uint64),
		chnEvent:           make(chan *ReportEvent),
		filters:            make(map[filter.IFilter]*entity.Track),
		endpoints:          make([]Endpoint, 0),
		lastBlockNumber:    chain.MineFromBlock,
		miningBlockNumber:  0,
		abis:               map[ContractName]IABI{},
		contractAddress:    map[ContractName]ContractAddress{},
		directContractTool: map[ContractName]*ContractTool{},
		quantityControl:    map[string]int{},
		originEndpoint:     map[string]Endpoint{},
		toolLabels:         make([]string, 0),

		syncWindow: map[uint64]*Window{},
	}
	__campmap[chain.Name] = camp

	camp.toolLabels = append(camp.toolLabels, chain.Name)

	for origin, endpoint := range chain.Endpoints {

		camp.AddEndpoint(origin, Endpoint(endpoint))
	}

	goworker.AddToolWithControl(chain.Name, &ClientBlackSmith{
		Campain: camp,
	}, chain.NumWorker)

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
	fmt.Println("add endpoint", origin, endpoint)
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
	fmt.Println("enpoint not found", origin, campain.originEndpoint)
	return ""
}

func (campain *Campain) EndpointQuantityReport(origin string, meta interface{}, failCount int) {
	fmt.Println("fail count", origin, failCount)
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
	name := ContractName(contract.Name)
	address := ContractAddress(contract.Address)
	if contract.AbiName != "" {
		if campain.chainName == "bsc" {
			abiObj, err := NewEthereumABI(contract.AbiName, contract.Address)
			if err != nil {

				return err
			} else {
				//abiObj.Info()
				campain.abis[name] = abiObj
				campain.contractAddress[name] = address

			}
		} else if campain.chainName == "kai" {
			abiObj, err := NewKaiABI(contract.AbiName, address)
			if err != nil {

				return err
			} else {
				//abiObj.Info()
				campain.abis[name] = abiObj
				campain.contractAddress[name] = address
			}
		}

		contractTool, err := NewContractTool(campain, name, campain.endpoints[0])
		if err == nil {
			campain.directContractTool[name] = contractTool
		} else {
			fmt.Println("", err)
		}
	}

	return nil
}

func (campain *Campain) SetupContractWorker(contract *entity.Contract, numWorker int) {

	labelContract := fmt.Sprintf("%s.%s", campain.chainName, contract.Name)
	labelTrans := fmt.Sprintf("%s.%s.trans", campain.chainName, contract.Name)
	name := ContractName(contract.Name)

	goworker.AddToolWithControl(labelContract, &ContractBlackSmith{
		Campain:      campain,
		ContractName: name,
	}, numWorker)

	goworker.AddToolWithControl(labelTrans, &TransactionBlackSmith{
		Campain:      campain,
		ContractName: name,
	}, numWorker)

	campain.toolLabels = append(campain.toolLabels, labelTrans)
	campain.toolLabels = append(campain.toolLabels, labelContract)
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
		latestBlockNumber := <-campain.chnBlockNumber

		fmt.Println("process block", latestBlockNumber)
		if latestBlockNumber < campain.lastBlockNumber {

			continue
		}

		campain.lastBlockNumber = latestBlockNumber

		latestWindow := latestBlockNumber / uint64(windowSize)

		if latestWindow <= campain.highWindow {
			//window had been planned to sync
			continue
		}
		if campain.lowWindow == 0 {

			campain.lowWindow = latestWindow
		}

		if latestWindow-campain.lowWindow > 10 {

			continue
		}

		beginBlock := latestWindow * uint64(windowSize)
		toBlock := (latestWindow + 1) * uint64(windowSize)

		campain.highWindow = latestWindow
		campain.syncWindow[latestWindow] = newWindowSync()

		fmt.Println("plan sync window", latestWindow, beginBlock, toBlock)
		//plan to sync the window
		for i := beginBlock; i < toBlock; i++ {

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
		if len(trans) == 0 {
			continue
		}
		blockNumber, err := strconv.ParseUint(trans[0].BlockNumber, 10, 64)

		if err != nil {
			panic(err)
		}

		window := blockNumber / uint64(windowSize)
		beginBlock := window * uint64(windowSize)
		index := int(blockNumber - beginBlock)

		if syncWindow, ok := campain.syncWindow[window]; ok {

			if syncWindow.finishBlock(index) {

				delete(campain.syncWindow, window)
				fmt.Println("finish window", window)
				low := campain.highWindow
				for windowIndex, _ := range campain.syncWindow {
					if windowIndex < low {
						low = windowIndex
					}
				}
				campain.lowWindow = low
			}
		}

		for _, tran := range trans {

			//fmt.Println("transaction ", len(campain.filters), trans.BlockNumber)
			for filter, track := range campain.filters {
				if filter.Match(tran) {
					event := map[string]interface{}{}

					fmt.Println(campain.chainName, "found transaction:", tran.Hash)
					fmt.Println("\tfrom:", tran.From)
					fmt.Println("\tto:", tran.To)
					event["from"] = tran.From
					event["to"] = tran.To
					event["hash"] = tran.Hash
					if track.ContractName != "" {

						if abiObj, ok := campain.abis[ContractName(track.ContractName)]; ok {

							method, args, err := abiObj.GetMethod(tran.Input)
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
							transaction: tran,
							track:       track,
						}
						go goworker.AddTask(&task)
					}
				}
			}
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
	for _, toolLabel := range campain.toolLabels {

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

func (campain *Campain) GetDirectContractTool(contractName ContractName) *ContractTool {
	if tool, ok := campain.directContractTool[contractName]; ok {
		return tool
	}
	return nil
}
