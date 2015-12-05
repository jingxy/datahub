package daemon

import (
	"encoding/json"
	"fmt"
	"github.com/asiainfoLDP/datahub/cmd"
	"github.com/asiainfoLDP/datahub/ds"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
)

func epGetHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	msg := ds.MsgResp{}

	if len(EntryPoint) == 0 {
		EntryPoint = getEntryPoint()
		if len(EntryPoint) == 0 {
			msg.Msg = "you don't have any entrypoint."
		} else {
			msg.Msg = EntryPoint + " " + EntryPointStatus
		}

	} else {
		msg.Msg = EntryPoint + " " + EntryPointStatus
	}
	resp, _ := json.Marshal(&msg)
	fmt.Fprintln(w, string(resp))
	return
}

func epPostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(reqBody))
	ep := cmd.FormatEp{}
	if err := json.Unmarshal(reqBody, &ep); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	EntryPoint = ep.Ep
	saveEntryPoint(EntryPoint)

	msg := ds.MsgResp{Msg: "OK. your entrypoint is: " + EntryPoint}

	resp, _ := json.Marshal(&msg)
	fmt.Fprintln(w, string(resp))
	return
}

func epDeleteHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	EntryPoint = ""
	delEntryPoint()
	msg := ds.MsgResp{Msg: "OK. your entrypoint has been removed"}

	resp, _ := json.Marshal(&msg)
	fmt.Fprintln(w, string(resp))
	return
}
