package task

import (
	"github.com/robfig/cron/v3"
	"kmpi-go/log"
	"kmpi-go/service"
)

func init() {
	c := cron.New(cron.WithSeconds())
	taskCron := "0 0 0/1 * * ? "
	//taskCron := "0/5 * * * * ? "
	_, err := c.AddFunc(taskCron, func() {
		checkAndDelete()
	})
	if err != nil {
		log.Error("AddFunc error:", err.Error())
		return
	}

	c.Start()
}
func checkAndDelete() {
	count, err := service.LogCount() //modify the history count and delete the count
	if err != nil {
		log.Error("LogCount error", err.Error())
		return
	}
	if *count > 610*10000 {
		service.LogDeleteLast(10 * 10000)
	}
	log.Info("history count is %d", *count)
}
