package daemon

import (
	"encoding/json"
	"fmt"
	"github.com/asiainfoLDP/datahub/cmd"
	"github.com/asiainfoLDP/datahub/ds"
	log "github.com/asiainfoLDP/datahub/utils/clog"
	"github.com/asiainfoLDP/datahub/utils/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func dpPostOneHandler(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()

	rw.WriteHeader(http.StatusOK)

	if r.Method == "POST" {
		result, _ := ioutil.ReadAll(r.Body)
		reqJson := cmd.FormatDpCreate{}
		err := json.Unmarshal(result, &reqJson)
		if err != nil {
			fmt.Printf("%T\n%s\n%#v\n", err, err, err)
			fmt.Fprintln(rw, "Invalid argument.")
		}
		if len(reqJson.Name) == 0 {
			fmt.Fprintln(rw, "Invalid argument.")
		} else {
			msg := &ds.MsgResp{}
			var sdpDirName string
			if len(reqJson.Conn) == 0 {
				reqJson.Conn = g_strDpPath
				sdpDirName = g_strDpPath

			} else if reqJson.Conn[0] != '/' {
				sdpDirName = g_strDpPath + reqJson.Conn
				reqJson.Conn = sdpDirName
				/*if reqJson.Conn[len(reqJson.Conn)-1] == '/' {
					sdpDirName = sdpDirName + reqJson.Name
				} else {
					sdpDirName = sdpDirName + "/" + reqJson.Name
				}*/
			} else {
				sdpDirName = reqJson.Conn
			}

			dpexist := CheckDataPoolExist(reqJson.Name)
			if dpexist {
				msg.Msg = fmt.Sprintf("Datapool %s is already exist, please use another name!", reqJson.Name)
				resp, _ := json.Marshal(msg)
				fmt.Fprintln(rw, string(resp))
				return
			}
			if err := os.MkdirAll(sdpDirName, 0777); err != nil {
				log.Println(sdpDirName)
				msg.Msg = err.Error()
			} else {
				msg.Msg = fmt.Sprintf("dp create success. dp:%s total path:%s", reqJson.Name, sdpDirName)
				reqJson.Conn = strings.TrimRight(reqJson.Conn, "/")
				sql_dp_insert := fmt.Sprintf(`insert into DH_DP (DPID, DPNAME, DPTYPE, DPCONN, STATUS)
					values (null, '%s', '%s', '%s', 'A')`, reqJson.Name, reqJson.Type, reqJson.Conn)
				if _, err := g_ds.Insert(sql_dp_insert); err != nil {
					os.Remove(sdpDirName)
					msg.Msg = err.Error()
				}
			}
			resp, _ := json.Marshal(msg)
			respStr := string(resp)
			fmt.Fprintln(rw, respStr)
		}

	}

}

func dpGetAllHandler(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	rw.WriteHeader(http.StatusOK)
	//msg := &ds.MsgResp{}
	//msg.Msg = "OK."
	dps := []cmd.FormatDp{}
	result := &ds.Result{Code: cmd.ResultOK, Data: &dps} //must use a pointer dps to initial Data
	onedp := cmd.FormatDp{}
	sqlDp := fmt.Sprintf(`SELECT DPNAME, DPTYPE FROM DH_DP WHERE STATUS = 'A'`)
	rows, err := g_ds.QueryRows(sqlDp)
	if err != nil {
		result.Msg = err.Error()
		result.Code = cmd.ErrorSqlExec
		resp, _ := json.Marshal(result)
		rw.Write(resp)
	}
	defer rows.Close()
	bresultflag := false
	for rows.Next() {
		bresultflag = true
		rows.Scan(&onedp.Name, &onedp.Type)
		dps = append(dps, onedp)
	}
	if bresultflag == false {
		//fmt.Println(bresultflag)
		result.Msg = "There isn't any datapool."
		result.Code = cmd.ErrorNoRecord
		resp, _ := json.Marshal(result)
		fmt.Println(string(resp))
		rw.Write(resp)
		return
	}

	resp, _ := json.Marshal(result)
	log.Println(string(resp))
	fmt.Fprintln(rw, string(resp))

}

func dpGetOneHandler(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	rw.WriteHeader(http.StatusOK)
	dpname := ps.ByName("dpname")

	/*In future, we need to get dptype in Json to surpport FILE\ DB\ SDK\ API datapool
	result, _ := ioutil.ReadAll(r.Body)
	reqJson := cmd.FormatDp{}
	err := json.Unmarshal(result, &reqJson)
	if err != nil {
		fmt.Printf("%T\n%s\n%#v\n", err, err, err)
		fmt.Println(rw, "invalid argument.")
	}*/

	onedp := cmd.FormatDpDetail{}
	result := &ds.Result{Code: cmd.ResultOK, Data: &onedp}

	sqlTotal := fmt.Sprintf(`SELECT COUNT(*) FROM DH_DP 
		WHERE STATUS = 'A' AND DPNAME = '%s'`, string(dpname))
	row, err := g_ds.QueryRow(sqlTotal)
	if err != nil {
		SqlExecError(rw, result, err.Error())
		return
	}
	var total int
	row.Scan(&total)
	if total == 0 {
		result.Msg = fmt.Sprintf("Datapool %v not found.", dpname)
		result.Code = cmd.ErrorNoRecord
		resp, _ := json.Marshal(result)
		fmt.Fprintln(rw, string(resp))
		return
	}

	sqlDp := fmt.Sprintf(`SELECT DPID, DPNAME, DPTYPE, DPCONN FROM DH_DP 
		WHERE STATUS = 'A' AND DPNAME = '%s'`, dpname)
	rowdp, err := g_ds.QueryRow(sqlDp)
	if err != nil {
		SqlExecError(rw, result, err.Error())
		return
	}

	var dpid int
	onedp.Items = make([]cmd.Item, 0, 16)
	rowdp.Scan(&dpid, &onedp.Name, &onedp.Type, &onedp.Conn)
	if dpid > 0 {
		//Use "left out join" to get repository/dataitem records, whether it has tags or not.
		//B.STATUS='A'
		sqlTag := fmt.Sprintf(`SELECT A.REPOSITORY, A.DATAITEM, A.ITEMDESC, A.PUBLISH ,strftime(A.CREATE_TIME), B.TAGNAME, B.DETAIL,strftime(B.CREATE_TIME)
				FROM DH_DP_RPDM_MAP A LEFT JOIN DH_RPDM_TAG_MAP B
				ON (A.RPDMID = B.RPDMID)
				WHERE A.DPID = %v AND A.STATUS='A' `, dpid)
		tagrows, err := g_ds.QueryRows(sqlTag)
		if err != nil {
			SqlExecError(rw, result, err.Error())
			return
		}
		defer tagrows.Close()
		for tagrows.Next() {
			item := cmd.Item{}
			var repoitemtime string
			tagrows.Scan(&item.Repository, &item.DataItem, &item.ItemDesc, &item.Publish, &repoitemtime, &item.Tag, &item.TagDetail, &item.Time)
			if len(item.Time) == 0 {
				item.Time = repoitemtime
			}
			//log.Println(item.Repository, item.DataItem, item.Tag, item.Time, item.Publish)
			onedp.Items = append(onedp.Items, item)
		}
	}
	resp, _ := json.Marshal(result)
	log.Println(string(resp))
	fmt.Fprintln(rw, string(resp))

}

func SqlExecError(rw http.ResponseWriter, result *ds.Result, msg string) {
	result.Msg = msg
	result.Code = cmd.ErrorSqlExec
	resp, _ := json.Marshal(result)
	fmt.Fprintln(rw, string(resp))
}

func dpDeleteOneHandler(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()

	dpname := ps.ByName("dpname")
	msg := &ds.MsgResp{}

	sql_dp_rm := fmt.Sprintf(`SELECT DPID, DPTYPE FROM DH_DP WHERE DPNAME ='%s' AND STATUS='A'`, dpname)
	dprows, err := g_ds.QueryRows(sql_dp_rm)
	if err != nil {
		msg.Msg = err.Error()
	}

	bresultflag := false

	dpid_type := make([]strc_dp, 0, 8)
	strcone := strc_dp{}
	for dprows.Next() {
		dprows.Scan(&strcone.Dpid, &strcone.Dptype)
		dpid_type = append(dpid_type, strcone)
	}
	dprows.Close()

	for _, v := range dpid_type {
		var dpid = v.Dpid
		var dptype = v.Dptype
		bresultflag = true
		//dprow.Scan(&dpid, &dptype)
		sql_dp_item := fmt.Sprintf("SELECT PUBLISH FROM DH_DP_RPDM_MAP WHERE DPID = %v ", dpid)
		row, err := g_ds.QueryRow(sql_dp_item)
		if err != nil {
			msg.Msg = err.Error()
		}
		//time.Sleep(60*time.Second)
		var sPublish string
		row.Scan(&sPublish)
		if sPublish == "Y" {
			msg.Msg = fmt.Sprintf(`Datapool %s with type:%s can't be removed , it contains published DataItem !`, dpname, dptype)
		} else {
			sql_update := fmt.Sprintf("UPDATE DH_DP SET STATUS = 'N' WHERE DPID = %v", dpid)
			_, err := g_ds.Update(sql_update)
			if err != nil {
				msg.Msg = err.Error()
			} else {
				msg.Msg = fmt.Sprintf("Datapool %s with type:%s removed successfully!", dpname, dptype)
			}
		}
		rw.WriteHeader(http.StatusOK)
		resp, _ := json.Marshal(msg)
		respStr := string(resp)
		fmt.Fprintln(rw, respStr)
	}
	if bresultflag == false {
		rw.WriteHeader(http.StatusOK)
		msg.Msg = fmt.Sprintf("Datapool %s not found.", dpname)
		resp, _ := json.Marshal(msg)
		fmt.Fprintln(rw, string(resp))
	}
}
