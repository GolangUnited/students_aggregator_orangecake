package handlers

import (
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

const (
	TestDataFolder = "./test_data"
)

var testServer httptest.Server
var testDataURL string

func TestMain(m *testing.M) {
	mux := http.NewServeMux()
	fServer := http.FileServer(http.Dir(TestDataFolder))
	mux.Handle("/", fServer)
	testServer = *httptest.NewServer(mux)

	testDataURL = testServer.URL + "/"
	os.Exit(m.Run())
}

func TestGetArticlesListAppliedGo(t *testing.T) {
	var lReceivedLinksList []string
	var lErr error
	lExpectedLinksList := []string{
		"https://appliedgo.net/rich/",
		"https://appliedgo.net/generictree/",
		"https://appliedgo.net/instantgo/",
		"https://appliedgo.net/mantil/",
		"https://appliedgo.net/auxin/",
	}

	lReceivedLinksList, lErr = ParseAppliedGoMain(testDataURL + "AppliedGoMain.htm")
	assert.ElementsMatch(t, lExpectedLinksList, lReceivedLinksList)
	assert.Equal(t, nil, lErr)
}

func TestArticleScraping(t *testing.T) {
	var lReceivedData core.Article
	var lErr error
	lReceivedData, lErr = ParseAppliedGoArticle(testDataURL + "AppliedGoArticle.htm")
	lExpectedData := core.Article{
		Title:       "How I used Go to make my radio auto-switch to AUX-IN when a Raspi plays music - Applied Go",
		Author:      "",
		Link:        "https://appliedgo.net/auxin/",
		PublishDate: time.Date(2022, time.August, 20, 0, 0, 0, 0, time.UTC),
		Description: "How Go code detects music output on a Raspberry and switches a 3sixty radio to AUX-IN via Frontier Silicon API",
	}
	assert.Equal(t, lExpectedData, lReceivedData)
	assert.Equal(t, nil, lErr)
}
