package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetArticlesNotRegistered(t *testing.T) {
	lBuilder := initTestBuilder()
	const TEST_URL = "/test"

	lFailed := true
	RegisterGetArticles(TEST_URL, lBuilder, &lFailed)

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

func TestGetArticlesRegistered(t *testing.T) {
	lBuilder := initTestBuilder()
	const TEST_URL = "/test"

	lFailed := false
	RegisterGetArticles(TEST_URL, lBuilder, &lFailed)

	// create test request
	lRequest, lErr := http.NewRequest("GET", TEST_URL, nil)
	if lErr != nil {
		t.Errorf("Cannot create http request. %s", lErr.Error())
		return
	}

	lResponse := httptest.NewRecorder()
	lBuilder.ServeMux().ServeHTTP(lResponse, lRequest)
	assert.Equal(t, http.StatusOK, lResponse.Code, "Response status")
	//TODO: check response data
}

func TestGetArticlesInvalidFromDate(t *testing.T) {
	lBuilder := initTestBuilder()
	const TEST_URL = "/test"

	lFailed := false
	RegisterGetArticles(TEST_URL, lBuilder, &lFailed)

	// create test request
	lRequest, lErr := http.NewRequest("GET", TEST_URL+"?from=2022-23-24", nil)
	if lErr != nil {
		t.Errorf("Cannot create http request. %s", lErr.Error())
		return
	}

	lResponse := httptest.NewRecorder()
	lBuilder.ServeMux().ServeHTTP(lResponse, lRequest)
	assert.Equal(t, http.StatusBadRequest, lResponse.Code, "Response status")
}

func TestGetArticlesInvalidToDate(t *testing.T) {
	lBuilder := initTestBuilder()
	const TEST_URL = "/test"

	lFailed := false
	RegisterGetArticles(TEST_URL, lBuilder, &lFailed)

	// create test request
	lRequest, lErr := http.NewRequest("GET", TEST_URL+"?to=2022-23-24", nil)
	if lErr != nil {
		t.Errorf("Cannot create http request. %s", lErr.Error())
		return
	}

	lResponse := httptest.NewRecorder()
	lBuilder.ServeMux().ServeHTTP(lResponse, lRequest)
	assert.Equal(t, http.StatusBadRequest, lResponse.Code, "Response status")
}
