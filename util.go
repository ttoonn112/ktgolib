package ktgolib

// Shortcut Function
// ใช้งานเพื่อให้ Code สั้นลง

import (
	"strings"
	"fmt"
	"os"
	"time"
	"github.com/kr/pretty"
)

// Check ว่ามี key อยู่ m หรือไม่
func Has(m map[string]interface{}, key string) bool {
	_, exists := m[key]
	return exists
}

// ดึงข้อมูลจาก mapData เฉพาะ key ที่อยู่ใน keyArr
func GetMask(mapData map[string]interface{}, keyArr []string) map[string]interface{} {
  result := map[string]interface{}{}
	for _, k := range keyArr {
		if strings.HasSuffix(k, "_") {
			info := map[string]interface{}{}
			for km, v := range mapData {
				if strings.HasPrefix(km, k) {
					key := strings.TrimPrefix(km, k)
					info[key] = v
				}
			}
			result[strings.TrimSuffix(k, "_")] = info
		}else if strings.HasPrefix(k, "list_") {
			if arr, ok := mapData[k].([]map[string]interface{}); ok {
				result[strings.TrimPrefix(k, "list_")] = CompressArray(arr)
			}
		}else{
			if v, ok := mapData[k]; ok {
				result[k] = v
			}
		}
  }
	return result
}

// Optional chainning
func OC(condition bool, a string, b string) string {
  if condition {
    return a
  }else{
    return b
  }
}

func FirstXChar(str string, num int) string{
	if len(str)-num >= 0 {
		return str[:num]
	}else{
		return str
	}
}

func LastXChar(str string, num int) string{
	if len(str)-num >= 0 {
		return str[len(str)-num:]
	}
	return ""
}

func TryCatch(callback func(errStr string)) {
    if r := recover(); r != nil {
        errStr := ""
        if err, ok := r.(error); ok {
            errStr = err.Error()
        } else if errS, ok := r.(string); ok {
            errStr = errS
        }
        if callback != nil {
            callback(errStr)
        }
    }
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

func Println(object interface{}){
	pretty.Println(object)
}

func Print(object interface{}){
	pretty.Print(object)
}
