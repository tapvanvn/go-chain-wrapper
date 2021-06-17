package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/tapvanvn/go-jsonrpc-wrapper/campain"
	"github.com/tapvanvn/go-jsonrpc-wrapper/entity"
	"github.com/tapvanvn/go-jsonrpc-wrapper/route"
	"github.com/tapvanvn/go-jsonrpc-wrapper/system"
	"github.com/tapvanvn/go-jsonrpc-wrapper/utility"
	"github.com/tapvanvn/gorouter/v2"
	"github.com/tapvanvn/goworker"
)

type Handles []gorouter.RouteHandle
type Endpoint gorouter.EndpointDefine

func initWorker() {
	//GetNumWorker
	numWorker := system.Config.GetNumWorker()
	goworker.OrganizeWorker(numWorker)

	fmt.Println("num worker:", numWorker)

	for _, chain := range system.Config.Chains {
		fmt.Println("add tool for chain ", chain.Name)

		camp := campain.NewCampain(chain.Name, time.Second*5)

		for _, export := range chain.Exports {
			err := camp.AddExport(&export)
			if err != nil {
				panic(err)
			}
		}

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

func main() {

	var port = utility.MustGetEnv("PORT")
	rootPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	system.RootPath = rootPath
	//MARK: init system config
	jsonFile2, err := os.Open(rootPath + "/config/config.json")

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
