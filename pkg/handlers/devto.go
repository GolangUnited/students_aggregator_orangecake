package handlers

import (
	"fmt"
	"github.com/gocolly/colly"
	"time"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
)

type DevtoHandler struct {
	URL      string
	Articles []core.Article
}

func NewHandler(aURL string) DevtoHandler{

	return DevtoHandler{URL: aURL, Articles: make([]core.Article, 0)}
}

func Run() []core.Article {
	
	h := NewHandler("https://dev.to/t/go")

	return h.Scrap()
}

func (aHandler DevtoHandler) Scrap() []core.Article {

	lScrapper := colly.NewCollector()

	lArticle := core.Article{}

	lScrapper.OnHTML("div.substories", func(e *colly.HTMLElement) {

		e.ForEach("div.crayons-story", func(i int, h *colly.HTMLElement) {

			lArticle.Author = h.ChildText("button.profile-preview-card__trigger")
			lArticle.Title = h.ChildText("h2.crayons-story__title")
			lArticle.Link = "https://dev.to" + h.ChildAttr("a", "href")

			lDate := h.ChildAttr("time", "datetime")
			lArticle.PublishDate = parseDateDevto(lDate)

			lScrapper.Visit(e.Request.AbsoluteURL(lArticle.Link))

			aHandler.Articles = append(aHandler.Articles, lArticle)
		})

	})

	lScrapper.OnHTML("div.crayons-article__body p:nth-of-type(1)", func(e *colly.HTMLElement) {
		lArticle.Description = e.Text
	})

	lErr := lScrapper.Visit(aHandler.URL)
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
