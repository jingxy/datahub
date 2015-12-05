package cmd

import (
	"fmt"
	"github.com/asiainfoLDP/datahub/utils/mflag"
	"io/ioutil"
	"net/http"
)

func Job(needLogin bool, args []string) (err error) {

	f := mflag.NewFlagSet("datahub job", mflag.ContinueOnError)
	fListall := f.Bool([]string{"-all", "a"}, false, "list all jobs")

	path := "/job"
	if len(args) > 0 && len(args[0]) > 0 && args[0][0] != '-' {
		path += "/" + args[0]
	} else {
		if err := f.Parse(args); err != nil {
			return err
		}
		if *fListall {
			path += "?opt=all"
		}
	}

	resp, err := commToDaemon("GET", path, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		showResponse(resp)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
	}

	return err
}

func JobRm(needLogin bool, args []string) (err error) {
	fmt.Println("job rm called.")
	f := mflag.NewFlagSet("datahub job rm", mflag.ContinueOnError)
	fForce := f.Bool([]string{"-force", "f"}, false, "force cancel a pulling job.")

	path := "/job"
	if len(args) > 0 && len(args[0]) > 0 && args[0][0] != '-' {
		path += "/" + args[0]

	}
	if len(args) > 1 && len(args[1]) > 0 {
		if err := f.Parse(args[1:]); err != nil {
			return err
		}
		if *fForce {
			path += "?opt=force"
		}
	}

	resp, err := commToDaemon("DELETE", path, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		showResponse(resp)
	} else {
		//showjobResp(resp)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
	}
	return err
}

func jobUsage() {
	fmt.Println("Usage: datahub job [-a]")
}
