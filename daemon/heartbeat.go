package daemon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/asiainfoLDP/datahub/ds"
	log "github.com/asiainfoLDP/datahub/utils/clog"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Beatbody struct {
	Daemonid   string   `json:"daemonid"`
	Entrypoint []string `json:"entrypoint"`
}

var (
	EntryPoint       string
	EntryPointStatus string
	DaemonID         string
	heartbeatTimeout = 5 * time.Second
)

func HeartBeat() {
	for {
		heartbeatbody := Beatbody{Daemonid: DaemonID}
		heartbeatbody.Entrypoint = append(heartbeatbody.Entrypoint, EntryPoint)
		jsondata, err := json.Marshal(heartbeatbody)
		url := DefaultServer + "/heartbeat"
		log.Println("connecting to", url, string(jsondata))
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsondata))
		/*
			if len(loginAuthStr) > 0 {
				req.Header.Set("Authorization", loginAuthStr)
			}
		*/
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Error(err.Error())
			time.Sleep(10 * time.Second)
			continue
		}

		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("HeartBeat http statuscode:%v,  http body:%s", resp.StatusCode, body)

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

func getdaemonid() (id string) {
	fmt.Println("TODO get daemonid from db.")
	return

}
func saveDaemonID(id string) {
	fmt.Println("TODO save daemonid to db when srv returns code 0.")
}

func init() {
	EntryPoint = os.Getenv("DAEMON_ENTRYPOINT")
}
