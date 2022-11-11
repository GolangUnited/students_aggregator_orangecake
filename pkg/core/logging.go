package core

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"time"
)

type Logger interface {
	Trace(msg string, keyvals ...Value)
	Info(msg string, keyvals ...Value)
	Debug(msg string, keyvals ...Value)
	Warning(msg string, keyvals ...Value)
	Error(msg string, err error, keyvals ...Value)
}

type gLogger struct {
	Log zerolog.Logger
}

type Value struct {
	Key     string
	Payload interface{}
}

func newGLogger(writers ...io.Writer) gLogger {
	multi := zerolog.MultiLevelWriter(writers...)
	w := zerolog.ConsoleWriter{
		Out:        multi,
		TimeFormat: time.RFC822,
		PartsOrder: []string{"level", "time", "message"},
	}
	return gLogger{log.Output(w)}
}

func (g *gLogger) Trace(msg string, keyvals ...Value) {
	event := g.Log.Trace()
	addFields(event, keyvals...)
	event.Msg(msg)
}

func (g *gLogger) Info(msg string, keyvals ...Value) {
	event := g.Log.Info()
	addFields(event, keyvals...)
	event.Msg(msg)
}

func (g *gLogger) Debug(msg string, keyvals ...Value) {
	event := g.Log.Debug()
	addFields(event, keyvals...)
	event.Msg(msg)
}

func (g *gLogger) Warning(msg string, keyvals ...Value) {
	event := g.Log.Warn()
	addFields(event, keyvals...)
	event.Msg(msg)
}

func (g *gLogger) Error(msg string, err error, keyvals ...Value) {
	event := g.Log.Error().Err(err)
	addFields(event, keyvals...)
	event.Msg(msg)
}

func addFields(e *zerolog.Event, keyvals ...Value) {
	for _, value := range keyvals {
		e = e.Interface(fmt.Sprintf("%s", value.Key), value.Payload)
	}
}

var logger = newGLogger(os.Stdout)
