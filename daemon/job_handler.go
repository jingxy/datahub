package daemon

import (
	"crypto/rand"
	"fmt"
	"github.com/asiainfoLDP/datahub/ds"
	log "github.com/asiainfoLDP/datahub/utils/clog"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var DatahubJob = make(map[string]ds.JobInfo) //job[id]=JobInfo

func jobHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	log.Trace("from", req.RemoteAddr, req.Method, req.URL.RequestURI(), req.Proto)

	var joblist []ds.JobInfo
	for _, job := range DatahubJob {
		joblist = append(joblist, job)
	}
	r, _ := buildResp(0, "ok", joblist)
	w.WriteHeader(http.StatusOK)
	w.Write(r)

	//http.Error(w, log.Info(req.URL.RequestURI()), http.StatusNotImplemented)

	return

}

func jobDetailHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	log.Trace("from", req.RemoteAddr, req.Method, req.URL.RequestURI(), req.Proto)
	http.Error(w, log.Info(req.URL.RequestURI()), http.StatusNotImplemented)

	return

}

func jobRmHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	log.Trace("from", req.RemoteAddr, req.Method, req.URL.RequestURI(), req.Proto)
	http.Error(w, log.Info(req.URL.RequestURI()), http.StatusNotImplemented)

	return

}

func genJobID() (id string, err error) {
	c := 8
	b := make([]byte, c)
	_, err = rand.Read(b)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	return fmt.Sprintf("%x", b), nil
}
