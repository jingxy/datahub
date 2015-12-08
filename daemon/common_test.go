package daemon

import (
	"fmt"
	"github.com/asiainfoLDP/datahub/ds"
	"os"
	"testing"
)

var (
	Name  = "DatahubUnitTest"
	Type  = "file"
	Conn  = "/var/lib/datahub/datahubUnitTest"
	testP = ds.DsPull{
		Tag:        "tagTest",
		Datapool:   Name,
		DestName:   "tagTestD.txt",
		Repository: "repotest",
		Dataitem:   "itemtest",
		ItemDesc:   "dirrepoitemtest",
	}
)

func TestCheckDataPoolExist(t *testing.T) {
	b := CheckDataPoolExist("DatahubUnitTest")
	if b {
		t.Errorf("1.test CheckDataPoolExist fail--------")
	} else {
		t.Log("1.test CheckDataPoolExist success--------")
	}

	sqldpinsert := fmt.Sprintf(`insert into DH_DP (DPID, DPNAME, DPTYPE, DPCONN, STATUS)
					values (null, '%s', '%s', '%s', 'A')`, Name, Type, Conn)
	if _, err := g_ds.Insert(sqldpinsert); err != nil {
		os.Remove("/var/lib/datahub/datahubUnitTest")
		t.Error(err)
	}
	b = CheckDataPoolExist("DatahubUnitTest")
	if !b {
		t.Error("2.test CheckDataPoolExist fail--------")
	} else {
		t.Log("2.test CheckDataPoolExist success--------")
	}
}

func TestGetDataPoolDpconn(t *testing.T) {
	conn := GetDataPoolDpconn("DatahubUnitTest")
	if conn == Conn {
		t.Logf("1.GetDataPoolDpconn %s success--------\n", conn)
	} else {
		t.Errorf("1.test GetDataPoolDpconn %s fail--------\n", conn)
	}
}

func TestGetDataPoolDpid(t *testing.T) {
	if id := GetDataPoolDpid(Name); id == 0 {
		t.Errorf("1.GetDataPoolDpid fail--------dpid:%d\n", id)
	} else {
		t.Log("1.GetDataPoolDpid success--------dpid:", id)
	}
}

func TestInsertTagToDb(t *testing.T) {
	if e := InsertTagToDb(true, testP); e != nil {
		t.Error("1.InsertTagToDb fail--------", e)
	} else {
		t.Log("1.InsertTagToDb success--------")
	}
}

func TestGetRepoItemId(t *testing.T) {
	if id := GetRepoItemId(testP.Repository, testP.Dataitem); id == 0 {
		t.Errorf("1.GetRepoItemId fail-------- repo:%s, item:%s\n", testP.Repository, testP.Dataitem)
	} else {
		t.Logf("1.GetRepoItemId success-------- repo:%s, item:%s\n", testP.Repository, testP.Dataitem)
	}
}

func TestInsertItemToDb(t *testing.T) {
	if e := DeleteRepoItemHard(testP.Repository, testP.Dataitem); e != nil {
		t.Errorf("Recover db for TestInsertItemToDb test fail. Delete %s, %s error, %v", testP.Repository, testP.Dataitem, e)
	}
	if e := DeleteTagHard(testP.Tag); e != nil {
		t.Errorf("Recover db for TestInsertItemToDb test fail. Delete %s error, %v", testP.Tag, e)
	}
	err := InsertItemToDb(testP.Repository, testP.Dataitem, testP.Datapool, testP.ItemDesc)
	if err != nil {
		t.Error("1.InsertItemToDb fail-------- ", err)
	} else {
		t.Log("1.InsertItemToDb success---------")
	}
	if e := InsertTagToDb(true, testP); e != nil {
		t.Error("InsertTagToDb fail--------", e)
	}

	err = InsertItemToDb(testP.Repository, testP.Dataitem, "No_this_datapool", testP.ItemDesc)
	if err != nil {
		t.Log("2.InsertItemToDb success:insert item to a non-existent datapool---------")
	} else {
		t.Error("2.InsertItemToDb fail:insert item to a non-existent datapool---------")
	}
}

func TestGetRpdmidDpidItemdesc(t *testing.T) {
	r, d, desc := GetRpdmidDpidItemdesc(testP.Repository, testP.Dataitem)
	if r == 0 || d == 0 {
		t.Errorf("1.TestGetRpdmidDpidItemdesc fail-------- rpdmid:%d, dpid:%d, desc:%s\n", r, d, desc)
	} else {
		t.Logf("1.TestGetRpdmidDpidItemdesc success-------- rpdmid:%d, dpid:%d, desc:%s\n", r, d, desc)
	}

	r, d, desc = GetRpdmidDpidItemdesc("No_this_repo", testP.Dataitem)
	if r == 0 && d == 0 {
		t.Logf("2.TestGetRpdmidDpidItemdesc success-------- rpdmid:%d, dpid:%d, desc:%s\n", r, d, desc)
	} else {
		t.Errorf("2.TestGetRpdmidDpidItemdesc fail-------- rpdmid:%d, dpid:%d, desc:%s\n", r, d, desc)
	}
}

func TestCheckTagExist(t *testing.T) {
	exist, e := CheckTagExist(testP.Repository, testP.Dataitem, testP.Tag)
	if exist == false || e != nil {

	}
}

func TestRecover(t *testing.T) {
	if e := DeleteDpHard(Name); e != nil {
		t.Errorf("Recover db for common.go test fail. Delete %s error, %v", Name, e)
	}
	if e := DeleteRepoItemHard(testP.Repository, testP.Dataitem); e != nil {
		t.Errorf("Recover db for common.go test fail. Delete %s, %s error, %v", testP.Repository, testP.Dataitem, e)
	}
	if e := DeleteTagHard(testP.Tag); e != nil {
		t.Errorf("Recover db for common.go test fail. Delete %s error, %v", testP.Tag, e)
	}
}

func DeleteRepoItemHard(repo, item string) (e error) {
	sql := fmt.Sprintf("DELETE FROM DH_DP_RPDM_MAP WHERE REPOSITORY='%s' AND DATAITEM='%s'", repo, item)
	_, e = g_ds.Delete(sql)
	return e
}

func DeleteTagHard(tag string) (e error) {
	sql := fmt.Sprintf("DELETE FROM DH_RPDM_TAG_MAP WHERE TAGNAME='%s'", tag)
	_, e = g_ds.Delete(sql)
	return e
}
