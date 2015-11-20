package daemon

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	EntryPoint string
	DaemonID   string
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
			fmt.Println(err.Error())
			time.Sleep(10 * time.Second)
			continue
		}

		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("HeartBeat http statuscode:%v,  http body:%s\n", resp.StatusCode, body)

		time.Sleep(30 * time.Second)
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
