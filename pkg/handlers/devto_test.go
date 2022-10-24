package handlers

import (
	"net/http"
	"testing"
    "errors"
    "time"

	"github.com/gocolly/colly"
    "github.com/stretchr/testify/assert"
    "github.com/indikator/aggregator_orange_cake/pkg/core"
)

const (
	TestDataFolder = "./test_data/devto/"
)

func NewTestHandler(aURL string) DevtoHandler {
    // We don't need to add anything to link during tests
    Devto_URL = ""

    lScrapper := colly.NewCollector()

    // Switching Colly's data source from network to local files
    lTransport := &http.Transport{}
    lTransport.RegisterProtocol("file", http.NewFileTransport(http.Dir(TestDataFolder)))
    lScrapper.WithTransport(lTransport)

    return DevtoHandler{URL: aURL, Articles: make([]core.Article, 0), Colly: lScrapper}
}

func TestScrapDevto(t *testing.T) {

    lHandle := NewTestHandler("file://./Devto.html")
    lHandle.Scrap()

	lExpectedArticles := 
        []core.Article(
            []core.Article{
                {
                    Title:"Restful API with Golang practical approach",
                    Author:"Firdavs Kasymov",
                    Link:"file://./Article1.html",
                    PublishDate:time.Date(2022, time.October, 17, 0, 0, 0, 0, time.UTC),
                    Description:"In this tutorial, we would be creating a Restful API with a practical approach of clean architecture and native Golang without any frameworks.",
                },
                {
                    Title:"Learn Go in Minutes",
                    Author:"Ayoub Ali",
                    Link:"file://./Article2.html",
                    PublishDate:time.Date(2022, time.October, 16, 0, 0, 0, 0, time.UTC),
                    Description:"\nGo is an open source programming language supported by Google.\nEasy to learn and get started with.\nBuilt-in concurrency and a robust standard library.\nGrowing ecosystem of partners, communities, and tools.\n",
                },
            },
        )

    assert.Equal(t, lExpectedArticles, lHandle.Articles, "ArticlesEqual")
}

func TestScrapDevtoEmptyNotRequiredFields(t *testing.T) {

    lHandle := NewTestHandler("file://./DevtoEmptyFields.html")
    lHandle.Scrap()

	lExpectedArticles := 
        []core.Article(
            []core.Article{
                {
                    Title:"Restful API with Golang practical approach",
                    Author:"",
                    Link:"file://./Article1.html",
                    PublishDate:time.Date(2022, time.October, 17, 0, 0, 0, 0, time.UTC),
                    Description:"In this tutorial, we would be creating a Restful API with a practical approach of clean architecture and native Golang without any frameworks.",
                },
                {
                    Title:"Learn Go in Minutes",
                    Author:"Ayoub Ali",
                    Link:"file://./Article2.html",
                    PublishDate: core.NormalizeDate(time.Now()),
                    Description:"\nGo is an open source programming language supported by Google.\nEasy to learn and get started with.\nBuilt-in concurrency and a robust standard library.\nGrowing ecosystem of partners, communities, and tools.\n",
                },
            },
        )

    assert.Equal(t, lExpectedArticles, lHandle.Articles, "ArticlesEqual")
}

func TestScrapDevtoEmptyTitle(t *testing.T) {

    var lErr error
	var lExpectedErr = errors.New("no title found for an article - quitting")

    lHandle := NewTestHandler("file://./DevtoEmptyTitle.html")
    lErr = lHandle.Scrap()

	assert.Equal(t, lExpectedErr, lErr, "ErrorsEqual")
}

func TestScrapDevtoEmptyLink(t *testing.T) {

    var lErr error
	var lExpectedErr = errors.New("no link found for an article - quitting")

    lHandle := NewTestHandler("file://./DevtoEmptyLink.html")
    lErr = lHandle.Scrap()

	assert.Equal(t, lExpectedErr, lErr, "ErrorsEqual")
}

func TestScrapDevtoBadUrl(t *testing.T) {

    lBadUrl := "https://de.vto/t/go"

    var lErr error
	var lExpectedErr = errors.New("error scrapping " + lBadUrl)

    lHandle := NewHandler(lBadUrl)
    lErr = lHandle.Scrap()

	assert.Equal(t, lExpectedErr, lErr, "ErrorsEqual")
}