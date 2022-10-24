package handlers

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
)

const Go_URL = "https://dev.to/t/go"
var Devto_URL = "https://dev.to" // To be able to redefine later for testing purposes

const Substories_class = "div.substories"
const Story_class = "div.crayons-story"
const Author_class = "button.profile-preview-card__trigger"
const Title_class = "h2.crayons-story__title"
const Article_class = "div.crayons-article__body p:nth-of-type(1)"

type DevtoHandler struct {
    URL      string
    Articles []core.Article
    Colly    *colly.Collector
    Err      error
}

func NewHandler(aURL string) DevtoHandler {
    return DevtoHandler{URL: aURL, Articles: make([]core.Article, 0), Colly: colly.NewCollector()}
}

// Initializes new Handler and runs scrapping of URL provided
func Run() []core.Article {
    lH := NewHandler(Go_URL)
    lH.Scrap()

    return lH.Articles
}

// Runs scrapping and fills Handler field with Articles info
func (aHandler *DevtoHandler) Scrap() error {

    lArticle := core.Article{}

    aHandler.Colly.OnHTML(Substories_class, func(e *colly.HTMLElement) {

        e.ForEachWithBreak(Story_class, func(i int, h *colly.HTMLElement) bool {

                lArticle.Author = strings.TrimSpace(h.ChildText(Author_class))
                if len(lArticle.Author) <= 0 {
                    fmt.Printf("No author for an Article\n")
                }

                // Title is a requiered field
                lArticle.Title = strings.TrimSpace(h.ChildText(Title_class))
                if len(lArticle.Title) <= 0 {
                    aHandler.Err = errors.New("no title found for an article - quitting")
                    return false
                }

                lDate := h.ChildAttr("time", "datetime")
                var lErr error
                lArticle.PublishDate, lErr = core.ParseDate(time.RFC3339, lDate)
                if lErr != nil {
                    fmt.Printf("date cannot be parsed: %s\n", lErr.Error())
                }

                // Link is a required field
                lArticle.Link = Devto_URL + h.ChildAttr("a", "href")
                if len(lArticle.Link) <= 0 {
                    aHandler.Err = errors.New("no link found for an article - quitting")
                    return false
                }

                // Following link to an Article itself to get the description
                aHandler.Colly.Visit(e.Request.AbsoluteURL(lArticle.Link))

                aHandler.Articles = append(aHandler.Articles, lArticle)
                return true
        })
    })

    // Scrapping Artile.Description from the first paragraph of text
    aHandler.Colly.OnHTML(Article_class, func(e *colly.HTMLElement) {
        lArticle.Description = e.Text
    })

    lErr := aHandler.Colly.Visit(aHandler.URL)
    if lErr != nil {
        return errors.New("error scrapping " + aHandler.URL)
    }
    
    // Checking Handler for any Articles errors
    if aHandler.Err != nil {
        return aHandler.Err
    }

    return nil
}
