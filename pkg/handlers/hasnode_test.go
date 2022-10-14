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

func NewTestServer() *httptest.Server {
	lServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, test.HashnodeTestData)
	}))

	return lServer
}

func TestScrapUrl(t *testing.T) {

	lServer := NewTestServer()
	defer lServer.Close()

	lTestCases := []core.Article{
		{
			Title:       "Working with strings effectively with Golang",
			Author:      "Ayomide Ayanwola",
			Link:        "https://thecodeway.hashnode.dev/working-with-strings-effectively-with-golang",
			PublishDate: time.Date(2022, time.October, 9, 0, 0, 0, 0, time.UTC),
			Description: "What are Strings? The way strings is stored in memory makes them immutable, making it difficult to perform simple operations like changing the value of an index in a string. For example you can't perform an index assignment operation on …",
		},
		{
			Title:       "How to create a static GO-lang server",
			Author:      "Ayush Bajpai",
			Link:        "https://gitayush.hashnode.dev/how-to-create-a-static-go-lang-server",
			PublishDate: time.Date(2022, time.September, 11, 0, 0, 0, 0, time.UTC),
			Description: "Open Powershell in your windows and then direct to GO directory using the command cd go Once inside GO directory, go to the src folder using the command cd src. Then you need to make a folder/directory…",
		},
	}

	lHS := hashnodeScraper{URL: lServer.URL}
	lHS.ScrapUrl()

	for idx, val := range lHS.Articles {
		assert.Equal(t, lTestCases[idx], val)
	}
}
