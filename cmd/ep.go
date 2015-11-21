package cmd

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type FormatEp struct {
	Ep string `json:"entrypoint"`
}

func Ep(needLogin bool, args []string) (err error) {

	//d := FormatEp{}
	if len(args) > 1 {
		epUsage()
		return
	}

	jdata := []byte(nil)
	method := "GET"

	if len(args) == 1 {

		ep := parseEp(args[0])
		if len(ep) == 0 {
			epUsage()
			return
		}

		fmt.Println("entrypoint:", ep)
		j := FormatEp{}
		j.Ep = ep

		jdata, err = json.Marshal(j)
		if err != nil {
			fmt.Println(err)
			return err
		}

		method = "POST"

	}

	resp, err := commToDaemon(method, "/ep", jdata)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	showResponse(resp)

	return err
}

func epUsage() {
	fmt.Printf("usage: %s ep [http://HOST:PORT]\n", os.Args[0])
}

func parseEp(s string) (ep string) {
	ep = ""
	u, err := url.Parse(s)
	if err != nil {
		fmt.Println(err)
		return
	}

	host, port, _ := net.SplitHostPort(u.Host)

	if len(u.Scheme) == 0 || len(host) == 0 || len(port) == 0 || strings.ToLower(u.Scheme) != "http" {
		return
	} else {
		if _, err = strconv.Atoi(port); err == nil {
			ep = u.Scheme + "://" + host + ":" + port
		}
	}

	return

}
