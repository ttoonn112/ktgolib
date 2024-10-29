package ktgolib

import (
  "sync"
  "golang.org/x/time/rate"
)

var userOperationLimiters = make(map[string]map[string]*rate.Limiter)
var mu sync.Mutex

func Operation_IsLimitedExceed(username string, operation string, key string) bool {
  mu.Lock()
  defer mu.Unlock()

  if userOperationLimiters[username] == nil {
    userOperationLimiters[username] = make(map[string]*rate.Limiter)
  }

  limiter, exists := userOperationLimiters[username][operation]
  if !exists {
    limiter = rate.NewLimiter(1, 1)       // กำหนด 1 request per second
    userOperationLimiters[username][operation] = limiter
  }

  result := !limiter.Allow()
  if result {
    LogHidden(operation, username, key, "", "OperationLimitedExceed")
  }

  return result
}
