package ktgolib

import (
	"fmt"
	"testing"
	"unicode/utf8"
)

// FirstXChar รุ่นก่อน 2026-10-20 (ตัดเป็นไบต์) — ใช้พิสูจน์ว่าข้อความ ASCII ได้ผลเหมือนเดิมทุกไบต์
func firstXCharBytes(str string, num int) string {
	if len(str)-num >= 0 {
		return str[:num]
	}
	return str
}

func firstXCharSafe(f func(string, int) string, s string, n int) (out string) {
	defer func() {
		if r := recover(); r != nil {
			out = fmt.Sprintf("PANIC(%v)", r)
		}
	}()
	return f(s, n)
}

// ข้อความ ASCII ทุกความยาว × ทุกจำนวนที่ขอ (รวมติดลบ ศูนย์ และเกินความยาว) = ผลเหมือนรุ่นเดิม รวมถึงการ panic
func TestFirstXChar_AsciiSameAsBefore(t *testing.T) {
	inputs := []string{"", "a", "abc", "2026-10-20 13:45:59", "Bearer abc.def", "GJ2610/00006", "line1\nline2\t\x00end"}
	ascii := make([]byte, 128)
	for i := range ascii {
		ascii[i] = byte(i)
	}
	inputs = append(inputs, string(ascii))
	for _, s := range inputs {
		for n := -2; n <= len(s)+3; n++ {
			want := firstXCharSafe(firstXCharBytes, s, n)
			got := firstXCharSafe(FirstXChar, s, n)
			if got != want {
				t.Errorf("FirstXChar(%q, %d) = %q · รุ่นเดิม %q", s, n, got, want)
			}
		}
	}
}

// ภาษาไทยตัดตามตัวอักษร ได้ข้อความที่ถูกต้องเสมอ (รุ่นเดิมตัดกลางตัว ขึ้น "��")
func TestFirstXChar_Thai(t *testing.T) {
	s := "คืนเงินสำรองจ่ายค่าตรายางบริษัท"
	for n := 0; n <= utf8.RuneCountInString(s)+2; n++ {
		got := FirstXChar(s, n)
		want := n
		if c := utf8.RuneCountInString(s); want > c {
			want = c
		}
		if !utf8.ValidString(got) || utf8.RuneCountInString(got) != want || s[:len(got)] != got {
			t.Errorf("FirstXChar(%q, %d) = %q", s, n, got)
		}
	}
	if got := FirstXChar("GJ ค่าตรายาง 690.00", 6); got != "GJ ค่า" {
		t.Errorf("ข้อความผสม = %q", got)
	}
	// ตัวเดิมกับข้อความนี้ = ไม่ใช่ UTF-8 (จับได้ว่าเทสนี้ล้มบนรุ่นเดิม)
	if utf8.ValidString(firstXCharBytes(s, 80)) {
		t.Fatalf("ข้อมูลทดสอบไม่พิสูจน์อะไร: รุ่นเดิมก็ได้ข้อความที่ถูกต้อง")
	}
}

// ไบต์ที่ไม่ใช่ UTF-8 คงไว้ตามต้นฉบับ นับไบต์ละ 1 ตัว (ไม่แปลงเป็น U+FFFD)
func TestFirstXChar_InvalidBytesKept(t *testing.T) {
	if got := FirstXChar("\xff\xfeab", 3); got != "\xff\xfea" {
		t.Errorf("got %q", got)
	}
}
