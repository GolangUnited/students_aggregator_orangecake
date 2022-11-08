package handlers

import (
    "fmt"
    "strings"
    "time"

    "github.com/gocolly/colly"
    "github.com/indikator/aggregator_orange_cake/pkg/core"
)

var Devto_URL = "https://dev.to" // To be able to redefine later for testing purposes

const (
    DEVTO_GO_URL           = "https://dev.to/t/go"
    DEVTO_SUBSTORIES_CLASS = "div.substories"
    DEVTO_STORY_CLASS      = "div.crayons-story"
    DEVTO_AUTHOR_CLASS     = "button.profile-preview-card__trigger"
    DEVTO_TITLE_CLASS      = "h2.crayons-story__title"
    DEVTO_ARTICLE_CLASS    = "div.crayons-article__body p:nth-of-type(1)"

    FILE_PREFIX = "file"
)

type DevtoHandler struct {
    url      string
    articles []core.Article
    colly    *colly.Collector
    err      error
}

func NewDevtoHandler(aURL string) DevtoHandler {
    return DevtoHandler{url: aURL, articles: make([]core.Article, 0), colly: colly.NewCollector()}
}

// Runs scrapping of URL provided
func (aHandler *DevtoHandler) Run() []core.Article {
    aHandler.Scrap()

    return aHandler.articles
}

// Runs scrapping and fills Handler field with Articles info
func (aHandler *DevtoHandler) Scrap() error {

    lArticle := core.Article{}

    aHandler.colly.OnHTML(DEVTO_SUBSTORIES_CLASS, func(e *colly.HTMLElement) {

        e.ForEachWithBreak(DEVTO_STORY_CLASS, func(i int, h *colly.HTMLElement) bool {

            // Title is a required field
            lArticle.Title = strings.TrimSpace(h.ChildText(DEVTO_TITLE_CLASS))
            if len(lArticle.Title) == 0 {
                aHandler.err = core.RequiredFieldError{ErrorType: core.ErrFieldIsEmpty, Field: core.TitleFieldName}
                return false
            }

            // Link is a required field
            lLink := h.ChildAttr("a", "href")
            if len(lLink) == 0 {
                aHandler.err = core.RequiredFieldError{ErrorType: core.ErrFieldIsEmpty, Field: core.LinkFieldName}
                return false
            }
            if !strings.HasPrefix(lLink, "/") && !strings.HasPrefix(lLink, FILE_PREFIX) {
                lLink = "/" + lLink
            }
            lArticle.Link = Devto_URL + lLink
            
            lArticle.Author = strings.TrimSpace(h.ChildText(DEVTO_AUTHOR_CLASS))
            if len(lArticle.Author) == 0 {
                // TODO: log
                fmt.Printf("%s\n", core.EmptyFieldError{Field: core.AuthorFieldName})
            }

            lDate := h.ChildAttr("time", "datetime")
            var lErr error
            lArticle.PublishDate, lErr = core.ParseDate(time.RFC3339, lDate)
            if lErr != nil {
                // TODO: log
                fmt.Printf("date cannot be parsed: %s\n", lErr.Error())
            }

            // Following link to an Article itself to get the description
            aHandler.colly.Visit(e.Request.AbsoluteURL(lArticle.Link))

            aHandler.articles = append(aHandler.articles, lArticle)
            return true
        })
    })

    // Scrapping Article.Description from the first paragraph of text
    aHandler.colly.OnHTML(DEVTO_ARTICLE_CLASS, func(e *colly.HTMLElement) {
        lArticle.Description = e.Text
    })

    lErr := aHandler.colly.Visit(aHandler.url)
    if lErr != nil {
        return core.ErrHTMLAccess
    }

    // Checking Handler for any Articles errors
    if aHandler.err != nil {
        // TODO: log
        return aHandler.err
    }

    return nil
}