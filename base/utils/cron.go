package utils

import (
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
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
