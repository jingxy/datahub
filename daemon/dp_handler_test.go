package daemon

import (
	"database/sql"
	"encoding/json"
	"github.com/asiainfoLDP/datahub/cmd"
	"github.com/asiainfoLDP/datahub/ds"
	"github.com/julienschmidt/httprouter"
	//"log"
	//"net"
	//"errors"
	//"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	g_dbfileU string = "/var/lib/datahub/datahubUtest.db"
)

type Context struct {
	desc        string
	requestBody string
	ps          httprouter.Params
	rw          *httptest.ResponseRecorder
	out         string
	d           interface{}
}

func init() {

	fmt.Println("connect to db sqlite3----unit-test")
	db, err := sql.Open("sqlite3", g_dbfileU)
	if err != nil {
		fmt.Println(err)
	}
	//defer db.Close()
	chk(err)
	g_ds.Db = db
	g_ds.Create(ds.Create_dh_dp)
	g_ds.Create(ds.Create_dh_dp_repo_ditem_map)
	g_ds.Create(ds.Create_dh_repo_ditem_tag_map)
}

func Test_dpPostOneHandler(t *testing.T) {

	con := []Context{
		Context{
			desc: "1.create datapool--------",
			requestBody: `{
									"dpname":"datahub_unit_test",
									"dptype":"file",
									"dpconn":"/var/lib/datahub/datahub-Unit-Test"						
							 }`,
			d: cmd.FormatDpCreate{
				Name: "datahub_unit_test",
				Type: "file",
				Conn: "/var/lib/datahub/datahub-Unit-Test",
			},
			rw:  httptest.NewRecorder(),
			out: "dp create success. name:datahub_unit_test type:file path:/var/lib/datahub/datahub-Unit-Test",
		},
		Context{
			desc: "2.create dup datapool err--------",
			requestBody: `{
									"dpname":"datahub_unit_test",
									"dptype":"file",
									"dpconn":"/var/lib/datahub/datahub-Unit-Test"						
							 }`,
			d: cmd.FormatDpCreate{
				Name: "datahub_unit_test",
				Type: "file",
				Conn: "/var/lib/datahub/datahub-Unit-Test",
			},
			rw:  httptest.NewRecorder(),
			out: "The datapool datahub_unit_test is already exist, please use another name!",
		},
	}

	for _, v := range con {
		req, _ := http.NewRequest("POST", "/datapools", strings.NewReader(v.requestBody))
		dpPostOneHandler(v.rw, req, v.ps)
		if !Expect(t, v.rw, v.out) {
			t.Logf("%s fail!", v.desc)
		} else {
			t.Logf("%s success.", v.desc)
		}
	}
	for _, v := range con {
		dpc, ok := v.d.(cmd.FormatDpCreate)
		if ok {
			if e := DeleteDpHard(dpc.Name); e != nil {
				t.Logf("DeleteDpHard error. %v", e)
			}
			fmt.Println("-----delete----", dpc.Name)
		}
	}
}

func DeleteDpHard(dp string) (e error) {
	sqlDelDp := fmt.Sprintf("DELETE FROM DH_DP WHERE DPNAME='%s'", dp)
	_, e = g_ds.Delete(sqlDelDp)
	return e
}

func Expect(t *testing.T, rw *httptest.ResponseRecorder, out string) bool {
	msg := &ds.MsgResp{}
	body, _ := ioutil.ReadAll(rw.Body)

	if err := json.Unmarshal(body, &msg); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(msg.Msg)
	}
	if msg.Msg != out {
		t.Errorf("expected http.Msg(%s) != return http.Msg(%s)", out, msg.Msg)
		return false
	}
	return true
}
