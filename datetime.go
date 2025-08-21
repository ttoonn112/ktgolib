package ktgolib

import (
	"math"
	"time"
	"strings"
	"fmt"
)

func Now() string {
	return DateTimeString(time.Now())
}

func NowDate() string {
	nowstr := Now()
	return nowstr[:10]
}

func NextDate(day int) string {
	nowstr := Now()
	return DateTimeAddString(nowstr, day*3600*24)[:10]
}

func AddDate(datestr string, day int) string {
	if len(datestr) == 10 { datestr = datestr+" 00:00:00" }
	t, err := DateTimeValue(datestr)
	if err != nil {
		return ""
	}
	return DateTimeString(t.Add(time.Second * time.Duration(day*3600*24)))[:10]
}

func DateTimeValueDiffSec(less time.Time, more time.Time) int64{
  diff := more.Sub(less)
  return int64(diff/1000/time.Millisecond)
}

func DateTimeAddString(datestr string, sec int) string{
  t, err := DateTimeValue(datestr)
	if err != nil {
		return ""
	}
	return DateTimeString(t.Add(time.Second * time.Duration(sec)))
}

func DateTimeAddYMDString(datestr string, year int, month int, day int) string{
  t, err := DateTimeValue(datestr)
	if err != nil {
		return ""
	}
	return DateTimeString(t.AddDate(year, month, day))
}

func DateTimeValue(datestr string) (time.Time, error){
	if len(datestr) != 10 && len(datestr) != 19 && len(datestr) != 16 {
		return time.Parse("2006-01-02 15:04:05", "1900-01-01 00:00:00")
	}
  return time.Parse("2006-01-02 15:04:05", datestr)
}

func DateTimeString(t time.Time) string{
  return t.Format("2006-01-02 15:04:05")
}

func DateTimeDiff(less string, more string) int64{
	if len(less) != 10 && len(less) != 19 && len(less) != 16 {return 0}
	if len(more) != 10 && len(more) != 19 && len(more) != 16 {return 0}
	if less == "" || less == "0000-00-00 00:00:00" || more == "" || more == "0000-00-00 00:00:00" {return 0}
  t1, err1 := DateTimeValue(less)
  t2, err2 := DateTimeValue(more)
  if err1 != nil { panic(err1) }
	if err2 != nil { panic(err2) }
	return DateTimeValueDiff(t1, t2)
}

func DateTimeValueDiff(t1 time.Time, t2 time.Time) int64{
  diff := t2.Sub(t1)
  return int64(diff/1000/time.Millisecond)
}

func DateDiff(less string, more string) int64{
  return DateTimeDiff(less+" 00:00:00", more+" 00:00:00")/(3600 * 24)
}

// แปลงจาก "YYYY-MM-DD HH:MM:SS" (ตีความเป็นเวลาโซนของเครื่อง) -> UTC/ISO8601 (RFC3339)
func ToRFC3339(src string) string {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", src, time.Local) // ใช้โซนของเครื่อง
	if err != nil {
		return "INVALID"
	}
	return t.UTC().Format(time.RFC3339) // ออกเป็น UTC เสมอ
}

// แปลงจาก UTC/ISO8601 (RFC3339) -> "YYYY-MM-DD HH:MM:SS" (ฟอร์แมตตามโซนของเครื่อง)
func FromRFC3339(src string) string {
	// ลอง parse แบบ RFC3339 หากไม่ผ่านจะลอง RFC3339Nano เผื่อมี nanoseconds
	t, err := time.Parse(time.RFC3339, src)
	if err != nil {
		t, err = time.Parse(time.RFC3339Nano, src)
		if err != nil {
			return "INVALID"
		}
	}
	return t.In(time.Local).Format("2006-01-02 15:04:05") // แสดงตามโซนของเครื่อง
}

func DateTimeFormat(dtstr string, format string, lang string) string{
	if len(dtstr) != 10 && len(dtstr) != 19 && len(dtstr) != 16 {return ""}
	if lang == "" {lang = "en"}
	if len(dtstr) == 10 {
		if format == "" { format = "DisplayDate"}
	}else{
		if format == "" { format = "DisplayDateTime"}
	}
	syear := ""; smonth := ""; sday := ""; shour := ""; smin := ""; ssec := ""
	if strings.Contains(dtstr, "/") {
		sday = dtstr[0:2]
		smonth = dtstr[3:5]
		syear = dtstr[6:10]
	}else{
		syear = dtstr[0:4]
		smonth = dtstr[5:7]
		sday = dtstr[8:10]
	}
	if len(dtstr) == 19 {
		shour = dtstr[11:13]
		smin = dtstr[14:16]
		ssec = dtstr[17:19]
	}
	if len(dtstr) == 16 {
		shour = dtstr[11:13]
		smin = dtstr[14:16]
		ssec = "00"
	}

	datestr := ""
	switch dformat := lang+format; dformat {
		case "enDisplayDateTime", "thDisplayDateTime":
			datestr = sday+"/"+smonth+"/"+syear+" "+shour+":"+smin+":"+ssec
			if datestr == "00/00/0000 00:00:00" {datestr = ""}
		case "enDisplayDateTimeMin", "thDisplayDateTimeMin":
			datestr = sday+"/"+smonth+"/"+syear+" "+shour+":"+smin
			if datestr == "00/00/0000 00:00" {datestr = ""}
		case "enShortDisplayDateTimeMin", "thShortDisplayDateTimeMin":
			datestr = sday+"/"+smonth+"/"+LastXChar(syear,2)+" "+shour+":"+smin
			if datestr == "00/00/00 00:00" {datestr = ""}
		case "enDisplayDate", "thDisplayDate":
			datestr = sday+"/"+smonth+"/"+syear
			if datestr == "00/00/0000" {datestr = ""}
		case "enShortDisplayDate", "thShortDisplayDate":
			datestr = sday+"/"+smonth+"/"+LastXChar(syear,2)
			if datestr == "00/00/00" {datestr = ""}
		case "enDBDateTime", "thDBDateTime":
			datestr = syear+"-"+smonth+"-"+sday+" "+shour+":"+smin+":"+ssec
		case "enDBDate", "thDBDate":
			datestr = syear+"-"+smonth+"-"+sday
		case "enTime", "thTime":
			datestr = shour+":"+smin+":"+ssec
		case "enTimeIfHave", "thTimeIfHave":
			datestr = shour+":"+smin+":"+ssec
			if datestr == "00:00:00" {datestr = ""}
		case "enTimeMin", "thTimeMin":
			datestr = shour+":"+smin
			if datestr == "00:00" {datestr = ""}
		case "enTimeMinF", "thTimeMinF":
			datestr = shour+":"+smin
		case "enDTABBR", "thDTABBR":
			datestr = syear[2:]+smonth+sday+shour+smin+ssec
		case "enLogDate", "thLogDate":
			datestr = syear[2:]+smonth+sday
		case "enLogDateTime", "thLogDateTime":
			datestr = sday+smonth+syear+""+shour+""+smin+""+ssec
		case "enTimestamp", "thTimestamp":
			datestr = syear+smonth+sday+""+shour+""+smin+""+ssec
		case "enDDMMYYYY":
			datestr = sday+smonth+syear
			if datestr == "00000000" {datestr = ""}
		case "enYYYYMMDD", "thYYYYMMDD":
			datestr = syear+smonth+sday
		case "enHHMMSS", "thHHMMSS":
			datestr = shour+smin+ssec
	}

  return datestr
}

func GetTimeHourFromSecondNoBlank(s int64) string{
	result := GetTimeHourFromSecond(s)
	if result == "00:00" {
		result = ""
	}
	return result
}

func GetTimeHourFromSecond(s int64) string{
	second := float64(0)
	min := float64(0)
	hour := float64(0)
	sec := float64(s)
	if(sec < 60){
		second = sec
	}else if(sec < 3600){
		second = math.Mod(sec,60)
		min = math.Floor(sec/60)
	}else if(sec < 3600 * 24){
		second = math.Mod(sec,60)
		min = math.Floor(math.Mod(sec,3600)/60)
		hour = math.Floor(sec/3600)
	}else {
		second = math.Mod(sec,60)
		min = math.Floor(math.Mod(sec,3600)/60)
		hour = math.Floor(sec/3600)
	}
	if second >= 30 {
		min += 1
		if min >= 60 {
			min = min - 60
			hour += 1
		}
	}
	hourStr := F64_S(hour,0)
	if(hour < 10){
		hourStr = "0"+hourStr
	}
	minStr := "00"+F64_S(min,0)
	minStr = minStr[len(minStr)-2:]
	return hourStr+":"+minStr
}

func GetSecondFromHHmmss(s string) (int, error) {
	t, err := time.Parse("15:04:05", s)
	if err != nil {
		return 0, err
	}
	return t.Hour()*3600 + t.Minute()*60 + t.Second(), nil
}

func FormatDuration(seconds int) string {
	days := seconds / 86400
	seconds %= 86400

	hours := seconds / 3600
	seconds %= 3600

	minutes := seconds / 60
	seconds %= 60

	result := ""

	switch {
	case days > 0:
		result += fmt.Sprintf("%d วัน", days)
		if hours > 0 {
			result += fmt.Sprintf(" %d ชม.", hours)
		}

	case hours > 0:
		result += fmt.Sprintf("%d ชม.", hours)
		if minutes > 0 {
			result += fmt.Sprintf(" %d น.", minutes)
		}

	case minutes > 0:
		result += fmt.Sprintf("%d น.", minutes)
		if seconds > 0 {
			result += fmt.Sprintf(" %d วิ", seconds)
		}

	default:
		result += fmt.Sprintf("%d วิ", seconds)
	}

	return result
}
