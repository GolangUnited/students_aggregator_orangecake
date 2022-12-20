package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenApiUINotRegistered(t *testing.T) {
	lBuilder := initTestBuilder()
	const TEST_URL = "/test"

	lFailed := true
	RegisterOpenApiUI(TEST_URL, lBuilder, &lFailed)

	// create test request
	lRequest, lErr := http.NewRequest("GET", TEST_URL, nil)
	if lErr != nil {
		t.Errorf("Cannot create http request. %s", lErr.Error())
		return
	}

	lResponse := httptest.NewRecorder()
	lBuilder.ServeMux().ServeHTTP(lResponse, lRequest)
	assert.Equal(t, http.StatusNotFound, lResponse.Code, "Response status")
}

func TestOpenApiUIRegistered(t *testing.T) {
	lBuilder := initTestBuilder()
	const TEST_URL = "/test"

	lFailed := false
	RegisterOpenApiUI(TEST_URL, lBuilder, &lFailed)

	// create test request
	lRequest, lErr := http.NewRequest("GET", TEST_URL, nil)
	if lErr != nil {
		t.Errorf("Cannot create http request. %s", lErr.Error())
		return
	}

	lResponse := httptest.NewRecorder()
	lBuilder.ServeMux().ServeHTTP(lResponse, lRequest)
	assert.Equal(t, http.StatusOK, lResponse.Code, "Response status")
}
