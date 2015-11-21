package daemon

import (
	"errors"
	"fmt"
	"github.com/asiainfoLDP/datahub/ds"
	log "github.com/asiainfoLDP/datahub/utils/clog"
)

func CheckDataPoolExist(datapoolname string) (bexist bool) {
	sqlcheck := fmt.Sprintf("SELECT COUNT(1) FROM DH_DP WHERE DPNAME='%s' AND STATUS='A'", datapoolname)
	row, err := g_ds.QueryRow(sqlcheck)
	//fmt.Println(sqlcheck)
	if err != nil {
		log.Println("CheckDataPoolExist QueryRow error:", err.Error())
		return
	} else {
		var num int
		row.Scan(&num)
		//fmt.Println("num:", num)
		if num == 0 {
			return false
		} else {
			return true
		}
	}
}

func GetDataPoolDpconn(datapoolname string) (dpconn string) {
	sqlgetdpconn := fmt.Sprintf("SELECT DPCONN FROM DH_DP WHERE DPNAME='%s'  AND STATUS='A'", datapoolname)
	//fmt.Println(sqlgetdpconn)
	row, err := g_ds.QueryRow(sqlgetdpconn)
	if err != nil {
		log.Println("GetDataPoolDpconn QueryRow error:", err.Error())
		return
	} else {
		row.Scan(&dpconn)
		return dpconn
	}
}

func GetDataPoolDpid(datapoolname string) (dpid int) {
	sqlgetdpid := fmt.Sprintf("SELECT DPID FROM DH_DP WHERE DPNAME='%s'  AND STATUS='A'", datapoolname)
	//fmt.Println(sqlgetdpid)
	row, err := g_ds.QueryRow(sqlgetdpid)
	if err != nil {
		log.Println("GetDataPoolDpid QueryRow error:", err.Error())
		return
	} else {
		row.Scan(&dpid)
		return
	}
}

func InsertTagToDb(dpexist bool, p ds.DsPull) (err error) {
	if dpexist == false {
		return
	}
	DpId := GetDataPoolDpid(p.Datapool)
	if DpId == 0 {
		return
	}
	rpdmid := GetRepoItemId(p.Repository, p.Dataitem)
	//fmt.Println("GetRepoItemId1", rpdmid, DpId)
	if rpdmid == 0 {
		sqlInsertRpdm := fmt.Sprintf(`INSERT INTO DH_DP_RPDM_MAP(RPDMID ,REPOSITORY , DATAITEM, 
        	DPID  , PUBLISH ,CREATE_TIME ) VALUES (null, '%s', '%s', %d, 'N', datetime('now'))`,
			p.Repository, p.Dataitem, DpId)
		g_ds.Insert(sqlInsertRpdm)
		rpdmid = GetRepoItemId(p.Repository, p.Dataitem)
		//fmt.Println("GetRepoItemId2", rpdmid, DpId)
	}
	sqlInsertTag := fmt.Sprintf(`INSERT INTO DH_RPDM_TAG_MAP(TAGNAME ,RPDMID ,DETAIL,CREATE_TIME) 
		VALUES ('%s', '%d', '%s', datetime('now'))`,
		p.Tag, rpdmid, p.DestName)
	log.Println(sqlInsertTag)
	_, err = g_ds.Insert(sqlInsertTag)
	return err
}

func GetRepoItemId(repository, dataitem string) (rpdmid int) {
	sqlgetrpdmId := fmt.Sprintf("SELECT RPDMID FROM DH_DP_RPDM_MAP WHERE REPOSITORY='%s' AND DATAITEM='%s'",
		repository, dataitem)
	row, err := g_ds.QueryRow(sqlgetrpdmId)
	if err != nil {
		log.Println("GetRepoItemId QueryRow error:", err.Error())
		return
	} else {
		row.Scan(&rpdmid)
		return
	}
}

func InsertItemToDb(repo, item, datapool string) (err error) {
	dpid := GetDataPoolDpid(datapool)
	if dpid > 0 {
		sqlInsertItem := fmt.Sprintf(`INSERT INTO DH_DP_RPDM_MAP (RPDMID, REPOSITORY, DATAITEM, DPID, PUBLISH, CREATE_TIME)
			VALUES (null, '%s', '%s', %d, 'Y',  datetime('now'))`, repo, item, dpid)
		_, err = g_ds.Insert(sqlInsertItem)
		log.Println(sqlInsertItem)

	} else {
		err = errors.New("dpid is not found")
	}
	return err
}

func GetDataPoolStatusByID(dpid int) (status string) {
	sqlGetDpStatus := fmt.Sprintf("SELECT STATUS FROM DH_DP WHERE DPID=%d", dpid)
	fmt.Println(sqlGetDpStatus)
	row, err := g_ds.QueryRow(sqlGetDpStatus)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	row.Scan(&status)
	log.Println("GetDataPoolStatusByID status:", status)
	return
}

func GetRpdmIdAndDpId(repo, item string) (rpdmid, dpid int) {
	sqlGetRpdmIdAndDpId := fmt.Sprintf("SELECT RPDMID, DPID FROM DH_DP_RPDM_MAP WHERE REPOSITORY='%s' AND DATAITEM='%s'", repo, item)
	row, err := g_ds.QueryRow(sqlGetRpdmIdAndDpId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	row.Scan(&rpdmid, &dpid)
	status := GetDataPoolStatusByID(dpid)
	if rpdmid == 0 || dpid == 0 {
		log.Println("GetRpdmIdAndDpId rpdmid, dpid: ", rpdmid, dpid)
		log.Println("datapool status:", status)
	}
	if status != "A" {
		return 0, 0
	}

	return
}

func CheckTagExist(repo, item, tag string) (exits bool, err error) {
	rpdmid, dpid := GetRpdmIdAndDpId(repo, item)
	if rpdmid == 0 || dpid == 0 {
		fmt.Println("rpdmid, dpid ", rpdmid, dpid)
		return false, errors.New("repo and dataitem not exist")
	}
	sqlCheckTag := fmt.Sprintf("SELECT COUNT(1) FROM DH_RPDM_TAG_MAP WHERE RPDMID='%d' AND TAGNAME='%s'", rpdmid, tag)
	row, err := g_ds.QueryRow(sqlCheckTag)
	var count int
	row.Scan(&count)
	if count > 0 {
		return true, nil
	}
	return
}

func GetDpNameAndDpConn(repo, item, tag string) (dpname, dpconn string) {
	_, dpid := GetRpdmIdAndDpId(repo, item)
	if dpid == 0 {
		log.Println("GetDpNameAndDpConn dpid==0")
		return "", ""
	}
	dpname, dpconn = GetDpnameDpconnByDpidAndStatus(dpid, "A")
	return
}

func GetDpnameDpconnByDpidAndStatus(dpid int, status string) (dpname, dpconn string) {
	sqlgetdpconn := fmt.Sprintf("SELECT DPNAME ,DPCONN FROM DH_DP WHERE DPID='%d'  AND STATUS='%s'", dpid, status)
	//fmt.Println(sqlgetdpconn)
	row, err := g_ds.QueryRow(sqlgetdpconn)
	if err != nil {
		log.Println("GetDpnameDpconnByDpidAndStatus QueryRow error:", err.Error())
		return
	} else {
		row.Scan(&dpname, &dpconn)
		return
	}
	return
}

func InsertPubTagToDb(repo, item, tag, FileName string) (err error) {
	rpdmid := GetRepoItemId(repo, item)
	if rpdmid == 0 {
		return errors.New("Dataitem is not found which need to be published before publishing tag. ")
	}
	sqlInsertTag := fmt.Sprintf("INSERT INTO DH_RPDM_TAG_MAP (TAGNAME, RPDMID, DETAIL, CREATE_TIME) VALUES ('%s', %d, '%s', datetime('now'))",
		tag, rpdmid, FileName)
	log.Println(sqlInsertTag)
	_, err = g_ds.Insert(sqlInsertTag)
	if err != nil {
		return err
	}
	return
}
