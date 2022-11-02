package handlers

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"net/http"
	"strings"
)

const (
	GOLANG_ORG_URL = "https://tip.golang.org"
)

type GolangOrgHandler struct {
	url      string
	articles []core.Article
	err      error
}

func NewGolangOrgHandler(aUrl string) GolangOrgHandler {
	return GolangOrgHandler{url: aUrl, articles: make([]core.Article, 0)}
}

func (h *GolangOrgHandler) getDescription(aSelection *goquery.Selection) string {
	lDesc := aSelection.Next()
	if len(lDesc.Nodes) == 1 {
		lClass, lOk := lDesc.Attr("class")
		if lOk && lClass == "blogsummary" {
			return strings.TrimSpace(lDesc.Text())
		}
	}

	return ""
}

func (h *GolangOrgHandler) Run() []core.Article {
	h.GolangOrgScraper()
	return h.articles
}

// GolangOrgScraper takes data from tip.golang.com/blog/all and converts it into a structure of json.
func (h *GolangOrgHandler) GolangOrgScraper() ([]core.Article, error) {

	lArticle := core.Article{}
	lArticles := make([]core.Article, 0, 0)

	resp, lErr := http.Get(h.url)
	if lErr != nil {
		// TODO: write error to log
		h.err = fmt.Errorf("http request returns an error: %w", lErr)
		return nil, h.err
	}

	defer resp.Body.Close()

	if resp.StatusCode > 400 {
		h.err = fmt.Errorf("Status code: %d", resp.StatusCode)
		return nil, h.err
	}

	doc, lErr := goquery.NewDocumentFromReader(resp.Body)
	if lErr != nil {
		// TODO: write error to log
		h.err = fmt.Errorf("goquery.NewDocumentFromReader returns an error: %w", lErr)
		return nil, h.err
	}

	// doc.Find("p.blogtitle").Each(func(aIndex int, aSelection *goquery.Selection) {
	doc.Find("p.blogtitle").Each(func(aIndex int, aSelection *goquery.Selection) {
		lOk := true
		lLink, _ := aSelection.Find("a").Attr("href")
		lT, lErr := core.ParseDate("_2 January 2006", aSelection.Find("span.date").Text())
		if lErr != nil {
			// TODO: write error to log
			h.err = errors.New("error of ParseDate function: ")
		}

		//title is a required field
		lArticle.Title = aSelection.Find("a[href]").Text()
		if len(lArticle.Title) == 0 {
			// TODO: write error to log
			lOk = false
			h.err = errors.New("title not found in article ==> exit")
			//how to return h.err
		}

		lArticle.Link = GOLANG_ORG_URL + lLink
		if len(lLink) == 0 {
			// TODO: write error to log
			lOk = false
			h.err = errors.New("link not found in article ==> exit")
			//how to return h.err
		}

		lArticle.Author = aSelection.Find("span.author").Text()
		if len(lArticle.Author) == 0 {
			// TODO: write error to log
			h.err = errors.New("author not found in article")
		}

		lArticle.PublishDate = lT

		lArticle.Description = h.getDescription(aSelection)
		if len(lArticle.Description) == 0 {
			// TODO: write error to log
			h.err = errors.New("description not found in article")
		}

		if lOk {
			lArticles = append(lArticles, lArticle)
		}
	})

	return lArticles, nil

}
