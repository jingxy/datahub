package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/asiainfoLDP/datahub/ds"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"syscall"
)

const GstrDpPath string = "/var/lib/datahub"

type UserInfo struct {
	userName string
	password string
	b64      string
}

var (
	User     = UserInfo{}
	UnixSock = "/var/run/datahub.sock"
	Logged   = false
	pidFile  = "/var/run/datahub.pid"
)

type Command struct {
	Name      string
	SubCmd    []Command
	Handler   func(login bool, args []string) error
	Desc      string
	NeedLogin bool
}

type MsgResp struct {
	Msg string `json:"msg"`
}

const (
	ResultOK         = 0
	ErrorInvalidPara = iota + 4000
	ErrorNoRecord
	ErrorSqlExec
	ErrorInsertItem
	ErrorUnmarshal
	ErrorMarshal
	ErrorServiceUnavailable
	ErrorFileNotExist
	ErrorTagAlreadyExist
	ErrorDatapoolNotExits
)

var Cmd = []Command{
	{
		Name:    "dp",
		Handler: Dp,
		SubCmd: []Command{
			{
				Name:    "create",
				Handler: DpCreate,
			},
			{
				Name:    "rm",
				Handler: DpRm,
			},
		},
		Desc: "datapools management",
	},
	{
		Name:    "ep",
		Handler: Ep,
		SubCmd: []Command{
			{
				Name:    "rm",
				Handler: EpRm,
			},
		},
		Desc: "entrypoint management",
	},
	{
		Name:    "job",
		Handler: Job,
		SubCmd: []Command{
			{
				Name:    "rm",
				Handler: JobRm,
			},
		},
		Desc: "datahub transfer job management",
	},
	{
		Name:      "login",
		Handler:   Login,
		Desc:      "login to hub.dataos.io",
		NeedLogin: true,
	},

	{
		Name:      "pub",
		Handler:   Pub,
		Desc:      "publish dataitem or tag",
		NeedLogin: true,
	},

	{
		Name:      "pull",
		Handler:   Pull,
		Desc:      "pull dataitem from peer",
		NeedLogin: true,
	},
	{
		Name:      "repo",
		Handler:   Repo,
		Desc:      "repostories mangement",
		NeedLogin: true,
	},
	{
		Name:      "subs",
		Handler:   Subs,
		Desc:      "subscription of items",
		NeedLogin: true,
	},
	{
		Name:    "version",
		Handler: Version,
		Desc:    "datahub version infomation",
	},
}

func login(interactive bool) {
	if Logged {
		if interactive {
			fmt.Println("you are already logged in as", User.userName)
		}
		return
	}

}

func commToDaemon(method, path string, jsonData []byte) (resp *http.Response, err error) {
	//fmt.Println(method, path, string(jsonData))

	req, err := http.NewRequest(strings.ToUpper(method), path, bytes.NewBuffer(jsonData))

	if len(User.userName) > 0 {
		req.SetBasicAuth(User.userName, User.password)
	}

	/* else {
		req.Header.Set("Authorization", "Basic "+os.Getenv("DAEMON_USER_AUTH_INFO"))
	}
	*/
	conn, err := net.Dial("unix", UnixSock)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Datahub Daemon not running? use 'datahub --daemon' to start daemon.")
		os.Exit(2)
	}
	//client := &http.Client{}
	client := httputil.NewClientConn(conn, nil)
	return client.Do(req)
	/*
		defer resp.Body.Close()
		response = *resp
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
	*/
}

func printDash(n int) {
	for i := 0; i < n; i++ {
		fmt.Printf("-")
	}
	fmt.Println()
}

func showResponse(resp *http.Response) {
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(resp.StatusCode, string(body))
		return
	}

	msg := MsgResp{}
	body, _ := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &msg); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(msg.Msg)
	}
}

func showError(resp *http.Response) {

	body, _ := ioutil.ReadAll(resp.Body)
	result := ds.Result{}
	err := json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("ERROR[%v] %v\n", result.Code, result.Msg)
	}

}

func StopP2P() error {

	data, err := ioutil.ReadFile(pidFile)

	if err != nil {
		if os.IsNotExist(err) {
			err = fmt.Errorf("datahub is not running.")
		}
	} else {
		if pid, err := strconv.Atoi(string(data)); err == nil {
			return syscall.Kill(pid, syscall.SIGQUIT)
		}
	}
	return err
	//commToDaemon("get", "/stop", nil)
}

func ShowUsage() {
	fmt.Println("Usage:\tdatahub COMMAND [arg...]")
	fmt.Println("\tdatahub COMMAND [ --help ]")
	fmt.Println("\tdatahub help [COMMAND]\n")
	fmt.Println("A client for Datahub to publish and pull data\n")
	fmt.Println("Commands:")
	for _, v := range Cmd {
		fmt.Printf("    %-10s%s\n", v.Name, v.Desc)
	}
	fmt.Printf("\nrun '%s COMMAND --help' for details on a command.\n", os.Args[0])
}
