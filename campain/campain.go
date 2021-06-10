package campain

import (
	"fmt"
	"sync"
	"time"

	"github.com/tapvanvn/go-jsonrpc-wrapper/entity"
	"github.com/tapvanvn/go-jsonrpc-wrapper/filter"
	"github.com/tapvanvn/go-jsonrpc-wrapper/utility"

	"github.com/tapvanvn/goworker"
)

type Campain struct {
	mux             sync.Mutex
	timeRange       time.Duration
	isRun           bool
	chainName       string
	ChnTransactions chan entity.Transaction
	ChnBlockNumber  chan int64
	filters         map[filter.IFilter]*entity.Track
	lastBlockNumber int64
	abis            map[string]IABI
}

func NewCampain(chain string, timeRange time.Duration) *Campain {

	camp := &Campain{
		timeRange:       timeRange,
		isRun:           false,
		chainName:       chain,
		ChnTransactions: make(chan entity.Transaction),
		ChnBlockNumber:  make(chan int64),
		filters:         make(map[filter.IFilter]*entity.Track),
		lastBlockNumber: 0,
		abis:            map[string]IABI{},
	}

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

	campain.mux.Unlock()
	return nil
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

			fmt.Println("block:", i)
			cmd := CreateCmdTransactionsOfBlock(i)
			cmd.Init()
			task := NewTask(campain.chainName, cmd)
			goworker.AddTask(task)
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
				//isFilted = false
				fmt.Println("found transaction:", trans.Hash)
				fmt.Println("\tfrom:", trans.From)
				fmt.Println("\tto:", trans.To)
				if track.ContractName != "" {
					if abiObj, ok := campain.abis[track.ContractName]; ok {
						method, args, err := abiObj.GetMethod(trans.Input)
						if err == nil {
							fmt.Println("\tmethod:", method, args)
						}
					}
				}
			}
		}

		campain.mux.Unlock()
	}
}

func (campain *Campain) run() {

	cmd := CreateCmdLatestBlockNumber()
	cmd.Init()
	task := NewTask(campain.chainName, cmd)
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
