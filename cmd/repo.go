package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/asiainfoLDP/datahub/ds"
	"io/ioutil"
	"os"
)

func Repo(login bool, args []string) (err error) {

	itemDetail := false
	if len(args) > 1 {
		fmt.Println("invalid argument..")
		repoUsage()
		return
	}

	uri := "/repositories"
	if len(args) == 1 {
		uri = uri + "/" + args[0]
		itemDetail = true
	}

	resp, err := commToDaemon("get", uri, nil)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		if len(args) == 1 {
			repoResp(itemDetail, body, args[0])
		} else {
			repoResp(itemDetail, body, "")
		}
	} else if resp.StatusCode == 401 {
		//fmt.Println(resp.StatusCode, "returned....")
		if err := Login(false, nil); err == nil {
			Repo(login, args)
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(resp.StatusCode)
	}

	return err
}

func repoUsage() {
	fmt.Printf("usage: %s repo [[URL]/[REPO]/[ITEM]\n", os.Args[0])
}

func repoResp(detail bool, respbody []byte, repoitem string) {
	//fmt.Println(string(respbody))

	if detail {
		subs := ds.Data{}
		result := &ds.Result{Data: &subs}
		err := json.Unmarshal(respbody, &result)
		if err != nil {
			panic(err)
		}
		n, _ := fmt.Printf("%s\t%s\n", "REPOSITORY/ITEM[:TAG]", "UPDATETIME")
		printDash(n + 12)
		for _, tag := range subs.Tags {
			fmt.Printf("%s:%-8s\t%s\n", repoitem, tag.Tag, tag.Optime)
		}
	} else {
		subs := []ds.Data{}
		result := &ds.Result{Data: &subs}
		err := json.Unmarshal(respbody, &result)
		if err != nil {
			panic(err)
		}
		n, _ := fmt.Printf("%s/%-8s\t%s\n", "REPOSITORY", "ITEM", "TYPE")
		printDash(n + 5)

		for _, item := range subs {
			fmt.Printf("%s/%-8s\t%s\n", item.Repository_name, item.Dataitem_name, "File")
		}

	}

}
