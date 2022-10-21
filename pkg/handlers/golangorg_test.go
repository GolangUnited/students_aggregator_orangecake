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
	TestDataFolder = "./test_data"
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

func TestGolangorg(t *testing.T) {

	testServer := newTestServer()
	defer testServer.Close()

	// create expected data
	lExpectedData := []core.Article{
		{
			Title:  "Go runtime: 4 years later",
			Author: "Michael Knyszek",
			Link:   "https://tip.golang.org/blog/all/go119runtime",
			PublishDate: func() time.Time {
				t, _ := time.Parse("_2 January 2006", "26 September 2022")
				return t.UTC()
			}(),
			Description: "A check-in on the status of Go runtime development",
		},
		{
			Title:  "Go Developer Survey 2022 Q2 Results",
			Author: "Todd Kulesza",
			Link:   "https://tip.golang.org/blog/all/survey2022-q2-results",
			PublishDate: func() time.Time {
				t, _ := time.Parse("_2 January 2006", "8 September 2022")
				return t.UTC()
			}(),
			Description: "An analysis of the results from the 2022 Q2 Go Developer Survey.",
		},
		{
			Title:  "Vulnerability Management for Go",
			Author: "Julie Qiu, for the Go security team",
			Link:   "https://tip.golang.org/blog/all/vuln",
			PublishDate: func() time.Time {
				t, _ := time.Parse("_2 January 2006", "6 September 2022")
				return t.UTC()
			}(),
			Description: "Announcing vulnerability management for Go, to help developers learn about known vulnerabilities in their dependencies.",
		},
		{
			Title:  "Go 1.19 is released!",
			Author: "The Go Team",
			Link:   "https://tip.golang.org/blog/all/go1.19",
			PublishDate: func() time.Time {
				t, _ := time.Parse("_2 January 2006", "2 August 2022")
				return t.UTC()
			}(),
			Description: "Go 1.19 adds richer doc comments, performance improvements, and more.",
		},
		{
			Title:  "Share your feedback about developing with Go",
			Author: "Todd Kulesza",
			Link:   "https://tip.golang.org/blog/all/survey2022-q2",
			PublishDate: func() time.Time {
				t, _ := time.Parse("_2 January 2006", "1 June 2022")
				return t.UTC()
			}(),
			Description: "Help shape the future of Go by sharing your thoughts via the Go Developer Survey",
		},
	}

	lGot, lErr := GolangorgScraper(testServer.URL + File)
	if lErr != nil {
		fmt.Println("function Scraping return the error: ", lErr)
		lGotErr = lErr
	}

	lWant := lExpectedData

	//compare data
	if !reflect.DeepEqual(lGot, lWant) {
		t.Errorf("Mismatch between expected and actual data")
	}

	//compare errors
	if lErr != lGotErr {
		t.Errorf("Mismatch between expected and actual errors")
	}
}
