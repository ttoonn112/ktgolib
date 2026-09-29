package ktgolib

// รายการเคสที่ถูกล็อกไว้ (คู่กับ behavior_lock_test.go)
//
// เก็บเฉพาะฟังก์ชันที่ผลลัพธ์ขึ้นกับ input อย่างเดียว — ตัวที่อ่านเวลาปัจจุบัน สุ่มค่า
// ต่อเน็ต อ่าน/เขียนไฟล์ หรือจำสถานะไว้ในตัว ล็อกแบบนี้ไม่ได้ (ดูรายการท้ายไฟล์)
//
// เพิ่มฟังก์ชันใหม่เข้าไลบรารี = เพิ่มเคสที่นี่ด้วย แล้วสร้าง golden ใหม่

import (
	"sort"
	"time"
)

// sortedByCode — เรียงผลลัพธ์ให้คงที่ก่อนเทียบ สำหรับฟังก์ชันที่วนบน map
func sortedByCode(in []map[string]interface{}) []map[string]interface{} {
	out := append([]map[string]interface{}{}, in...)
	sort.Slice(out, func(i, j int) bool { return T(out[i], "code") < T(out[j], "code") })
	return out
}

var (
	fixedT1 = time.Date(2026, 9, 28, 10, 30, 0, 0, time.Local)
	fixedT2 = time.Date(2026, 9, 28, 12, 0, 0, 0, time.Local)

	sampleMap = map[string]interface{}{
		"code": "C001", "name": "บริษัท ทดสอบ", "qty": 5, "rate": 1.5, "flag": true, "none": nil,
	}
	sampleList = []map[string]interface{}{
		{"code": "b", "grp": "x", "n": "2"},
		{"code": "a", "grp": "x", "n": "10"},
		{"code": "c", "grp": "y", "n": "1"},
	}
	// ค่าที่ผู้ใช้พิมพ์แล้วเคยทำระบบพัง — ต้องล็อกผลไว้ทุกตัว
	nastyValues = map[string]string{
		"plain":     "สมชาย",
		"quote":     `O'Brien`,
		"dquote":    `say "hi"`,
		"backslash": `C:\new`,
		"bs_end":    `path\`,
		"amp":       "A & B",
		"newline":   "line1\nline2",
		"sentinel":  "ราคา u0026 ส่วนลด",
		"attack":    "x' or 1=1 --",
		"empty":     "",
		"all":       "All",
		"multi":     "A|N",
	}
)

func lockedCases() []lockCase {
	var cs []lockCase
	add := func(key string, run func() interface{}) { cs = append(cs, lockCase{key, run}) }

	// ── sql.go — ตัวที่เคยเปลี่ยนพฤติกรรมแล้วทำ 3 repo พัง ต้องล็อกทุกค่า
	for n, v := range nastyValues {
		v := v
		add("SqlStr/"+n, func() interface{} { return SqlStr(v) })
		add("AddSqlFilter/"+n, func() interface{} { return AddSqlFilter("A.col", v) })
		add("AddSqlLikeFilter/"+n, func() interface{} { return AddSqlLikeFilter("A.col", v) })
		add("AddSqlLikePrefixFilter/"+n, func() interface{} { return AddSqlLikePrefixFilter("A.col", v) })
		add("AddSqlLikeSuffixFilter/"+n, func() interface{} { return AddSqlLikeSuffixFilter("A.col", v) })
		add("AddSqlMultipleFilter/"+n, func() interface{} { return AddSqlMultipleFilter("A.col", v) })
		add("AddSqlNotInMultipleFilter/"+n, func() interface{} { return AddSqlNotInMultipleFilter("A.col", v) })
		add("GetSqlMultipleFilter/"+n, func() interface{} { return GetSqlMultipleFilter(v) })
		add("AddSqlDateRangeFilter/"+n, func() interface{} { return AddSqlDateRangeFilter("A.d", "2026-01-01", v) })
	}
	add("AddSqlDateRangeFilter/normal", func() interface{} { return AddSqlDateRangeFilter("A.d", "2026-01-01", "2026-12-31") })
	add("AddSqlDateRangeFilter/zerodate", func() interface{} { return AddSqlDateRangeFilter("A.d", "0000-00-00", "2026-12-31") })

	// ── json.go — รูปแบบข้อมูลที่เก็บถาวรในฐานข้อมูลของ consumer
	for n, v := range nastyValues {
		v := v
		add("EncloseText/"+n, func() interface{} { return EncloseText(v) })
		add("UncloseText/"+n, func() interface{} { return UncloseText(v) })
		add("RoundTrip.Compress.Extract/"+n, func() interface{} {
			return T(Extract(Compress(map[string]interface{}{"x": v})), "x")
		})
		add("RoundTrip.CompressArray.ExtractArray/"+n, func() interface{} {
			return ExtractArray(CompressArray([]map[string]interface{}{{"x": v}}))
		})
		add("Compress/"+n, func() interface{} { return Compress(map[string]interface{}{"x": v}) })
		add("CompressRaw/"+n, func() interface{} { return CompressRaw(map[string]interface{}{"x": v}) })
	}
	// ชื่อ field ที่มีอักขระพิเศษ — ล็อกไว้ว่าตอนนี้ "ไม่" ถูกเข้ารหัส (ช่องโหว่ที่รู้แล้ว ดู AGENTS.md)
	add("Compress/keyWithQuote", func() interface{} { return Compress(map[string]interface{}{`a'b`: "x"}) })
	add("CompressArray/keyWithQuote", func() interface{} {
		return CompressArray([]map[string]interface{}{{`a'b`: "x"}})
	})
	add("Compress/zeroValuesDropped", func() interface{} {
		return Compress(map[string]interface{}{"i": 0, "f": 0.0, "s": "", "keep": "v"})
	})
	add("CompressRaw/zeroValuesKept", func() interface{} {
		return CompressRaw(map[string]interface{}{"i": 0, "f": 0.0, "s": "", "keep": "v"})
	})
	add("Compress/mixed", func() interface{} { return Compress(sampleMap) })
	add("CompressArray/empty", func() interface{} { return CompressArray(nil) })
	add("Extract/invalidJson", func() interface{} { return Extract(`{"a":`) })
	add("ExtractArray/invalidJson", func() interface{} { return ExtractArray(`[{"a":`) })
	add("Extract/legacyData", func() interface{} { return Extract(`{"n":"Au0026#39;B"}`) })
	add("ExtractArray/legacyData", func() interface{} { return ExtractArray(`[{"n":"A u0026 B"}]`) })

	// ── convert.go
	add("MapToString/sample", func() interface{} { return MapToString(sampleMap) })
	add("MapToString/secrets", func() interface{} {
		return MapToString(map[string]interface{}{"username": "u", "password": "p1", "setup_password": "p2"})
	})
	add("MapToString/emptyPassword", func() interface{} {
		return MapToString(map[string]interface{}{"password": "", "username": "u"})
	})
	add("ArrayOfMapToString/sample", func() interface{} { return ArrayOfMapToString(sampleList) })
	add("RawStringToPureASCII/thai", func() interface{} { return RawStringToPureASCII("บริษัท ABC\n\tx") })
	add("CopyMap/sample", func() interface{} { return CopyMap(sampleMap) })
	add("ListToMap/sample", func() interface{} { return ListToMap(sampleList, "code") })
	add("ListToMapOfArray/sample", func() interface{} { return ListToMapOfArray(sampleList, "grp") })
	// สองตัวนี้วนบน map ซึ่ง Go สุ่มลำดับทุกครั้ง — ล็อก "สมาชิกที่ได้" ไม่ใช่ "ลำดับ"
	// (ลำดับไม่ใช่สัญญาของฟังก์ชัน ล็อกไปจะกลายเป็นเทสที่แดงสลับเขียวเอง)
	add("MapToList/sample", func() interface{} { return sortedByCode(MapToList(ListToMap(sampleList, "code"))) })
	add("MapOfArrayToList/sample", func() interface{} {
		return sortedByCode(MapOfArrayToList(ListToMapOfArray(sampleList, "grp")))
	})

	// ── fm.go / util.go — ตัวแปลงค่าและตัวจัดรูปแบบที่ consumer ใช้ทุกหน้า
	add("T/present", func() interface{} { return T(sampleMap, "code") })
	add("T/missing", func() interface{} { return T(sampleMap, "nope") })
	add("T/number", func() interface{} { return T(sampleMap, "qty") })
	add("T/float", func() interface{} { return T(sampleMap, "rate") })
	add("T/bool", func() interface{} { return T(sampleMap, "flag") })
	add("T/nil", func() interface{} { return T(sampleMap, "none") })
	add("MapVStr/default", func() interface{} { return MapVStr(sampleMap, "nope", "ค่าเริ่มต้น") })
	add("M/missing", func() interface{} { return M(sampleMap, "nope") })
	add("Has/yes", func() interface{} { return Has(sampleMap, "code") })
	add("Has/no", func() interface{} { return Has(sampleMap, "nope") })
	add("SI/float", func() interface{} { return SI(sampleMap, "rate") })
	add("SI64/number", func() interface{} { return SI64(sampleMap, "qty") })
	add("SF64/number", func() interface{} { return SF64(sampleMap, "rate") })
	add("NObj/nested", func() interface{} {
		return NObj(map[string]map[string]interface{}{"a": {"b": "c"}}, "a", "b")
	})
	add("GetMask/plain", func() interface{} { return GetMask(sampleMap, []string{"code", "qty"}) })
	add("GetMask/prefix", func() interface{} {
		return GetMask(map[string]interface{}{"info_a": "1", "info_b": "2", "x": "3"}, []string{"info_"})
	})
	add("GetMask/prefixQuoteKey", func() interface{} {
		return GetMask(map[string]interface{}{`info_a'b`: "1"}, []string{"info_"})
	})
	add("GetMask/array", func() interface{} {
		return GetMask(map[string]interface{}{"rows": []interface{}{map[string]interface{}{"k": "v"}}}, []string{"rows!"})
	})
	for _, s := range []string{"", "0", "12", "-3", "3.7", "abc", " 5 "} {
		s := s
		add("S_I/"+s, func() interface{} { return S_I(s) })
		add("S_I64/"+s, func() interface{} { return S_I64(s) })
		add("S_I32/"+s, func() interface{} { return S_I32(s) })
		add("S_F64/"+s, func() interface{} { return S_F64(s) })
	}
	add("I_S", func() interface{} { return I_S(42) })
	add("I64_S", func() interface{} { return I64_S(-42) })
	add("I32_S", func() interface{} { return I32_S(7) })
	for _, f := range []float64{0, 1.005, 2.5, -2.5, 1234.5678, 1e6} {
		f := f
		add("F64_S/2/"+F64_S(f, 6), func() interface{} { return F64_S(f, 2) })
		add("F64_S_AUTO/"+F64_S(f, 6), func() interface{} { return F64_S_AUTO(f) })
		add("FormatNumber/2/"+F64_S(f, 6), func() interface{} { return FormatNumber(f, 2, true) })
		add("FormatNumber/noZero/"+F64_S(f, 6), func() interface{} { return FormatNumber(f, 2, false) })
		add("ToFixed/2/"+F64_S(f, 6), func() interface{} { return ToFixed(f, 2) })
		add("Round/"+F64_S(f, 6), func() interface{} { return Round(f) })
		add("RoundToInt64/"+F64_S(f, 6), func() interface{} { return RoundToInt64(f) })
	}
	add("ZeroString/pad", func() interface{} { return ZeroString(42, 6) })
	add("ZeroString/overflow", func() interface{} { return ZeroString(1234567, 3) })
	add("FirstXChar/short", func() interface{} { return FirstXChar("abc", 10) })
	add("FirstXChar/cut", func() interface{} { return FirstXChar("abcdef", 3) })
	add("LastXChar/cut", func() interface{} { return LastXChar("abcdef", 3) })
	add("OC/true", func() interface{} { return OC(true, "a", "b") })
	add("OC/false", func() interface{} { return OC(false, "a", "b") })
	add("Includes/hit", func() interface{} { return Includes("0001", "0000|0001", "|") })
	add("Includes/miss", func() interface{} { return Includes("0009", "0000|0001", "|") })
	add("HasIncludes/hit", func() interface{} { return HasIncludes([]string{"a", "b"}, []string{"b", "c"}) })
	add("IsLetterOnly/yes", func() interface{} { return IsLetterOnly("abc123") })
	add("IsLetterOnly/no", func() interface{} { return IsLetterOnly("abc-123") })
	add("HasRegexSpecialChars/yes", func() interface{} { return HasRegexSpecialChars("a.*b") })
	add("HasRegexSpecialChars/no", func() interface{} { return HasRegexSpecialChars("abc") })
	for n, v := range nastyValues {
		v := v
		add("DoLetterOnly/"+n, func() interface{} { return DoLetterOnly(v) })
	}
	add("DivMod", func() interface{} { a, b := DivMod(17, 5); return []int64{a, b} })
	add("Median/odd", func() interface{} { return Median([]float64{3, 1, 2}) })
	add("Median/even", func() interface{} { return Median([]float64{4, 1, 2, 3}) })
	add("FormatDuration/secs", func() interface{} { return FormatDuration(3725) })
	add("FormatDuration/zero", func() interface{} { return FormatDuration(0) })

	// ── datetime.go (ตรึงโซนเวลาไว้ที่ +07:00 ใน behavior_lock_test.go)
	add("DateTimeFormat/date", func() interface{} { return DateTimeFormat("2026-09-28 10:30:00", "Date", "th") })
	add("DateTimeFormat/datetime", func() interface{} { return DateTimeFormat("2026-09-28 10:30:00", "DateTime", "th") })
	add("DateTimeFormat/timestamp", func() interface{} { return DateTimeFormat("2026-09-28 10:30:00", "Timestamp", "") })
	add("DateTimeFormat/en", func() interface{} { return DateTimeFormat("2026-09-28 10:30:00", "Date", "en") })
	add("DateTimeFormat/invalid", func() interface{} { return DateTimeFormat("ไม่ใช่วันที่", "Date", "th") })
	add("AddDate/plus", func() interface{} { return AddDate("2026-09-28", 5) })
	add("AddDate/minus", func() interface{} { return AddDate("2026-03-01", -1) })
	add("DateTimeAddString", func() interface{} { return DateTimeAddString("2026-09-28 10:30:00", 90) })
	add("DateTimeAddYMDString", func() interface{} { return DateTimeAddYMDString("2026-09-28 10:30:00", 1, 2, 3) })
	add("DateDiff", func() interface{} { return DateDiff("2026-09-01", "2026-09-28") })
	add("DateTimeDiff", func() interface{} { return DateTimeDiff("2026-09-28 10:00:00", "2026-09-28 12:30:00") })
	add("DateTimeString", func() interface{} { return DateTimeString(fixedT1) })
	add("DateTimeValue/ok", func() interface{} {
		v, err := DateTimeValue("2026-09-28 10:30:00")
		return []string{DateTimeString(v), render(err == nil)}
	})
	add("DateTimeValue/bad", func() interface{} {
		_, err := DateTimeValue("ไม่ใช่วันที่")
		return err != nil
	})
	add("DateTimeValueDiff", func() interface{} { return DateTimeValueDiff(fixedT1, fixedT2) })
	add("DateTimeValueDiffSec", func() interface{} { return DateTimeValueDiffSec(fixedT1, fixedT2) })
	add("GetTimeHourFromSecond", func() interface{} { return GetTimeHourFromSecond(3725) })
	add("GetTimeHourFromSecond/zero", func() interface{} { return GetTimeHourFromSecond(0) })
	add("GetTimeHourFromSecondNoBlank/zero", func() interface{} { return GetTimeHourFromSecondNoBlank(0) })
	add("GetSecondFromHHmmss", func() interface{} {
		s, err := GetSecondFromHHmmss("01:02:05")
		return []string{I_S(s), render(err == nil)}
	})
	add("ToRFC3339", func() interface{} { return ToRFC3339("2026-09-28 10:30:00") })
	add("FromRFC3339", func() interface{} { return FromRFC3339("2026-09-28T10:30:00+07:00") })
	add("IsOverlap/yes", func() interface{} {
		return IsOverlap("2026-09-01", "2026-09-10", "2026-09-05", "2026-09-20")
	})
	add("IsOverlap/no", func() interface{} {
		return IsOverlap("2026-09-01", "2026-09-10", "2026-09-11", "2026-09-20")
	})

	// ── list.go
	add("ListSort/asc", func() interface{} { return ListSort(copyList(sampleList), "code", "asc") })
	add("ListSort/desc", func() interface{} { return ListSort(copyList(sampleList), "code", "desc") })
	add("ListSortNumber/asc", func() interface{} { return ListSortNumber(copyList(sampleList), "n", "asc") })
	add("ListReverse", func() interface{} { return ListReverse(copyList(sampleList)) })
	add("SortedListToArrayOfArray", func() interface{} { return SortedListToArrayOfArray(copyList(sampleList), "grp") })
	add("ListToArrayOfArraySortByKey", func() interface{} {
		return ListToArrayOfArraySortByKey(copyList(sampleList), "grp", "code", "asc")
	})

	return cs
}

func copyList(in []map[string]interface{}) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(in))
	for _, m := range in {
		out = append(out, CopyMap(m))
	}
	return out
}

// ฟังก์ชันที่ล็อกแบบนี้ไม่ได้ (ผลไม่ขึ้นกับ input อย่างเดียว) — ยังห้ามเปลี่ยนพฤติกรรมเหมือนกัน
// ต้องตรวจด้วยมือตอนแก้:
//   เวลาปัจจุบัน : Now · NowDate · NextDate
//   สุ่ม         : GenerateRandomString · GenerateRandomNumberString
//   ต่อเน็ต      : Http_Get · Http_PostForm · Http_PostFormWithHeader · Http_PostJson · Notify_*
//   ไฟล์         : CsvLog · CsvRead · Log · LogHidden · LogWithDuration · LogHiddenWithDuration
//   จำสถานะ/เวลา : Operation_IsLimitExceeded* · StartTask · StartScheduleTask · Attempt · TryCatch
//   พิมพ์จอ      : Print · Println
