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
	"github.com/asiainfoLDP/datahub/utils/julienschmidt/httprouter"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	g_ds = new(ds.Ds)

	wg sync.WaitGroup
)

const (
	g_dbfile    string = "/var/lib/datahub/datahub.db"
	g_strDpPath string = cmd.GstrDpPath
)

type StoppableListener struct {
	*net.UnixListener          //Wrapped listener
	stop              chan int //Channel used only to indicate listener should shutdown
}

type StoppabletcpListener struct {
	*net.TCPListener          //Wrapped listener
	stop             chan int //Channel used only to indicate listener should shutdown
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

	var RetDhRpdmTagMap string
	row, err := g_ds.QueryRow(ds.SQLIsExistRpdmTagMap)
	if err != nil {
		log.Error("Get Dh_Rpdm_Tag_Map error!")
		return
	}
	row.Scan(&RetDhRpdmTagMap)
	if len(RetDhRpdmTagMap) == 0 {
		g_ds.Create(ds.Create_dh_dp)
		g_ds.Create(ds.Create_dh_dp_repo_ditem_map)
		g_ds.Create(ds.Create_dh_repo_ditem_tag_map)
	} else {
		if false == strings.Contains(RetDhRpdmTagMap, "TAGID") {
			UpdateSql04To05()
		}
	}

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
	unixL, ok := l.(*net.UnixListener)

	if !ok {
		return nil, errors.New("Cannot wrap listener")
	}

	retval := &StoppableListener{}
	retval.UnixListener = unixL
	retval.stop = make(chan int)

	return retval, nil
}
func tcpNew(l net.Listener) (*StoppabletcpListener, error) {
	tcpL, ok := l.(*net.TCPListener)

	if !ok {
		return nil, errors.New("Cannot wrap listener")
	}

	retval := &StoppabletcpListener{}
	retval.TCPListener = tcpL
	retval.stop = make(chan int)

	return retval, nil
}

var StoppedError = errors.New("Listener stopped")
var sl = new(StoppableListener)
var p2psl = new(StoppabletcpListener)

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

func (tcpsl *StoppabletcpListener) Accept() (net.Conn, error) {

	for {
		//Wait up to one second for a new connection
		tcpsl.SetDeadline(time.Now().Add(time.Second))

		newConn, err := tcpsl.TCPListener.Accept()

		//Check for the channel being closed
		select {
		case <-tcpsl.stop:
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

func (tcpsl *StoppabletcpListener) Stop() {

	close(tcpsl.stop)
}

func helloHttp(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	rw.WriteHeader(http.StatusOK)
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Fprintf(rw, "%s Hello HTTP!\n", req.URL.Path)
	fmt.Fprintf(rw, "%s \n", string(body))
}

func stopHttp(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
	//fmt.Fprintf(rw, "Hello HTTP!\n")
	sl.Close()
	p2psl.Close()
	log.Println("connect close")
}

func isDirExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		log.Println(err.Error())
		return os.IsExist(err)
	} else {
		//log.Println(fi.IsDir())
		return fi.IsDir()
	}
	//panic("not reached")
	return false
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
	//fmt.Println("Run daemon..")
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
		log.Fatal(err)
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

	router.GET("/ep", epGetHandler)
	router.POST("/ep", epPostHandler)
	router.DELETE("/ep", epDeleteHandler)

	router.GET("/repositories/:repo/:item/:tag", repoTagHandler)
	router.GET("/repositories/:repo/:item", repoItemHandler)
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

	go func() {

		stop := make(chan os.Signal)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		select {
		case signal := <-stop:
			log.Printf("Got signal:%v", signal)
		}

		sl.Stop()
		if len(DaemonID) > 0 {
			p2psl.Stop()
		}

	}()

	if len(DaemonID) > 0 {
		go startP2PServer()
		go HeartBeat()
		go datapoolMonitor()
	} else {
		log.Error("no daemonid specificed.")
		fmt.Println("You don't have a daemonid specificed.")
	}

	/*
		wg.Add(1)
		defer wg.Done()
	*/
	log.Info("starting daemon listener...")
	server.Serve(sl)
	log.Info("Stopping daemon listener...")

	if len(DaemonID) > 0 {
		wg.Wait()
	}

	daemonigo.UnlockPidFile()
	g_ds.Db.Close()

	log.Info("daemon exit....")
	log.CloseLogFile()

}

func startP2PServer() {
	p2pListener, err := net.Listen("tcp", ":35800")
	if err != nil {
		log.Fatal(err)
	}

	p2psl, err = tcpNew(p2pListener)
	if err != nil {
		log.Fatal(err)
	}

	P2pRouter := httprouter.New()
	P2pRouter.GET("/", sayhello)
	P2pRouter.GET("/pull/:repo/:dataitem/:tag", p2p_pull)
	P2pRouter.GET("/health", p2pHealthyCheckHandler)

	p2pserver := http.Server{Handler: P2pRouter}

	//stop := make(chan os.Signal)
	//signal.Notify(stop, syscall.SIGINT)

	wg.Add(1)
	defer wg.Done()

	log.Info("p2p server start")
	p2pserver.Serve(p2psl)
	log.Info("p2p server stop")

}

func p2pHealthyCheckHandler(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	healthbody := Beatbody{Daemonid: DaemonID}
	jsondata, err := json.Marshal(healthbody)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(jsondata)
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
	var stagdetail, sdpname, sdpconn, itemdesc string
	msg := &ds.MsgResp{}
	msg.Msg = "OK."

	sSqlGetRpdmidDpid := fmt.Sprintf(`SELECT DPID, RPDMID, ITEMDESC FROM DH_DP_RPDM_MAP 
    	WHERE REPOSITORY = '%s' AND DATAITEM = '%s' AND STATUS='A'`, sRepoName, sDataItem)
	row, err := g_ds.QueryRow(sSqlGetRpdmidDpid)
	if err != nil {
		msg.Msg = err.Error()
	}
	row.Scan(&idpid, &irpdmid, &itemdesc)
	if len(itemdesc) == 0 {
		itemdesc = sRepoName + "_" + sDataItem
	}
	log.Println("dpid:", idpid, "rpdmid:", irpdmid, "itemdesc:", itemdesc)

	sSqlGetTagDetail := fmt.Sprintf(`SELECT DETAIL FROM DH_RPDM_TAG_MAP 
        WHERE RPDMID = '%d' AND TAGNAME = '%s' AND STATUS='A'`, irpdmid, sTag)
	tagrow, err := g_ds.QueryRow(sSqlGetTagDetail)
	if err != nil {
		msg.Msg = err.Error()
	}
	tagrow.Scan(&stagdetail)
	log.Println("tagdetail", stagdetail)
	if len(stagdetail) == 0 {
		log.Warnf("%s(tag:%s) not found", stagdetail, sTag)
		http.Error(rw, sTag+" not found", http.StatusNotFound)
		return
	}

	sSqlGetDpconn := fmt.Sprintf(`SELECT DPNAME, DPCONN FROM DH_DP WHERE DPID='%d'`, idpid)
	dprow, err := g_ds.QueryRow(sSqlGetDpconn)
	if err != nil {
		msg.Msg = err.Error()
	}
	dprow.Scan(&sdpname, &sdpconn)
	log.Println("dpname:", sdpname, "dpconn:", sdpconn)

	filepathname := sdpconn + "/" + itemdesc + "/" + stagdetail
	log.Println("filename:", filepathname)
	if exists := isFileExists(filepathname); !exists {
		log.Error("1 file not found", filepathname)
		//filepathname = "/" + sdpconn + "/" + sdpname + "/" + sRepoName + "/" + sDataItem + "/" + stagdetail
		//if exists := isFileExists(filepathname); !exists {
		//filepathname = "/" + sdpconn + "/" + stagdetail
		//if exists := isFileExists(filepathname); !exists {
		//	filepathname = "/var/lib/datahub/" + sTag
		//	if exists := isFileExists(filepathname); !exists {
		//log.Error("2 filename not found:", filepathname)
		//http.NotFound(rw, r)
		msg.Msg = fmt.Sprintf("tag:%s not found", sTag)
		resp, _ := json.Marshal(msg)
		respStr := string(resp)
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(rw, respStr)
		return
		//}
	}
	log.Println("Tag file full path name :", filepathname)
	rw.Header().Set("Source-FileName", stagdetail)
	log.Info("transfering", filepathname)
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
	log.Println(string(body))

	return tkresp.Valid
}

func init() {
	if srv := os.Getenv("DATAHUB_SERVER"); len(srv) > 0 {
		DefaultServer = srv
	}

	log.SetLogLevel(log.LOG_LEVEL_INFO)
}
