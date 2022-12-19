package handlers

import (
	"io"
	"net/http"
	"testing"

	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"github.com/stretchr/testify/assert"
)

func checkHttpResponse(t *testing.T, aResponse *http.Response, aStatusCode int, aContentType string, aBody []byte) {
	lContentType := aResponse.Header.Get(HTTP_HEADER_CONTENT_TYPE)
	lBody, _ := io.ReadAll(aResponse.Body)

	assert.Equal(t, aStatusCode, aResponse.StatusCode, "HTTP Status")
	assert.Equal(t, aContentType, lContentType, HTTP_HEADER_CONTENT_TYPE)
	assert.Equal(t, string(aBody), string(lBody), "Body")
}

func getErrMsg(aError error) string {
	if aError != nil {
		return aError.Error()
	}
	return ""
}

func initCustomTestBuilder(aStorage *core.MockDBStorage) WebServerBuilder {

	lStorage := aStorage
	if lStorage == nil {
		lStorage = &core.MockDBStorage{}
	}

	lWebServer, _ := NewWebServer(newMockLogger(), lStorage, lStorage)

	return NewWebServerBuilder(lWebServer)
}

func initTestBuilder() WebServerBuilder {
	return initCustomTestBuilder(nil)
}
