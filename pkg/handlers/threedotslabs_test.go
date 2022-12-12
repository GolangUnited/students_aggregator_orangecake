package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"github.com/stretchr/testify/assert"
)

// TODO: move to test utils package
func buildTestDataPath() (string, error) {
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

// TODO: move to test utils package
func buildTestFileName(aFileName string) (string, error) {
	lPath, lErr := buildTestDataPath()
	if lErr != nil {
		return "", lErr
	}

	return filepath.Join(lPath, aFileName), nil
}

// TODO: move to test utils package
func buildTestFileServer(aDataPath string) (*httptest.Server, error) {
	lServeMux := http.NewServeMux()
	lFileServer := http.FileServer(http.Dir(aDataPath))
	lServeMux.Handle("/", lFileServer)
	return httptest.NewServer(lServeMux), nil
}

// TODO: move to test utils package
func buildTestDataPathServer() (*httptest.Server, error) {
	lDataPath, lErr := buildTestDataPath()
	if lErr != nil {
		return nil, lErr
	}

	return buildTestFileServer(lDataPath)
}

// TODO: move to test utils package
func readTestFile(aFileName string) (io.Reader, error) {
	lFileName, lErr := buildTestFileName(aFileName)
	if lErr != nil {
		return nil, lErr
	}

	lData, lErr := os.ReadFile(lFileName)
	if lErr != nil {
		return nil, lErr
	}

	return bytes.NewReader(lData), nil
}

// TODO: move to test utils
type threeDotsLabsErrorIOReader struct{ Error error }

func (e threeDotsLabsErrorIOReader) Read([]byte) (int, error) {
	return 0, e.Error
}

func TestThreeDotsLabsParseError(t *testing.T) {
	const READ_ERROR = "internal read error"
	lReader := threeDotsLabsErrorIOReader{Error: errors.New(READ_ERROR)}

	lHandler := NewThreeDotsLabsHandler("")
	lArticles, lWarnings, lErr := lHandler.ParseHtml(lReader)

	assert.EqualError(t, lErr, READ_ERROR)
	assert.Empty(t, lArticles, "Articles")
	assert.Empty(t, lWarnings, "Warnings")
}

func TestThreeDotsLabsData(t *testing.T) {
	lExpectedArticles := [...]core.Article{
		{Author: "Author", Title: "Title 1", Link: "LinkUrl", Description: "Summary", PublishDate: time.Date(2022, time.February, 1, 0, 0, 0, 0, time.UTC)},
		{Author: /*e*/ "", Title: "Title 2", Link: "LinkUrl", Description: "Summary", PublishDate: core.NormalizeDate(time.Now())},
		{Author: /*e*/ "", Title: "Title 3", Link: "LinkUrl", Description: "Summary", PublishDate: core.NormalizeDate(time.Now())},
		{Author: "Author", Title: "Title 4", Link: "LinkUrl", Description: "Summary", PublishDate: core.NormalizeDate(time.Now())},
		{Author: /*e*/ "", Title: "Title 5", Link: "LinkUrl", Description: "Summary", PublishDate: core.NormalizeDate(time.Now())},
		{Author: /*e*/ "", Title: "Title 6", Link: "LinkUrl", Description: /*em*/ "", PublishDate: core.NormalizeDate(time.Now())},
	}

	lExpectedWarnings := [...]core.Warning{
		"Warning[1,0]: article Header is empty",
		"Warning[2,0]: invalid article Header format",
		"Warning[3,0]: cannot parse article date 'Feb 31, 2022'. invalid Date format",
		"Warning[4,0]: article Author is empty",
		"Warning[4,1]: cannot parse article date ''. empty Date",
		"Warning[5,0]: article Header node not found",
		"Warning[5,1]: article Description is empty",
		"Error[6]: article Link node not found",
		"Error[7]: article Title is empty",
		"Error[8]: article Link URL not found",
		"Error[9]: article Link URL is empty",
	}

	lTestData, lErr := readTestFile("threedotslabs_test.html")
	if lErr != nil {
		t.Error(lErr.Error())
		return
	}

	lHandler := NewThreeDotsLabsHandler("")
	lArticles, lWarnings, lErr := lHandler.ParseHtml(lTestData)
	if lErr != nil {
		t.Error(lErr.Error())
		return
	}

	if assert.Equal(t, len(lExpectedArticles), len(lArticles), "Invalid article count") {
		for i, lArticle := range lArticles {
			assert.Equal(t, lExpectedArticles[i], lArticle, fmt.Sprintf("Article %d", i))
		}
	}

	if assert.Equal(t, len(lExpectedWarnings), len(lWarnings), "Invalid warning count") {
		for i, lWarning := range lWarnings {
			assert.Equal(t, lExpectedWarnings[i], lWarning, fmt.Sprintf("Warning %d", i))
		}
	}
}

func TestThreeDotsLabsUrl(t *testing.T) {
	// test we can parse the actual web page, it should contain 10 articles and no warnings
	const EXPECTED_COUNT = 10

	lHandler := NewThreeDotsLabsHandler(THREE_DOTS_LABS_URL)

	lArticles, lWarnings, lErr := lHandler.ParseArticles()

	assert.Nil(t, lErr, "Error")
	assert.Equal(t, EXPECTED_COUNT, len(lArticles), "ArticleCount")
	assert.Empty(t, lWarnings, "Warnings")
}

func TestThreeDotsLabsEmptyUrl(t *testing.T) {
	// test HTTP Error (empty URL)
	lHandler := NewThreeDotsLabsHandler("")

	lArticles, lWarnings, lErr := lHandler.ParseArticles()

	assert.Empty(t, lArticles, "Articles")
	assert.Empty(t, lWarnings, "Warnings")
	assert.EqualError(t, lErr, "Get \"\": unsupported protocol scheme \"\"")
}

func TestThreeDotsLabsHttpStatusError(t *testing.T) {
	lTestServer, lErr := buildTestDataPathServer()
	if !assert.Nil(t, lErr, "TestServer") {
		return // no need to continue
	}
	defer lTestServer.Close()

	lHandler := NewThreeDotsLabsHandler(lTestServer.URL + "/not_found.html")
	lArticles, lWarnings, lErr := lHandler.ParseArticles()

	assert.Nil(t, lArticles, "Articles")
	assert.Nil(t, lWarnings, "Warnings")
	assert.EqualError(t, lErr, "status code error: 404 404 Not Found")
}
