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
	Info(msg string, keyvals ...Value)
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

func (g *gLogger) Info(msg string, keyvals ...Value) {
	event := g.Log.Info()
	addFields(event, keyvals...)
	event.Msg(msg)
}

func addFields(e *zerolog.Event, keyvals ...Value) {
	for _, value := range keyvals {
		e = e.Interface(fmt.Sprintf("%s", value.Key), value.Payload)
	}
}

var logger = newGLogger(os.Stdout)
