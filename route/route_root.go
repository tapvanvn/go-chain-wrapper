package route

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/tapvanvn/go-jsonrpc-wrapper/campain"
	"github.com/tapvanvn/go-jsonrpc-wrapper/form"
	"github.com/tapvanvn/go-jsonrpc-wrapper/route/request"
	"github.com/tapvanvn/go-jsonrpc-wrapper/route/response"
	"github.com/tapvanvn/gorouter/v2"
	goworker "github.com/tapvanvn/goworker/v2"
)

//Unhandle handle unhandling route
func Unhandle(context *gorouter.RouteContext) {

	fmt.Println("cannot handle:", context.Path)
	responseData := response.Data{Success: false,
		ErrorCode: 0,
		Message:   "Route to nowhere",
		Data:      nil}

	if data, err := json.Marshal(responseData); err == nil {

		context.W.Write(data)
	}
	context.Handled = true
}

//Root handle root
func Root(context *gorouter.RouteContext) {

	fmt.Println(context.Action)
	if context.Action == "call_contract" {

		callContract(context)
	}
}

func callContract(context *gorouter.RouteContext) {

	frm := &form.FormCallContract{}
	err := request.FromRequest(frm, context.R)
	if err != nil {
		response.BadRequest(context, 0, err.Error(), nil)
		return
	}
	err = frm.IsValid()
	if err != nil {
		response.BadRequest(context, 0, err.Error(), nil)
		return
	}
	inputs := make([]interface{}, 0)
	for _, param := range frm.Params {
		if param.Type == "uint256" {
			input := &big.Int{}
			num, ok := input.SetString(string(param.Value), 10)
			if !ok {
				fmt.Println("parse param error")
				response.BadRequest(context, 0, "cannot parse:"+param.Value, nil)
				return
			}
			input = num
			inputs = append(inputs, input)
		}
	}
	var call *campain.ContractCall = campain.CreateContractCall(frm.Name, inputs, frm.ReportName, frm.Topic)
	if _, ok := context.R.Header["Direct"]; ok {
		parts := strings.Split(frm.ChainName, ".")
		if len(parts) != 2 {
			response.BadRequest(context, 0, err.Error(), nil)
			return
		}
		camp := campain.GetCampain(parts[0])
		if camp == nil {
			fmt.Println("cam not found", parts)
			response.NotFound(context)
			return
		}
		if tool := camp.GetDirectContractTool(parts[1]); tool != nil {

			err := tool.Process(call)
			if err != nil {
				response.BadRequest(context, 0, err.Error(), nil)
				return
			}
			response.Success(context, call)
			return
		} else {
			fmt.Println("tool not found", parts)
		}

	} else {
		fmt.Println(context.R.Header)
		//TODO: check if chain name is valid
		task := campain.NewContractTask(frm.ChainName, call)
		go goworker.AddTask(task)
		response.Success(context, "planed")
	}
}

func callBlockSync(context *gorouter.RouteContext) {
	frm := &form.FormBlockSync{}
	err := request.FromRequest(frm, context.R)
	if err != nil {
		response.BadRequest(context, 0, err.Error(), nil)
		return
	}
	err = frm.IsValid()
	if err != nil {
		response.BadRequest(context, 0, err.Error(), nil)
		return
	}
	blockNumber := &big.Int{}

	if frm.BlockNum.Type == "uint256" {
		num, ok := blockNumber.SetString(string(frm.BlockNum.Value), 10)
		blockNumber = num
		if !ok {
			fmt.Println("parse param error", err)
			response.BadRequest(context, 0, err.Error(), nil)
			return
		}
	} else {
		response.BadRequest(context, 0, "block number format is invalid", nil)
		return
	}
	camp := campain.GetCampain(frm.ChainName)
	if camp == nil {
		response.NotFound(context)
		return
	}
	cmd := campain.CreateCmdTransactionsOfBlock(blockNumber.Uint64())
	cmd.Init()
	task := campain.NewClientTask(frm.ChainName, cmd)
	go goworker.AddTask(task)

}
