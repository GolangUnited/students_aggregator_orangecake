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

func NewTestServer(data string) *httptest.Server {
	lServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, data)
	}))

	return lServer
}

func TestScrapUrl(t *testing.T) {

	lServer := NewTestServer(test.HashnodeOkTestData)
	defer lServer.Close()

	lTestCases := []core.Article{
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
			Link:        "Link 1",
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
	lHS := NewHashnodeScraper(lServer.URL, log)

	lHS.ScrapUrl()

	for idx, val := range lHS.Articles {
		assert.Equal(t, lTestCases[idx], val)
	}
}
