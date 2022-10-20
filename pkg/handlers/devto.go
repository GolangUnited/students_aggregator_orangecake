package handlers

import (
    "fmt"
    "time"
    "net/http"

    "github.com/gocolly/colly"
    "github.com/indikator/aggregator_orange_cake/pkg/core"
)

const Go_URL = "https://dev.to/t/go"
var Devto_URL = "https://dev.to"

const Substories_class = "div.substories"
const Story_class = "div.crayons-story"
const Author_class = "button.profile-preview-card__trigger"
const Title_class = "h2.crayons-story__title"
const Article_class = "div.crayons-article__body p:nth-of-type(1)"

type DevtoHandler struct {
    URL      string
    Articles []core.Article
    Colly    *colly.Collector
}

func NewHandler(aURL string) DevtoHandler {

    return DevtoHandler{URL: aURL, Articles: make([]core.Article, 0), Colly: colly.NewCollector()}
}

func NewTestHandler(aURL string) DevtoHandler {
    
    Devto_URL = ""

    lScrapper := colly.NewCollector()

    lTransport := &http.Transport{}
    lTransport.RegisterProtocol("file", http.NewFileTransport(http.Dir("./test_data/")))
    lScrapper.WithTransport(lTransport)

    return DevtoHandler{URL: aURL, Articles: make([]core.Article, 0), Colly: lScrapper}
}

func Run() []core.Article {
    
    h := NewHandler(Go_URL)

    return h.Scrap()
}

func (aHandler DevtoHandler) Scrap() []core.Article {

    lArticle := core.Article{}

    aHandler.Colly.OnHTML(Substories_class, func(e *colly.HTMLElement) {

        e.ForEach(Story_class, func(i int, h *colly.HTMLElement) {

            lArticle.Author = h.ChildText(Author_class)
            lArticle.Title = h.ChildText(Title_class)
            lArticle.Link = Devto_URL + h.ChildAttr("a", "href")

            lDate := h.ChildAttr("time", "datetime")
            lArticle.PublishDate = parseDateDevto(lDate)

            aHandler.Colly.Visit(e.Request.AbsoluteURL(lArticle.Link))

            aHandler.Articles = append(aHandler.Articles, lArticle)
        })

    })

    aHandler.Colly.OnHTML(Article_class, func(e *colly.HTMLElement) {
        lArticle.Description = e.Text
    })

    lErr := aHandler.Colly.Visit(aHandler.URL)
    if lErr != nil {
        fmt.Printf("Error scrapping %s: %s\n\n", aHandler.URL, lErr.Error())
    }

    return aHandler.Articles
}

func parseDateDevto(aDate string) time.Time {

    lDate, lErr := time.Parse(time.RFC3339, aDate)
    if lErr != nil {
        fmt.Printf("Error: %s\n\n", lErr.Error())
        lDate = time.Now()
    }

    return lDate.UTC()
}
