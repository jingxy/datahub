package daemon

import (
	//"errors"
	//"fmt"
	//"github.com/asiainfoLDP/datahub/cmd"
	//"github.com/asiainfoLDP/datahub/daemon/daemonigo"
	//"github.com/asiainfoLDP/datahub/ds"
	log "github.com/asiainfoLDP/datahub/utils/clog"
	//"github.com/asiainfoLDP/datahub/utils/logq"
	"github.com/julienschmidt/httprouter"
	//"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	//"os"
	//"os/signal"
	"strings"
	//"sync"
	//"syscall"
	//"time"
)

func init() {
	DefaultServer = "http://localhost:35888"
	go testserver()
}

func testserver() {
	var tsl = new(StoppabletcpListener)

	Listener, err := net.Listen("tcp", ":35888")
	if err != nil {
		log.Fatal(err)
	}

	tsl, err = tcpNew(Listener)
	if err != nil {
		log.Fatal(err)
	}

	tRouter := httprouter.New()
	tRouter.GET("/", sayhello)
	tRouter.GET("/pull/:repo/:dataitem/:tag", p2p_pull)
	tRouter.GET("/health", p2pHealthyCheckHandler)
	//tRouter.GET("/", helloHttp)
	tRouter.POST("/datapools", dpPostOneHandler)
	tRouter.GET("/datapools", dpGetAllHandler)
	tRouter.GET("/datapools/:dpname", dpGetOneHandler)
	tRouter.DELETE("/datapools/:dpname", dpDeleteOneHandler)

	tRouter.GET("/ep", epGetHandler)
	tRouter.POST("/ep", epPostHandler)
	tRouter.DELETE("/ep", epDeleteHandler)

	tRouter.GET("/repositories/:repo/:item/:tag", repoTagHandler)
	tRouter.GET("/repositories/:repo/:item", repoItemHandler)
	tRouter.GET("/repositories/:repo", repoRepoNameHandler)
	tRouter.GET("/repositories", repoHandler)
	tRouter.GET("/subscriptions", subsHandler)

	tRouter.POST("/repositories/:repo/:item", pubItemHandler)
	tRouter.POST("/repositories/:repo/:item/:tag", pubTagHandler)

	tRouter.POST("/subscriptions/:repo/:item/pull", pullHandler)

	tRouter.GET("/job", jobHandler)
	tRouter.GET("/job/:id", jobDetailHandler)
	tRouter.DELETE("/job/:id", jobRmHandler)

	http.Handle("/", tRouter)
	http.HandleFunc("/stop", stopHttp)
	http.HandleFunc("/users/auth", loginHandler)

	server := http.Server{Handler: tRouter}

	log.Info("p2p server start")
	server.Serve(tsl)
	log.Info("p2p server stop")
}

func Test_commToServer(t *testing.T) {
	w := httptest.NewRecorder()
	commToServer("get", "/", nil, w)
}

func Test_loginHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", strings.NewReader(""))
	w := httptest.NewRecorder()
	loginHandler(w, req)
}
