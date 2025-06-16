package ktgolib

import (
  "sync"
  "time"
  "golang.org/x/time/rate"
)

var userOperationLimiters = make(map[string]map[string]*rate.Limiter)
var mu sync.Mutex

func Operation_IsLimitExceeded(username string, operation string, key string) bool {
  mu.Lock()
  defer mu.Unlock()

  if userOperationLimiters[username] == nil {
    userOperationLimiters[username] = make(map[string]*rate.Limiter)
  }

  limiter, exists := userOperationLimiters[username][operation]
  if !exists {
    limiter = rate.NewLimiter(1, 1)       // 1 request ต่อวินาที, เก็บ Token ได้สูงสุด 1
    userOperationLimiters[username][operation] = limiter
  }

  result := !limiter.Allow()
  if result {
    LogHidden("Operation_IsLimitExceeded", username, key, operation, "SECURITY_RISK")
  }

  return result
}

// sec = "จำนวนวินาทีที่ต้องรอ ก่อนทำ Operation ครั้งถัดไป"
func Operation_IsLimitExceededBySec(username string, operation string, key string, sec int) bool {
  mu.Lock()
  defer mu.Unlock()

  if userOperationLimiters[username] == nil {
    userOperationLimiters[username] = make(map[string]*rate.Limiter)
  }

  limiter, exists := userOperationLimiters[username][operation]
  if !exists {
    limiter = rate.NewLimiter(rate.Every(time.Duration(sec) * time.Second), 1)       // จะสามารถทำงานได้ เพียง 1 ครั้งทุก ๆ [sec] วินาที โดยไม่มีการสะสม Token เพื่อทำงานหลายครั้งต่อเนื่อง
    userOperationLimiters[username][operation] = limiter
  }

  result := !limiter.Allow()
  if result {
    LogHidden("Operation_IsLimitExceededBySec", username, key, operation, "SECURITY_RISK")
  }

  return result
}

// sec = "จำนวนวินาทีที่ต้องรอ ก่อนทำ Operation ครั้งถัดไป"
// ระบบสามารถรับคำขอได้ ทันทีสูงสุด [burst] requests อย่างรวดเร็ว
// หลังจากนั้น คำขอใหม่จะต้องรอให้ "โทเค็น" เพิ่มขึ้นตามอัตรา [sec] requests/second
func Operation_IsLimitExceededBySecToken(username string, operation string, key string, sec int, burst int) bool {
  mu.Lock()
  defer mu.Unlock()

  if userOperationLimiters[username] == nil {
    userOperationLimiters[username] = make(map[string]*rate.Limiter)
  }

  limiter, exists := userOperationLimiters[username][operation]
  if !exists {
    limiter = rate.NewLimiter(rate.Every(time.Duration(sec) * time.Second), burst)       // จะสามารถทำงานได้ เพียง 1 ครั้งทุก ๆ [sec] วินาที โดยไม่มีการสะสม Token เพื่อทำงานหลายครั้งต่อเนื่อง
    userOperationLimiters[username][operation] = limiter
  }

  result := !limiter.Allow()
  if result {
    LogHidden("Operation_IsLimitExceededBySecToken", username, key, operation, "SECURITY_RISK")
  }

  return result
}
