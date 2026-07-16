package ktgolib

// Shortcut Function
// ใช้งานเพื่อให้ Code สั้นลง

import (
	"strings"
)

// SqlStr — escape ค่าที่จะเอาไปวางใน single-quoted string literal ของ SQL (กัน SQL injection)
// คืน "ค่าที่ escape แล้ว" ไม่รวม quote ครอบ (ผู้เรียกเป็นคนใส่ '...' เอง)
// เทียบเท่า mysql_real_escape_string: neutralize \ ' " NUL \n \r Ctrl-Z
// ทำแบบ byte-wise ปลอดภัยกับ UTF-8 (multi-byte มี high bit เสมอ ไม่ชนกับอักขระ ASCII เหล่านี้)
// ⚠️ ใช้กับ "ค่า" เท่านั้น — ห้ามใช้กับชื่อ column/table (identifier ไม่ได้อยู่ใน quote)
func SqlStr(value string) string {
	var b strings.Builder
	b.Grow(len(value) + 8)
	for i := 0; i < len(value); i++ {
		switch c := value[i]; c {
		case '\\':
			b.WriteString(`\\`)
		case '\'':
			b.WriteString(`\'`)
		case '"':
			b.WriteString(`\"`)
		case 0:
			b.WriteString(`\0`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case 0x1a:
			b.WriteString(`\Z`)
		default:
			b.WriteByte(c)
		}
	}
	return b.String()
}

func AddSqlFilter(field_name string, value string) string{
	if value != "All" && value != "" {
		filter := " and "+field_name+" = '"+SqlStr(value)+"' "
		return filter
	}
	return ""
}

func AddSqlLikePrefixFilter(field_name string, value string) string{
	if value != "All" && value != "" {
		filter := " and "+field_name+" like '%"+SqlStr(value)+"' "
		return filter
	}
	return ""
}

func AddSqlLikeSuffixFilter(field_name string, value string) string{
	if value != "All" && value != "" {
		filter := " and "+field_name+" like '"+SqlStr(value)+"%' "
		return filter
	}
	return ""
}

func AddSqlLikeFilter(field_name string, value string) string{
	if value != "All" && value != "" {
		filter := " and "+field_name+" like '%"+SqlStr(value)+"%' "
		return filter
	}
	return ""
}

func AddSqlDateRangeFilter(field_name string, start_date string, end_date string) string{
	if start_date != "" && FirstXChar(start_date,10) != "0000-00-00" && end_date != "" && FirstXChar(end_date,10) != "0000-00-00" {
		filter := " and "+field_name+" between '"+SqlStr(start_date)+"' and '"+SqlStr(end_date)+"' "
		return filter
	}
	return ""
}

func AddSqlMultipleFilter(field_name string, value string) string{
	if value != "All" && value != "" {
		values := strings.Split(value,"|")
		if len(values) > 0 {
			filter := " and "+field_name+" in ("
			for _, skey := range values {
				filter += "'"+SqlStr(skey)+"',"
			}
			filter = filter[:len(filter)-1]
			filter += ")"
			return filter
		}
	}
	return ""
}

func AddSqlNotInMultipleFilter(field_name string, value string) string{
	if value != "All" && value != "" {
		values := strings.Split(value,"|")
		if len(values) > 0 {
			filter := " and "+field_name+" not in ("
			for _, skey := range values {
				filter += "'"+SqlStr(skey)+"',"
			}
			filter = filter[:len(filter)-1]
			filter += ")"
			return filter
		}
	}
	return ""
}

func GetSqlMultipleFilter(value string) string{
	if value != "All" && value != "" {
		values := strings.Split(value,"|")
		if len(values) > 0 {
			filter := "("
			for _, skey := range values {
				filter += "'"+SqlStr(skey)+"',"
			}
			filter = filter[:len(filter)-1]
			filter += ")"
			return filter
		}
	}
	return ""
}
