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
	TestDataPath  = "./testdata"
	GolangOrgHtml = "/golangorg.html"
)

var lGotErr error

// newGolangOrgTestServer create a new server
func newGolangOrgTestServer() *httptest.Server {
	mux := http.NewServeMux()
	fServer := http.FileServer(http.Dir(TestDataPath))
	mux.Handle("/", fServer)
	return httptest.NewServer(mux)
}

func TestGolangOrgData(t *testing.T) {

	testServer := newGolangOrgTestServer()
	defer testServer.Close()

	// create expected data
	lExpectedData := []core.Article{

		{
			Title:       "Go 1.19 is released!",
			Author:      "", //empty author
			Link:        "https://tip.golang.org/blog/go1.19",
			PublishDate: time.Date(2022, time.August, 2, 0, 0, 0, 0, time.UTC),
			Description: "Go 1.19 adds richer doc comments, performance improvements, and more.",
		},
		{
			Title:       "Share your feedback about developing with Go",
			Author:      "Todd Kulesza",
			Link:        "https://tip.golang.org/blog/survey2022-q2",
			PublishDate: core.NormalizeDate(time.Now()), //empty date
			Description: "",                             //empty description
		},
		{
			Title:       "Go: What's New in March 2010",
			Author:      "", //empty author and author attribute
			Link:        "https://tip.golang.org/blog/hello-world",
			PublishDate: time.Date(2010, time.March, 18, 0, 0, 0, 0, time.UTC),
			Description: "First post!",
		},
		{
			Title:       "Third-party libraries: goprotobuf and beyond",
			Author:      "Andrew Gerrand",
			Link:        "https://tip.golang.org/blog/protobuf",
			PublishDate: time.Date(2010, time.April, 20, 0, 0, 0, 0, time.UTC),
			Description: "", //empty desc and desc attribute
		},
	}

	lExpectedWarnings := []core.Warning{
		"Warning[0,0]: article's author is empty",
		"Warning[1,0]: cannot parse article date ''. empty Date",
		"Warning[1,1]: article description is empty",
		"Warning[2,0]: article's author is empty",
		"Warning[3,0]: article description node not found",
		"Warning[3,1]: article description is empty",
		"Error[4]: article's title is empty",
		"Error[5]: article's link is empty",
	}

	h := NewGolangOrgHandler(testServer.URL + GolangOrgHtml)
	lArticles, lWarnings, lErr := h.ParseArticles()
	if lErr != nil {
		t.Error(lErr.Error())
		return
	}

	for i, lExpectedArticle := range lExpectedData {
		if !reflect.DeepEqual(lArticles[i], lExpectedArticle) {
			t.Errorf("Expected %s, but got %s", lArticles[i], lExpectedArticle)
		}
	}

	for i, lExpectedWarning := range lExpectedWarnings {
		if !(reflect.DeepEqual(lWarnings[i], lExpectedWarning)) {
			t.Errorf("Expected %s, but got %s", lExpectedWarning, lWarnings[i])
		}
	}
}

func TestGolangOrgHandler_EmptyUrl(t *testing.T) {
	golangOrgHandler := NewGolangOrgHandler("")
	lArticles, lWarnings, lErr := golangOrgHandler.ParseArticles()
	if lArticles != nil && lWarnings != nil {
		t.Errorf("articles and warnings must be nil")
	}

	if lErr == nil {
		t.Errorf("error must be not nil")
	}
}
