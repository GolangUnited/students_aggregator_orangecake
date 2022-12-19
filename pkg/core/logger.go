package core

import (
	"io"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type zeroLogger struct {
	log zerolog.Logger
}

func NewZeroLogger(aWriters ...io.Writer) Logger {
	lMulti := zerolog.MultiLevelWriter(aWriters...)
	w := zerolog.ConsoleWriter{
		Out:        lMulti,
		TimeFormat: time.RFC822,
		PartsOrder: []string{"level", "time", "message"},
	}
	return zeroLogger{log.Output(w)}
}

func NewDebugZeroLogger(aWriters ...io.Writer) Logger {
	lMulti := zerolog.MultiLevelWriter(aWriters...)
	w := zerolog.ConsoleWriter{
		Out:        lMulti,
		PartsOrder: []string{"level", "message"},
		NoColor:    true,
	}
	return zeroLogger{log.Output(w)}
}

func (zero zeroLogger) Trace(aMessage string, aValues ...interface{}) {
	lEvent := zero.log.Trace()
	addFields(lEvent, aValues...)
	lEvent.Msg(aMessage)
}

func (zero zeroLogger) Info(aMessage string, aValues ...interface{}) {
	lEvent := zero.log.Info()
	addFields(lEvent, aValues...)
	lEvent.Msg(aMessage)
}

func (zero zeroLogger) Debug(aMessage string, aValues ...interface{}) {
	lEvent := zero.log.Debug()
	addFields(lEvent, aValues...)
	lEvent.Msg(aMessage)
}

func (zero zeroLogger) Warn(aMessage string, aValues ...interface{}) {
	lEvent := zero.log.Warn()
	addFields(lEvent, aValues...)
	lEvent.Msg(aMessage)
}

func (zero zeroLogger) Error(aMessage string, aValues ...interface{}) {
	lEvent := zero.log.Error()
	addFields(lEvent, aValues...)
	lEvent.Msg(aMessage)
}

func addFields(e *zerolog.Event, aValues ...interface{}) {
	for i, value := range aValues {
		e = e.Interface("v"+strconv.Itoa(i), value)
	}
}
