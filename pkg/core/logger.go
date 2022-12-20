package core

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"time"
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

func (zero zeroLogger) Trace(aMessage string, aValues ...interface{}) {
	lFormattedMessage := fmt.Sprintf(aMessage, aValues...)
	zero.log.Trace().Msg(lFormattedMessage)
}

func (zero zeroLogger) Info(aMessage string, aValues ...interface{}) {
	lFormattedMessage := fmt.Sprintf(aMessage, aValues...)
	zero.log.Info().Msg(lFormattedMessage)
}

func (zero zeroLogger) Debug(aMessage string, aValues ...interface{}) {
	lFormattedMessage := fmt.Sprintf(aMessage, aValues...)
	zero.log.Debug().Msg(lFormattedMessage)
}

func (zero zeroLogger) Warn(aMessage string, aValues ...interface{}) {
	lFormattedMessage := fmt.Sprintf(aMessage, aValues...)
	zero.log.Warn().Msg(lFormattedMessage)
}

func (zero zeroLogger) Error(aMessage string, aValues ...interface{}) {
	lFormattedMessage := fmt.Sprintf(aMessage, aValues...)
	zero.log.Error().Msg(lFormattedMessage)
}
