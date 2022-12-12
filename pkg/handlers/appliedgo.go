package handlers

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"net/http"
	"strings"
)

const APPLIED_GO_URL = "https://appliedgo.net/"

const (
	appliedGoArticlePath = "div[class=\"article-inner\"]"
	appliedGoLinkPath    = "header[class=\"article-header\"] > a"
	appliedGoTitlePath   = "h1[class=\"article-title\"]"
	appliedGoDatePath    = "div[class=\"article-date\"] > time"
	appliedGoDateLayout  = "2006-01-02 15:04:05 -0700 MST"
	appliedGoDescrPath   = "div[class=\"summary doc \"] > p"
)

type appliedGoScraper struct {
	url string
}

func newAppliedGoScraper(url string) core.ArticleScraper {
	return &appliedGoScraper{
		url: url,
	}
}

func (s *appliedGoScraper) ParseArticles() ([]core.Article, []core.Warning, error) {
	lArticles := make([]core.Article, 0)
	lWarnings := make([]core.Warning, 0)

	var lArticleCollector = colly.NewCollector()
	lArticleCollector.OnHTML(appliedGoArticlePath, func(e *colly.HTMLElement) {
		lNewArticle := core.Article{}

		//Link (required field)
		aLink := e.ChildAttr(appliedGoLinkPath, "href")
		if aLink == "" {
			//TODO
			lWarnings = append(lWarnings, "error: link field is empty")
			return
		}
		lNewArticle.Link = aLink

		//Title (required field)
		aTitle := e.ChildText(appliedGoTitlePath)
		if aTitle == "" {
			//TODO
			lWarnings = append(lWarnings, "error: title field is empty")
			return
		}
		lNewArticle.Title = aTitle

		//Date
		aPublishDateStr := e.ChildAttr(appliedGoDatePath, "datetime")
		aPublishDate, aDateErr := core.ParseDate(appliedGoDateLayout, strings.TrimSpace(aPublishDateStr))
		if aDateErr != nil {
			lWarnings = append(lWarnings, core.Warning(fmt.Sprintf("cannot parse article date '%s', %s", aPublishDate, aDateErr.Error())))
		}
		lNewArticle.PublishDate = aPublishDate

		//Description
		aDescription := e.ChildText(appliedGoDescrPath)
		if aDescription == "" {
			//TODO
			//Link core.RequiredFieldError{ErrorType: core.ErrFieldIsEmpty, Field: core.LinkFieldName}
			lWarnings = append(lWarnings, "warning: description field is empty")
		}
		lNewArticle.Description = aDescription

		lArticles = append(lArticles, lNewArticle)
	})

	_, lCallErr := http.Get(s.url)
	if lCallErr != nil {
		return nil, nil, lCallErr
	}

	lVisitErr := lArticleCollector.Visit(s.url)
	lArticleCollector.Wait()

	if len(lArticles) == 0 && len(lWarnings) == 0 {
		return nil, nil, core.ErrNoArticles
	}

	return lArticles, lWarnings, lVisitErr
}
