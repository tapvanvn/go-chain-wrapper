package route

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/tapvanvn/go-jsonrpc-wrapper/campain"
	"github.com/tapvanvn/go-jsonrpc-wrapper/form"
	"github.com/tapvanvn/go-jsonrpc-wrapper/route/request"
	"github.com/tapvanvn/go-jsonrpc-wrapper/route/response"
	"github.com/tapvanvn/gorouter/v2"
	"github.com/tapvanvn/goworker"
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
			err := json.Unmarshal(param.Value, input)
			if err != nil {
				fmt.Println("parse param error", err)
				response.BadRequest(context, 0, err.Error(), nil)
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

		err := json.Unmarshal(frm.BlockNum.Value, blockNumber)
		if err != nil {
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
