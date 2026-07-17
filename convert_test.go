package ktgolib

import (
	"strings"
	"testing"
)

// T-002: MapToString ต้อง redact password/setup_password (ใช้เฉพาะสร้าง log string)
func TestMapToStringRedactsSecrets(t *testing.T) {
	out := MapToString(map[string]interface{}{
		"username":       "somchai",
		"password":       "S3cr3t!pw",
		"setup_password": "newpass99",
		"role":           "admin",
	})
	// รหัสจริงต้องไม่โผล่
	for _, secret := range []string{"S3cr3t!pw", "newpass99"} {
		if strings.Contains(out, secret) {
			t.Errorf("MapToString ยัง leak ความลับ %q: %s", secret, out)
		}
	}
	// ต้อง redact เป็น ***
	if !strings.Contains(out, `"password":"***"`) {
		t.Errorf("password ไม่ถูก redact เป็น ***: %s", out)
	}
	if !strings.Contains(out, `"setup_password":"***"`) {
		t.Errorf("setup_password ไม่ถูก redact เป็น ***: %s", out)
	}
	// field ปกติต้องอยู่ครบ ไม่ถูกแตะ
	if !strings.Contains(out, `"username":"somchai"`) || !strings.Contains(out, `"role":"admin"`) {
		t.Errorf("field ปกติหาย/เพี้ยน: %s", out)
	}
}

// ค่าว่างไม่ควรกลายเป็น *** (ไม่มีความลับให้ซ่อน)
func TestMapToStringEmptyPasswordNotRedacted(t *testing.T) {
	out := MapToString(map[string]interface{}{"password": "", "username": "u"})
	if strings.Contains(out, "***") {
		t.Errorf("password ว่างไม่ควรเป็น ***: %s", out)
	}
}

// key ปกติที่ไม่ sensitive ต้องผ่านเหมือนเดิม (ไม่ redact เกิน)
func TestMapToStringNormalKeysUnchanged(t *testing.T) {
	out := MapToString(map[string]interface{}{"code": "C001", "name": "บริษัท ทดสอบ", "qty": 5})
	for _, want := range []string{`"code":"C001"`, `"qty":5`} {
		if !strings.Contains(out, want) {
			t.Errorf("field ปกติเพี้ยน (ต้องมี %s): %s", want, out)
		}
	}
	if strings.Contains(out, "***") {
		t.Errorf("ไม่ควรมี *** กับ field ปกติ: %s", out)
	}
}
