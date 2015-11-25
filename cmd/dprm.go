package cmd

func DpRm(needLogin bool, args []string) (err error) {
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
