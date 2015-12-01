package daemon

import (
	"fmt"
	log "github.com/asiainfoLDP/datahub/utils/clog"
	"github.com/asiainfoLDP/datahub/utils/go-fsnotify/fsnotify"
	"github.com/asiainfoLDP/datahub/utils/logq"
)

var monitList map[string]string

func datapoolMonitor() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		l := log.Error(err)
		logq.LogPutqueue(l)
	}
	defer watcher.Close()

	initMonitList()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Debug("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					l := log.Warn("modified file:", event.Name)
					logq.LogPutqueue(l)
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					l := log.Warn("deleted file:", event.Name)
					logq.LogPutqueue(l)
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					l := log.Warn("renamed file:", event.Name)
					logq.LogPutqueue(l)
				}
			case err := <-watcher.Errors:
				l := log.Error("error:", err)
				logq.LogPutqueue(l)
			}
		}
	}()

	for tag, filecheck := range monitList {
		err = watcher.Add(filecheck)
		if err != nil {
			l := log.Error(tag, filecheck, err)
			logq.LogPutqueue(l)
		}
	}

	<-done
}

func initMonitList() {
	fmt.Println("TODO INIT MONIT LIST.")
	monitList = make(map[string]string)
	if e := GetTagDetails(&monitList); e != nil {
		log.Errorf("GetTagDetails error. %v", e)
	}
	log.Println(monitList)
}
