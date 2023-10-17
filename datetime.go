package ktgolib

import (
	"time"
)

func DateTimeValueDiffSec(less time.Time, more time.Time) int64{
  diff := more.Sub(less)
  return int64(diff/1000/time.Millisecond)
}
