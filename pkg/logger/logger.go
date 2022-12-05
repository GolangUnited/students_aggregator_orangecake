package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"time"
)

type gLogger struct {
	Log zerolog.Logger
}

type Field struct {
	Key   string
	Value interface{}
}

func NewGLogger(aWriters ...io.Writer) gLogger {
	lMulti := zerolog.MultiLevelWriter(aWriters...)
	w := zerolog.ConsoleWriter{
		Out:        lMulti,
		TimeFormat: time.RFC822,
		PartsOrder: []string{"level", "time", "message"},
	}
	return gLogger{log.Output(w)}
}

func (g gLogger) Trace(aMsg string, aKeyvals ...Field) {
	lEvent := g.Log.Trace()
	addFields(lEvent, aKeyvals...)
	lEvent.Msg(aMsg)
}

func (g gLogger) Info(aMsg string, aKeyvals ...Field) {
	lEvent := g.Log.Info()
	addFields(lEvent, aKeyvals...)
	lEvent.Msg(aMsg)
}

func (g gLogger) Debug(aMsg string, aKeyvals ...Field) {
	lEvent := g.Log.Debug()
	addFields(lEvent, aKeyvals...)
	lEvent.Msg(aMsg)
}

func (g gLogger) Warning(aMsg string, aKeyvals ...Field) {
	lEvent := g.Log.Warn()
	addFields(lEvent, aKeyvals...)
	lEvent.Msg(aMsg)
}

func (g gLogger) Error(aMsg string, aErr error, aKeyvals ...Field) {
	lEvent := g.Log.Error().Err(aErr)
	addFields(lEvent, aKeyvals...)
	lEvent.Msg(aMsg)
}

func addFields(e *zerolog.Event, aKeyvals ...Field) {
	for _, field := range aKeyvals {
		e = e.Interface(fmt.Sprintf("%s", field.Key), field.Value)
	}
}
