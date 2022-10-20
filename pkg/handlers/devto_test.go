package handlers

import (
	"net/http"
	"testing"
    "time"

	"github.com/gocolly/colly"
    "github.com/stretchr/testify/assert"
    "github.com/indikator/aggregator_orange_cake/pkg/core"
)

func TestGetArticlesListAppliedGo(t *testing.T) {
    lTransport := &http.Transport{}
    lTransport.RegisterProtocol("file", http.NewFileTransport(http.Dir("./test_data/")))

    lArticleCollectorTester := colly.NewCollector()
    lArticleCollectorTester.WithTransport(lTransport)

    handle := NewTestHandler("file://./Devto.html")

    lArticles := handle.Scrap()

	lExpectedArticles := 
        []core.Article(
            []core.Article{
                {
                    Title:"Restful API with Golang practical approach",
                    Author:"Firdavs Kasymov",
                    Link:"file://./Article1.html",
                    PublishDate:time.Date(2022, time.October, 17, 8, 42, 7, 0, time.UTC),
                    Description:"In this tutorial, we would be creating a Restful API with a practical approach of clean architecture and native Golang without any frameworks.",
                },
                {
                    Title:"Learn Go in Minutes",
                    Author:"Ayoub Ali",
                    Link:"file://./Article2.html",
                    PublishDate:time.Date(2022, time.October, 16, 14, 42, 58, 0, time.UTC),
                    Description:"\nGo is an open source programming language supported by Google.\nEasy to learn and get started with.\nBuilt-in concurrency and a robust standard library.\nGrowing ecosystem of partners, communities, and tools.\n",
                },
            },
        )

    assert.Equal(t, lExpectedArticles, lArticles, "ArticlesEqual")
}