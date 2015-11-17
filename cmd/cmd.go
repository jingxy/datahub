package cmd

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

const GstrDpPath string = "/var/lib/datahub/"

type UserInfo struct {
	userName string
	password string
	b64      string
}

var (
	User     = UserInfo{}
	UnixSock = "/var/run/datahub.sock"
	//DefaultServer = "http://10.1.235.98:8080"
	Logged = false
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
		Desc: "list all of datapools.",
	},
	{
		Name:      "repo",
		Handler:   Repo,
		Desc:      "Repostories mangement",
		NeedLogin: true,
	},
	{
		Name:      "subs",
		Handler:   Subs,
		Desc:      "subscription of item.",
		NeedLogin: true,
	},
	{
		Name:      "pull",
		Handler:   Pull,
		Desc:      "pull item from peer.",
		NeedLogin: true,
	},
	{
		Name:      "login",
		Handler:   Login,
		Desc:      "login in to dataos.io.",
		NeedLogin: true,
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
	} else {
		req.Header.Set("Authorization", "Basic "+os.Getenv("DAEMON_USER_AUTH_INFO"))
	}
	conn, err := net.Dial("unix", UnixSock)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Datahub Daemon not running? Or you are not root?")
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
