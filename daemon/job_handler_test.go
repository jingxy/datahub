package daemon

import (
	"github.com/asiainfoLDP/datahub/ds"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

/*type Context struct {
	desc        string
	requestBody string
	ps          httprouter.Params
	rw          *httptest.ResponseRecorder
	out         string
	code        int
	d           interface{}
}*/

func Test_jobDetailHandler(t *testing.T) {
	DatahubJob = append(DatahubJob,
		ds.JobInfo{
			ID:     "qwert358",
			Tag:    "ttttTag",
			Path:   "/var/lib/datahub/datahubUnitTest",
			Stat:   "transfering",
			Dlsize: 0},
		ds.JobInfo{
			ID:     "ewlxsoso",
			Tag:    "2tttTag",
			Path:   "/var/lib/datahub/datahubUnitTest",
			Stat:   "download",
			Dlsize: 702})

	contexts := []Context{
		{
			requestBody: "",
			rw:          httptest.NewRecorder(),
			out:         "ok",
			code:        0,
			ps:          httprouter.Params{{"id", "qwerty358"}},
		},
		{
			requestBody: "",
			rw:          httptest.NewRecorder(),
			out:         "ok",
			code:        0,
			ps:          httprouter.Params{{"id", "12345678"}},
		},
	}

	for _, v := range contexts {
		req, _ := http.NewRequest("Get", "/job", strings.NewReader(v.requestBody))
		jobDetailHandler(v.rw, req, v.ps)
		if !ExpectResult(t, v.rw, v.out, v.code) {
			t.Log("jobDetail-------- fail", v.code)
		} else {
			t.Log("jobDetail-------- success")
		}
	}
}

func Test_jobHandler(t *testing.T) {

	contexts := []Context{
		{
			requestBody: "",
			rw:          httptest.NewRecorder(),
			out:         "ok",
			code:        0,
		},
	}

	for _, v := range contexts {
		req, _ := http.NewRequest("Get", "/job", strings.NewReader(v.requestBody))
		jobHandler(v.rw, req, v.ps)
		if !ExpectResult(t, v.rw, v.out, v.code) {
			t.Log("job-------- fail", v.code)
		} else {
			t.Log("job-------- success")
		}
	}
}
