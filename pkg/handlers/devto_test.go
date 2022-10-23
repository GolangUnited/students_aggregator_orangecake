package handlers

import (
	"net/http"
	"testing"
    "time"

	"github.com/gocolly/colly"
    "github.com/stretchr/testify/assert"
    "github.com/indikator/aggregator_orange_cake/pkg/core"
)

const (
	TestDataFolder = "./test_data/"
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

func TestGetArticlesListAppliedGo(t *testing.T) {

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
