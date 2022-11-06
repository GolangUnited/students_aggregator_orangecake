package handlers

import (
	"errors"
	"github.com/gocolly/colly"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"net/http"
	"strings"
	"time"
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

type appliedGoArticle struct {
	article  core.Article
	warnings []string
	errors   []error
}

func newAppliedGoArticle() appliedGoArticle {
	return appliedGoArticle{
		article:  core.Article{},
		warnings: make([]string, 0),
		errors:   make([]error, 0),
	}
}

type appliedGoParser struct {
	articles []appliedGoArticle
	warnings []string
	errors   []error
}

func newAppliedGoParser() appliedGoParser {
	return appliedGoParser{
		articles: make([]appliedGoArticle, 0),
		warnings: make([]string, 0),
		errors:   make([]error, 0),
	}
}

func (p *appliedGoParser) ParseAppliedGo(link string) error {
	var lArticleCollector = colly.NewCollector()
	lArticleCollector.OnHTML(appliedGoArticlePath, func(e *colly.HTMLElement) {
		lNewArticle := newAppliedGoArticle()
		lCorrect := true

		//Link
		aLink := e.ChildAttr(appliedGoLinkPath, "href")
		if aLink == "" {
			//TODO
			lNewArticle.errors = append(lNewArticle.errors, errors.New("error: link field is empty"))
			lCorrect = false
		}

		//Title
		aTitle := e.ChildText(appliedGoTitlePath)
		if aTitle == "" {
			//TODO
			lNewArticle.errors = append(lNewArticle.errors, errors.New("error: title field is empty"))
			lCorrect = false
		}

		//Date
		var aPublicationDate time.Time
		var dateErr error
		aPublicationDateStr := e.ChildAttr(appliedGoDatePath, "datetime")
		aPublicationDate, dateErr = core.ParseDate(appliedGoDateLayout, strings.TrimSpace(aPublicationDateStr))
		if dateErr != nil {
			lNewArticle.errors = append(lNewArticle.errors, dateErr)
		}

		//Description
		aDescription := e.ChildText(appliedGoDescrPath)
		if aDescription == "" {
			//TODO
			lNewArticle.warnings = append(lNewArticle.warnings, "warning: description field is empty")
		}

		if lCorrect {
			lNewArticle.article.Title = aTitle
			lNewArticle.article.Description = aDescription
			lNewArticle.article.Link = aLink
			lNewArticle.article.PublishDate = aPublicationDate
		}
		p.articles = append(p.articles, lNewArticle)
	})

	_, lCallErr := http.Get(link)
	if lCallErr != nil {
		return lCallErr
	}

	lVisitErr := lArticleCollector.Visit(link)
	lArticleCollector.Wait()

	if len(p.articles) == 0 {
		return core.ErrNoArticles
	}

	return lVisitErr
}

//Link core.RequiredFieldError{ErrorType: core.ErrFieldIsEmpty, Field: core.LinkFieldName}
