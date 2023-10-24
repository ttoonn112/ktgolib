package ktgolib

import (
	"time"
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

func DateTimeValue(datestr string) (time.Time, error){
	if len(datestr) != 10 && len(datestr) != 19 && len(datestr) != 16 {
		return time.Parse("2006-01-02 15:04:05", "1900-01-01 00:00:00")
	}
  return time.Parse("2006-01-02 15:04:05", datestr)
}

func DateTimeString(t time.Time) string{
  return t.Format("2006-01-02 15:04:05")
}
