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
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

const (
	TestDataAddress = "127.0.0.1:8080"
	TestDataURL     = "http://127.0.0.1:8080/"
)

func buildTestDataPathAppliedGo() (string, error) {
	_, lFileName, _, lOK := runtime.Caller(0)
	if !lOK {
		return "", errors.New("cannot get TestDataPath")
	}

	lFullFileName, lErr := filepath.Abs(lFileName)
	if lErr != nil {
		return "", lErr
	}
	lFullFileName = filepath.Join(filepath.Dir(lFullFileName), "testdata")

	return lFullFileName, nil
}

func TestMain(m *testing.M) {
	// Create mux handler
	mux := http.NewServeMux()
	testDataFolder, _ := buildTestDataPathAppliedGo()
	fServer := http.FileServer(http.Dir(testDataFolder))
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
	lExpectedLinks := []string{
		"https://appliedgo.net/rich/",
		"https://appliedgo.net/generictree/",
		"https://appliedgo.net/instantgo/",
		"https://appliedgo.net/mantil/",
		"https://appliedgo.net/auxin/",
	}
	lReceivedLinks := newAppliedGoMainParser()
	lReceivedErr := lReceivedLinks.ParseAppliedGoMain(TestDataURL + "AppliedGoMain.htm")
	assert.ElementsMatch(t, lExpectedLinks, lReceivedLinks.Links)
	assert.Equal(t, nil, lReceivedErr)
}

func TestAppliedGoMainIncorrectUrlProtocol(t *testing.T) {
	var badUrl = ""
	var lExpectedErr = fmt.Errorf("unsupported protocol scheme %q", badUrl)

	lReceivedLinks := newAppliedGoMainParser()
	lReceivedErr := lReceivedLinks.ParseAppliedGoMain(badUrl)
	assert.Equal(t, []string{}, lReceivedLinks.Links)
	assert.Equal(t, lExpectedErr, errors.Unwrap(lReceivedErr))
}

func TestAppliedGoScrapeSingleArticle(t *testing.T) {
	lReceivedArticle := newAppliedGoArticleParser()
	lReceivedErr := lReceivedArticle.ParseAppliedGoArticle(TestDataURL + "AppliedGoArticle.htm")
	lExpectedArticle := core.Article{
		Title:       "How I used Go to make my radio auto-switch to AUX-IN when a Raspi plays music - Applied Go",
		Author:      "",
		Link:        TestDataURL + "AppliedGoArticle.htm",
		PublishDate: time.Date(2022, time.August, 20, 0, 0, 0, 0, time.UTC),
		Description: "How Go code detects music output on a Raspberry and switches a 3sixty radio to AUX-IN via Frontier Silicon API",
	}
	assert.Equal(t, lExpectedArticle, lReceivedArticle.Article)
	assert.Equal(t, []string{}, lReceivedArticle.Warnings)
	assert.Equal(t, nil, lReceivedErr)
}

func TestAppliedGoScrapeEmptyFields(t *testing.T) {
	lReceivedArticle := newAppliedGoArticleParser()
	lReceivedErr := lReceivedArticle.ParseAppliedGoArticle(TestDataURL + "AppliedGoArticleEmptyFields.htm")
	lExpectedArticle := core.Article{
		Title:       "",
		Author:      "",
		Link:        TestDataURL + "AppliedGoArticleEmptyFields.htm",
		PublishDate: time.Date(1970, time.January, 20, 0, 0, 0, 0, time.UTC),
		Description: "",
	}
	lExpectedWarnings := []string{
		"error: title field is empty",
		"warning: description field is empty",
		"warning: date field is empty",
	}
	assert.Equal(t, lExpectedArticle, lReceivedArticle.Article)
	assert.ElementsMatch(t, lExpectedWarnings, lReceivedArticle.Warnings)
	assert.Equal(t, nil, lReceivedErr)
}

func TestAppliedGoScrapeInvalidDate(t *testing.T) {
	lReceivedArticle := newAppliedGoArticleParser()
	lReceivedErr := lReceivedArticle.ParseAppliedGoArticle(TestDataURL + "AppliedGoArticleInvalidDate.htm")

	assert.Equal(t, 1, len(lReceivedArticle.Warnings))
	assert.Equal(t, "invalid date format: invalid Date format", lReceivedArticle.Warnings[0])
	assert.Equal(t, nil, lReceivedErr)
}

func TestAppliedGoArticleIncorrectUrlProtocol(t *testing.T) {
	var badUrl = ""
	var lExpectedArticle = core.Article{
		Title:  "",
		Author: "", Link: "",
		PublishDate: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		Description: ""}
	var lExpectedErr = fmt.Errorf("unsupported protocol scheme %q", badUrl)

	lReceivedArticle := newAppliedGoArticleParser()
	lReceivedErr := lReceivedArticle.ParseAppliedGoArticle(badUrl)
	assert.Equal(t, lExpectedArticle, lReceivedArticle.Article)
	assert.Equal(t, lExpectedErr, errors.Unwrap(lReceivedErr))
}

func TestAppliedGoGetCombinedResults(t *testing.T) {
	var lReceivedLinksList []core.Article
	var lReceivedWarningsList []string
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

	lReceivedLinksList, lReceivedWarningsList, lErr = ParseAppliedGo(TestDataURL + "AppliedGoMainCombined.htm")
	assert.ElementsMatch(t, lExpectedLinksList, lReceivedLinksList)
	assert.ElementsMatch(t, []error(nil), lReceivedWarningsList)
	assert.Equal(t, nil, lErr)

}
