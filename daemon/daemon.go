package daemon

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asiainfoLDP/datahub/cmd"
	"github.com/asiainfoLDP/datahub/daemon/daemonigo"
	"github.com/asiainfoLDP/datahub/ds"
	log "github.com/asiainfoLDP/datahub/utils/clog"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

var g_ds = new(ds.Ds)

const g_dbfile string = "/var/lib/datahub/datahub.db"

const g_strDpPath string = cmd.GstrDpPath

type StoppableListener struct {
	*net.UnixListener          //Wrapped listener
	stop              chan int //Channel used only to indicate listener should shutdown
}

type strc_dp struct {
	Dpid   int
	Dptype string
}

func dbinit() {
	log.Println("connect to db sqlite3")
	db, err := sql.Open("sqlite3", g_dbfile)
	//defer db.Close()
	chk(err)
	g_ds.Db = db

	g_ds.Create(ds.Create_dh_dp)
	g_ds.Create(ds.Create_dh_dp_repo_ditem_map)
	g_ds.Create(ds.Create_dh_repo_ditem_tag_map)

}

func Trace() {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	fmt.Printf("%s:%d %s()\n", file, line, f.Name())
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
func get(err error) {
	if err != nil {
		log.Println(err)
	}
}

func New(l net.Listener) (*StoppableListener, error) {
	tcpL, ok := l.(*net.UnixListener)

	if !ok {
		return nil, errors.New("Cannot wrap listener")
	}

	retval := &StoppableListener{}
	retval.UnixListener = tcpL
	retval.stop = make(chan int)

	return retval, nil
}

var StoppedError = errors.New("Listener stopped")
var sl = new(StoppableListener)

func (sl *StoppableListener) Accept() (net.Conn, error) {

	for {
		//Wait up to one second for a new connection
		sl.SetDeadline(time.Now().Add(time.Second))

		newConn, err := sl.UnixListener.Accept()

		//Check for the channel being closed
		select {
		case <-sl.stop:
			return nil, StoppedError
		default:
			//If the channel is still open, continue as normal
		}

		if err != nil {
			netErr, ok := err.(net.Error)

			//If this is a timeout, then continue to wait for
			//new connections
			if ok && netErr.Timeout() && netErr.Temporary() {
				continue
			}
		}

		return newConn, err
	}
}

func (sl *StoppableListener) Stop() {
	close(sl.stop)
}

func helloHttp(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	rw.WriteHeader(http.StatusOK)
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Fprintf(rw, "%s Hello HTTP!\n", req.URL.Path)
	fmt.Fprintf(rw, "%s \n", string(body))
}

func stopHttp(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintf(rw, "Hello HTTP!\n")
	sl.Close()
	log.Println("connect close")
}

func isDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
	panic("not reached")
}
func isFileExists(file string) bool {
	fi, err := os.Stat(file)
	if err == nil {
		log.Println("exist", file)
		return !fi.IsDir()
	}
	return os.IsExist(err)
}

func RunDaemon() {
	log.Println("Run daemon..")
	// Daemonizing echo server application.
	switch isDaemon, err := daemonigo.Daemonize(); {
	case !isDaemon:
		return
	case err != nil:
		log.Fatal("main(): could not start daemon, reason -> %s", err.Error())
	}
	//fmt.Printf("server := http.Server{}\n")

	if false == isDirExists(g_strDpPath) {
		err := os.MkdirAll(g_strDpPath, 0755)
		if err != nil {
			log.Printf("mkdir %s error! %v ", g_strDpPath, err)
		}

	}

	dbinit()

	if len(DaemonID) == 40 {
		log.Println("daemonid", DaemonID)
		saveDaemonID(DaemonID)
	} else {
		log.Println("get daemonid from db")
		DaemonID = getdaemonid()
	}

	os.Chdir(g_strDpPath)
	originalListener, err := net.Listen("unix", cmd.UnixSock)
	if err != nil {
		panic(err)
	} else {
		if err = os.Chmod(cmd.UnixSock, os.ModePerm); err != nil {
			log.Error(err)
		}
	}

	sl, err = New(originalListener)
	if err != nil {
		panic(err)
	}

	router := httprouter.New()
	router.GET("/", helloHttp)
	router.POST("/datapools", dpPostOneHandler)
	router.GET("/datapools", dpGetAllHandler)
	router.GET("/datapools/:dpname", dpGetOneHandler)
	router.DELETE("/datapools/:dpname", dpDeleteOneHandler)

	router.GET("/repositories/:repo/:item/:tag", repoTagHandler)
	router.GET("/repositories/:repo/:item", repoDetailHandler)
	router.GET("/repositories/:repo", repoRepoNameHandler)
	router.GET("/repositories", repoHandler)
	router.GET("/subscriptions", subsHandler)

	router.POST("/repositories/:repo/:item", pubItemHandler)
	router.POST("/repositories/:repo/:item/:tag", pubTagHandler)

	router.POST("/subscriptions/:repo/:item/pull", pullHandler)

	http.Handle("/", router)
	http.HandleFunc("/stop", stopHttp)
	http.HandleFunc("/users/auth", loginHandler)

	server := http.Server{}

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer wg.Done()
		server.Serve(sl)
	}()

	//p2p server
	P2pRouter := httprouter.New()
	P2pRouter.GET("/", sayhello)
	P2pRouter.GET("/pull/:repo/:dataitem/:tag", p2p_pull)
	P2pRouter.GET("/health", p2pHealthyCheckHandler)

	go func() {
		wg.Add(1)
		defer wg.Done()
		http.ListenAndServe(":35800", P2pRouter)
	}()
	go HeartBeat()

	log.Printf("Serving HTTP\n")
	select {
	case signal := <-stop:
		log.Printf("Got signal:%v\n", signal)
	}
	log.Printf("Stopping listener\n")
	sl.Stop()
	log.Printf("Waiting on server\n")
	wg.Wait()
	daemonigo.UnlockPidFile()
	g_ds.Db.Close()

}
func p2pHealthyCheckHandler(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	rw.WriteHeader(http.StatusOK)
}

/*pull parses filename and target IP from HTTP GET method, and start downloading routine. */
func p2p_pull(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("p2p pull...", r.URL.Path)
	r.ParseForm()
	sRepoName := ps.ByName("repo")
	sDataItem := ps.ByName("dataitem")
	sTag := ps.ByName("tag")

	tokenValid := false

	token := r.Form.Get("token")
	username := r.Form.Get("username")
	if len(token) > 0 && len(username) > 0 {
		log.Println(r.URL.Path, "token:", token, "username:", username)
		url := "/transaction/" + sRepoName + "/" + sDataItem + "/" + sTag +
			"?cypt_accesstoken=" + token + "&username=" + username
		tokenValid = checkAccessToken(url)
	}

	if !tokenValid {
		http.Error(rw, "Bad Request", http.StatusBadRequest)
		return
	}

	log.Println(sRepoName, sDataItem, sTag)
	var irpdmid, idpid int
	var stagdetail, sdpname, sdpconn string
	msg := &ds.MsgResp{}
	msg.Msg = "OK."

	sSqlGetRpdmidDpid := fmt.Sprintf(`SELECT DPID, RPDMID FROM DH_DP_RPDM_MAP 
    	WHERE REPOSITORY = '%s' AND DATAITEM = '%s'`, sRepoName, sDataItem)
	row, err := g_ds.QueryRow(sSqlGetRpdmidDpid)
	if err != nil {
		msg.Msg = err.Error()
	}
	row.Scan(&idpid, &irpdmid)
	fmt.Println("dpid", idpid, "rpdmid", irpdmid)

	sSqlGetTagDetail := fmt.Sprintf(`SELECT DETAIL FROM DH_RPDM_TAG_MAP 
        WHERE RPDMID = '%d' AND TAGNAME = '%s'`, irpdmid, sTag)
	tagrow, err := g_ds.QueryRow(sSqlGetTagDetail)
	if err != nil {
		msg.Msg = err.Error()
	}
	tagrow.Scan(&stagdetail)
	log.Println("tagdetail", stagdetail)

	sSqlGetDpconn := fmt.Sprintf(`SELECT DPNAME, DPCONN FROM DH_DP WHERE DPID='%d'`, idpid)
	dprow, err := g_ds.QueryRow(sSqlGetDpconn)
	if err != nil {
		msg.Msg = err.Error()
	}
	dprow.Scan(&sdpname, &sdpconn)
	fmt.Println("dpname", sdpname, "dpconn", sdpconn)

	filepathname := "/" + sdpconn + "/" + sdpname + "/" + sRepoName + "/" + sDataItem + "/" + stagdetail
	fmt.Println(" filename:", filepathname)
	if exists := isFileExists(filepathname); !exists {
		filepathname = "/" + sdpconn + "/" + sdpname + "/" + sRepoName + "/" + sDataItem + "/" + sTag
		if exists := isFileExists(filepathname); !exists {
			filepathname = "/" + sdpconn + "/" + stagdetail
			if exists := isFileExists(filepathname); !exists {
				filepathname = "/var/lib/datahub/" + sTag
				if exists := isFileExists(filepathname); !exists {
					fmt.Println(" filename:", filepathname)
					//http.NotFound(rw, r)
					msg.Msg = fmt.Sprintf("tag:%s not found", sTag)
					resp, _ := json.Marshal(msg)
					respStr := string(resp)
					rw.WriteHeader(http.StatusNotFound)
					fmt.Fprintln(rw, respStr)
					return
				}
			}
		}
	}

	rw.Header().Set("Source-FileName", stagdetail)
	http.ServeFile(rw, r, filepathname)

	resp, _ := json.Marshal(msg)
	respStr := string(resp)
	fmt.Fprintln(rw, respStr)
	return
}

func sayhello(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	rw.WriteHeader(http.StatusOK)
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Fprintf(rw, "%s Hello p2p HTTP !\n", req.URL.Path)
	fmt.Fprintf(rw, "%s \n", string(body))
}

func checkAccessToken(tokenUrl string) bool {
	log.Println("daemon: connecting to", DefaultServer+tokenUrl)
	req, err := http.NewRequest("GET", DefaultServer+tokenUrl, nil)
	if len(loginAuthStr) > 0 {
		req.Header.Set("Authorization", loginAuthStr)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	type tokenDs struct {
		Valid bool `json:"valid"`
	}
	tkresp := tokenDs{}
	result := &ds.Result{Data: &tkresp}
	if err = json.Unmarshal(body, &result); err != nil {
		log.Println(err)
	}

	return tkresp.Valid
}

func init() {
	if srv := os.Getenv("DATAHUB_SERVER"); len(srv) > 0 {
		DefaultServer = srv
	}
}
