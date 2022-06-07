package utils

import (
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"time"
)

var c = cron.New()

func StartCron() {
	c.Start()
}

func AddCron(spec string, cmd func()) cron.EntryID {
	entryId, err := c.AddFunc(spec, cmd)
	if err != nil {
		err = errors.Wrap(err, "AddCron Failed")
	}
	return entryId
}

func RemoveCron(entryId cron.EntryID) {
	c.Remove(entryId)
}

func StopCron() {
	c.Stop()
}

// todo: Sleep这个主意由Copilot自动生成，本来想尝试time.timer，所以待验证是否可行
// 到期自动删除Cron任务
func RemoveCronByTime(entryID cron.EntryID, finalTime time.Time) {
	go func() {
		time.Sleep(finalTime.Sub(time.Now()))
		RemoveCron(entryID)
	}()
}
