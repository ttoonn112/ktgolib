package ktgolib

import (
  "sync"
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
    LogHidden(operation, username, key, "", "OperationLimitExceeded")
  }

  return result
}

func Operation_IsUserLimitExceeded(username string, operation string, key string) bool {
  mu.Lock()
  defer mu.Unlock()

  if userOperationLimiters[username] == nil {
    userOperationLimiters[username] = make(map[string]*rate.Limiter)
  }

  limiter, exists := userOperationLimiters[username][operation]
  if !exists {
    limiter = rate.NewLimiter(rate.Every(2 * time.Second), 1)       // จะสามารถทำงานได้ เพียง 1 ครั้งทุก ๆ 2 วินาที โดยไม่มีการสะสม Token เพื่อทำงานหลายครั้งต่อเนื่อง
    userOperationLimiters[username][operation] = limiter
  }

  result := !limiter.Allow()
  if result {
    LogHidden(operation, username, key, "", "OperationUserLimitExceeded")
  }

  return result
}
