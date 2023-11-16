package sms

import (
  "strings"
  "net/http"
  "io/ioutil"
)

func SendOTP_MKT(api_key string, secret_key string, project_key string, phone string) (string, boolean) {
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
  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
      return err.Error(), false
  }

  return string(body), true
}

func ValidateOTP_MKT(api_key string, secret_key string, token string, otp string) (string, boolean) {
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
      return
  }
  req.Header.Add("Content-Type", "application/json")
  req.Header.Add("api_key", api_key)
  req.Header.Add("secret_key", secret_key)
  res, err := client.Do(req)
  if err != nil {
      return err.Error(), false
      return
  }
  defer res.Body.Close()
  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
      return err.Error(), false
      return
  }

  return string(body), true
}
