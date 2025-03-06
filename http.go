package ktgolib

import (
  "net/http"
  "io/ioutil"
  "time"
  "encoding/json"
  "bytes"
  "strings"
  "net/url"
)

func Http_Get(apiurl string) map[string]interface{}{
  var msgObj map[string]interface{} = nil

  res, err := http.Get(apiurl)
  if err != nil {
    panic(err.Error())
  }
	body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    panic(err.Error())
  }

  err = json.Unmarshal([]byte(body), &msgObj)
  if err != nil {
		panic(err.Error())
  }

  return msgObj
}

func Http_PostJson(apiurl string, data map[string]interface{}) map[string]interface{}{
  var msgObj map[string]interface{} = nil

  if body := http_PostJson(apiurl, map[string]string{}, data); body != nil {
    if err := json.Unmarshal(body, &msgObj); err != nil {
      Log("Http_PostJson", "", apiurl, "JSON result parsing failed : "+err.Error(), "HTTP_ERROR")
    }
  }

  return msgObj
}

func http_PostJson(apiurl string, headers map[string]string, data map[string]interface{}) []byte{
  runTime := time.Now()

  resultByte,err := json.Marshal(data)
  if err != nil {
    LogHidden("http_PostJson", "", apiurl, "JSON request parsing failed : "+err.Error(), "HTTP_ERROR")
    panic("error.APIRequestFailed")
  }

  req, err := http.NewRequest("POST", apiurl, bytes.NewBuffer(resultByte))
  if err != nil {
    LogHidden("http_PostJson", "", apiurl, "Buffer failed : "+err.Error(), "HTTP_ERROR")
    panic("error.APIRequestFailed")
  }

  req.Header.Set("Content-Type", "application/json")
  for key, value := range headers {
    req.Header.Set(key, value)
  }

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    LogHidden("http_PostJson", "", apiurl, "Request failed : "+err.Error(), "HTTP_ERROR")
    panic("error.APIRequestFailed")
  }
  defer resp.Body.Close()

  if resp.Status != "200 OK" {
    body, err := ioutil.ReadAll(resp.Body)
    if err == nil {
      LogHidden("http_PostJson", "", apiurl, "Status "+resp.Status+" : "+string(body), "HTTP_ERROR")
    }else{
      LogHidden("http_PostJson", "", apiurl, "Status "+resp.Status, "HTTP_ERROR")
      panic("error.APIRequestFailed")
    }
  }

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    LogHidden("http_PostJson", "", apiurl, "Body failed : "+err.Error(), "HTTP_ERROR")
    panic("error.APIRequestFailed")
  }

  LogHiddenWithDuration("http_PostJson", "", apiurl, string(body), I64_S(DateTimeValueDiffSec(runTime, time.Now())) , "HTTP")

  return body
}

func Http_PostForm(apiurl string, data url.Values) []map[string]interface{}{

  resp, err := http.PostForm(apiurl,data)

  if nil != err {
    LogHidden("Http_PostForm", "", apiurl, "Error : "+err.Error(), "HTTP_ERROR")
    return nil
  }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)

  if nil != err {
    LogHidden("Http_PostForm", "", apiurl, "body failed : "+err.Error(), "HTTP_ERROR")
    return nil
  }

  msgObj := make([]map[string]interface{},0)
  if err := json.Unmarshal(body, &msgObj); err != nil {
    LogHidden("Http_PostForm", "", apiurl, "parsing json failed : "+err.Error(), "HTTP_ERROR")
    return nil
  }

  return msgObj
}

func Http_PostFormWithHeader(apiurl string, headers map[string]string, data url.Values) string {

  req, err := http.NewRequest("POST", apiurl, strings.NewReader(data.Encode()) )
  if err != nil {
    LogHidden("Http_PostFormWithHeader", "", apiurl, "Creating request failed : "+err.Error(), "HTTP_ERROR")
    return ""
  }

  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  for key, value := range headers {
    req.Header.Set(key, value)
  }

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    LogHidden("Http_PostFormWithHeader", "", apiurl, "Request failed : "+err.Error(), "HTTP_ERROR")
    return ""
  }
  defer resp.Body.Close()

  if resp.Status != "200 OK" {
    LogHidden("Http_PostFormWithHeader", "", apiurl, "Status code : "+resp.Status, "HTTP_ERROR")
    return ""
  }

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    LogHidden("Http_PostFormWithHeader", "", apiurl, "Body failed : "+err.Error(), "HTTP_ERROR")
    return ""
  }

  return string(body)
}
