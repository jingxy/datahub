package ds

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

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
}

type Tag struct {
	Dataitem_id int64  `json:"dataitem_id,omitempty"`
	Tag         string `json:"tag,omitempty"`
	Filename    string `json:"filename,omitempty"`
	Optime      string `json:"optime,omitempty"`
	Comment     string `json:"comment,omitempty"`
}

type Data struct {
	Repository_name string `json:"repname,omitempty"`
	Dataitem_name   string `json:"itemname,omitempty"`
	//Usage *DataItemUsage `json:"statis,omitempty"`
	Tagsnum int   `json:"tags",omitempty`
	Taglist []Tag `json:"taglist,omitempty"`
}

type Repositories struct {
	RepositoryName string `json:"repname,omitempty"`
	Comment        string `json:"comment,omitempty"`
	Optime         string `json:"optime,omitempty"`
}

type Repository struct {
	DataItems []string `json:"dataitems, omitempty"`
}

type PubPara struct {
	Datapool   string `json:"datapool, omitempty"`
	Detail     string `json:"detail, omitempty"`
	Accesstype string `json:"itemaccesstype, omitempty"`
	Comment    string `json:"comment, omitempty"`
}

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
//TAGID        INTEGER PRIMARY KEY AUTOINCREMENT,
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
