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
	filters         []filter.IFilter
	lastBlockNumber int64
}

func NewCampain(chain string, timeRange time.Duration) *Campain {

	return &Campain{
		timeRange:       timeRange,
		isRun:           false,
		chainName:       chain,
		ChnTransactions: make(chan entity.Transaction),
		ChnBlockNumber:  make(chan int64),
		filters:         make([]filter.IFilter, 0),
		lastBlockNumber: 0,
	}
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
		isFilted := true
		for _, filter := range campain.filters {
			if filter.Match(&trans) {
				isFilted = false
			}
		}
		if !isFilted {
			fmt.Println("found transaction:", trans.BlockHash)
		}
		campain.mux.Unlock()
	}
}

func (campain *Campain) AddFilter(filter filter.IFilter) {
	campain.mux.Lock()
	campain.filters = append(campain.filters, filter)
	campain.mux.Unlock()
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
