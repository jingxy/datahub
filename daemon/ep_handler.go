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
		msg.Msg = "you don't have any entrypoint."
	} else {
		msg.Msg = EntryPoint + " " + EntryPointStatus
	}
	resp, _ := json.Marshal(&msg)
	fmt.Fprintln(w, string(resp))
	return
}

func epPostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	ep := cmd.FormatEp{}
	if err := json.Unmarshal(reqBody, &ep); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	EntryPoint = ep.Ep

	msg := ds.MsgResp{Msg: "OK. your entrypoint is: " + EntryPoint}

	resp, _ := json.Marshal(&msg)
	fmt.Fprintln(w, string(resp))
	return
}
