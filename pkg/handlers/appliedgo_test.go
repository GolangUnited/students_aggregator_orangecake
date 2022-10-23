package handlers

import (
	"errors"
	"fmt"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

const (
	TestDataFolder  = "./test_data"
	TestDataAddress = "127.0.0.1:8080"
	TestDataURL     = "http://127.0.0.1:8080/"
)

func TestMain(m *testing.M) {
	// Create mux handler
	mux := http.NewServeMux()
	fServer := http.FileServer(http.Dir(TestDataFolder))
	mux.Handle("/", fServer)

	// Create the server with certain address
	lNewListener, _ := net.Listen("tcp", TestDataAddress)
	testServer := *httptest.NewUnstartedServer(mux)
	testServer.Listener = lNewListener

	// Start the server and run the tests
	testServer.Start()
	defer testServer.Close()
	os.Exit(m.Run())
}

func TestAppliedGoGetArticlesList(t *testing.T) {
	var lReceivedLinksList []string
	var lErr error
	lExpectedLinksList := []string{
		"https://appliedgo.net/rich/",
		"https://appliedgo.net/generictree/",
		"https://appliedgo.net/instantgo/",
		"https://appliedgo.net/mantil/",
		"https://appliedgo.net/auxin/",
	}

	lReceivedLinksList, lErr = ParseAppliedGoMain(TestDataURL + "AppliedGoMain.htm")
	assert.ElementsMatch(t, lExpectedLinksList, lReceivedLinksList)
	assert.Equal(t, nil, lErr)
}

func TestAppliedGoMainIncorrectUrlProtocol(t *testing.T) {
	var badUrl = ""
	var lReceivedLinksList []string
	var lErr error
	var lExpectedErr = fmt.Errorf("unsupported protocol scheme %q", badUrl)

	lReceivedLinksList, lErr = ParseAppliedGoMain(badUrl)
	assert.Equal(t, []string(nil), lReceivedLinksList)
	assert.Equal(t, lExpectedErr, errors.Unwrap(lErr))
}

func TestAppliedGoScrapeSingleArticle(t *testing.T) {
	var lReceivedData core.Article
	var lErr []error
	lReceivedData, lErr = ParseAppliedGoArticle(TestDataURL + "AppliedGoArticle.htm")
	lExpectedData := core.Article{
		Title:       "How I used Go to make my radio auto-switch to AUX-IN when a Raspi plays music - Applied Go",
		Author:      "",
		Link:        TestDataURL + "AppliedGoArticle.htm",
		PublishDate: time.Date(2022, time.August, 20, 0, 0, 0, 0, time.UTC),
		Description: "How Go code detects music output on a Raspberry and switches a 3sixty radio to AUX-IN via Frontier Silicon API",
	}
	assert.Equal(t, lExpectedData, lReceivedData)
	assert.Equal(t, []error(nil), lErr)
}

func TestAppliedGoScrapeEmptyFields(t *testing.T) {
	var lReceivedData core.Article
	var lErr []error
	lReceivedData, lErr = ParseAppliedGoArticle(TestDataURL + "AppliedGoArticleEmptyFields.htm")
	lExpectedData := core.Article{
		Title:       "",
		Author:      "",
		Link:        TestDataURL + "AppliedGoArticleEmptyFields.htm",
		PublishDate: time.Date(1970, time.January, 20, 0, 0, 0, 0, time.UTC),
		Description: "",
	}
	lExpectedErr := []error{
		errors.New("error: title field is empty"),
		errors.New("warning: description field is empty"),
		errors.New("warning: date field is empty"),
	}
	assert.Equal(t, lExpectedData, lReceivedData)
	assert.ElementsMatch(t, lExpectedErr, lErr)
}

func TestAppliedGoScrapeInvalidDate(t *testing.T) {
	var lErr []error
	_, lErr = ParseAppliedGoArticle(TestDataURL + "AppliedGoArticleInvalidDate.htm")

	assert.Equal(t, 1, len(lErr))
	assert.Equal(t, fmt.Errorf("invalid date format: invalid Date format"), lErr[0])
}

func TestAppliedGoArticleIncorrectUrlProtocol(t *testing.T) {
	var badUrl = ""
	var lExpectedArticle = core.Article{
		Title:  "",
		Author: "", Link: "",
		PublishDate: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		Description: ""}
	var lReceivedArticle core.Article
	var lExpectedErr = fmt.Errorf("unsupported protocol scheme %q", badUrl)
	var lErr []error

	lReceivedArticle, lErr = ParseAppliedGoArticle(badUrl)
	assert.Equal(t, lExpectedArticle, lReceivedArticle)
	assert.Equal(t, 1, len(lErr))
	assert.Equal(t, lExpectedErr, errors.Unwrap(lErr[0]))
}

func TestAppliedGoGetCombinedResults(t *testing.T) {
	var lReceivedLinksList []core.Article
	var lReceivedErrorsList []error
	var lErr error
	lExpectedLinksList := []core.Article{
		{
			Title:       "How I used Go to make my radio auto-switch to AUX-IN when a Raspi plays music - Applied Go",
			Author:      "",
			Link:        TestDataURL + "AppliedGoArticleCombined.htm",
			PublishDate: time.Date(2022, time.August, 20, 0, 0, 0, 0, time.UTC),
			Description: "How Go code detects music output on a Raspberry and switches a 3sixty radio to AUX-IN via Frontier Silicon API",
		},
	}

	lReceivedLinksList, lReceivedErrorsList, lErr = ParseAppliedGo(TestDataURL + "AppliedGoMainCombined.htm")
	assert.ElementsMatch(t, lExpectedLinksList, lReceivedLinksList)
	assert.ElementsMatch(t, []error(nil), lReceivedErrorsList)
	assert.Equal(t, nil, lErr)

}
