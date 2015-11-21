package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asiainfoLDP/datahub/ds"
	"github.com/asiainfoLDP/datahub/utils/mflag"
	"io/ioutil"
	//"net/url"
	//"os"
	"log"
	"path/filepath"
	"strings"
)

const (
	PRIVATE = "private"
	PUBLIC  = "public"
)

func Pub(needlogin bool, args []string) (err error) {
	usage := "usage: sudo datahub pub repository/dataitem --datapool=? \n\tsudo datahub pub repository/dataitem:tag --detail=?"
	if len(args) == 0 {
		fmt.Println(usage)
		return errors.New("args len error!")
	}
	pub := ds.PubPara{}
	var largs []string
	var repo, item, tag, src string
	f := mflag.NewFlagSet("dp create", mflag.ContinueOnError)
	f.StringVar(&pub.Datapool, []string{"-datapool"}, "", "datapool name")
	f.StringVar(&pub.Accesstype, []string{"-accesstype", "t"}, "", "dataitem accesstype, private or public")
	f.StringVar(&pub.Comment, []string{"-comment", "m"}, "", "comments")
	f.StringVar(&pub.Detail, []string{"-detail"}, "", "tag detail ,for example file name")

	if len(args) > 0 && args[0][0] != '-' {
		src = args[0]
		largs = args[1:]
	}

	if err = f.Parse(largs); err != nil {
		fmt.Println("parse parameter error")
		return err
	}
	if len(f.Args()) > 0 {
		fmt.Printf("invalid argument.\nSee '%s --help'.\n", f.Name())
		return errors.New("invalid argument")
	}

	src = strings.Trim(src, "/")
	sp := strings.Split(src, "/")
	if len(sp) != 2 {
		fmt.Println(usage)
		return errors.New("invalid repo/item")
	}
	repo = sp[0]
	sptag := strings.Split(sp[1], ":")
	l := len(sptag)
	if l == 1 {
		item = sptag[0]
		err = PubItem(repo, item, pub, args)
	} else if l == 2 {
		item = sptag[0]
		tag = sptag[1]
		err = PubTag(repo, item, tag, pub, args)
	} else {
		fmt.Printf("invalid argument.\nSee '%s --help'.\n", f.Name())
		return errors.New("invalid argument")
	}

	return err

}

func PubItem(repo, item string, p ds.PubPara, args []string) (err error) {
	url := repo + "/" + item
	if len(p.Accesstype) == 0 {
		p.Accesstype = PRIVATE
	}
	if len(p.Datapool) == 0 {
		fmt.Println("Publishing dataitem requires a parameter \"--datapool=???\" .")
		return
	}
	jsonData, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Mrashal pubdata error while publishing dateitem.")
		return err
	}
	err = pubResp(url, jsonData, args)
	return err
}

func PubTag(repo, item, tag string, p ds.PubPara, args []string) (err error) {
	url := repo + "/" + item + "/" + tag
	if len(p.Detail) == 0 {
		fmt.Println("Publishing tag requires a parameter \"--detail=???\" to ")
		return
	}
	if p.Detail[0] != '/' && strings.Contains(p.Detail, "/") {
		p.Detail, err = filepath.Abs(p.Detail)
		if err != nil {
			log.Print(err.Error())
			return
		}
	}
	jsonData, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Mrashal pubdata error while publishing tag.")
		return err
	}
	err = pubResp(url, jsonData, args)

	return err
}

func pubResp(url string, jsonData []byte, args []string) (err error) {
	resp, err := commToDaemon("POST", "/repositories/"+url, jsonData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		result := ds.Result{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			fmt.Println("PubItem error 1.")
			return err
		} else {
			if result.Code == 0 {
				fmt.Println("pub success")
			} else {
				fmt.Println("Error code: ", result.Code, " Msg: ", result.Msg)
			}
		}
	} else if resp.StatusCode == 401 {
		if err := Login(false, nil); err == nil {
			Pub(true, args)
		} else {
			fmt.Println(err)
		}
	} else {
		result := ds.Result{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			fmt.Println("PubItem error .")
			return err
		} else {
			fmt.Println("Error code: ", resp.StatusCode, "  ", result.Code, " Msg: ", result.Msg)
		}
	}
	return err
}
