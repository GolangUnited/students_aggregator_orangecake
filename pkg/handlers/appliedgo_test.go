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
	lExpectedArticles := appliedGoParser{
		articles: []core.Article{
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
		},
		warnings: []string{},
		errors:   []error{}}

	lReceivedArticles := newAppliedGoParser()
	lReceivedErr := lReceivedArticles.ParseAppliedGo(TestDataURL + "AppliedGoMain.htm")

	assert.Equal(t, true, reflect.DeepEqual(lExpectedArticles, lReceivedArticles), "invalid articles content")
	assert.Equal(t, nil, lReceivedErr, "invalid error message")
}

func TestAppliedGoIncorrectUrlProtocol(t *testing.T) {
	var badUrl = ""
	var lExpectedErr = fmt.Errorf("unsupported protocol scheme %q", badUrl)

	lReceivedArticles := newAppliedGoParser()
	lReceivedErr := lReceivedArticles.ParseAppliedGo(badUrl)

	assert.Equal(t, []core.Article{}, lReceivedArticles.articles, "invalid articles content")
	assert.Equal(t, lExpectedErr, errors.Unwrap(lReceivedErr), "invalid error message")
}

func TestAppliedGoScrapeEmptyRequiredFields(t *testing.T) {
	lExpectedArticle := appliedGoParser{
		articles: []core.Article{},
		errors: []error{
			errors.New("error: link field is empty")},
		warnings: []string{},
	}

	lReceivedArticle := newAppliedGoParser()
	lReceivedErr := lReceivedArticle.ParseAppliedGo(TestDataURL + "AppliedGoArticleEmptyRequiredFields.htm")

	assert.Equal(t, lExpectedArticle, lReceivedArticle, "invalid articles content")
	assert.Equal(t, nil, lReceivedErr, "invalid error message")
}

func TestAppliedGoScrapeEmptyNonRequiredFields(t *testing.T) {
	lDate := time.Now()
	lExpectedArticle := appliedGoParser{
		articles: []core.Article{
			{
				Title:       "How I used Go to make my radio auto-switch to AUX-IN when a Raspi plays music",
				Link:        "https://appliedgo.net/auxin/",
				Description: "",
				PublishDate: time.Date(lDate.Year(), lDate.Month(), lDate.Day(), 0, 0, 0, 0, time.UTC),
			},
		},
		errors: []error{
			core.ErrEmptyDate,
		},
		warnings: []string{
			"warning: description field is empty",
		}}

	lReceivedArticle := newAppliedGoParser()
	lReceivedErr := lReceivedArticle.ParseAppliedGo(TestDataURL + "AppliedGoArticleEmptyNonRequiredFields.htm")

	assert.Equal(t, lExpectedArticle, lReceivedArticle, "invalid articles content")
	assert.Equal(t, nil, lReceivedErr, "invalid error message")
}

func TestAppliedGoScrapeInvalidDate(t *testing.T) {
	lDate := time.Now()
	lExpectedArticle := appliedGoParser{
		articles: []core.Article{
			{
				Title:       "How I used Go to make my radio auto-switch to AUX-IN when a Raspi plays music",
				Link:        "https://appliedgo.net/auxin/",
				Description: "Ok, so your radio lacks AirPlay support but has an auxiliary input and can be remote-controlled via the Frontier Silicon API. Fetch a Raspberry Pi, put Shairport-sync and Raspotify on it, plug it into the AUX port, and glue everything together with some Go code. Et voilà - home automation in the small.",
				PublishDate: time.Date(lDate.Year(), lDate.Month(), lDate.Day(), 0, 0, 0, 0, time.UTC),
			},
		},
		errors: []error{
			core.ErrInvalidDateFormat,
		},
		warnings: []string{}}

	lReceivedArticle := newAppliedGoParser()
	lReceivedErr := lReceivedArticle.ParseAppliedGo(TestDataURL + "AppliedGoArticleInvalidDate.htm")

	assert.Equal(t, lExpectedArticle, lReceivedArticle, "invalid articles content")
	assert.Equal(t, nil, lReceivedErr, "invalid error message")
}

func TestAppliedGoNoArticles(t *testing.T) {
	lExpectedArticle := appliedGoParser{
		articles: []core.Article{},
		errors:   []error{},
		warnings: []string{}}

	lReceivedArticle := newAppliedGoParser()
	lReceivedErr := lReceivedArticle.ParseAppliedGo(TestDataURL + "AppliedGoNoArticles.htm")

	assert.Equal(t, lExpectedArticle, lReceivedArticle, "invalid articles content")
	assert.Equal(t, core.ErrNoArticles, lReceivedErr, "invalid error message")
}
