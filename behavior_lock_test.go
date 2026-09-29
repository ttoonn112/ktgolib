package ktgolib

// 🔒 ล็อกพฤติกรรม — ห้ามเปลี่ยนผลลัพธ์ของฟังก์ชันที่มีอยู่แล้ว
//
// ทำไมต้องมีเทสนี้
// ไลบรารีนี้ถูก 16 repo เรียกใช้ และแต่ละ repo ปักเวอร์ชันเอง การแก้พฤติกรรมของ
// ฟังก์ชันเดิมโดยไม่เปลี่ยนชื่อ = ระเบิดเวลา: repo ที่ยังไม่เลื่อนเวอร์ชันไม่รู้สึกอะไร
// แต่วันที่เลื่อน โค้ดจะ build ผ่านเหมือนเดิมทุกบรรทัด แล้วพฤติกรรมเปลี่ยนเงียบ ๆ
//
// เกิดจริงมาแล้ว 2 ครั้งจากการแก้ครั้งเดียว (2026-07-16 AddSql* เริ่ม escape ให้เอง):
//   · magbiz service/master.go:19 — ครอบ SqlStr ทับ ค้นค่าที่มี ' ไม่เจอ
//   · besttransport + smartfarmservice — หุ้ม AddSql* ไว้เองกัน injection
//     วันที่เลื่อนเวอร์ชันจะ escape สองชั้นทุกจุดเรียก (164 / 5 จุด)
//
// กติกา
//   ✅ ปรับปรุงได้: เพิ่มฟังก์ชันใหม่ · เร่งความเร็ว · จัดระเบียบโค้ด · เพิ่มคอมเมนต์
//      ตราบใดที่ผลลัพธ์ของ input เดิมยังเหมือนเดิมทุกตัวอักษร
//   ⛔ ห้าม: เปลี่ยนผลลัพธ์ของฟังก์ชันที่มีอยู่ แม้จะ "ถูกต้องกว่าเดิม"
//      ต้องการพฤติกรรมใหม่ → ตั้งชื่อใหม่ ปล่อยของเดิมไว้
//
// เทสนี้ล้มเมื่อผลลัพธ์เปลี่ยน — ห้ามสร้าง golden ใหม่เพื่อให้ผ่าน
// (วิธีสร้างใหม่และกรณีที่อนุญาต: ดู AGENTS.md §เปลี่ยนพฤติกรรม)

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
	"time"
)

const goldenPath = "testdata/behavior.golden"

// เวลาคงที่ — ฟังก์ชันวันเวลาอ่านโซนเวลาของเครื่อง ถ้าไม่ตรึงไว้ เทสจะได้ผลต่างกันคนละเครื่อง
func init() { time.Local = time.FixedZone("ICT", 7*3600) }

// render — แปลงค่าที่ฟังก์ชันคืนมาให้เป็นข้อความเทียบได้ (map เรียงคีย์ให้เองโดย json)
func render(v interface{}) string {
	switch x := v.(type) {
	case string:
		return fmt.Sprintf("%q", x)
	default:
		b, err := json.Marshal(x)
		if err != nil {
			return fmt.Sprintf("<marshal error: %v>", err)
		}
		return string(b)
	}
}

// safe — เก็บพฤติกรรมตอน panic ไว้ด้วย การเปลี่ยนจาก panic เป็นไม่ panic ก็คือเปลี่ยนพฤติกรรม
func safe(f func() interface{}) (out string) {
	defer func() {
		if r := recover(); r != nil {
			out = fmt.Sprintf("PANIC(%v)", r)
		}
	}()
	return render(f())
}

type lockCase struct {
	key string
	run func() interface{}
}

func TestBehaviorLocked(t *testing.T) {
	cases := lockedCases()

	seen := map[string]bool{}
	got := map[string]string{}
	for _, c := range cases {
		if seen[c.key] {
			t.Fatalf("ชื่อเคสซ้ำ: %s", c.key)
		}
		seen[c.key] = true
		got[c.key] = safe(c.run)
	}

	if os.Getenv("KTGOLIB_WRITE_GOLDEN") == "1" {
		writeGolden(t, got)
		t.Fatalf("สร้าง %s ใหม่แล้ว %d เคส — ตรวจ diff ทีละบรรทัดก่อน commit (ดู AGENTS.md)", goldenPath, len(got))
	}

	want := readGolden(t)

	for key, w := range want {
		g, ok := got[key]
		if !ok {
			t.Errorf("เคส %s หายไปจากเทส — ถ้าลบฟังก์ชัน/เคสทิ้ง แปลว่า consumer ที่เรียกอยู่จะพัง", key)
			continue
		}
		if g != w {
			t.Errorf("⛔ พฤติกรรมเปลี่ยน: %s\n    เดิม  %s\n    ใหม่  %s\n"+
				"    ไลบรารีนี้ใช้ร่วม 16 repo — ห้ามเปลี่ยนผลของฟังก์ชันเดิม\n"+
				"    ต้องการพฤติกรรมใหม่ให้ตั้งชื่อฟังก์ชันใหม่ แล้วปล่อยของเดิมไว้", key, w, g)
		}
	}
	for key := range got {
		if _, ok := want[key]; !ok {
			t.Errorf("เคสใหม่ %s ยังไม่มีใน golden — เพิ่มเคสให้ฟังก์ชันใหม่ได้ ให้สร้าง golden ใหม่ตาม AGENTS.md", key)
		}
	}
}

func readGolden(t *testing.T) map[string]string {
	t.Helper()
	b, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("อ่าน %s ไม่ได้: %v", goldenPath, err)
	}
	out := map[string]string{}
	for _, line := range strings.Split(string(b), "\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			t.Fatalf("บรรทัดในไฟล์ golden ผิดรูปแบบ: %q", line)
		}
		out[parts[0]] = parts[1]
	}
	if len(out) == 0 {
		t.Fatal("ไฟล์ golden ว่าง — เทสนี้จะกลายเป็นของประดับ")
	}
	return out
}

func writeGolden(t *testing.T, got map[string]string) {
	t.Helper()
	if err := os.MkdirAll("testdata", 0o755); err != nil {
		t.Fatal(err)
	}
	keys := make([]string, 0, len(got))
	for k := range got {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sb strings.Builder
	sb.WriteString("# ผลลัพธ์ที่ถูกล็อกไว้ของทุกฟังก์ชันที่ consumer เรียกใช้ — ห้ามแก้ด้วยมือ\n")
	sb.WriteString("# เปลี่ยนบรรทัดใดบรรทัดหนึ่ง = เปลี่ยนพฤติกรรมของไลบรารีที่ 16 repo ใช้ร่วมกัน\n")
	sb.WriteString("# สร้างใหม่: KTGOLIB_WRITE_GOLDEN=1 go test -run TestBehaviorLocked\n")
	for _, k := range keys {
		sb.WriteString(k + "\t" + got[k] + "\n")
	}
	if err := os.WriteFile(goldenPath, []byte(sb.String()), 0o644); err != nil {
		t.Fatal(err)
	}
}
