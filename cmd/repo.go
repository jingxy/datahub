package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asiainfoLDP/datahub/ds"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

const (
	Repos = iota
	ReposReponame
	ReposReponameDataItem
	ReposReponameDataItemTag
)

func Repo(login bool, args []string) (err error) {
	var icmd int
	if len(args) > 1 {
		fmt.Println("invalid argument..")
		repoUsage()
		return
	}
	var repo, item, tag, uri string
	if len(args) == 0 {
		uri = "/repositories"
		icmd = Repos
	} else {
		u, err := url.Parse(args[0])
		if err != nil {
			panic(err)
		}
		source := u.Path
		if len(u.Path) > 0 && u.Path[0] == '/' {
			source = u.Path[1:]
		}

		urls := strings.Split(source, "/")
		lenth := len(urls)
		//fmt.Println(lenth, urls)
		if lenth == 0 {
			uri = "/repositories"
			icmd = Repos
			//fmt.Println(uri, icmd)
		} else if lenth == 1 || (lenth == 2 && len(urls[1]) == 0) {
			uri = "/repositories/" + urls[0]
			icmd = ReposReponame
			repo = urls[0]
			//fmt.Println(uri, icmd)
		} else if lenth == 2 || (lenth == 3 && len(urls[2]) == 0) {
			target := strings.Split(urls[1], ":")
			tarlen := len(target)
			if tarlen == 1 || (tarlen == 2 && len(target[1]) == 0) {
				uri = "/repositories/" + urls[0] + "/" + target[0]
				icmd = ReposReponameDataItem
				repo = urls[0]
				item = target[0]
				//fmt.Println(uri, icmd)
			} else if tarlen == 2 {
				uri = "/repositories/" + urls[0] + "/" + target[0] + "/" + target[1]
				icmd = ReposReponameDataItemTag
				repo = urls[0]
				item = target[0]
				tag = target[1]
				//fmt.Println(uri, icmd)
			}
		} else {
			fmt.Println("The parameter after repo is in wrong format!")
			return errors.New("The parameter after repo is in wrong format!")
		}
	}
	//fmt.Println(uri)
	resp, err := commToDaemon("get", uri, nil)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		repoResp(icmd, body, repo, item, tag)
	} else if resp.StatusCode == 401 || resp.StatusCode == 400 {
		result := ds.Result{}
		err := json.Unmarshal(body, &result)
		if err != nil {
			panic(err)
		}
		if result.Code != 1400 {
			fmt.Println(result.Msg)
			return nil
		}
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

func repoResp(icmd int, respbody []byte, repo, item, tag string) {
	//fmt.Println(string(respbody))
	result := ds.Result{Code: ResultOK}
	if icmd == Repos {
		repos := []ds.Repositories{}
		result.Data = &repos
		err := json.Unmarshal(respbody, &result)
		if err != nil {
			panic(err)
		}
		n, _ := fmt.Printf("%-16s\t%-16s\t%-16s\n", "REPOSITORY", "UPDATETIME", "COMMENT")
		printDash(n + 12)
		for _, v := range repos {
			fmt.Printf("%-16s\t%-16s\t%-16s\n", v.RepositoryName, v.Optime, v.Comment)
		}
	} else if icmd == ReposReponame {
		onerepo := ds.Repository{}
		result.Data = &onerepo
		err := json.Unmarshal(respbody, &result)
		if err != nil {
			panic(err)
		}
		n, _ := fmt.Printf("REPOSITORY/DATAITEM\n")
		printDash(n + 12)
		for _, v := range onerepo.DataItems {
			fmt.Printf("%s/%s\n", repo, v)
		}

	} else if icmd == ReposReponameDataItem {
		repoitemtags := ds.Data{}
		result.Data = &repoitemtags
		err := json.Unmarshal(respbody, &result)
		if err != nil {
			panic(err)
		}
		n, _ := fmt.Printf("%s\t%s\n", "REPOSITORY/ITEM:TAG", "UPDATETIME")
		printDash(n + 12)
		for _, v := range repoitemtags.Taglist {
			fmt.Printf("%s/%s:%s\t%s\n", repo, item, v.Tag, v.Optime)
		}
	} else if icmd == ReposReponameDataItemTag {
		onetag := ds.Tag{}
		result.Data = &onetag
		err := json.Unmarshal(respbody, &result)
		if err != nil {
			panic(err)
		}
		n, _ := fmt.Printf("%s\t%s\t%s\n", "REPOSITORY/ITEM:TAG", "UPDATETIME", "COMMENT")
		printDash(n + 12)
		fmt.Printf("%s/%s:%s\t%s\t%s\n", repo, item, tag, onetag.Optime, onetag.Comment)
	}
}
