package core

import (
	"github.com/indikator/aggregator_orange_cake/pkg/logger"
	"io"
)

type Logger interface {
	Trace(msg string, keyvals ...Field)
	Info(msg string, keyvals ...Field)
	Debug(msg string, keyvals ...Field)
	Warning(msg string, keyvals ...Field)
	Error(msg string, err error, keyvals ...Field)
}
type Field = logger.Field

var _ Logger = (*logger.GLogger)(nil)

func NewLogger(writers ...io.Writer) Logger {
	return logger.NewGLogger(writers...)
}
