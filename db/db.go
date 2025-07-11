package db

import (
	"os"
	"fmt"
	"time"
	"strconv"
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
    Log("db.Open", conn_name, err.Error(), "", "DB_ERROR")
    panic("error.DBOperationFailed")
  }
	return db
}

func Close(conn_name string, db mysql.Conn) {
	if db != nil && db.IsConnected() {
		 err := db.Close()
		 db = nil
		 if err != nil {
       Log("db.Close", conn_name, err.Error(), "", "DB_ERROR")
			 panic("error.DBOperationFailed")
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
    Log("db.Execute", conn_name, err.Error(), sql, "DB_ERROR")
    panic("error.DBOperationFailed")
  }
	return false
}

func Query(conn_name string, sql string) []map[string]interface{}{
	runTime := time.Now()
  db := Open(conn_name)
	defer func() {
    Close(conn_name, db)
	}()

	records := []map[string]interface{}{}

	rows, res, err := db.Query(sql)
	if err != nil {
    Log("db.Query", conn_name, err.Error(), sql, "DB_ERROR")
    panic("error.DBOperationFailed")
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

	if DateTimeValueDiff(runTime, time.Now()) > 30 {
		LogHiddenWithDuration("Query", "", "", sql, I64_S(DateTimeValueDiff(runTime, time.Now()))+"s", "DB_SLOWQUERY")
	}

	return records
}

func OpenTrans(conn_name string) (Transaction){
	conn := Open(conn_name)
	tr, err := conn.Begin()
	if err != nil {
    Log("db.OpenTrans", conn_name, err.Error(), "", "DB_ERROR")
    panic("error.DBOperationFailed")
  }
	return Transaction{conn:conn, tr:tr}
}

func (trans *Transaction) Close() {
  if trans != nil && trans.conn.IsConnected() {
		 err := trans.conn.Close()
		 trans = nil
		 if err != nil {
       Log("trans.Close", "", err.Error(), "", "DB_ERROR")
     }
	}
}

func (trans *Transaction) SetTimeout(timeout time.Duration) {
	trans.conn.SetTimeout(timeout)
}

func (trans *Transaction) Commit(){
	if trans != nil && trans.conn.IsConnected() {
		err := trans.tr.Commit()
	  if err != nil {
	    Log("trans.Commit", "", err.Error(), "", "DB_ERROR")
	    panic("error.DBOperationFailed")
	  }
	}
}

func (trans *Transaction) Rollback(){
	if trans != nil && trans.conn.IsConnected() {
		err := trans.tr.Rollback()
	  if err != nil {
	    Log("trans.Rollback", "", err.Error(), "", "DB_ERROR")
	    panic("error.DBOperationFailed")
	  }
	}
}

func (trans *Transaction) Execute(sql string) {
	_, err := trans.tr.Start(sql)
  if err != nil {
    Log("trans.Execute", "", err.Error(), sql, "DB_ERROR")
    panic("error.DBOperationFailed")
  }
}

func (trans *Transaction) Query(sql string) []map[string]interface{} {

	runTime := time.Now()

	sel, err := trans.conn.Prepare(sql)
  if err != nil {
    Log("trans.Query (Prepare)", "", err.Error(), sql, "DB_ERROR")
    panic("error.DBOperationFailed")
  }

	rows, res, err := trans.tr.Do(sel).Exec()
  if err != nil {
    Log("trans.Query", "", err.Error(), sql, "DB_ERROR")
    panic("error.DBOperationFailed")
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

	if DateTimeValueDiff(runTime, time.Now()) > 30 {
		LogHiddenWithDuration("trans.Query", "", "", sql, I64_S(DateTimeValueDiff(runTime, time.Now()))+"s", "DB_SLOWQUERY")
	}

	return records
}

func writeLog(operation string, username string, key string, msg string, duration string, logfilename string, showDisplay bool){
	t := time.Now()

  if _, err := os.Stat("logs/"); os.IsNotExist(err) {
	   os.Mkdir("logs/", os.ModePerm)
	}

	logdatepath := "logs/"+t.Format("060102")
	if _, err := os.Stat(logdatepath); os.IsNotExist(err) {
	    os.Mkdir(logdatepath, os.ModePerm)
	}

	file, err := os.OpenFile(logdatepath+"/"+logfilename+".txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
  if err != nil {
    //log.Fatal("Cannot create file", err)
  }
  defer file.Close()
	if showDisplay {
		fmt.Print("Log|o="+operation+"|u="+username+"|k="+key+"|d="+duration+"|m=["+msg+"] => "+logfilename+"\r\n")
	}
  fmt.Fprintf(file, "t="+t.Format("15:04:05.000")+"|o="+operation+"|u="+username+"|k="+key+"|d="+duration+"|m=["+msg+"]\r\n")
}

func Log(operation string, username string, key string, msg string, logfilename string){
	writeLog(operation, username, key, msg, "", logfilename, true)
}

func LogHidden(operation string, username string, key string, msg string, logfilename string){
	writeLog(operation, username, key, msg, "", logfilename, false)
}

func LogHiddenWithDuration(operation string, username string, key string, msg string, duration string, logfilename string){
	writeLog(operation, username, key, msg, duration, logfilename, false)
}

func LogWithDuration(operation string, username string, key string, msg string, duration string, logfilename string){
	writeLog(operation, username, key, msg, duration, logfilename, true)
}

func DateTimeValueDiff(t1 time.Time, t2 time.Time) int64{
	diff := t2.Sub(t1)
	return int64(diff/1000/time.Millisecond)
}

func I64_S(value int64) string {
	return strconv.FormatInt(value, 10)
}