package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestJsonResponseSuccess(t *testing.T) {
	type TestData struct {
		Key   string //`json:"key"`
		Value string //`json:"value"`
	}

	lTestData := TestData{Key: "test key", Value: "test value"}
	lExpectedBody, _ := json.Marshal(lTestData)

	lLog := new(strings.Builder)
	lLogger := core.NewDebugZeroLogger(lLog)
	lResponseWriter := httptest.NewRecorder()

	internalSetJsonResponse(lResponseWriter, lLogger, http.StatusOK, lTestData, false)
	checkHttpResponse(t, lResponseWriter.Result(), http.StatusOK, HTTP_CONTENT_TYPE_JSON_UTF8, lExpectedBody)

	assert.Equal(t, lLog.String(), "", "Log")
}

func TestJsonResponseMarshalError(t *testing.T) {
	const EXPECTED_LOG = "ERR Cannot convert complex64 to json\n"
	const EXPECTED_BODY = `{"code":500,"message":"Cannot convert complex64 to json"}`

	lLog := new(strings.Builder)
	lLogger := core.NewDebugZeroLogger(lLog)
	lResponseWriter := httptest.NewRecorder()

	// complex64 type is not supported by json.Marshal
	internalSetJsonResponse(lResponseWriter, lLogger, http.StatusOK, complex64(1), false)
	checkHttpResponse(t, lResponseWriter.Result(), http.StatusInternalServerError, HTTP_CONTENT_TYPE_JSON_UTF8, []byte(EXPECTED_BODY))

	assert.Equal(t, lLog.String(), EXPECTED_LOG, "Log")
}

func TestJsonResponseApiErrorMarshalError(t *testing.T) {
	const EXPECTED_LOG = "ERR Cannot convert complex64 to json\n"
	const EXPECTED_BODY = "\n"

	lLog := new(strings.Builder)
	lLogger := core.NewDebugZeroLogger(lLog)
	lResponseWriter := httptest.NewRecorder()

	// complex64 type is not supported by json.Marshal
	internalSetJsonResponse(lResponseWriter, lLogger, http.StatusOK, complex64(1), true)
	checkHttpResponse(t, lResponseWriter.Result(), http.StatusInternalServerError, HTTP_CONTENT_TYPE_TEXT_UTF8, []byte(EXPECTED_BODY))

	assert.Equal(t, lLog.String(), EXPECTED_LOG, "Log")
}

func TestSetApiErrorResponse(t *testing.T) {
	const EXPECTED_LOG = ""
	const ERROR_TEMPLATE = "test error %d"
	const EXPECTED_BODY = `{"code":400,"message":"test error 125"}`

	lLog := new(strings.Builder)
	lLogger := core.NewDebugZeroLogger(lLog)
	lResponseWriter := httptest.NewRecorder()

	SetApiErrorResponsef(lResponseWriter, lLogger, http.StatusBadRequest, ERROR_TEMPLATE, 125)
	checkHttpResponse(t, lResponseWriter.Result(), http.StatusBadRequest, HTTP_CONTENT_TYPE_JSON_UTF8, []byte(EXPECTED_BODY))

	assert.Equal(t, lLog.String(), EXPECTED_LOG, "Log")
}
