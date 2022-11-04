package handlers

import (
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

const (
	TestDataFolder = "./testdata"
	File           = "/golangorg.html"
)

var lGotErr error

// newGolangOrgTestServer create a new server
func newGolangOrgTestServer() *httptest.Server {
	mux := http.NewServeMux()
	fServer := http.FileServer(http.Dir(TestDataFolder))
	mux.Handle("/", fServer)
	return httptest.NewServer(mux)
}

func TestGolangOrgData(t *testing.T) {

	testServer := newGolangOrgTestServer()
	defer testServer.Close()

	// create expected data
	lExpectedData := []core.Article{
		{
			Title:       "Vulnerability Management for Go",
			Author:      "Julie Qiu, for the Go security team",
			Link:        "https://tip.golang.org/blog/vuln",
			PublishDate: time.Date(2022, time.September, 6, 0, 0, 0, 0, time.UTC),
			Description: "",
		},
		{
			Title:       "Go 1.19 is released!",
			Author:      "The Go Team",
			Link:        "https://tip.golang.org/blog/go1.19",
			PublishDate: time.Date(2022, time.August, 2, 0, 0, 0, 0, time.UTC),
			Description: "Go 1.19 adds richer doc comments, performance improvements, and more.",
		},
		{
			Title:       "Share your feedback about developing with Go",
			Author:      "Todd Kulesza",
			Link:        "https://tip.golang.org/blog/survey2022-q2",
			PublishDate: time.Date(2022, time.June, 1, 0, 0, 0, 0, time.UTC),
			Description: "Help shape the future of Go by sharing your thoughts via the Go Developer Survey",
		},
	}

	h := NewGolangOrgHandler(testServer.URL + File)
	lArticles, lErr := h.GolangOrgScraper()
	if lErr != nil {
		t.Error(lErr.Error())
		return
	}

	for i, lExpectedArticle := range lExpectedData {
		if !reflect.DeepEqual(lArticles[i], lExpectedArticle) {
			t.Errorf("Expected %s, but got %s", lArticles[i], lExpectedArticle)
		}
	}

	if lErr != lGotErr {
		t.Errorf("Mismatch between expected and actual errors")
	}
}

func TestGolangOrgHandler_EmptyUrl(t *testing.T) {
	golangOrgHandler := NewGolangOrgHandler("")
	lArticles, lErr := golangOrgHandler.GolangOrgScraper()
	if lArticles != nil {
		t.Errorf("articles and warnings must be nil")
	}

	if lErr == nil {
		t.Errorf("error must be not nil")
	}
}
