package ktgolib

// 🔒 ล็อกรายชื่อฟังก์ชันสาธารณะ — ห้ามลบ ห้ามเปลี่ยนชื่อ ห้ามเปลี่ยนพารามิเตอร์
//
// ทำไมต้องมีเทสนี้
// การลบหรือเปลี่ยน signature ของฟังก์ชันที่ consumer เรียกอยู่ จะทำให้ repo นั้น
// build ไม่ผ่านทันทีที่เลื่อนเวอร์ชัน และคนแก้ไลบรารีจะไม่รู้เลย เพราะในนี้ยัง build ผ่าน
// (คู่กับ behavior_lock_test.go ที่ล็อก "ผลลัพธ์" — ตัวนี้ล็อก "หน้าตา")
//
// ✅ เพิ่มฟังก์ชันใหม่ได้ (นั่นคือวิธีที่ถูกต้องเวลาต้องการพฤติกรรมใหม่)
//    เพิ่มแล้วเทสจะล้มพร้อมบอกให้เอาชื่อใหม่ไปต่อท้ายไฟล์รายชื่อ — ตั้งใจให้เป็นขั้นตอนที่ต้องรู้ตัว
// ⛔ ลบ/เปลี่ยนชื่อ/เปลี่ยนพารามิเตอร์ = ของเดิมพัง ต้องไม่ทำ

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

const apiGoldenPath = "testdata/api.golden"

func TestPublicApiSurfaceLocked(t *testing.T) {
	got := publicApi(t)
	if len(got) < 50 {
		t.Fatalf("อ่านฟังก์ชันได้แค่ %d ตัว — น้อยผิดปกติ ตัวอ่านน่าจะพัง ไม่ใช่โค้ดหาย", len(got))
	}

	if os.Getenv("KTGOLIB_WRITE_GOLDEN") == "1" {
		writeLines(t, apiGoldenPath,
			"# รายชื่อฟังก์ชันสาธารณะที่ consumer เรียกได้ — ห้ามลบ/เปลี่ยนชื่อ/เปลี่ยนพารามิเตอร์",
			"# เพิ่มฟังก์ชันใหม่ = เพิ่มบรรทัดที่นี่ (อนุญาต) · สร้างใหม่: KTGOLIB_WRITE_GOLDEN=1 go test",
			got)
		t.Fatalf("สร้าง %s ใหม่แล้ว %d รายการ — ตรวจ diff ก่อน commit", apiGoldenPath, len(got))
	}

	want := readLines(t, apiGoldenPath)
	wantSet := map[string]bool{}
	for _, w := range want {
		wantSet[w] = true
	}
	gotSet := map[string]bool{}
	for _, g := range got {
		gotSet[g] = true
	}

	for _, w := range want {
		if !gotSet[w] {
			t.Errorf("⛔ ฟังก์ชันหายไปหรือเปลี่ยนหน้าตา: %s\n"+
				"    consumer ที่เรียกอยู่จะ build ไม่ผ่านทันทีที่เลื่อนเวอร์ชัน\n"+
				"    ต้องการรูปแบบใหม่ให้เพิ่มฟังก์ชันใหม่ แล้วปล่อยของเดิมไว้", w)
		}
	}
	for _, g := range got {
		if !wantSet[g] {
			t.Errorf("ฟังก์ชันใหม่ที่ยังไม่อยู่ในรายชื่อ: %s\n"+
				"    เพิ่มของใหม่ทำได้ — ให้สร้างรายชื่อใหม่ตาม AGENTS.md แล้วเพิ่มเคสใน behavior_lock_cases_test.go ด้วย", g)
		}
	}
}

// publicApi — ชื่อ + พารามิเตอร์ + ชนิดที่คืน ของทุกฟังก์ชันสาธารณะในแพ็กเกจนี้
func publicApi(t *testing.T) []string {
	t.Helper()
	fset := token.NewFileSet()
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, path := range files {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			t.Fatalf("อ่าน %s ไม่ได้: %v", path, err)
		}
		for _, d := range f.Decls {
			fd, ok := d.(*ast.FuncDecl)
			if !ok || fd.Recv != nil || !fd.Name.IsExported() {
				continue
			}
			out = append(out, fd.Name.Name+typeString(fset, fd.Type))
		}
	}
	sort.Strings(out)
	return out
}

func typeString(fset *token.FileSet, ft *ast.FuncType) string {
	var sb strings.Builder
	ft2 := *ft
	ft2.Func = token.NoPos
	if err := printer.Fprint(&sb, fset, &ft2); err != nil {
		return "<?>"
	}
	// printer คืนมาเป็น "func(a int) string" — ตัดคำว่า func ออกให้เหลือแค่พารามิเตอร์กับชนิดที่คืน
	return strings.TrimPrefix(strings.Join(strings.Fields(sb.String()), " "), "func")
}

func readLines(t *testing.T, path string) []string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("อ่าน %s ไม่ได้: %v", path, err)
	}
	var out []string
	for _, l := range strings.Split(string(b), "\n") {
		if l != "" && !strings.HasPrefix(l, "#") {
			out = append(out, l)
		}
	}
	if len(out) == 0 {
		t.Fatalf("%s ว่าง — เทสนี้จะกลายเป็นของประดับ", path)
	}
	return out
}

func writeLines(t *testing.T, path, header1, header2 string, lines []string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	body := fmt.Sprintf("%s\n%s\n%s\n", header1, header2, strings.Join(lines, "\n"))
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}
