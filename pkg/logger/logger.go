package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"time"
)

type GLogger struct {
	Log zerolog.Logger
}

type Field struct {
	Key   string
	Value interface{}
}

func NewGLogger(writers ...io.Writer) GLogger {
	multi := zerolog.MultiLevelWriter(writers...)
	w := zerolog.ConsoleWriter{
		Out:        multi,
		TimeFormat: time.RFC822,
		PartsOrder: []string{"level", "time", "message"},
	}
	return GLogger{log.Output(w)}
}

func (g GLogger) Trace(msg string, keyvals ...Field) {
	event := g.Log.Trace()
	addFields(event, keyvals...)
	event.Msg(msg)
}

func (g GLogger) Info(msg string, keyvals ...Field) {
	event := g.Log.Info()
	addFields(event, keyvals...)
	event.Msg(msg)
}

func (g GLogger) Debug(msg string, keyvals ...Field) {
	event := g.Log.Debug()
	addFields(event, keyvals...)
	event.Msg(msg)
}

func (g GLogger) Warning(msg string, keyvals ...Field) {
	event := g.Log.Warn()
	addFields(event, keyvals...)
	event.Msg(msg)
}

func (g GLogger) Error(msg string, err error, keyvals ...Field) {
	event := g.Log.Error().Err(err)
	addFields(event, keyvals...)
	event.Msg(msg)
}

func addFields(e *zerolog.Event, keyvals ...Field) {
	for _, field := range keyvals {
		e = e.Interface(fmt.Sprintf("%s", field.Key), field.Value)
	}
}
