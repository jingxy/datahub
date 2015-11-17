package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asiainfoLDP/datahub/ds"
	"io/ioutil"
)

type FormatDp struct {
	Name string `json:"dpname"`
	Type string `json:"dptype"`
	//Conn string `json:"dpconn"`
}

type Item struct {
	Repository string `json:"repository"`
	DataItem   string `json:"dataitem"`
	Tag        string `json:"tag"`
	Time       string `json:"time"`
	Publish    string `json:"publish"`
}
type FormatDpDetail struct {
	Name  string `json:"dpname"`
	Type  string `json:"dptype"`
	Conn  string `json:"dpconn"`
	Items []Item `json:"items"`
}

func Dp(needLogin bool, args []string) (err error) {
	if needLogin && !Logged {
		login(false)
	}

	if len(args) == 0 {
		resp, _ := commToDaemon("GET", "/datapools", nil)
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode == 200 {
			dpResp(false, body)
		} else {
			fmt.Println(string(body))
			fmt.Println(resp.StatusCode)
			err = errors.New(string(resp.StatusCode))
		}

	} else {
		//support: dp name1 name2 name3
		for _, v := range args {
			if v[0] != '-' {
				strdp := fmt.Sprintf("/datapools/%s", v)
				resp, _ := commToDaemon("GET", strdp, nil)
				defer resp.Body.Close()
				body, _ := ioutil.ReadAll(resp.Body)
				if resp.StatusCode == 200 {
					dpResp(true, body)
				} else {
					fmt.Println(resp.StatusCode)
					err = errors.New(string(resp.StatusCode))
				}
			}
		}
	}
	return err
}

func dpResp(bDetail bool, RespBody []byte) {
	if bDetail == false {
		strcDps := []FormatDp{}
		result := &ds.Result{Data: &strcDps}
		err := json.Unmarshal(RespBody, result)
		if err != nil {
			fmt.Println("Get /datapools  dpResp json.Unmarshal error!")
			return
		}
		if result.Code == ResultOK {
			n, _ := fmt.Printf("%-16s    %-8s\n", "DATAPOOL", "TYPE")
			printDash(n - 5)
			for _, dp := range strcDps {
				fmt.Printf("%-16s    %-8s\n", dp.Name, dp.Type)
			}
		} else {
			fmt.Println("Result code:", result.Code, " Msg:", result.Msg)
		}
	} else {
		strcDp := FormatDpDetail{}
		result := &ds.Result{Data: &strcDp}
		err := json.Unmarshal(RespBody, &result)
		if err != nil {
			fmt.Println("Get /datapools/:dpname  dpResp json.Unmarshal error!")
			return
		}
		if result.Code == ResultOK {
			n, _ := fmt.Printf("%s%-16s\t%-16s\t%-16s\n", "DATAPOOL:", strcDp.Name, strcDp.Type, strcDp.Conn)
			for _, item := range strcDp.Items {
				RepoItemTag := item.Repository + "/" + item.DataItem + ":" + item.Tag
				if item.Publish == "Y" {
					fmt.Printf("%-32s\t%-20s\t%-5s\n", RepoItemTag, item.Time, "pub")
				} else {
					fmt.Printf("%-32s\t%-20s\t%-5s\n", RepoItemTag, item.Time, "pull")
				}
			}
			printDash(n)
		} else {
			fmt.Println("Result code:", result.Code, " Msg:", result.Msg)
		}
	}
}

func GetResultMsg(RespBody []byte, bprint bool) (sMsgResp string) {
	result := &ds.Result{}
	err := json.Unmarshal(RespBody, result)
	if err != nil {
		sMsgResp = "Get /datapools  dpResp json.Unmarshal error!"
	} else {
		sMsgResp = "Result code:" + string(result.Code) + " Msg:" + string(result.Msg)
		if bprint == true {
			fmt.Println(sMsgResp)
		}
	}
	return sMsgResp
}

func ShowMsgResp(RespBody []byte, bprint bool) (sMsgResp string) {
	msg := MsgResp{}
	err := json.Unmarshal(RespBody, &msg)
	if err != nil {
		sMsgResp = err.Error() + " " + "ShowMsgResp unmarshal error!"
		fmt.Println(sMsgResp)
	} else {
		sMsgResp = msg.Msg
		if bprint == true {
			fmt.Println(sMsgResp)
		}
	}
	return sMsgResp
}
