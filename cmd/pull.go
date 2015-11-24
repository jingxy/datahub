package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/asiainfoLDP/datahub/ds"
	"github.com/asiainfoLDP/datahub/utils/mflag"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

func Pull(login bool, args []string) (err error) {
	f := mflag.NewFlagSet("pull", mflag.ContinueOnError)
	f.Usage = pullUsage
	if err = f.Parse(args); err != nil {
		return err
	}
	if len(args) != 2 {
		fmt.Println("invalid argument..")
		pullUsage()
		return
	}
	u, err := url.Parse(args[0])
	if err != nil {
		panic(err)
	}
	source := u.Path
	if u.Path[0] == '/' {
		source = u.Path[1:]
	}

	var repo, item string
	dstruc := ds.DsPull{}
	if url := strings.Split(source, "/"); len(url) != 2 {
		fmt.Println("invalid argument..")
		pullUsage()
		return
	} else {
		target := strings.Split(url[1], ":")
		if len(target) == 1 {
			target = append(target, "latest")
		} else if len(target[1]) == 0 {
			target[1] = "latest"
		}
		//uri = fmt.Sprintf("%s/%s:%s", url[0], target[0], target[1])
		repo = url[0]
		item = target[0]
		dstruc.Tag = target[1]
		dstruc.DestName = dstruc.Tag
		dstruc.Datapool = args[1]
	}

	//fmt.Println("uri:", uri)

	jsonData, err := json.Marshal(dstruc)
	if err != nil {
		return
	}

	resp, err := commToDaemon("post", "/subscriptions/"+repo+"/"+item+"/pull", jsonData)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		//body, _ := ioutil.ReadAll(resp.Body)
		//fmt.Println(body)
		body, _ := ioutil.ReadAll(resp.Body)
		ShowMsgResp(body, true)
		//fmt.Printf("%s/%s:%s will be download to %s\n.", repo, item, ds.Tag, ds.Datapool)

	} else if resp.StatusCode == 401 {
		if err := Login(false, nil); err == nil {
			Pull(login, args)
		} else {
			fmt.Println(err)
			return err
		}
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		result := ds.Result{}
		err := json.Unmarshal(body, &result)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		fmt.Println(resp.StatusCode, result.Msg, ", please ensure you have subscribed the repository/dataitem.")
		return nil
	}
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(body)

	return nil // dl(uri)
	//return nil
}

func pullUsage() {
	fmt.Printf("usage: %s pull [[URL]/[REPO]/[ITEM][:TAG]] [DATAPOOL]\n", os.Args[0])
}
