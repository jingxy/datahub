package cmd

import (
	"fmt"
	"github.com/asiainfoLDP/datahub/utils/mflag"
)

func DpRm(needLogin bool, args []string) (err error) {
	f := mflag.NewFlagSet("dp rm", mflag.ContinueOnError)
	f.Usage = dprUseage
	if err = f.Parse(args); err != nil {
		return err
	}
	if needLogin && !Logged {
		login(false)
	}

	if len(args) > 0 && args[0][0] != '-' {
		for _, v := range args {
			dp := v
			if v[0] != '-' {

				resp, _ := commToDaemon("DELETE", "/datapools/"+dp, nil)
				defer resp.Body.Close()
				showResponse(resp)
			}
		}
	}
	return nil
}

func dprUseage() {
	fmt.Println("Usage of datahub dp rm:")
	fmt.Println("  datahub dp rm DATAPOOL")
	fmt.Println("Remove a datapool")
}
