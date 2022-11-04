package handlers

import (
    "net/http"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/indikator/aggregator_orange_cake/pkg/core"
)

const TestDataFolder = "./test_data/devto/"

func NewTestDevtoHandler(aURL string) DevtoHandler {
    // We don't need to add anything to link during tests
    Devto_URL = ""

    // Switching Colly's data source from network to local files
    lTransport := &http.Transport{}
    lTransport.RegisterProtocol("file", http.NewFileTransport(http.Dir(TestDataFolder)))

    lHandler := NewDevtoHandler(aURL)
    lHandler.colly.WithTransport(lTransport)

    return lHandler
}

func TestScrapDevto(t *testing.T) {

    lHandle := NewTestDevtoHandler("file://./Devto.html")
    lHandle.Scrap()

    lExpectedArticles := 
        []core.Article(
            []core.Article{
                {
                    Title:"Restful API with Golang practical approach",
                    Author:"Firdavs Kasymov",
                    Link:"file://./DevtoArticle1.html",
                    PublishDate:time.Date(2022, time.October, 17, 0, 0, 0, 0, time.UTC),
                    Description:"In this tutorial, we would be creating a Restful API with a practical approach of clean architecture and native Golang without any frameworks.",
                },
                {
                    Title:"Learn Go in Minutes",
                    Author:"Ayoub Ali",
                    Link:"file://./DevtoArticle2.html",
                    PublishDate:time.Date(2022, time.October, 16, 0, 0, 0, 0, time.UTC),
                    Description:"\nGo is an open source programming language supported by Google.\nEasy to learn and get started with.\nBuilt-in concurrency and a robust standard library.\nGrowing ecosystem of partners, communities, and tools.\n",
                },
            },
        )

    assert.Equal(t, lExpectedArticles, lHandle.articles, "ArticlesEqual")
}

func TestScrapDevtoEmptyNotRequiredFields(t *testing.T) {

    lHandle := NewTestDevtoHandler("file://./DevtoEmptyFields.html")
    lHandle.Scrap()

    lExpectedArticles := 
        []core.Article(
            []core.Article{
                {
                    Title:"Restful API with Golang practical approach",
                    Author:"",
                    Link:"file://./DevtoArticle1.html",
                    PublishDate:time.Date(2022, time.October, 17, 0, 0, 0, 0, time.UTC),
                    Description:"In this tutorial, we would be creating a Restful API with a practical approach of clean architecture and native Golang without any frameworks.",
                },
                {
                    Title:"Learn Go in Minutes",
                    Author:"Ayoub Ali",
                    Link:"file://./DevtoArticle2.html",
                    PublishDate: core.NormalizeDate(time.Now()),
                    Description:"\nGo is an open source programming language supported by Google.\nEasy to learn and get started with.\nBuilt-in concurrency and a robust standard library.\nGrowing ecosystem of partners, communities, and tools.\n",
                },
            },
        )

    assert.Equal(t, lExpectedArticles, lHandle.articles, "ArticlesEqual")
}

func TestScrapDevtoEmptyTitle(t *testing.T) {

    var lErr error
    var lExpectedErr = core.RequiredFieldError{ErrorType: core.ErrFieldIsEmpty, Field: core.TitleFieldName}

    lHandle := NewTestDevtoHandler("file://./DevtoEmptyTitle.html")
    lErr = lHandle.Scrap()

    assert.Equal(t, lExpectedErr, lErr, "ErrorsEqual")
}

func TestScrapDevtoEmptyLink(t *testing.T) {

    var lErr error
    var lExpectedErr = core.RequiredFieldError{ErrorType: core.ErrFieldIsEmpty, Field: core.LinkFieldName}

    lHandle := NewTestDevtoHandler("file://./DevtoEmptyLink.html")
    lErr = lHandle.Scrap()

    assert.Equal(t, lExpectedErr, lErr, "ErrorsEqual")
}

func TestScrapDevtoBadUrl(t *testing.T) {

    lBadUrl := "https://de.vto/t/go"

    var lErr error
    var lExpectedErr = core.ErrHTMLAccess

    lHandle := NewDevtoHandler(lBadUrl)
    lErr = lHandle.Scrap()

    assert.Equal(t, lExpectedErr, lErr, "ErrorsEqual")
}
