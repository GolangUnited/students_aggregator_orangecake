package handlers

import (
    "errors"
    "fmt"
    "strings"
    "time"

    "github.com/gocolly/colly"
    "github.com/indikator/aggregator_orange_cake/pkg/core"
)

var Devto_URL = "https://dev.to" // To be able to redefine later for testing purposes

const (
    GO_URL           = "https://dev.to/t/go"
    SUBSTORIES_CLASS = "div.substories"
    STORY_CLASS      = "div.crayons-story"
    AUTHOR_CLASS     = "button.profile-preview-card__trigger"
    TITLE_CLASS      = "h2.crayons-story__title"
    ARTICLE_CLASS    = "div.crayons-article__body p:nth-of-type(1)"
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

    aHandler.colly.OnHTML(SUBSTORIES_CLASS, func(e *colly.HTMLElement) {

        e.ForEachWithBreak(STORY_CLASS, func(i int, h *colly.HTMLElement) bool {

            // Title is a required field
            lArticle.Title = strings.TrimSpace(h.ChildText(TITLE_CLASS))
            if len(lArticle.Title) <= 0 {
                // TODO: switch to custom errors
                aHandler.err = errors.New("no title found for an article - quitting")
                return false
            }

            // Link is a required field
            lLink := h.ChildAttr("a", "href")
            if len(lLink) == 0 {
                // TODO: switch to custom errors
                aHandler.err = errors.New("no link found for an article - quitting")
                return false
            }
            if !strings.HasPrefix(lLink, "/") && !strings.HasPrefix(lLink, "file") {
                lLink = "/" + lLink
            }
            lArticle.Link = Devto_URL + lLink
            
            lArticle.Author = strings.TrimSpace(h.ChildText(AUTHOR_CLASS))
            if len(lArticle.Author) <= 0 {
                // TODO: log
                fmt.Printf("No author for an Article\n")
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
    aHandler.colly.OnHTML(ARTICLE_CLASS, func(e *colly.HTMLElement) {
        lArticle.Description = e.Text
    })

    lErr := aHandler.colly.Visit(aHandler.url)
    if lErr != nil {
        // TODO: switch to custom errors
        return errors.New("error scrapping " + aHandler.url)
    }

    // Checking Handler for any Articles errors
    if aHandler.err != nil {
        // TODO: log
        return aHandler.err
    }

    return nil
}
