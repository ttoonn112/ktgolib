package ktgolib

// Shortcut Function
// ใช้งานเพื่อให้ Code สั้นลง

import (
	"strings"
	"fmt"
	"os"
	"time"
	"reflect"
	"runtime"
	"github.com/kr/pretty"
)

// Check ว่ามี key อยู่ m หรือไม่
func Has(m map[string]interface{}, key string) bool {
	_, exists := m[key]
	return exists
}

// ดึงข้อมูลจาก mapData เฉพาะ key ที่อยู่ใน keyArr
// ถ้า key ลงท้ายด้วย _ เป็นข้อมูลประเภท Json
// ถ้า key ลงท้ายด้วย ! เป็นข้อมูลประเภท Json Array
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
		}else if strings.HasSuffix(k, "!") {
			thekey := strings.TrimSuffix(k, "!")
			if arrI, ok := mapData[thekey].([]interface{}); ok {
				items := []map[string]interface{}{}
				for _, itemI := range arrI {
					if item, ok := itemI.(map[string]interface{}); ok {
						items = append(items, item)
					}
				}
				result[thekey] = items
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

func CopyMap(item map[string]interface{}) map[string]interface{} {
	nitem := map[string]interface{}{}
	for k,v := range item {
		nitem[k] = v
	}
	return nitem
}

func Includes(arr1_str string, arr2_str string, sep string) bool{
	arr1 := strings.Split(arr1_str, sep)
	arr2 := strings.Split(arr2_str, sep)
	for _, a1 := range arr1 {
    for _, a2 := range arr2 {
      if a1 == a2 {
        return true
      }
    }
  }
	return false
}

func HasIncludes(arr1 []string, arr2 []string) bool{
	for _, a1 := range arr1 {
    for _, a2 := range arr2 {
      if a1 == a2 {
        return true
      }
    }
  }
	return false
}

// ไม่สามารถใช้ร่วมกับการ return ค่า กรณี TryCatch จับการ Error ได้
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

// ###################### Begin Attempt ########################
// ฟังก์ชัน Attempt รับฟังก์ชันที่ต้องการรันซ้ำ, จำนวนครั้งที่จำกัด และพารามิเตอร์แบบ variadic
func Attempt(operation interface{}, maxAttempts int, delaySec int, params ...interface{}) string {
	errMsgReturn := ""

	for i := 0; i < maxAttempts; i++ {
		runTime := time.Now()

		errMsg := callFunction(i+1, operation, params...)

		if errMsg != "" {
			errMsgReturn = errMsg
			LogHiddenWithDuration(getFunctionName(operation), "", I_S(i+1), errMsg, I64_S(DateTimeValueDiff(runTime, time.Now()))+"s", "Attempt")
		}else{
			return ""
		}

		time.Sleep(time.Duration(delaySec)*time.Second) // ถ้าไม่สำเร็จ รอช่วงเวลาที่กำหนดก่อนพยายามใหม่
	}

	return errMsgReturn
}

// callFunction ใช้ reflect เพื่อเรียกใช้ฟังก์ชันพร้อมกับพารามิเตอร์ที่กำหนด (ถูกใช้งานที่ Attempt)
func callFunction(numAttempt int, fn interface{}, params ...interface{}) (errMsg string) {

	defer func(){
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
					errMsg = err.Error()
			} else if errS, ok := r.(string); ok {
					errMsg = errS
			}
		}
	}()

	fnValue := reflect.ValueOf(fn)
	if fnValue.Kind() != reflect.Func {
		errMsg = "Operation is not a function"
	}

	args := make([]reflect.Value, len(params))
	for i, param := range params {
		args[i] = reflect.ValueOf(param)
	}

	fnValue.Call(args)		// Can return []reflect.Value from fnValue.Call to get the result(s) back

	return
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
// ###################### End Attempt ########################

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

func Println(object interface{}){
	pretty.Println(object)
}

func Print(object interface{}){
	pretty.Print(object)
}
