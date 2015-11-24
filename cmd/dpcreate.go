package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/asiainfoLDP/datahub/utils/mflag"
	"io/ioutil"
	//"log"
	//"path/filepath"
	"strings"
)

type FormatDpCreate struct {
	Name string `json:"dpname"`
	Type string `json:"dptype"`
	Conn string `json:"dpconn"`
}

func DpCreate(needLogin bool, args []string) (err error) {
	f := mflag.NewFlagSet("dp create", mflag.ContinueOnError)
	d := FormatDpCreate{}
	f.StringVar(&d.Type, []string{"-type", "t"}, "file", "datapool type")
	f.StringVar(&d.Conn, []string{"-conn"}, "", "datapool connection info")

	if len(args) > 0 && args[0][0] != '-' {
		d.Name = args[0]
		args = args[1:]
	}

	if len(args) == 0 {
		d.Conn = GstrDpPath
		d.Type = "file"
		//fmt.Printf("missing argument.\nSee '%s --help'.\n", f.Name())
		//return
	}

	if err = f.Parse(args); err != nil {
		//fmt.Println(err.Error())
		//fmt.Println("parse parameter error")
		return
	}

	if len(f.Args()) > 0 {
		fmt.Printf("invalid argument.\nSee '%s --help'.\n", f.Name())
		return

	}

	dptype := strings.ToLower(d.Type)
	if dptype != "file" && dptype != "db" && dptype != "hadoop" && dptype != "api" && dptype != "storm" {
		fmt.Println("Datapool type need to be :file,db,hadoop,api,storm")
		return
	}

	if d.Conn[0] != '/' {
		d.Conn = GstrDpPath + d.Conn
		//log.Println(d.Conn)
	}

	jsonData, err := json.Marshal(d)
	if err != nil {
		return err
	}

	if needLogin && !Logged {
		login(false)
	}

	resp, err := commToDaemon("POST", "/datapools", jsonData)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	ShowMsgResp(body, true)

	return err
}
