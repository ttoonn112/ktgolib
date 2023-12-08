package ktgolib

// Shortcut Function
// ใช้งานเพื่อให้ Code สั้นลง

import (
	"strings"
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
