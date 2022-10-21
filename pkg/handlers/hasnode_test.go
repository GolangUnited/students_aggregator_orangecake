package handlers

import (
	"fmt"
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

type TestCase struct {
	TestName string
	TestData string
	Err      error
	URL      string
}

func NewTestServer(data string) *httptest.Server {
	lServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, data)
	}))

	return lServer
}

func TestOkScrapUrl(t *testing.T) {

	lServer := NewTestServer(test.OkTestData)
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
			PublishDate: time.Date(2022, time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC),
			Description: "",
		},
	}

	log := log.New(os.Stdout, "HS", log.Flags())
	lHS := NewHashnodeScraper(log)
	lHS.URL = lServer.URL

	lHS.ScrapUrl()

	assert.Equal(t, len(lOkTestCases), len(lHS.Articles))

	for idx, val := range lOkTestCases {
		assert.Equal(t, val, lHS.Articles[idx])

	}
}

func TestErrorsScrapUrl(t *testing.T) {

	TestCases := []TestCase{
		{
			TestName: "No Articles Err",
			TestData: test.NoArticlesTestData,
			Err:      fmt.Errorf("no correct articles"),
		},
		{
			TestName: "No Fields Err",
			TestData: test.NoFieldsTestData,
			Err:      fmt.Errorf("unable to find required field"),
		},
		{
			TestName: "URL Visit Err",
			TestData: "",
			Err:      fmt.Errorf("visit error"),
			URL:      "http://127.0.0.1",
		},
	}

	for _, tCase := range TestCases {

		lServer := NewTestServer(tCase.TestData)
		defer lServer.Close()

		log := log.New(os.Stdout, "HS", log.Flags())
		lHS := NewHashnodeScraper(log)
		lHS.URL = lServer.URL

		if tCase.URL != "" {
			lHS.URL = tCase.URL
		}

		err := lHS.ScrapUrl()

		assert.Contains(t, err.Error(), tCase.Err.Error())

	}

}
