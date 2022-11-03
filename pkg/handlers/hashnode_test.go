package handlers

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
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
	url      string
}

func newHTTPTestServer(data string) *httptest.Server {
	lServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, data)
	}))

	return lServer
}

func TestOkScrapUrl(t *testing.T) {

	lServer := newHTTPTestServer(test.OK_TEST_DATA)
	defer lServer.Close()

	lOkTestCases := []core.Article{
		{
			Title:       "Title 1",
			Author:      "Author 1",
			Link:        "Link 1",
			PublishDate: time.Date(2022, time.October, 9, 0, 0, 0, 0, time.UTC),
			Description: "Text 1…",
		},
		{
			Title:       "Title 2",
			Author:      "Author 2",
			Link:        "Link 2",
			PublishDate: time.Date(2022, time.September, 8, 0, 0, 0, 0, time.UTC),
			Description: "Text 2…",
		},
		{
			Title:       "Title 3",
			Author:      "Author 3",
			Link:        "Link 3",
			PublishDate: core.NormalizeDate(time.Now()),
			Description: "",
		},
	}

	log := log.New(os.Stdout, "HS", log.Flags())
	lHS := NewHashnodeScraper(log, lServer.URL)

	lHS.ScrapUrl()

	assert.Equal(t, len(lOkTestCases), len(lHS.Articles))

	for idx, val := range lOkTestCases {
		assert.Equal(t, val, lHS.Articles[idx])

	}
}

func TestErrorsScrapUrl(t *testing.T) {

	lTestCases := []hashnodeTestCase{
		{
			testName: "No Articles Err",
			testData: test.NO_ARTICLES_TEST_DATA,
			err:      core.ErrArticlesNotFound,
		},
		{
			testName: "URL Visit Err",
			testData: "",
			err:      core.ErrUrlVisit,
			url:      "http://127.0.0.1",
		},
		{
			testName: "No Correct Articles Err",
			testData: test.NO_CORRECT_ARTICLES_TEST_DATA,
			err:      core.ErrNoArticles,
		},
	}

	for _, tCase := range lTestCases {

		lServer := newHTTPTestServer(tCase.testData)
		defer lServer.Close()

		log := log.New(os.Stdout, "HS", log.Flags())
		lHS := NewHashnodeScraper(log, lServer.URL)

		if tCase.url != "" {
			lHS.URL = tCase.url
		}

		err := lHS.ScrapUrl()

		assert.Contains(t, err.Error(), tCase.err.Error())

	}

}
