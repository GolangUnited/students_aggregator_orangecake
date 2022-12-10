package handlers

import (
	"errors"
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

type AppliedGoParser struct {
	articles []core.Article
	warnings []string
	errors   []error
}

func NewAppliedGoParser() AppliedGoParser {
	return AppliedGoParser{
		articles: make([]core.Article, 0),
		warnings: make([]string, 0),
		errors:   make([]error, 0),
	}
}

func (p *AppliedGoParser) ParseAppliedGo(link string) error {
	var lArticleCollector = colly.NewCollector()
	lArticleCollector.OnHTML(appliedGoArticlePath, func(e *colly.HTMLElement) {
		lNewArticle := core.Article{}

		//Link
		aLink := e.ChildAttr(appliedGoLinkPath, "href")
		if aLink == "" {
			//TODO
			p.errors = append(p.errors, errors.New("error: link field is empty"))
			return
		}
		lNewArticle.Link = aLink

		//Title
		aTitle := e.ChildText(appliedGoTitlePath)
		if aTitle == "" {
			//TODO
			p.errors = append(p.errors, errors.New("error: title field is empty"))
			return
		}
		lNewArticle.Title = aTitle

		//Date
		aPublishDateStr := e.ChildAttr(appliedGoDatePath, "datetime")
		aPublishDate, dateErr := core.ParseDate(appliedGoDateLayout, strings.TrimSpace(aPublishDateStr))
		if dateErr != nil {
			p.errors = append(p.errors, dateErr)
		}
		lNewArticle.PublishDate = aPublishDate

		//Description
		aDescription := e.ChildText(appliedGoDescrPath)
		if aDescription == "" {
			//TODO
			p.warnings = append(p.warnings, "warning: description field is empty")
		}
		lNewArticle.Description = aDescription

		p.articles = append(p.articles, lNewArticle)
	})

	_, lCallErr := http.Get(link)
	if lCallErr != nil {
		return lCallErr
	}

	lVisitErr := lArticleCollector.Visit(link)
	lArticleCollector.Wait()

	if len(p.articles) == 0 && len(p.errors) == 0 {
		return core.ErrNoArticles
	}

	return lVisitErr
}

//Link core.RequiredFieldError{ErrorType: core.ErrFieldIsEmpty, Field: core.LinkFieldName}
