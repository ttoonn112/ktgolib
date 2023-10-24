package ktgolib

import (
  "time"
  "strings"
  "strconv"
  "regexp"
	"math/rand"
)

func MapVStr(m map[string]interface{}, k string, default_value interface{}) string {
	if v, present := m[k]; present {
	 switch value := v.(type) {
	 case float64:
			return F64_S(value,-1)
	 case int:
		 return I_S(value)
	 case int64:
			return I64_S(value)
	 case bool:
			return OC(value,"true","false")
	 case string:
			return value
	 default:
			 return default_value.(string)
	 }
 }
 return default_value.(string)
}

func S_I(str string) int{
  i,err := strconv.Atoi(str)
  if err != nil {
    return 0
  }
  return int(i)
}

func S_I64(str string) int64{
  i,err := strconv.Atoi(str)
  if err != nil {
    return 0
  }
  return int64(i)
}

func S_I32(str string) int32{
  i,err := strconv.Atoi(str)
  if err != nil {
    return 0
  }
  return int32(i)
}

func S_F64(str string) float64{
  i, err := strconv.ParseFloat(str, 64)
  if err != nil {
      return 0
  }
  return float64(i)
}

func I_S(value int) string {
	var i64 int64
	i64 = int64(value)
	return strconv.FormatInt(i64, 10)
}

func I64_S(value int64) string {
	return strconv.FormatInt(value, 10)
}

func I32_S(value int32) string {
	return strconv.FormatInt(int64(value), 10)
}

func F64_S(value float64, decimal int) string {
	return strconv.FormatFloat(value, 'f', decimal, 64)
}

func F64_S_AUTO(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func T(m map[string]interface{}, k string) string {
	return MapVStr(m, k, "")
}

func SF64(m map[string]interface{}, k string) float64 {
	return S_F64(MapVStr(m, k, ""))
}

func SI64(m map[string]interface{}, k string) int64 {
	return S_I64(MapVStr(m, k, ""))
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789$@!*.+_-"
func GenerateRandomString(n int) string {		//RandASCIIBytes
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[seededRand.Intn(len(letterBytes))]
	}
	return string(b)
}

const letterNumBytes = "0123456789"
func GenerateRandomNumberString(n int) string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = letterNumBytes[seededRand.Intn(len(letterNumBytes))]
	}
	return string(b)
}

func DoLetterOnly(str_ string) string{
	reg, err := regexp.Compile("[^a-zA-Z0-9]$@!*.+_-")
  if err != nil {return str_}
	return reg.ReplaceAllString(str_, "")
}

func IsLetterOnly(s string) bool {
   for _, char := range s {
      if !strings.Contains(letterBytes, string(char)) {
         return false
      }
   }
   return true
}
