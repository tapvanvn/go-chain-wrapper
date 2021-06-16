package route

import (
	"encoding/json"
	"fmt"

	"github.com/tapvanvn/go-jsonrpc-wrapper/campain"
	"github.com/tapvanvn/go-jsonrpc-wrapper/route/form"
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
	call := campain.CreateContractCall(frm.Name, frm.Params, frm.ReportName, frm.Topic)

	//TODO: check if chain name is valid
	task := campain.NewContractTask(frm.ChainName, call)
	goworker.AddTask(task)
}
