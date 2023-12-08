package ktgolib

import (
  "time"
  "github.com/jasonlvhit/gocron"	// cron job
)

type scheduletask func()

func StartScheduleTask(repeatsec uint64,runfirst bool,task scheduletask){
  defer TryCatch(func(errStr string) {
    Log("lib.StartScheduleTask", "", "", errStr, "Task")
  })
  if(runfirst){ task() }
  c := make(chan int)
  defer close(c)
  go doScheduleTask(c,repeatsec,task)
}

func StartTask(nextsec uint64, task scheduletask){
  defer TryCatch(func(errStr string) {
    Log("lib.StartTask", "", "", errStr, "Task")
  })
  c := make(chan int)
  defer close(c)
  go doTask(c,nextsec,task)
}

func doScheduleTask(work chan int,repeatsec uint64,task scheduletask){
  defer TryCatch(func(errStr string) {
    Log("lib.doScheduleTask", "", "", errStr, "Task")
  })
  s := gocron.NewScheduler()
	s.Every(repeatsec).Seconds().Do(task)
	<- s.Start()
}

func doTask(work chan int,nextsec uint64,task scheduletask){
  defer TryCatch(func(errStr string) {
    Log("lib.doTask", "", "", errStr, "Task")
  })
  time.Sleep(time.Duration(nextsec)*time.Second)
  task()
}
