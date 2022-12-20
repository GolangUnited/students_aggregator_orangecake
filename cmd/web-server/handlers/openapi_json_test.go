package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenApiJsonNotRegistered(t *testing.T) {
	lBuilder := initTestBuilder()
	const TEST_URL = "/test"

	lFailed := true
	RegisterOpenApiJson(TEST_URL, lBuilder, &lFailed)

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

func TestOpenApiJsonRegistered(t *testing.T) {
	lBuilder := initTestBuilder()
	const TEST_URL = "/test"

	lFailed := false
	RegisterOpenApiJson(TEST_URL, lBuilder, &lFailed)

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

func TestMarshalJsonError(t *testing.T) {
	lBuilder := initTestBuilder()
	// add a complex64 to generate an errror in aOpenApi.Spec.MarshalJSON()
	lBuilder.OpenApi().Spec.WithMapOfAnythingItem("bad", complex64(1))
	const TEST_URL = "/test"

	lFailed := false
	RegisterOpenApiJson(TEST_URL, lBuilder, &lFailed)
	assert.True(t, lFailed)
}
