package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"github.com/indikator/aggregator_orange_cake/test"
	"github.com/stretchr/testify/assert"
)

type hashnodeTestCase struct {
	testName string
	testData string
	err      error
	warnings []core.Warning
	url      string
}

func newHTTPTestServer(data string) *httptest.Server {
	lServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, data)
	}))

	return lServer
}

func TestHashnodeOkScrapUrl(t *testing.T) {

	lServer := newHTTPTestServer(test.OK_TEST_DATA)
	defer lServer.Close()

	lExpectedArticles := []core.Article{
		{
			Title:       "Title 1",
			Author:      "Author 1",
			Link:        "Link 1",
			PublishDate: time.Date(2022, time.October, 9, 0, 0, 0, 0, time.UTC),
			Description: "Text 1…",
		},
		{
			Title:       "Title 2",
			Author:      "",
			Link:        "Link 2",
			PublishDate: time.Date(2022, time.September, 8, 0, 0, 0, 0, time.UTC),
			Description: "Text 2…",
		},
		{
			Title:       "Title 3",
			Author:      "Author 3",
			Link:        "Link 3",
			PublishDate: time.Date(2022, time.September, 8, 0, 0, 0, 0, time.UTC),
			Description: "",
		},
	}

	lExpectedWarnings := []core.Warning{}
	var lExpectedError error = nil

	log := core.NewZeroLogger(io.Discard)
	lHS := NewHashnodeScraper(lServer.URL, log)

	lGotArticles, lGotWarnings, lGotError := lHS.ParseArticles()

	assert.Equal(t, lExpectedError, lGotError)
	assert.Equal(t, lExpectedWarnings, lGotWarnings)
	assert.Equal(t, len(lExpectedArticles), len(lGotArticles))

	for idx, val := range lExpectedArticles {
		assert.Equal(t, val, lGotArticles[idx])
	}
}

func TestHashnodeErrorsScrapUrl(t *testing.T) {

	lTestCases := []hashnodeTestCase{
		{
			testName: "URL Visit Err",
			testData: "",
			err:      core.ErrUrlVisit,
			warnings: make([]core.Warning, 0),
			url:      "http://",
		},
		{
			testName: "No Articles Err",
			testData: test.NO_CORRECT_ARTICLES_TEST_DATA,
			err:      core.ErrNoArticles,
			warnings: []core.Warning{
				"error: field is empty, field: Link",
				"error: attribute doesn't exists, field: Link",
				"error: node not found, field: Title",
				"error: field is empty, field: Title",
			},
		},
	}

	for _, tCase := range lTestCases {

		lServer := newHTTPTestServer(tCase.testData)

		var lHS *HashnodeScraper

		log := core.NewZeroLogger(io.Discard)

		if tCase.url != "" {
			lHS = NewHashnodeScraper(tCase.url, log)
		} else {
			lHS = NewHashnodeScraper(lServer.URL, log)
		}

		_, lGotWarnings, lGotError := lHS.ParseArticles()

		assert.Contains(t, lGotError.Error(), tCase.err.Error())
		assert.Equal(t, tCase.warnings, lGotWarnings)

		lServer.Close()
	}
}

func TestHashnodeWithWarningsScrapUrl(t *testing.T) {

	lServer := newHTTPTestServer(test.ARTICLES_WITH_WARNINGS)
	defer lServer.Close()

	lExpectedArticles := []core.Article{
		{
			Title:       "Title 1",
			Author:      "Author 1",
			Link:        "Link 1",
			PublishDate: time.Date(2022, time.October, 9, 0, 0, 0, 0, time.UTC),
			Description: "",
		},
		{
			Title:       "Title 2",
			Author:      "",
			Link:        "Link 2",
			PublishDate: time.Date(2022, time.October, 9, 0, 0, 0, 0, time.UTC),
			Description: "Text 2…",
		},
		{
			Title:       "Title 3",
			Author:      "Author 3",
			Link:        "Link 3",
			PublishDate: core.NormalizeDate(time.Now()),
			Description: "Text 3…",
		},
	}

	lExpectedWarnings := []core.Warning{
		"For article Title 1, Link 1 field Description not found",
		"For article Title 2, Link 2 field Author not found",
		"For article Title 3, Link 3 field Data not found",
		"For article Title 3, Link 3 DateErr: empty Date",
	}
	var lExpectedError error = nil

	log := core.NewZeroLogger(io.Discard)
	lHS := NewHashnodeScraper(lServer.URL, log)

	lGotArticles, lGotWarnings, lGotError := lHS.ParseArticles()

	assert.Equal(t, lExpectedError, lGotError)
	assert.Equal(t, lExpectedWarnings, lGotWarnings)
	assert.Equal(t, len(lExpectedArticles), len(lGotArticles))

	for idx, val := range lExpectedArticles {
		assert.Equal(t, val, lGotArticles[idx])
	}
}
