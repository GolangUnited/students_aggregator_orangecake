package handlers

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"net/http"
	"strings"
	"time"
)

type appliedGoArticleParser struct {
	Article  core.Article
	Warnings []string
}

func newAppliedGoArticleParser() appliedGoArticleParser {
	return appliedGoArticleParser{
		Article:  core.Article{},
		Warnings: make([]string, 0),
	}
}

type appliedGoMainParser struct {
	Links []string
}

func newAppliedGoMainParser() appliedGoMainParser {
	return appliedGoMainParser{
		Links: make([]string, 0),
	}
}

// ParseAppliedGo returns the list of parsed articles from appliedgo.net site.
func ParseAppliedGo(aStartLink string) (aArticles []core.Article, aWarnings []string, aError error) {
	lReceivedLinks := newAppliedGoMainParser()
	lReceivedErr := lReceivedLinks.ParseAppliedGoMain(aStartLink)
	if lReceivedErr != nil {
		return nil, nil, lReceivedErr
	}

	var lArticlesList []core.Article
	var lWarningsList []string
	for aIndex, value := range lReceivedLinks.Links {
		lNewArticle := newAppliedGoArticleParser()
		lParseErr := lNewArticle.ParseAppliedGoArticle(value)
		if lParseErr != nil {
			lWarningsList = append(lWarningsList, fmt.Sprintf("Error[%d]: %s", aIndex, lParseErr.Error()))
		}
		lArticlesList = append(lArticlesList, lNewArticle.Article)
		for i, lWarning := range lNewArticle.Warnings {
			lWarningsList = append(lWarningsList, fmt.Sprintf("Warning[%d,%d]: %s", aIndex, i, lWarning))
		}

	}
	return lArticlesList, lWarningsList, nil
}

func (p *appliedGoMainParser) ParseAppliedGoMain(link string) error {
	var lArticleCollector = colly.NewCollector()
	lArticleCollector.OnHTML("header[class=\"article-header\"] > a", func(e *colly.HTMLElement) {
		p.Links = append(p.Links, e.Attr("href"))
	})
	_, lCallErr := http.Get(link)
	if lCallErr != nil {
		return lCallErr
	}
	lVisitErr := lArticleCollector.Visit(link)
	return lVisitErr
}

func (p *appliedGoArticleParser) ParseAppliedGoArticle(link string) error {
	var lArticleParser = colly.NewCollector()
	lArticleParser.OnHTML("head", func(e *colly.HTMLElement) {

		aTitle := e.ChildAttr("meta[property=\"og:title\"]", "content")
		if aTitle == "" {
			p.Warnings = append(p.Warnings, "error: title field is empty")
		}

		aDescription := e.ChildAttr("meta[property=\"og:description\"]", "content")
		if aDescription == "" {
			p.Warnings = append(p.Warnings, "warning: description field is empty")
		}

		var aPublicationDate time.Time
		var dateErr error
		aPublicationDateStr := e.ChildAttr("meta[property=\"article:published_time\"]", "content")
		if aPublicationDateStr == "" {
			p.Warnings = append(p.Warnings, "warning: date field is empty")
			aPublicationDate = time.Date(1970, time.January, 20, 0, 0, 0, 0, time.UTC)
		} else {
			aPublicationDate, dateErr = core.ParseDate(time.RFC3339, strings.TrimSpace(aPublicationDateStr))
			if dateErr != nil {
				p.Warnings = append(p.Warnings, fmt.Sprintf("invalid date format: %s", dateErr))
			}
		}

		p.Article.Title = aTitle
		p.Article.Description = aDescription
		p.Article.Link = link
		p.Article.PublishDate = aPublicationDate
	})

	_, lCallErr := http.Get(link)
	if lCallErr != nil {
		return lCallErr
	}
	lVisitErr := lArticleParser.Visit(link)
	return lVisitErr
}
