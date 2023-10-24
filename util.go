package ktgolib

// Shortcut Function
// ใช้งานเพื่อให้ Code สั้นลง

import (
	
)

// Check ว่ามี key อยู่ m หรือไม่
func Has(m map[string]interface{}, key string) bool {
	_, exists := m[key]
	return exists
}

// ดึงข้อมูลจาก mapData เฉพาะ key ที่อยู่ใน keyArr
func GetMask(mapData map[string]interface{}, keyArr []string) map[string]interface{} {
  result := map[string]interface{}{}
	for k1, v := range mapData {
    for _, k2 := range keyArr {
      if k1 == k2 {
        result[k1] = v
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
