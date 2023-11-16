package sms

import (
  "fmt"
  "os"
  "time"
  "github.com/kr/pretty"
  "strings"
  "net/http"
  "io/ioutil"
)

func SendOTP_MKT(api_key string, secret_key string, phone string) (string, boolean) {
  url := "https://portal-otp.smsmkt.com/api/otp-send"
  method := "POST"
  payload := strings.NewReader(`{
  }`)
  client := &http.Client {    }
  req, err := http.NewRequest(method, url, payload)
  if err != nil {
    LogHidden("SendOTP_MKT", "", secret_key+":"+phone, "Result : "+err.Error(), "OTP")
    LogHidden("SendOTP_MKT", "", secret_key+":"+phone, "Result : "+err.Error(), "OTPFailed")
    return err.Error(), false
  }
  req.Header.Add("Content-Type", "application/json")
  req.Header.Add("api_key", api_key)
  req.Header.Add("secret_key", secret_key)
  res, err := client.Do(req)
  if err != nil {
    LogHidden("SendOTP_MKT", "", secret_key+":"+phone, "Result : "+err.Error(), "OTP")
    LogHidden("SendOTP_MKT", "", secret_key+":"+phone, "Result : "+err.Error(), "OTPFailed")
    return err.Error(), false
  }
  defer res.Body.Close()
  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
      LogHidden("SendOTP_MKT", "", secret_key+":"+phone, "Result : "+err.Error(), "OTP")
      LogHidden("SendOTP_MKT", "", secret_key+":"+phone, "Result : "+err.Error(), "OTPFailed")
      return err.Error(), false
  }

  LogHidden("SendOTP_MKT", "", project+":"+phone, "Result : "+string(body), "OTP")

  return string(body), true
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
		fmt.Print("Log|o="+operation+"|u="+username+"|k="+key+"|d="+duration+"|m=["+msg+"]\r\n")
	}
  fmt.Fprintf(file, "t="+t.Format("15:04:05.000")+"|o="+operation+"|u="+username+"|k="+key+"|d="+duration+"|m=["+msg+"]\r\n")
}

func LogHidden(operation string, username string, key string, msg string, logfilename string){
	writeLog(operation, username, key, msg, "", logfilename, false)
}
