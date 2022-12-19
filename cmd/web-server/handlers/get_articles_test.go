package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/indikator/aggregator_orange_cake/pkg/core"
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

	lResponseWriter := httptest.NewRecorder()
	lBuilder.ServeMux().ServeHTTP(lResponseWriter, lRequest)
	assert.Equal(t, http.StatusNotFound, lResponseWriter.Code, "Response status")
}

func TestGetArticlesRegistered(t *testing.T) {
	const EXPECTED_BODY = "[]"

	lStorage := &core.MockDBStorage{}
	lStorage.ReadArticlesByDateRangeFunc = func(aMin, aMax time.Time) ([]core.ArticleDB, error) {
		return make([]core.ArticleDB, 0), nil
	}

	lBuilder := initCustomTestBuilder(lStorage)
	const TEST_URL = "/test"

	lFailed := false
	RegisterGetArticles(TEST_URL, lBuilder, &lFailed)

	// create test request
	lRequest, lErr := http.NewRequest("GET", TEST_URL, nil)
	if lErr != nil {
		t.Errorf("Cannot create http request. %s", lErr.Error())
		return
	}

	lResponseWriter := httptest.NewRecorder()
	lBuilder.ServeMux().ServeHTTP(lResponseWriter, lRequest)
	checkHttpResponse(t, lResponseWriter.Result(), http.StatusOK, HTTP_CONTENT_TYPE_JSON_UTF8, []byte(EXPECTED_BODY))
}

func TestGetArticlesReadError(t *testing.T) {
	const EXPECTED_BODY = `{"code":500,"message":"Cannot read articles. not implemented (MockDBStorage.ReadArticlesByDateRange)"}`

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

	lResponseWriter := httptest.NewRecorder()
	lBuilder.ServeMux().ServeHTTP(lResponseWriter, lRequest)
	checkHttpResponse(t, lResponseWriter.Result(), http.StatusInternalServerError, HTTP_CONTENT_TYPE_JSON_UTF8, []byte(EXPECTED_BODY))
}

func TestGetArticlesInvalidFromDate(t *testing.T) {
	const EXPECTED_BODY = `{"code":400,"message":"Invalid From parameter. parsing time \"2022-23-24\": month out of range"}`

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

	lResponseWriter := httptest.NewRecorder()
	lBuilder.ServeMux().ServeHTTP(lResponseWriter, lRequest)
	checkHttpResponse(t, lResponseWriter.Result(), http.StatusBadRequest, HTTP_CONTENT_TYPE_JSON_UTF8, []byte(EXPECTED_BODY))
}

func TestGetArticlesInvalidToDate(t *testing.T) {
	const EXPECTED_BODY = `{"code":400,"message":"Invalid To parameter. parsing time \"2022-23-24\": month out of range"}`

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

	lResponseWriter := httptest.NewRecorder()
	lBuilder.ServeMux().ServeHTTP(lResponseWriter, lRequest)
	checkHttpResponse(t, lResponseWriter.Result(), http.StatusBadRequest, HTTP_CONTENT_TYPE_JSON_UTF8, []byte(EXPECTED_BODY))
}
