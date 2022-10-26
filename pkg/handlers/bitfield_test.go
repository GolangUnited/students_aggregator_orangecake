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
	TEST_DATA_DIR = "./testdata"
)

func newTestServer() *httptest.Server {
	mux := http.NewServeMux()
	fServer := http.FileServer(http.Dir(TEST_DATA_DIR))
	mux.Handle("/", fServer)
	return httptest.NewServer(mux)
}

func TestBitfieldHandler_GetArticles(t *testing.T) {
	testServer := newTestServer()
	defer testServer.Close()

	lDefaultLink, _ := resolveURL(BITFIELD_URL, "/link")
	lDefaultDate := time.Date(2022, time.October, 4, 0, 0, 0, 0, time.UTC)
	expectedData := []core.Article{
		{Author: "author", Title: "Title", Link: lDefaultLink, Description: "summary", PublishDate: lDefaultDate},
		{Author: "author", Title: "Title", Link: lDefaultLink, Description: "summary", PublishDate: core.NormalizeDate(time.Now())},
		{Author: "author", Title: "Title", Link: lDefaultLink, Description: "summary", PublishDate: core.NormalizeDate(time.Now())},
		{Author: "author", Title: "Title", Link: lDefaultLink, Description: "summary", PublishDate: core.NormalizeDate(time.Now())},
		{Author: "author", Title: "Title", Link: lDefaultLink, Description: "summary", PublishDate: core.NormalizeDate(time.Now())},
		{Author: "", Title: "Title", Link: lDefaultLink, Description: "", PublishDate: lDefaultDate},
		{Author: "", Title: "Title", Link: lDefaultLink, Description: "", PublishDate: lDefaultDate},
	}

	lExpectedWarnings := []string{
		"Warning[1,0]: article date attribute not exists",
		"Warning[2,0]: article date node not found",
		"Warning[3,0]: cannot parse article date ''. empty Date",
		"Warning[4,0]: cannot parse article date 'Oct 4, 2022'. invalid Date format",
		"Warning[5,0]: article description node not found",
		"Warning[5,1]: article author node not found",
		"Warning[6,0]: article description is empty",
		"Warning[6,1]: article author is empty",
		"Error[7]: article title and url node not found",
		"Error[8]: article link attribute not found",
		"Error[9]: article link is empty",
		"Error[10]: article title is empty",
	}

	bitfieldHandler := NewBitfieldScrapper(testServer.URL + "/bitfield_test.html")
	lArticles, lWarnings, lErr := bitfieldHandler.GetArticles()
	if lErr != nil {
		t.Error(lErr.Error())
		return
	}

	for i, lExpectedArticle := range expectedData {
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

func TestBitfieldHandler_EmptyUrl(t *testing.T) {
	bitfieldHandler := NewBitfieldScrapper("")
	lArticles, lWarnings, lErr := bitfieldHandler.GetArticles()
	if lArticles != nil && lWarnings != nil {
		t.Errorf("articles and warnings must be nil")
	}

	if lErr == nil {
		t.Errorf("error must be not nil")
	}
}
