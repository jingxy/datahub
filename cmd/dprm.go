package cmd

import (
	"fmt"
	"io/ioutil"
)

func DpRm(needLogin bool, args []string) (err error) {
	var strdpname string
	if needLogin && !Logged {
		login(false)
	}
	if len(args) > 0 && args[0][0] != '-' {
		for _, v := range args {
			strdpname = v
			if v[0] != '-' {
				strDpurl := fmt.Sprintf("/datapools/%s", strdpname)

				resp, _ := commToDaemon("DELETE", strDpurl, nil)
				defer resp.Body.Close()
				body, _ := ioutil.ReadAll(resp.Body)
				ShowMsgResp(body, true)
			}
		}
	}
	return nil
}
