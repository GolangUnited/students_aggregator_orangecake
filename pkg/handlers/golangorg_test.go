package handlers

import (
	"fmt"
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

// newTestServer create a new server
func newTestServer() *httptest.Server {
	mux := http.NewServeMux()
	fServer := http.FileServer(http.Dir(TestDataFolder))
	mux.Handle("/", fServer)
	return httptest.NewServer(mux)
}

func TestGolangOrg(t *testing.T) {

	testServer := newTestServer()
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

	for i, lArticle := range lArticles {
		fmt.Printf("Node %d: %s\n", i, lArticle.Title)
		fmt.Printf("  Author: %s\n", lArticle.Author)
		fmt.Printf("  Date: %s\n", lArticle.PublishDate.Format("Jan _2, 2006"))
		fmt.Printf("  URL: %s\n", lArticle.Link)
		fmt.Printf("  Description:\n    %s\n\n", lArticle.Description)
	}

	fmt.Println()
	fmt.Println()
	fmt.Printf("  %d Articles detected.", len(lArticles))
	fmt.Println()

	lWant := lExpectedData

	//compare data
	if !reflect.DeepEqual(lArticles, lWant) {
		t.Errorf("Mismatch between expected and actual data")
	}

	//compare errors
	if lErr != lGotErr {
		t.Errorf("Mismatch between expected and actual errors")
	}
}
