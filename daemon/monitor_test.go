package daemon

import (
	log "github.com/asiainfoLDP/datahub/utils/clog"
	"github.com/asiainfoLDP/datahub/utils/logq"
	"github.com/go-fsnotify/fsnotify"
	"testing"
)

func Test_initMonitList(t *testing.T) {
	initMonitList()
}

func Test_datapoolMonitor(t *testing.T) {
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		l := log.Error(err)
		logq.LogPutqueue(l)
	}

	defer watcher.Close()
	AddtoMonitor("/var/lib/datahub/datahubUtest.db", "repounittest/itemunittest:tagunittestdatahubUtest.db")
}
