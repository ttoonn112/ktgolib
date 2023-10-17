package log

import (
  "fmt"
  "os"
  "time"
  "github.com/kr/pretty"
)

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
		fmt.Print("Log|o="+operation+"|u="+username+"|k="+key+"|d="+duration+"|m=["+msg+"]\r\n")
	}
  fmt.Fprintf(file, "t="+t.Format("15:04:05.000")+"|o="+operation+"|u="+username+"|k="+key+"|d="+duration+"|m=["+msg+"]\r\n")
}

func LogHidden(operation string, username string, key string, msg string, logfilename string){
	writeLog(operation, username, key, msg, "", logfilename, false)
}

func Log(operation string, username string, key string, msg string, logfilename string){
	writeLog(operation, username, key, msg, "", logfilename, true)
}

func LogHiddenWithDuration(operation string, username string, key string, msg string, duration string, logfilename string){
	writeLog(operation, username, key, msg, duration, logfilename, false)
}

func LogWithDuration(operation string, username string, key string, msg string, duration string, logfilename string){
	writeLog(operation, username, key, msg, duration, logfilename, true)
}

func Println(object interface{}){
	pretty.Println(object)
}

func Print(object interface{}){
	pretty.Print(object)
}
