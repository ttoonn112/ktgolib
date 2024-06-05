package ktgolib

// Shortcut Function
// ใช้งานเพื่อให้ Code สั้นลง

import (
	"strings"
)

func AddSqlFilter(field_name string, value string) string{
	if value != "All" && value != "" {
		filter := " and "+field_name+" = '"+value+"' "
		return filter
	}
	return ""
}

func AddSqlLikeFilter(field_name string, value string) string{
	if value != "All" && value != "" {
		filter := " and "+field_name+" like '%"+value+"%' "
		return filter
	}
	return ""
}

func AddSqlDateRangeFilter(field_name string, start_date string, end_date string) string{
	if start_date != "" && FirstXChar(start_date,10) != "0000-00-00" && end_date != "" && FirstXChar(end_date,10) != "0000-00-00" {
		filter := " and "+field_name+" between '"+start_date+"' and '"+end_date+"' "
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
				filter += "'"+skey+"',"
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
				filter += "'"+skey+"',"
			}
			filter = filter[:len(filter)-1]
			filter += ")"
			return filter
		}
	}
	return ""
}
