package db

import (
	"time"
  log "../log"
	"github.com/ziutek/mymysql/mysql"			// Mysql
	_ "github.com/ziutek/mymysql/native" 	// Mysql Native engine
)

var DB_Connections map[string]map[string]string

type Transaction struct {
	conn mysql.Conn
	tr mysql.Transaction
}

func Open(conn_name string) mysql.Conn{
	var db mysql.Conn
	db = mysql.New("tcp", "",
      DB_Connections[conn_name]["server"],
      DB_Connections[conn_name]["user"],
      DB_Connections[conn_name]["pass"],
      DB_Connections[conn_name]["dbname"])
	db.SetTimeout(15 * time.Second)
	db.Register("set names utf8")
	err := db.Connect()
	if err != nil {
    log.Log("db.Open", conn_name, err.Error(), "", "SQL")
    panic("error.db_connection_failed")
  }
	return db
}

func Close(conn_name string, db mysql.Conn) {
	if db != nil && db.IsConnected() {
		 err := db.Close()
		 db = nil
		 if err != nil {
       log.Log("db.Close", conn_name, err.Error(), "", "SQL")
     }
	}
}

func Execute(conn_name string, sql string) bool{
	db := Open(conn_name)
	defer func() {
    Close(conn_name, db)
	}()
	_, _, err := db.Query(sql)
	if err != nil {
    log.Log("db.Execute", conn_name, err.Error(), sql, "SQL")
    panic(err.Error())
  }
	return false
}

func Query(conn_name string, sql string) []map[string]interface{}{
  db := Open(conn_name)
	defer func() {
    Close(conn_name, db)
	}()

	records := []map[string]interface{}{}

	rows, res, err := db.Query(sql)
	if err != nil {
    log.Log("db.Query", conn_name, err.Error(), sql, "SQL")
    panic(err.Error())
  }

	for _, row := range rows {
		record := make(map[string]interface{})
		for _,field := range res.Fields() {
			switch field.Type {
			case 0x1, 0x3, 0x8:			//tinyint, int, bigint
				record[field.Name] = row.Int64(res.Map(field.Name))
			case 0x4, 0x5:					//float, double
				record[field.Name] = row.Float(res.Map(field.Name))
			default:
				record[field.Name] = row.Str(res.Map(field.Name))
			}
		}
		records = append(records, record)
	}

	for res != nil {
		if res, err = res.NextResult(); err != nil {
			return records
		}
	}

	return records
}

func OpenTrans(conn_name string) (Transaction){
	conn := Open(conn_name)
	tr, err := conn.Begin()
	if err != nil {
    log.Log("db.OpenTrans", conn_name, err.Error(), "", "SQL")
    panic(err.Error())
  }
	return Transaction{conn:conn, tr:tr}
}

func (trans *Transaction) Close() {
  if trans != nil && trans.conn.IsConnected() {
		 err := trans.conn.Close()
		 trans = nil
		 if err != nil {
       log.Log("trans.Close", "", err.Error(), "", "SQL")
     }
	}
}

func (trans *Transaction) SetTimeout(timeout time.Duration) {
	trans.conn.SetTimeout(timeout)
}

func (trans *Transaction) Commit(){
	err := trans.tr.Commit()
  if err != nil {
    log.Log("trans.Commit", "", err.Error(), "", "SQL")
    panic(err.Error())
  }
}

func (trans *Transaction) Rollback(){
	err := trans.tr.Rollback()
  if err != nil {
    log.Log("trans.Rollback", "", err.Error(), "", "SQL")
    panic(err.Error())
  }
}

func (trans *Transaction) Execute(sql string) {
	_, err := trans.tr.Start(sql)
  if err != nil {
    log.Log("trans.Execute", "", err.Error(), sql, "SQL")
    panic(err.Error())
  }
}

func (trans *Transaction) Query(sql string) []map[string]interface{} {

	sel, err := trans.conn.Prepare(sql)
  if err != nil {
    log.Log("trans.Query", "", err.Error(), sql, "SQL")
    panic(err.Error())
  }

	rows, res, err := trans.tr.Do(sel).Exec()
  if err != nil {
    log.Log("trans.Query", "", err.Error(), sql, "SQL")
    panic(err.Error())
  }

  records := []map[string]interface{}{}
	for _, row := range rows {
		record := make(map[string]interface{})
		for _,field := range res.Fields() {
			switch field.Type {
			case 0x1, 0x3, 0x8:			//tinyint, int, bigint
				record[field.Name] = row.Int64(res.Map(field.Name))
			case 0x4, 0x5:					//float, double
				record[field.Name] = row.Float(res.Map(field.Name))
			default:
				record[field.Name] = row.Str(res.Map(field.Name))
			}
		}
		records = append(records, record)
	}

	for res != nil {
		if res, err = res.NextResult(); err != nil {
			return records
		}
	}

	return records
}
