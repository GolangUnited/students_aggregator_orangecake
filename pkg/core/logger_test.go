package core

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWritingToMultipleWriters(t *testing.T) {
	buffer1 := new(bytes.Buffer)
	buffer2 := new(bytes.Buffer)
	logger := NewZeroLogger(buffer1, buffer2)

	logger.Info("Some important info...")

	expectedMsg := fmt.Sprintf("\x1b[32mINF\x1b[0m \x1b[90m%s\x1b[0m Some important info...\n", time.Now().Format(time.RFC822))
	// INF 11 Dec 22 13:59 +03 Some important info...
	assert.Equal(t, expectedMsg, buffer1.String(), "Wrong logged message.")
	assert.Equal(t, buffer1.String(), buffer2.String(), "Log results in writer1 and writer2 are not the same.")
}

func TestLoggingWithFormatValues(t *testing.T) {
	buffer := new(bytes.Buffer)
	logger := NewZeroLogger(buffer)

	logger.Info("Article %d: %s.", 5, "some article name")

	expectedMsg := fmt.Sprintf(
		"\x1b[32mINF\x1b[0m \x1b[90m%s\x1b[0m Article 5: some article name.\n",
		time.Now().Format(time.RFC822))
	// INF 11 Dec 22 15:19 +03 Article 5: some article name.
	assert.Equal(t, expectedMsg, buffer.String(), "Values given to logger are shown incorrect.")
}

func TestLoggingTraceLevel(t *testing.T) {
	buffer := new(bytes.Buffer)
	logger := NewZeroLogger(buffer)

	logger.Trace("Trace...")

	expectedMsg := fmt.Sprintf("\x1b[35mTRC\x1b[0m \x1b[90m%s\x1b[0m Trace...\n", time.Now().Format(time.RFC822))
	// TRC 11 Dec 22 13:59 +03 Trace...
	assert.Equal(t, expectedMsg, buffer.String(), "Wrong logged Trace string.")
}

func TestLoggingInfoLevel(t *testing.T) {
	buffer := new(bytes.Buffer)
	logger := NewZeroLogger(buffer)

	logger.Info("Some important info...")

	expectedMsg := fmt.Sprintf("\x1b[32mINF\x1b[0m \x1b[90m%s\x1b[0m Some important info...\n", time.Now().Format(time.RFC822))
	// INF 11 Dec 22 13:59 +03 Some important info...
	assert.Equal(t, expectedMsg, buffer.String(), "Wrong logged Info string.")
}

func TestLoggingDebugLevel(t *testing.T) {
	buffer := new(bytes.Buffer)
	logger := NewZeroLogger(buffer)

	logger.Debug("Debug...")

	expectedMsg := fmt.Sprintf("\x1b[33mDBG\x1b[0m \x1b[90m%s\x1b[0m Debug...\n", time.Now().Format(time.RFC822))
	// DBG 11 Dec 22 13:59 +03 Debug...
	assert.Equal(t, expectedMsg, buffer.String(), "Wrong logged Debug string.")
}

func TestLoggingWarnLevel(t *testing.T) {
	buffer := new(bytes.Buffer)
	logger := NewZeroLogger(buffer)

	logger.Warn("Warn...")

	expectedMsg := fmt.Sprintf("\x1b[31mWRN\x1b[0m \x1b[90m%s\x1b[0m Warn...\n", time.Now().Format(time.RFC822))
	// WRN 11 Dec 22 13:59 +03 Warn...
	assert.Equal(t, expectedMsg, buffer.String(), "Wrong logged Warning string.")
}

func TestLoggingErrorLevel(t *testing.T) {
	buffer := new(bytes.Buffer)
	logger := NewZeroLogger(buffer)

	logger.Error("Error...")

	expectedMsg := fmt.Sprintf("\x1b[1m\x1b[31mERR\x1b[0m\x1b[0m \x1b[90m%s\x1b[0m Error...\n", time.Now().Format(time.RFC822))
	// ERR 11 Dec 22 13:59 +03 Error...
	assert.Equal(t, expectedMsg, buffer.String(), "Wrong logged Error string.")
}

func TestDebugLogger(t *testing.T) {
	buffer := new(bytes.Buffer)
	logger := NewDebugZeroLogger(buffer)

	logger.Error("Error...")
	assert.Equal(t, "ERR Error...\n", buffer.String(), "Wrong logged Error string.")
}
