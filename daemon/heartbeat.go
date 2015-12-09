package daemon

import (
	"bytes"
	"encoding/json"
	"github.com/asiainfoLDP/datahub/ds"
	log "github.com/asiainfoLDP/datahub/utils/clog"
	"github.com/asiainfoLDP/datahub/utils/logq"
	"io/ioutil"
	"net/http"
	"time"
)

type Beatbody struct {
	Daemonid   string   `json:"daemonid"`
	Entrypoint []string `json:"entrypoint"`
	Log        []string `json:"log,omitempty"`
}

var (
	EntryPoint       string
	EntryPointStatus = "not available"
	DaemonID         string
	heartbeatTimeout = 5 * time.Second
)

func HeartBeat() {
	getEp := false
	for {

		heartbeatbody := Beatbody{Daemonid: DaemonID}
		if getEp == false && len(EntryPoint) == 0 {
			EntryPoint = getEntryPoint()
			getEp = true
		}
		heartbeatbody.Entrypoint = append(heartbeatbody.Entrypoint, EntryPoint)

		logQueue := logq.LogGetqueue()
		if len(logQueue) > 0 {
			heartbeatbody.Log = logQueue
		}

		jsondata, err := json.Marshal(heartbeatbody)
		if err != nil {
			l := log.Error(err)
			logq.LogPutqueue(l)
		}
		url := DefaultServer + "/heartbeat"
		log.Trace("connecting to", url, string(jsondata))
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsondata))
		/*
			if len(loginAuthStr) > 0 {
				req.Header.Set("Authorization", loginAuthStr)
			}
		*/
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			l := log.Error(err)
			logq.LogPutqueue(l)
			time.Sleep(10 * time.Second)
			continue
		}

		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		log.Tracef("HeartBeat http statuscode:%v,  http body:%s", resp.StatusCode, body)

		result := ds.Result{}
		if err := json.Unmarshal(body, &result); err == nil {
			if result.Code == 0 {
				EntryPointStatus = "available"
			} else {
				EntryPointStatus = "not available"
			}
		}

		time.Sleep(heartbeatTimeout)
	}
}
