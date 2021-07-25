package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/tapvanvn/go-chain-wrapper/campain"
	"github.com/tapvanvn/godashboard"
	engines "github.com/tapvanvn/godbengine"
	"github.com/tapvanvn/godbengine/engine"
	"github.com/tapvanvn/godbengine/engine/adapter"

	"github.com/tapvanvn/go-chain-wrapper/entity"
	"github.com/tapvanvn/go-chain-wrapper/export"
	"github.com/tapvanvn/go-chain-wrapper/form"
	"github.com/tapvanvn/go-chain-wrapper/route"
	"github.com/tapvanvn/go-chain-wrapper/system"
	"github.com/tapvanvn/go-chain-wrapper/utility"
	"github.com/tapvanvn/gopubsubengine"
	"github.com/tapvanvn/gorouter/v2"
	goworker "github.com/tapvanvn/goworker/v2"
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
			num, ok := input.SetString(string(param.Value), 10)

			if !ok {
				fmt.Println("parse param error", err)

				return
			}
			input = num
			inputs = append(inputs, input)
		}
	}

	call := campain.CreateContractCall(frm.Name, inputs, frm.ReportName, frm.Topic)

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
		num, ok := blockNumber.SetString(string(frm.BlockNum.Value), 10)

		if !ok {
			fmt.Println("parse param error", err)

			return
		}
		blockNumber = num
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

//Start start engine
func StartEngine(engine *engine.Engine) {

	//read redis define from env
	redisConnectString := utility.MustGetEnv("CONNECT_STRING_REDIS")
	fmt.Println("redis:", redisConnectString)
	redisPool := adapter.RedisPool{}

	err := redisPool.Init(redisConnectString)

	if err != nil {

		fmt.Println("cannot init redis")
	}

	engine.Init(nil, nil, nil)
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
		} else {
			panic("not support export type")
		}
	}

	for _, chain := range system.Config.Chains {

		camp := campain.AddCampain(&chain)

		if len(chain.Endpoints) == 0 {

			panic(errors.New("chain must has atleast 1 endpoint"))
		}

		camp.Run()
	}
}

func reportLive() {
	signal := &godashboard.Signal{ItemName: "chaininter." + system.NodeName,
		Params: campain.ReportLive(),
	}
	godashboard.Report(signal)
}

func main() {
	engines.InitEngineFunc = StartEngine
	_ = engines.GetEngine()

	var port = utility.MustGetEnv("PORT")
	rootPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	system.RootPath = rootPath
	configFile := utility.GetEnv("CONFIG")
	if configFile == "" {
		configFile = "config.jsonc"
	}
	system.NodeName = utility.GenVerifyCode(5)
	//MARK: init system config
	jsonFile2, err := os.Open(rootPath + "/config/" + configFile)

	if err != nil {

		panic(err)
	}

	defer jsonFile2.Close()
	bytes, err := ioutil.ReadAll(jsonFile2)
	bytes = utility.TripComment(bytes)
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
	jsonFile, err := os.Open(rootPath + "/config/route.jsonc")

	if err != nil {

		panic(err)
	}

	defer jsonFile.Close()

	bytes, _ = ioutil.ReadAll(jsonFile)
	bytes = utility.TripComment(bytes)
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

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
