package log

import (
	"fmt"
	"github.com/LightningFootball/backend/database/models/log"
	"github.com/fatih/color"
	"strings"

	"github.com/LightningFootball/backend/base"
	"github.com/LightningFootball/backend/base/event"
	"github.com/LightningFootball/backend/base/exit"
)

// Writers should only be used in this package.
// Other packages should use listeners to receive logs.
type _writer interface {
	log(log Log)
}

// Writes to the console.
type consoleWriter struct {
	Level
}

// Writes to the database for reading from web.
type databaseWriter struct {
	Level
	queue chan Log
}

// Calling log listeners.
type eventWriter struct{}

var (
	colors = []func(format string, a ...interface{}) string{
		FATAL:   color.MagentaString,
		ERROR:   color.RedString,
		WARNING: color.YellowString,
		INFO:    color.GreenString,
		DEBUG:   color.CyanString,
	}
)

func (w *consoleWriter) log(l Log) {
	if l.Level >= w.Level {
		fmt.Print(colors[l.Level]("[%s][%s] ▶ %s ",
			l.Time.Format("15:04:05"),
			strings.Replace(l.Caller, "github.com/LightningFootball/backend/", "", -1),
			l.Level.String()),
			l.Message,
			"\n")
	}
}

func (w *databaseWriter) log(l Log) {
	// avoid blocking the main thread.
	if l.Level >= w.Level && !strings.HasPrefix(l.Caller, "gorm.io/gorm") {
		select {
		case w.queue <- l:
		default:
		}
	}
}

func (w *databaseWriter) init() {
	w.queue = make(chan Log, 100)
	exit.QuitWG.Add(1)
	go func() {
		for {
			select {
			case l := <-w.queue:
				ll := int(l.Level)
				lm := log.Log{
					Level:     &ll,
					Message:   l.Message,
					Caller:    l.Caller,
					CreatedAt: l.Time,
				}
				base.DB.Create(&lm)
			case <-exit.BaseContext.Done():
				select {
				case l := <-w.queue:
					ll := int(l.Level)
					lm := log.Log{
						Level:     &ll,
						Message:   l.Message,
						Caller:    l.Caller,
						CreatedAt: l.Time,
					}
					base.DB.Create(&lm)
				default:
					exit.QuitWG.Done()
					return
				}
			}
		}
	}()
}

func (w *eventWriter) log(l Log) {
	go event.FireEvent("log", l)
}
