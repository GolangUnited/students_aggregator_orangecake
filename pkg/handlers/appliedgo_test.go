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
	"reflect"
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
	lExpectedArticles := []core.Article{
		{
			Title:       "How I used Go to make my radio auto-switch to AUX-IN when a Raspi plays music",
			Link:        "https://appliedgo.net/auxin/",
			Description: "Ok, so your radio lacks AirPlay support but has an auxiliary input and can be remote-controlled via the Frontier Silicon API. Fetch a Raspberry Pi, put Shairport-sync and Raspotify on it, plug it into the AUX port, and glue everything together with some Go code. Et voilà - home automation in the small.",
			PublishDate: time.Date(2022, time.August, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			Title:       "Rapid AWS Lambda development with Go and Mantil",
			Link:        "https://appliedgo.net/mantil/",
			Description: "If you need to develop an AWS Lambda function in Go, take a look at Mantil, a dev kit with staging and database connection included.",
			PublishDate: time.Date(2022, time.January, 28, 0, 0, 0, 0, time.UTC),
		},
	}
	lScraper := NewAppliedGoScraper(TestDataURL+"AppliedGoMain.htm", nil)
	lReceivedArticles, lReceivedWarnings, lReceivedErr := lScraper.ParseArticles()

	assert.Equal(t, true, reflect.DeepEqual(lExpectedArticles, lReceivedArticles), "invalid articles content")
	assert.Equal(t, make([]core.Warning, 0), lReceivedWarnings, "invalid warnings")
	assert.Equal(t, nil, lReceivedErr, "invalid error message")
}

func TestAppliedGoIncorrectUrlProtocol(t *testing.T) {
	var badUrl = ""
	var lExpectedErr = fmt.Errorf("unsupported protocol scheme %q", badUrl)

	lScraper := NewAppliedGoScraper(badUrl, nil)
	lReceivedArticles, lReceivedWarnings, lReceivedErr := lScraper.ParseArticles()

	assert.Equal(t, []core.Article(nil), lReceivedArticles, "there should be no articles")
	assert.Equal(t, []core.Warning(nil), lReceivedWarnings, "warnings should be nil")
	assert.Equal(t, lExpectedErr, errors.Unwrap(lReceivedErr), "invalid error message")
}

// Tests the situation when the required field is empty (such as link and title)
func TestAppliedGoScrapeEmptyRequiredFields(t *testing.T) {
	lExpectedWarnings := []core.Warning{
		"error: field is empty, field: Link",
		"error: field is empty, field: Title",
	}

	lScraper := NewAppliedGoScraper(TestDataURL+"AppliedGoArticleEmptyRequiredFields.htm", nil)
	lReceivedArticles, lReceivedWarnings, lReceivedErr := lScraper.ParseArticles()

	assert.Equal(t, []core.Article{}, lReceivedArticles, "there should be no articles")
	assert.Equal(t, lExpectedWarnings, lReceivedWarnings, "invalid warnings")
	assert.Equal(t, nil, lReceivedErr, "invalid error message")
}

// Tests the situation when the required field is empty (such as date and description)
func TestAppliedGoScrapeEmptyNonRequiredFields(t *testing.T) {
	lDate := time.Now().UTC()
	lExpectedDate := time.Date(lDate.Year(), lDate.Month(), lDate.Day(), 0, 0, 0, 0, time.UTC)
	lExpectedArticle := []core.Article{
		{
			Title:       "How I used Go to make my radio auto-switch to AUX-IN when a Raspi plays music",
			Link:        "https://appliedgo.net/auxin/",
			Description: "",
			PublishDate: lExpectedDate,
		},
	}
	lExpectedWarnings := []core.Warning{
		"error: field is empty, field: PublishDate",
		"error: field is empty, field: Description",
	}

	lScraper := NewAppliedGoScraper(TestDataURL+"AppliedGoArticleEmptyNonRequiredFields.htm", nil)
	lReceivedArticles, lReceivedWarnings, lReceivedErr := lScraper.ParseArticles()

	assert.Equal(t, lExpectedArticle, lReceivedArticles, "invalid articles content")
	assert.Equal(t, lExpectedWarnings, lReceivedWarnings, "invalid warnings")
	assert.Equal(t, nil, lReceivedErr, "invalid error message")
}

func TestAppliedGoScrapeInvalidDate(t *testing.T) {
	lDate := time.Now().UTC()
	lExpectedDate := time.Date(lDate.Year(), lDate.Month(), lDate.Day(), 0, 0, 0, 0, time.UTC)
	lExpectedArticle := []core.Article{
		{
			Title:       "How I used Go to make my radio auto-switch to AUX-IN when a Raspi plays music",
			Link:        "https://appliedgo.net/auxin/",
			Description: "Ok, so your radio lacks AirPlay support but has an auxiliary input and can be remote-controlled via the Frontier Silicon API. Fetch a Raspberry Pi, put Shairport-sync and Raspotify on it, plug it into the AUX port, and glue everything together with some Go code. Et voilà - home automation in the small.",
			PublishDate: lExpectedDate,
		},
	}
	lExpectedWarnings := []core.Warning{
		core.Warning(fmt.Sprintf("cannot parse article date '%s', invalid Date format", lExpectedDate)),
	}

	lScraper := NewAppliedGoScraper(TestDataURL+"AppliedGoArticleInvalidDate.htm", nil)
	lReceivedArticles, lReceivedWarnings, lReceivedErr := lScraper.ParseArticles()

	assert.Equal(t, lExpectedArticle, lReceivedArticles, "invalid articles content")
	assert.Equal(t, lExpectedWarnings, lReceivedWarnings, "invalid warnings")
	assert.Equal(t, nil, lReceivedErr, "invalid error message")
}

func TestAppliedGoNoArticles(t *testing.T) {
	lScraper := NewAppliedGoScraper(TestDataURL+"AppliedGoNoArticles.htm", nil)
	lReceivedArticles, lReceivedWarnings, lReceivedErr := lScraper.ParseArticles()

	assert.Equal(t, []core.Article(nil), lReceivedArticles, "there should be no articles")
	assert.Equal(t, []core.Warning(nil), lReceivedWarnings, "warnings should be nil")
	assert.Equal(t, core.ErrNoArticles, lReceivedErr, "invalid error message")
}

func TestScraperConformityToScraperInterface(t *testing.T) {
	lScraper := NewAppliedGoScraper("SomeUrl", nil)
	_, ok := lScraper.(core.ArticleScraper)
	assert.True(t, ok, "The scraper doesn't conform to the core.ArticleScraper interface.")
}
