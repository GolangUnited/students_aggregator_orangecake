package handlers

import (
	"strings"
	"testing"

	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"github.com/stretchr/testify/assert"
)

func newMockLogger() core.Logger {
	lLog := new(strings.Builder)
	return core.NewZeroLogger(lLog)
}

func TestNewWebServerWithoutLog(t *testing.T) {
	lServer, lErr := NewWebServer(nil, nil, nil)

	assert.Nil(t, lServer, "server")
	assert.ErrorIs(t, lErr, core.ErrLoggerNotAssigned)
}

func TestNewWebServerWithoutDBReader(t *testing.T) {
	lServer, lErr := NewWebServer(newMockLogger(), nil, nil)

	assert.Nil(t, lServer, "server")
	assert.ErrorIs(t, lErr, core.ErrDBReaderNotAssigned)
}

func TestNewWebServerWithoutDBWriter(t *testing.T) {
	lSrorage := core.MockDBStorage{}
	lServer, lErr := NewWebServer(newMockLogger(), lSrorage, nil)

	assert.Nil(t, lServer, "server")
	assert.ErrorIs(t, lErr, core.ErrDBWriterNotAssigned)
}

func TestNewWebServer(t *testing.T) {
	lSrorage := core.MockDBStorage{}
	lServer, lErr := NewWebServer(newMockLogger(), lSrorage, lSrorage)

	assert.Nil(t, lErr, "error")
	assert.NotNil(t, lServer, "server")
	assert.NotNil(t, lServer.Log(), "logger")
	assert.NotNil(t, lServer.DBReader(), "db reader")
	assert.NotNil(t, lServer.DBWriter(), "db writer")
}
