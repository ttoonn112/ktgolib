package ktgolib

import (
	"strings"
	"unicode"
	"encoding/json"
)

// เข้ารหัสอักขระที่ทำ SQL string literal พัง ให้เป็น sentinel — ถอดกลับด้วย UncloseText
// backslash ต้องมาก่อนตัวอื่น (mirror กับ UncloseText ที่ถอด #BS# ท้ายสุด)
func EncloseText(src string) string{
  s1 := strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(src,`\`,"#BS#",-1),"\r\n","\n",-1),"\n","#NL#",-1),"'","&#39;",-1),`"`,"#DQ#",-1)
	s1 = strings.Map(func(r rune) rune {
      if unicode.IsPrint(r) {
          return r
      }
      return -1
  }, s1)
	return s1
}

// ถอดรหัสจาก EncloseText — #BS# ต้องท้ายสุด (mirror กับ EncloseText ที่เข้ารหัส backslash ก่อน)
// u0026 = ร่องรอยของ data เก่าที่ & ถูก MySQL กิน backslash ไป ยังต้องรองรับต่อไป
func UncloseText(src string) string{
  str := strings.Replace( strings.Replace( strings.Replace(strings.Replace(strings.Replace(src,"#NL#","\n",-1),"u0026","&",-1) ,"&#39;",`'`,-1) ,"#DQ#",`"`,-1) ,"#BS#",`\`,-1)
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

// ถอด sentinel ทุก string field เหมือน Extract — ไม่งั้น '/"/&/newline/backslash ใน array element
// จะโผล่เป็น &#39;/#DQ#/u0026/#NL#/#BS# ให้ผู้ใช้เห็น (encode ผ่าน CompressArray แต่เดิมไม่มีตัวถอด)
func ExtractArray(str string) []map[string]interface{}{
  str = strings.Map(func(r rune) rune {
      if unicode.IsPrint(r) {
          return r
      }
      return -1
  }, str)
  var obj []map[string]interface{}
  if err := json.Unmarshal([]byte(str), &obj); err == nil {
    for _, row := range obj {
      for k, v := range row {
        if _, ok := v.(string); ok {
          row[k] = UncloseText(T(row, k))
        }
      }
    }
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
