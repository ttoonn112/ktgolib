package ktgolib

import "testing"

// เคสที่ต้อง round-trip ได้ครบ — backslash คือชุดที่ T-001 พบว่าพัง
var roundTripCases = []struct {
	name string
	val  string
}{
	{"apostrophe", "O'Brien"},
	{"ampersand", "A&B Co."},
	{"both", "O'Brien & Co."},
	{"dquote", `say "hi"`},
	{"thai", "บริษัท ทดสอบ จำกัด"},
	{"newline", "line1\nline2"},
	// T-001: backslash — โหมดพัง (ตัวถัดไปไม่ใช่ JSON escape)
	{"backslash_invalid_esc", `test\ok`},
	// T-001: backslash — โหมดเพี้ยนเงียบ (ตัวถัดไปเป็น JSON escape ที่ใช้ได้)
	{"backslash_valid_esc_n", `C:\new`},
	{"backslash_valid_esc_t", `C:\table`},
	{"backslash_valid_esc_r", `C:\report`},
	{"backslash_double", `a\\b`},
	{"backslash_only", `\`},
	{"backslash_end", `path\`},
	{"windows_path", `C:\new folder\temp.txt`},
	{"mixed", `O'Brien & C:\new "x"`},
}

// EncloseText -> UncloseText ต้องได้ค่าเดิมกลับมาเสมอ
func TestEncloseUncloseRoundTrip(t *testing.T) {
	for _, c := range roundTripCases {
		t.Run(c.name, func(t *testing.T) {
			got := UncloseText(EncloseText(c.val))
			if got != c.val {
				t.Errorf("round-trip เพี้ยน\n input   = %q\n enclose = %q\n got     = %q", c.val, EncloseText(c.val), got)
			}
		})
	}
}

// CRLF ถูก normalize เป็น LF ตั้งแต่ต้น (พฤติกรรมที่ตั้งใจ ไม่ใช่ round-trip)
func TestEncloseTextNormalizesCRLF(t *testing.T) {
	if got := UncloseText(EncloseText("line1\r\nline2")); got != "line1\nline2" {
		t.Errorf("CRLF ต้องถูก normalize เป็น LF: got %q", got)
	}
}

// EncloseText ต้องไม่ปล่อยอักขระที่ทำ SQL string literal / JSON พัง ออกไป
// (backslash คือ escape char ของ MySQL — หลุดไปแล้วโดนกินตอน parse)
func TestEncloseTextRemovesSqlBreakingChars(t *testing.T) {
	for _, c := range roundTripCases {
		enc := EncloseText(c.val)
		for _, bad := range []struct {
			ch   rune
			desc string
		}{
			{'\\', "backslash (MySQL escape char)"},
			{'\'', "single quote (ปิด string literal)"},
			{'"', "double quote (ทำ JSON พัง)"},
			{'\n', "newline"},
		} {
			for _, r := range enc {
				if r == bad.ch {
					t.Errorf("%s: EncloseText(%q) = %q ยังมี %s หลุดออกมา", c.name, c.val, enc, bad.desc)
					break
				}
			}
		}
	}
}

// Compress -> Extract ต้อง round-trip ได้ (เส้นทางจริงที่ app ใช้)
func TestCompressExtractRoundTrip(t *testing.T) {
	for _, c := range roundTripCases {
		t.Run(c.name, func(t *testing.T) {
			got := T(Extract(Compress(map[string]interface{}{"x": c.val})), "x")
			if got != c.val {
				t.Errorf("Compress/Extract เพี้ยน\n input = %q\n got   = %q", c.val, got)
			}
		})
	}
}

// backward compat: data เก่าที่เขียนก่อน T-001 fix ต้องยังอ่านได้เหมือนเดิม
// (& ถูก MySQL กิน backslash ของ & ไป เหลือ u0026 — UncloseText ต้องยังถอดให้)
func TestUncloseTextOldDataFormat(t *testing.T) {
	cases := []struct{ stored, want string }{
		{"Ou0026#39;Brien", "O'Brien"},           // ' ของ data เก่า
		{"Au0026B Co.", "A&B Co."},               // & ของ data เก่า
		{"Ou0026#39;Brien u0026 Co.", "O'Brien & Co."},
		{"say #DQ#hi#DQ#", `say "hi"`},
		{"line1#NL#line2", "line1\nline2"},
		{"บริษัท ทดสอบ จำกัด", "บริษัท ทดสอบ จำกัด"},
	}
	for _, c := range cases {
		if got := UncloseText(c.stored); got != c.want {
			t.Errorf("อ่าน data เก่าเพี้ยน: UncloseText(%q) = %q, ต้องได้ %q", c.stored, got, c.want)
		}
	}
}

// SqlStr ต้อง escape อักขระที่ทำ SQL พังครบ
func TestSqlStr(t *testing.T) {
	cases := []struct{ in, want string }{
		{`O'Brien`, `O\'Brien`},
		{`a"b`, `a\"b`},
		{`a\b`, `a\\b`},
		{"a\nb", `a\nb`},
		{"a\rb", `a\rb`},
		{"ปกติ", "ปกติ"},
	}
	for _, c := range cases {
		if got := SqlStr(c.in); got != c.want {
			t.Errorf("SqlStr(%q) = %q, ต้องได้ %q", c.in, got, c.want)
		}
	}
}
