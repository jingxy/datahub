package daemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asiainfoLDP/datahub/ds"
	log "github.com/asiainfoLDP/datahub/utils/clog"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type AccessToken struct {
	Accesstoken   string `json:"accesstoken,omitempty"`
	Remainingtime string `json:"remainingtime,omitempty"`
	Entrypoint    string `json:"entrypoint,omitempty"`
}

func pullHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println(r.URL.Path + "(pull)\n")
	result, _ := ioutil.ReadAll(r.Body)
	reqJson := ds.DsPull{}
	var strret string

	if err := json.Unmarshal(result, &reqJson); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	reqJson.Repository = ps.ByName("repo")
	reqJson.Dataitem = ps.ByName("item")
	if exist := CheckDataPoolExist(reqJson.Datapool); exist == false {
		strret = reqJson.Datapool + " not found. " + reqJson.Tag + " will be pull into " + g_strDpPath
	} else {
		strret = reqJson.Repository + "/" + reqJson.Dataitem + "/" + reqJson.Tag + " will be pull into " + reqJson.Datapool
	}

	url := "/transaction/" + ps.ByName("repo") + "/" + ps.ByName("item") + "/" + reqJson.Tag

	token, entrypoint, err := getAccessToken(url, w)
	if err != nil {
		log.Println(err)
		strret = err.Error()
		return
	} else {
		url = "/pull/" + ps.ByName("repo") + "/" + ps.ByName("item") + "/" + reqJson.Tag +
			"?token=" + token + "&username=" + gstrUsername
		go dl(url, entrypoint, reqJson)
	}

	log.Println(strret)
	msgret := ds.MsgResp{Msg: strret}
	resp, _ := json.Marshal(msgret)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

	/*
		url := "/transaction/" + ps.ByName("repo") + "/" + ps.ByName("item") + "/" + reqJson.Tag

		token, entrypoint, err := getAccessToken(url, w)
		if err != nil {
			return
		} else {
			url = "/pull/" + ps.ByName("repo") + "/" + ps.ByName("item") + "/" + reqJson.Tag +
				"?token=" + token + "?username=" + gstrUsername
		}
		//fmt.Fprintln(w, url)

		//url := "/pull/" + ps.ByName("repo") + "/" + ps.ByName("item") + "/" + reqJson.Tag
		//entrypoint := ""
	*/
	return

}

func dl(uri, entrypoint string, p ds.DsPull) error {
	var ip string
	if entrypoint == "" {
		ip = os.Getenv("DAEMON_IP_PEER")
	} else {
		ip = entrypoint
	}

	fmt.Println(ip)
	if len(ip) == 0 {
		ip = "http://54.223.58.0:35800"
	}

	target := ip + uri
	fmt.Println(target)
	n, err := download(target, p)
	if err != nil {
		log.Printf("[%d bytes returned.]\n", n)
		log.Println(err)
	}
	return err
}

/*download routine, supports resuming broken downloads.*/
func download(url string, p ds.DsPull) (int64, error) {
	log.Printf("we are going to download %s, save to dp=%s,name=%s\n", url, p.Datapool, p.DestName)

	var out *os.File
	var err error
	var destfilename string
	dpexist := CheckDataPoolExist(p.Datapool)
	if dpexist == false {
		destfilename = g_strDpPath + p.DestName
	} else {
		dpconn := GetDataPoolDpconn(p.Datapool)
		if len(dpconn) == 0 {
			destfilename = g_strDpPath + p.DestName
		} else {
			os.MkdirAll(dpconn+"/"+p.Datapool+"/"+p.Repository+"/"+p.Dataitem, 0755)
			destfilename = dpconn + "/" + p.Datapool + "/" + p.Repository + "/" + p.Dataitem + "/" + p.DestName
		}
	}
	fmt.Println(destfilename)
	out, err = os.OpenFile(destfilename, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		return 0, err
	}

	stat, err := out.Stat()
	if err != nil {
		out.Close()
		return 0, err
	}
	out.Seek(stat.Size(), 0)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "go-downloader")
	/* Set download starting position with 'Range' in HTTP header*/
	req.Header.Set("Range", "bytes="+strconv.FormatInt(stat.Size(), 10)+"-")
	log.Printf("%v bytes had already been downloaded.\n", stat.Size())

	resp, err := http.DefaultClient.Do(req)

	/*Save response body to file only when HTTP 2xx received. TODO*/
	if err != nil || (resp != nil && resp.StatusCode/100 != 2) {
		if resp != nil {
			log.Println("http status code:", resp.StatusCode, err)
			body, _ := ioutil.ReadAll(resp.Body)
			log.Println("response Body:", string(body))
		}
		filesize := stat.Size()
		out.Close()
		if filesize == 0 {
			os.Remove(destfilename)
		}
		return 0, err
	}
	defer resp.Body.Close()
	fname := resp.Header.Get("Source-Filename")
	if len(fname) > 0 {
		p.DestName = fname
	}

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		out.Close()
		return 0, err
	}
	out.Close()
	log.Printf("%d bytes downloaded.", n)
	InsertTagToDb(dpexist, p)
	return n, nil
}

func getAccessToken(url string, w http.ResponseWriter) (token, entrypoint string, err error) {
	//log.Println("can't get access token,direct download..")
	//return nil

	log.Println("daemon: connecting to", DefaultServer+url, "to get accesstoken")
	req, err := http.NewRequest("POST", DefaultServer+url, nil)
	if len(loginAuthStr) > 0 {
		req.Header.Set("Authorization", loginAuthStr)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		//w.WriteHeader(http.StatusServiceUnavailable)
		return "", "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(resp.StatusCode, string(body))
	if resp.StatusCode != 200 {

		w.WriteHeader(resp.StatusCode)
		w.Write(body)

		return "", "", errors.New(string(body))
	} else {

		t := AccessToken{}
		result := &ds.Result{Data: &t}
		if err = json.Unmarshal(body, result); err != nil {
			return "", "", err
		} else {
			if len(t.Accesstoken) > 0 {
				//w.WriteHeader(http.StatusOK)
				return t.Accesstoken, t.Entrypoint, nil
			}

		}

	}
	return "", "", errors.New("get access token error.")

}
