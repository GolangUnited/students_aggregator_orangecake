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
	url    string
	logger core.Logger
}

func NewAppliedGoScraper(url string, logger core.Logger) core.ArticleScraper {
	return &appliedGoScraper{
		url:    url,
		logger: logger,
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
			lError := core.RequiredFieldError{ErrorType: core.ErrFieldIsEmpty, Field: core.LinkFieldName}
			lWarnings = append(lWarnings, core.Warning(lError.Error()))
			return
		}
		lNewArticle.Link = aLink

		//Title (required field)
		aTitle := e.ChildText(appliedGoTitlePath)
		if aTitle == "" {
			lError := core.RequiredFieldError{ErrorType: core.ErrFieldIsEmpty, Field: core.TitleFieldName}
			lWarnings = append(lWarnings, core.Warning(lError.Error()))
			return
		}
		lNewArticle.Title = aTitle

		//Date
		aPublishDateStr := e.ChildAttr(appliedGoDatePath, "datetime")
		aPublishDate, aDateErr := core.ParseDate(appliedGoDateLayout, strings.TrimSpace(aPublishDateStr))
		if aDateErr != nil {
			if aDateErr.Error() == core.ErrEmptyDate.Error() {
				lError := core.EmptyFieldError{Field: core.PublishDateFieldName}
				lWarnings = append(lWarnings, core.Warning(lError.Error()))
			} else if aDateErr.Error() == core.ErrInvalidDateFormat.Error() {
				lWarnings = append(lWarnings, core.Warning(fmt.Sprintf("cannot parse article date '%s', %s", aPublishDate, aDateErr.Error())))
			}
		}
		lNewArticle.PublishDate = aPublishDate

		//Description
		aDescription := e.ChildText(appliedGoDescrPath)
		if aDescription == "" {
			lError := core.EmptyFieldError{Field: core.DescriptionFieldName}
			lWarnings = append(lWarnings, core.Warning(lError.Error()))
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
