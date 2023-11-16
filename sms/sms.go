package sms

import (
  "strings"
  "net/http"
  "io/ioutil"
  "encoding/json"
)

func SendOTP_MKT(api_key string, secret_key string, project_key string, phone string) (result string, ok bool) {

  defer func() {
		if r := recover(); r != nil {
      ok = false
   		if err,ok := r.(error); ok {
   			result = err.Error()
   		}else if errS,ok := r.(string); ok {
   			result = errS
   		}
		}
	}()

  url := "https://portal-otp.smsmkt.com/api/otp-send"
  method := "POST"
  payload := strings.NewReader(`{
        "project_key": "`+project_key+`",
        "phone": "`+phone+`"
    }`)
  client := &http.Client {    }
  req, err := http.NewRequest(method, url, payload)
  if err != nil {
    return err.Error(), false
  }
  req.Header.Add("Content-Type", "application/json")
  req.Header.Add("api_key", api_key)
  req.Header.Add("secret_key", secret_key)
  res, err := client.Do(req)
  if err != nil {
    return err.Error(), false
  }
  defer res.Body.Close()

  if res.Status != "200 OK" {
    return res.Status, false
  }

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
      return err.Error(), false
  }

  var obj map[string]interface{}
  if err := json.Unmarshal([]byte(string(body)), &obj); err != nil {
    return err.Error(), false
  }

  if result,ok := obj["result"].(map[string]interface{}); ok && result != nil {
    return result["token"].(string), true
  }else{
    return obj["code"].(string)+" - "+obj["detail"].(string), false
  }

}

func ValidateOTP_MKT(api_key string, secret_key string, token string, otp string) (result string, ok bool) {

  defer func() {
		if r := recover(); r != nil {
      ok = false
   		if err,ok := r.(error); ok {
   			result = err.Error()
   		}else if errS,ok := r.(string); ok {
   			result = errS
   		}
		}
	}()

  url := "https://portal-otp.smsmkt.com/api/otp-validate"
  method := "POST"
  payload := strings.NewReader(`{
      "token": "`+token+`",
      "otp_code": "`+otp+`"
  }`)
  client := &http.Client {    }
  req, err := http.NewRequest(method, url, payload)
  if err != nil {
      return err.Error(), false
  }
  req.Header.Add("Content-Type", "application/json")
  req.Header.Add("api_key", api_key)
  req.Header.Add("secret_key", secret_key)
  res, err := client.Do(req)
  if err != nil {
      return err.Error(), false
  }
  defer res.Body.Close()

  if res.Status != "200 OK" {
    return res.Status, false
  }

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
      return err.Error(), false
  }

  var obj map[string]interface{}
  if err := json.Unmarshal([]byte(string(body)), &obj); err != nil {
    return err.Error(), false
  }

  if result,ok := obj["result"].(map[string]interface{}); ok && result != nil {
    if result["status"].(bool) {
      return "", true
    }else{
      return obj["code"].(string)+" - "+obj["detail"].(string)+" => invalid code", false
    }
  }else{
    return obj["code"].(string)+" - "+obj["detail"].(string), false
  }

}
