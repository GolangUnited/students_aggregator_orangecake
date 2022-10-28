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
	GOLANG_ORG_FULL_URL = "https://tip.golang.org/blog/all"
	GOLANG_ORG_URL      = "https://tip.golang.org"
)

type GolangOrgHandler struct {
	url      string
	articles []core.Article
	err      error
}

func NewGolangOrgHandler(aURL string) GolangOrgHandler {
	return GolangOrgHandler{url: aURL, articles: make([]core.Article, 0)}
}

func (h *GolangOrgHandler) Run() []core.Article {
	h.GolangOrgScraper()
	return h.articles
}

// GolangOrgScraper takes data from tip.golang.com/blog/all and converts it into a structure of json.
func (h *GolangOrgHandler) GolangOrgScraper() ([]core.Article, error) {

	lArticle := core.Article{}
	lArticles := make([]core.Article, 0, 0)

	resp, lErr := http.Get(GOLANG_ORG_FULL_URL)
	if lErr != nil {
		// TODO: write error to log
		h.err = errors.New("http request returns an error: ")
	}

	defer resp.Body.Close()

	if resp.StatusCode > 400 {
		fmt.Println("Status code: ", resp.StatusCode)
	}

	//rewrite the code with goquery, without goColly
	doc, lErr := goquery.NewDocumentFromReader(resp.Body)
	if lErr != nil {
		// TODO: write error to log
		h.err = errors.New("goquery.NewDocumentFromReader returns an error: ")
	}

	doc.Find("p.blogtitle, p.blogsummary").Each(func(i int, selection *goquery.Selection) {
		if i%2 == 0 {
			lLink, _ := selection.Find("a").Attr("href")
			lT, lErr := core.ParseDate("_2 January 2006", selection.Find("span.date").Text())
			if lErr != nil {
				// TODO: write error to log
				h.err = errors.New("error of ParseDate function: ")
			}

			//title is a required field
			lArticle.Title = selection.Find("a[href]").Text()
			if len(lArticle.Title) == 0 {
				// TODO: write error to log
				h.err = errors.New("title not found in article ==> exit")
				//how to return h.err
			}

			lArticle.Link = GOLANG_ORG_URL + lLink
			if len(lLink) == 0 {
				// TODO: write error to log
				h.err = errors.New("link not found in article ==> exit")
				//how to return h.err
			}

			lArticle.Author = selection.Find("span.author").Text()
			if len(lArticle.Author) == 0 {
				// TODO: write error to log
				h.err = errors.New("author not found in article")
			}

			lArticle.PublishDate = lT

			lArticle.Description = strings.TrimSpace(selection.NextFiltered("p.blogsummary").Text())
			if len(lArticle.Description) == 0 {
				// TODO: write error to log
				h.err = errors.New("description not found in article")
			}

			lArticles = append(lArticles, lArticle)
		}
	})

	return lArticles, nil

}
