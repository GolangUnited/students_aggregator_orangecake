package handlers

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
    "io"

    "github.com/stretchr/testify/assert"
    "github.com/indikator/aggregator_orange_cake/pkg/core"
)

const TestDataFolder = "./testdata/devto/"

func newDevtoTestServer() *httptest.Server {
	lMux := http.NewServeMux()
	fServer := http.FileServer(http.Dir(TestDataFolder))
	lMux.Handle("/", fServer)
	return httptest.NewServer(lMux)
}

func TestScrapDevto(t *testing.T) {

    log := core.NewZeroLogger(io.Discard)

    testServer := newDevtoTestServer()
    Devto_URL = testServer.URL
    defer testServer.Close()

    lHandle := NewDevtoHandler(testServer.URL + "/Devto.html", log)
    lArticles, lWarnings, lErr := lHandle.ParseArticles()

    lExpectedArticles := 
        []core.Article(
            []core.Article{
                {
                    Title:"Restful API with Golang practical approach",
                    Author:"Firdavs Kasymov",
                    Link: testServer.URL + "/DevtoArticle1.html",
                    PublishDate:time.Date(2022, time.October, 17, 0, 0, 0, 0, time.UTC),
                    Description:"In this tutorial, we would be creating a Restful API with a practical approach of clean architecture and native Golang without any frameworks.",
                },
                {
                    Title:"Learn Go in Minutes",
                    Author:"Ayoub Ali",
                    Link: testServer.URL + "/DevtoArticle2.html",
                    PublishDate:time.Date(2022, time.October, 16, 0, 0, 0, 0, time.UTC),
                    Description:"\nGo is an open source programming language supported by Google.\nEasy to learn and get started with.\nBuilt-in concurrency and a robust standard library.\nGrowing ecosystem of partners, communities, and tools.\n",
                },
            },
        )

    assert.Equal(t, lExpectedArticles, lArticles, "ArticlesEqual")
    assert.Equal(t, []core.Warning(nil), lWarnings, "WarningsEqual")
    assert.Equal(t, nil, lErr, "ErrorsEqual")
}

func TestScrapDevtoEmptyNotRequiredFields(t *testing.T) {

    log := core.NewZeroLogger(io.Discard)

    testServer := newDevtoTestServer()
    Devto_URL = testServer.URL
    defer testServer.Close()

    lHandle := NewDevtoHandler(testServer.URL + "/DevtoEmptyFields.html", log)
    lArticles, lWarnings, lErr := lHandle.ParseArticles()

    lExpectedArticles := 
        []core.Article(
            []core.Article{
                {
                    Title:"Restful API with Golang practical approach",
                    Author:"",
                    Link: testServer.URL + "/DevtoArticle1.html",
                    PublishDate:time.Date(2022, time.October, 17, 0, 0, 0, 0, time.UTC),
                    Description:"In this tutorial, we would be creating a Restful API with a practical approach of clean architecture and native Golang without any frameworks.",
                },
                {
                    Title:"Learn Go in Minutes",
                    Author:"Ayoub Ali",
                    Link: testServer.URL + "/DevtoArticle2.html",
                    PublishDate: core.NormalizeDate(time.Now()),
                    Description:"\nGo is an open source programming language supported by Google.\nEasy to learn and get started with.\nBuilt-in concurrency and a robust standard library.\nGrowing ecosystem of partners, communities, and tools.\n",
                },
            },
        )

    assert.Equal(t, lExpectedArticles, lArticles, "ArticlesEqual")
    assert.Equal(t, []core.Warning{"error: field is empty, field: Author", "date cannot be parsed: empty Date"}, lWarnings, "WarningsEqual")
    assert.Equal(t, nil, lErr, "ErrorsEqual")
}

func TestScrapDevtoEmptyTitle(t *testing.T) {
    testServer := newDevtoTestServer()
    Devto_URL = testServer.URL
    defer testServer.Close()

    var lErr error
    var lExpectedErr = core.RequiredFieldError{ErrorType: core.ErrFieldIsEmpty, Field: core.TitleFieldName}

    log := core.NewZeroLogger(io.Discard)

    lHandle := NewDevtoHandler(testServer.URL + "/DevtoEmptyTitle.html", log)
    _, _, lErr = lHandle.ParseArticles()

    assert.Equal(t, lExpectedErr, lErr, "ErrorsEqual")
}

func TestScrapDevtoEmptyLink(t *testing.T) {
    testServer := newDevtoTestServer()
    Devto_URL = testServer.URL
    defer testServer.Close()

    var lErr error
    var lExpectedErr = core.RequiredFieldError{ErrorType: core.ErrFieldIsEmpty, Field: core.LinkFieldName}

    log := core.NewZeroLogger(io.Discard)

    lHandle := NewDevtoHandler(testServer.URL + "/DevtoEmptyLink.html", log)
    _, _, lErr = lHandle.ParseArticles()

    assert.Equal(t, lExpectedErr, lErr, "ErrorsEqual")
}

func TestScrapDevtoBadUrl(t *testing.T) {

    lBadUrl := "https://de.vto/t/go"

    var lErr error
    var lExpectedErr = core.ErrHTMLAccess

    log := core.NewZeroLogger(io.Discard)

    lHandle := NewDevtoHandler(lBadUrl, log)
    _, _, lErr = lHandle.ParseArticles()

    assert.Equal(t, lExpectedErr, lErr, "ErrorsEqual")
}
