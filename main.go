package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/tapvanvn/go-jsonrpc-wrapper/campain"
	"github.com/tapvanvn/godashboard"

	"github.com/tapvanvn/go-jsonrpc-wrapper/entity"
	"github.com/tapvanvn/go-jsonrpc-wrapper/export"
	"github.com/tapvanvn/go-jsonrpc-wrapper/form"
	"github.com/tapvanvn/go-jsonrpc-wrapper/route"
	"github.com/tapvanvn/go-jsonrpc-wrapper/system"
	"github.com/tapvanvn/go-jsonrpc-wrapper/utility"
	"github.com/tapvanvn/gopubsubengine"
	"github.com/tapvanvn/gorouter/v2"
	"github.com/tapvanvn/goworker"
)

type Handles []gorouter.RouteHandle
type Endpoint gorouter.EndpointDefine

var ContractCallSubcriber map[string]gopubsubengine.Subscriber = make(map[string]gopubsubengine.Subscriber)
var SyncBlockSubcriber map[string]gopubsubengine.Subscriber = make(map[string]gopubsubengine.Subscriber)

func OnContractCall(message string) {
	frm := &form.FormCallContract{}
	err := json.Unmarshal([]byte(message), frm)
	if err != nil {
		return
	}
	err = frm.IsValid()
	if err != nil {
		return
	}
	inputs := make([]interface{}, 0)
	for _, param := range frm.Params {
		if param.Type == "uint256" {
			input := &big.Int{}
			err := json.Unmarshal(param.Value, input)
			if err != nil {
				fmt.Println("parse param error", err)

				return
			}
			inputs = append(inputs, input)
		}
	}

	call := campain.CreateContractCall(frm.Name, inputs, frm.ReportName, frm.Topic)
	fmt.Println("call:", call)
	//TODO: check if chain name is valid
	task := campain.NewContractTask(frm.ChainName, call)
	go goworker.AddTask(task)
}

func OnSyncBlockCall(message string) {
	frm := &form.FormBlockSync{}
	err := json.Unmarshal([]byte(message), frm)
	if err != nil {

		return
	}
	err = frm.IsValid()
	if err != nil {

		return
	}
	blockNumber := &big.Int{}

	if frm.BlockNum.Type == "uint256" {

		err := json.Unmarshal(frm.BlockNum.Value, blockNumber)
		if err != nil {
			fmt.Println("parse param error", err)

			return
		}
	} else {

		return
	}
	camp := campain.GetCampain(frm.ChainName)
	if camp == nil {

		return
	}
	cmd := campain.CreateCmdTransactionsOfBlock(blockNumber.Uint64())
	cmd.Init()
	task := campain.NewClientTask(frm.ChainName, cmd)
	go goworker.AddTask(task)
}

func initWorker() {
	//GetNumWorker
	numWorker := system.Config.GetNumWorker()
	goworker.OrganizeWorker(numWorker)

	fmt.Println("num worker:", numWorker)
	for _, ex := range system.Config.Exports {
		err := export.AddExport(&ex)
		if err != nil {
			panic(err)
		}
		goworker.AddToolWithControl(ex.Name,
			&campain.ExportBlackSmith{
				ExportName: ex.Name,
			},
			1)
		if ex.Type == "wspubsub" {
			hub := export.GetPubSubHub(ex.Name)
			if hub != nil {
				subscriber, err := hub.SubscribeOn("contract.call")
				if err == nil {
					subscriber.SetProcessor(OnContractCall)
					ContractCallSubcriber[ex.Name] = subscriber
				}
				subscriber, err = hub.SubscribeOn("block.sync")
				if err == nil {
					subscriber.SetProcessor(OnSyncBlockCall)
					SyncBlockSubcriber[ex.Name] = subscriber
				}
			}
		}
	}

	for _, chain := range system.Config.Chains {
		fmt.Println("add tool for chain ", chain.Name)

		camp := campain.AddCampain(chain.Name)

		for _, contract := range chain.Contracts {
			err := camp.LoadContract(&contract)
			if err != nil {
				fmt.Println(err.Error())
				panic(err)
			}
			if chain.Name == "bsc" || chain.Name == "kai" {
				fmt.Println("create ethcontract blacksmith", chain.Name+"."+contract.Name, chain.NumWorker)
				goworker.AddToolWithControl(chain.Name+"."+contract.Name, &campain.EthContractBlackSmith{
					Campain:      camp,
					ContractName: contract.Name,
					BackendURLS:  chain.Endpoints,
				}, chain.NumWorker)
			}
		}

		for _, track := range chain.Tracking {

			camp.Tracking(track)
		}

		if chain.Name == "bsc" || chain.Name == "kai" {
			goworker.AddToolWithControl(chain.Name, &campain.ClientBlackSmith{
				Campain:     camp,
				BackendURLS: chain.Endpoints,
			}, chain.NumWorker)

		}
		camp.Run()
	}
}

func reportLive() {
	signal := &godashboard.Signal{ItemName: "chaininter." + system.NodeName}
	godashboard.Report(signal)
}
func main() {

	var port = utility.MustGetEnv("PORT")
	rootPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	system.RootPath = rootPath
	configFile := utility.GetEnv("CONFIG")
	if configFile == "" {
		configFile = "config.json"
	}
	system.NodeName = utility.GenVerifyCode(5)
	//MARK: init system config
	jsonFile2, err := os.Open(rootPath + "/config/" + configFile)

	if err != nil {

		panic(err)
	}

	defer jsonFile2.Close()
	bytes, err := ioutil.ReadAll(jsonFile2)
	systemConfig := entity.Config{}

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytes, &systemConfig)
	if err != nil {
		panic(err)
	}
	system.Config = &systemConfig

	if err != nil {

		panic(err)
	}

	go initWorker()

	for _, db := range system.Config.Dashboards {

		godashboard.AddDashboard(&db)
	}
	utility.Schedule(reportLive, time.Second*5)
	//MARK: init router
	jsonFile, err := os.Open(rootPath + "/config/route.json")

	if err != nil {

		panic(err)
	}

	defer jsonFile.Close()

	bytes, _ = ioutil.ReadAll(jsonFile)
	var handers = map[string]gorouter.EndpointDefine{

		"":         {Handles: Handles{route.Root}},
		"unhandle": {Handles: Handles{route.Unhandle}},
	}

	var router = gorouter.Router{}

	router.Init("v1/", string(bytes), handers)

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("i am ok"))
	})
	http.Handle("/v1/", router)

	fmt.Println("listen on port", port)

	/*for i := 0; i < 2; i++ {
		cmd := &command.CmdGetLatestBlockNumber{}
		cmd.Init()
		task := worker.NewTask("bsc", cmd)
		goworker.AddTask(task)
		//time.Sleep(10 * time.Second)
	}

	cmd := command.CreateCmdTransactionsOfBlock(-1)
	cmd.Init()
	task := worker.NewTask("bsc", cmd)
	goworker.AddTask(task)
	*/
	/*
		call := campain.ContractCall{
			FuncName: "totalSupply",
			Params:   nil,
		}
		task := campain.NewContractTask("bsc.pet", &call)
		goworker.AddTask(task)
	*/

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
