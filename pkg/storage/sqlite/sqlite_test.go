package sqlite

import (
	"fmt"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestNewSqliteConnection(t *testing.T) {
	logger := core.NewZeroLogger(io.Discard)

	testCases := map[string]struct {
		connectionString string
		expectedError    error
	}{
		"valid connection string": {
			connectionString: "file::memory:",
			expectedError:    nil,
		},
		"invalid connection string": {
			connectionString: "not::valid",
			expectedError:    fmt.Errorf("can't open database"),
		},
		"connection string read only": {
			connectionString: "file::memory:?mode=ro",
			expectedError:    fmt.Errorf("attempt to write a readonly database"),
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			_, gotError := NewSqliteConnection(testCase.connectionString, logger)

			if testCase.expectedError != gotError {
				assert.ErrorAs(t, testCase.expectedError, &gotError)
			}
		})
	}
}
