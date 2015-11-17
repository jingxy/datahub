package ds

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type FormatRepoList struct {
	RepoName string `json:"repoName"`
	ItemID   string `json:"itemID"`
}

type DsPull struct {
	Tag        string `json:"tag"`
	Datapool   string `json:"datapool"`
	DestName   string `json:"destname"`
	Repository string `json:"repository, omitempty"`
	Dataitem   string `jsong:"dataitem, omitempty"`
}
type Result struct {
	Code int         `json:"code,omitempty"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

const (
	DB_DML_INSERT = "insert"
	DB_DML_DELETE = "delete"
	DB_DML_UPDATE = "update"
	DB_DML_SELECT = "select"
	DB_DDL_CREATE = "create"
	DB_DDL_DROP   = "drop"
	TABLE_ORDER   = "order_t"
	TABLE_USER    = "user"
)

type MsgResp struct {
	Msg string `json:"msg"`
}

type DataItem struct {
	Repository_name string `json:"repname,omitempty"`
	Dataitem_name   string `json:"dataitem_name,omitempty"`
	/*
		User_id         int     `json:"user_id,omitempty"`
		Dataitem_id     int64   `json:"dataitem_id,omitempty"  xorm:"dataitem_id pk autoincr"`
		Ico_name        string  `json:"ico_name,omitempty"`
		Permit_type     int     `json:"permit_type,omitempty"`
		Key_words       string  `json:"key_words,omitempty"`
		Supply_style    int     `json:"supply_style,omitempty"`
		Priceunit_type  int     `json:"priceunit_type,omitempty"`
		Price           float32 `json:"price,omitempty"`
		Optime          string  `json:"optime,omitempty"`
		Data_format     int     `json:"data_format,omitempty"`
		Refresh_type    string  `json:"refresh_type,omitempty"`
		Refresh_num     int     `json:"refresh_num,omitempty"`
		Meta_filename   string  `json:"meta_filename,omitempty"`
		Sample_filename string  `json:"sample_filename,omitempty"`
		Comment         string  `json:"comment,omitempty"`
	*/
}

type DataItemUsage struct {
	Dataitem_id   int64  `json:"dataitem_id,omitempty"  xorm:"dataitem_id pk autoincr"`
	Dataitem_name string `json:"dataitem_name,omitempty"`
	Views         int    `json:"views"`
	Follows       int    `json:"follows"`
	//	Downloads     int    `json:"downloads"`
	Stars        int    `json:"stars"`
	Refresh_date string `json:"refresh_date,omitempty"`
	Usability    int    `json:"usability,omitempty"`
}

type Tag struct {
	Dataitem_id int64  `json:"dataitem_id,omitempty"`
	Tag         string `json:"tag,omitempty"`
	Filename    string `json:"filename,omitempty"`
	Optime      string `json:"optime,omitempty"`
}

type Data struct {
	Repository_name string `json:"repname,omitempty"`
	Dataitem_name   string `json:"itemname,omitempty"`
	//Usage *DataItemUsage `json:"statis,omitempty"`
	Tagsnum int   `json:"tags",omitempty`
	Tags    []Tag `json:"Tags,omitempty"`
}

/*
type ItemDetail struct {
	Item  *DataItem      `json:"item,omitempty"`
	Usage *DataItemUsage `json:"statis,omitempty"`
	Tags  []Tag          `json:"tags,omitempty"`
}
*/
type Ds struct {
	Db *sql.DB
}

const Create_dh_dp string = `CREATE TABLE IF NOT EXISTS 
    DH_DP ( 
       DPID    INTEGER PRIMARY KEY AUTOINCREMENT, 
       DPNAME  VARCHAR(32), 
       DPTYPE  VARCHAR(32), 
       DPCONN  VARCHAR(128), 
       STATUS  CHAR(2) 
    );`

//DH_DP STATUS : 'A' valid; 'N' invalid; 'P' contain dataitem published;

const Create_dh_dp_repo_ditem_map string = `CREATE TABLE IF NOT EXISTS 
    DH_DP_RPDM_MAP ( 
    	RPDMID       INTEGER PRIMARY KEY AUTOINCREMENT, 
        REPOSITORY   VARCHAR(128), 
        DATAITEM     VARCHAR(128), 
        DPID         INTEGER, 
        PUBLISH      CHAR(2), 
        CREATE_TIME  DATETIME 
    );`

//DH_DP_REPO_DITEM_MAP  PUBLISH: 'Y' the dataitem is published by you,
//'N' the dataitem is pulled by you

const Create_dh_repo_ditem_tag_map string = `CREATE TABLE IF NOT EXISTS 
    DH_RPDM_TAG_MAP ( 
        TAGNAME      VARCHAR(128),
        RPDMID       INTEGER,
        DETAIL       VARCHAR(128),
        CREATE_TIME  DATETIME
    );`

type Executer interface {
	Insert(cmd string) (interface{}, error)
	Delete(cmd string) (interface{}, error)
	Update(cmd string) (interface{}, error)
	QueryRaw(cmd string) (*sql.Rows, error)
	QueryRaws(cmd string) (*sql.Rows, error)

	Create(cmd string) (interface{}, error)
	Drop(cmd string) (interface{}, error)
}

func execute(p *Ds, cmd string) (interface{}, error) {
	tx, err := p.Db.Begin()
	if err != nil {
		return nil, err
	}
	var res sql.Result
	if res, err = tx.Exec(cmd); err != nil {
		log.Printf(`Exec("%s") err %s`, cmd, err.Error())
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return res, nil
}

func query(p *Ds, cmd string) (*sql.Row, error) {
	return p.Db.QueryRow(cmd), nil
}
func queryRows(p *Ds, cmd string) (*sql.Rows, error) {
	return p.Db.Query(cmd)
}

func (p *Ds) Insert(cmd string) (interface{}, error) {
	return execute(p, cmd)
}

func (p *Ds) Delete(cmd string) (interface{}, error) {
	return execute(p, cmd)
}

func (p *Ds) Update(cmd string) (interface{}, error) {
	return execute(p, cmd)
}

func (p *Ds) QueryRow(cmd string) (*sql.Row, error) {
	return query(p, cmd)
}

func (p *Ds) QueryRows(cmd string) (*sql.Rows, error) {
	return queryRows(p, cmd)
}
func (p *Ds) Create(cmd string) (interface{}, error) {
	return execute(p, cmd)
}

func (p *Ds) Drop(cmd string) (interface{}, error) {
	return execute(p, cmd)
}
