package ktgolib

import (
	"strings"
	"unicode"
	"encoding/json"
)

func EncloseText(src string) string{
  s1 := strings.Replace(strings.Replace(strings.Replace(strings.Replace(src,"\r\n","\n",-1),"\n","#NL#",-1),"'","&#39;",-1),`"`,"#DQ#",-1)
	s1 = strings.Map(func(r rune) rune {
      if unicode.IsPrint(r) {
          return r
      }
      return -1
  }, s1)
	return s1
}

func UncloseText(src string) string{
  str := strings.Replace( strings.Replace(strings.Replace(strings.Replace(src,"#NL#","\n",-1),"u0026","&",-1) ,"&#39;",`'`,-1) ,"#DQ#",`"`,-1)
  return str
}

func Extract(str string) map[string]interface{}{
  str = strings.Map(func(r rune) rune {
      if unicode.IsPrint(r) {
          return r
      }
      return -1
  }, str)
  var obj map[string]interface{}
  record := map[string]interface{}{}
  if err := json.Unmarshal([]byte(str), &obj); err == nil {
    for k, v := range obj {
      switch value := v.(type) {
       case string:
          record[k] = UncloseText(T(obj,k))
          break
       default:
          record[k] = value
          break
      }
    }
  }
  return record
}

func ExtractArray(str string) []map[string]interface{}{
  str = strings.Map(func(r rune) rune {
      if unicode.IsPrint(r) {
          return r
      }
      return -1
  }, str)
  var obj []map[string]interface{}
  if err := json.Unmarshal([]byte(str), &obj); err == nil {
    return obj
  }
  return nil
}

func Compress(payload map[string]interface{}) string{
  fields := map[string]interface{}{}
  for key, v := range payload {
    switch value := v.(type) {
     case string:
			 	if T(payload,key) != "" {
        	fields[key] = EncloseText(T(payload,key))
				}
        break
     default:
			 	if value != float64(0) && value != int64(0) && value != nil {
        	fields[key] = value
				}
        break
    }
  }
  fieldsByte, _ := json.Marshal(fields)
  str := string(fieldsByte)
  str = strings.Map(func(r rune) rune {
      if unicode.IsPrint(r) {
          return r
      }
      return -1
  }, str)
  return str
}

func CompressRaw(payload map[string]interface{}) string{
  fields := map[string]interface{}{}
  for key, v := range payload {
    switch value := v.(type) {
     case string:
			 	if T(payload,key) != "" {
        	fields[key] = EncloseText(T(payload,key))
				}
        break
     default:
			 	if value != nil {
        	fields[key] = value
				}
        break
    }
  }
  fieldsByte, _ := json.Marshal(fields)
  str := string(fieldsByte)
  str = strings.Map(func(r rune) rune {
      if unicode.IsPrint(r) {
          return r
      }
      return -1
  }, str)
  return str
}

func CompressArray(payload []map[string]interface{}) string{
  fieldValue := ""
  for _, obj := range payload {
    for k, v := range obj {
      switch v.(type) {
       case string:
         obj[k] = EncloseText(T(obj,k))
         break
       default:
         break
      }
    }
    jsonString, err := json.Marshal(obj);
    if err != nil {
      panic("error.JsonFailed")
    }
    fieldValue += ","+string(jsonString)
  }

  if fieldValue != "" {
    fieldValue = `[`+fieldValue[1:]+`]`
  }

  str := strings.Map(func(r rune) rune {
      if unicode.IsPrint(r) {
          return r
      }
      return -1
  }, fieldValue)

  return str
}
