package daemon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/asiainfoLDP/datahub/ds"
	log "github.com/asiainfoLDP/datahub/utils/clog"
	"github.com/asiainfoLDP/datahub/utils/logq"
	"io/ioutil"
	"net/http"
	"os"
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

	for {

		heartbeatbody := Beatbody{Daemonid: DaemonID}
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

func getDaemonid() (id string) {
	fmt.Println("TODO get daemonid from db.")
	s := `SELECT DAEMONID FROM DH_DAEMON;`
	row, e := g_ds.QueryRow(s)
	if e != nil {
		l := log.Error(s, "error.", e)
		logq.LogPutqueue(l)
		return
	}
	row.Scan(&id)
	log.Info("daemon id is", id)
	return id
}

func saveDaemonID(id string) {
	log.Println("TODO save daemonid to db when srv returns code 0.")
	count := `SELECT COUNT(*) FROM DH_DAEMON;`
	row, err := g_ds.QueryRow(count)
	if err != nil {
		l := log.Error(count, "error.", err)
		logq.LogPutqueue(l)
	}
	var c int
	row.Scan(&c)
	if c > 0 {
		Update := fmt.Sprintf(`UPDATE DH_DAEMON SET DAEMONID='%s';`, id)
		log.Debug(Update)
		if _, e := g_ds.Update(Update); e != nil {
			l := log.Error(Update, "error.", e)
			logq.LogPutqueue(l)
		}
	} else {
		Insert := fmt.Sprintf(`INSERT INTO DH_DAEMON (DAEMONID) VALUES ('%s');`, id)
		log.Debug(c, Insert)
		if _, e := g_ds.Insert(Insert); e != nil {
			l := log.Error(Insert, "error.", e)
			logq.LogPutqueue(l)
		}
	}
}

func init() {
	EntryPoint = os.Getenv("DAEMON_ENTRYPOINT")
}

func saveEntryPoint(ep string) {
	fmt.Println("TODO save ep to db")
}

func delEntryPoint() {
	fmt.Println("TODO remove ep from db.")
}
