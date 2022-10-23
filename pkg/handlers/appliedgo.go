package handlers

import (
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"net/http"
	"strings"
	"time"
)

// ParseAppliedGo returns ...
// ToDo about args
func ParseAppliedGo(aStartLink string) (aArticles []core.Article, aWarnings []error, aError error) {
	linksList, err := ParseAppliedGoMain(aStartLink)

	var lArticlesList []core.Article
	var lErrorsList []error
	for _, value := range linksList {
		lArticle, lErr := ParseAppliedGoArticle(value)
		lArticlesList = append(lArticlesList, lArticle)
		lErrorsList = append(lErrorsList, lErr...)
	}
	return lArticlesList, lErrorsList, err
}

func ParseAppliedGoMain(link string) ([]string, error) {
	var lLinksList []string
	var lArticleCollector = colly.NewCollector()
	lArticleCollector.OnHTML("header[class=\"article-header\"] > a", func(e *colly.HTMLElement) {
		lLinksList = append(lLinksList, e.Attr("href"))
	})
	_, err := http.Get(link)
	lArticleCollector.Visit(link)
	return lLinksList, err
}

func ParseAppliedGoArticle(link string) (core.Article, []error) {
	var lNewArticle core.Article
	var lParseErrors []error
	var lArticleParser = colly.NewCollector()
	lArticleParser.OnHTML("head", func(e *colly.HTMLElement) {

		aTitle := e.ChildAttr("meta[property=\"og:title\"]", "content")
		if aTitle == "" {
			lParseErrors = append(lParseErrors, errors.New("error: title field is empty"))
		}

		aDescription := e.ChildAttr("meta[property=\"og:description\"]", "content")
		if aDescription == "" {
			lParseErrors = append(lParseErrors, errors.New("warning: description field is empty"))
		}

		var aPublicationDate time.Time
		var dateErr error
		aPublicationDateStr := e.ChildAttr("meta[property=\"article:published_time\"]", "content")
		if aPublicationDateStr == "" {
			lParseErrors = append(lParseErrors, errors.New("warning: date field is empty"))
			aPublicationDate = time.Date(1970, time.January, 20, 0, 0, 0, 0, time.UTC)
		} else {
			aPublicationDate, dateErr = core.ParseDate(time.RFC3339, strings.TrimSpace(aPublicationDateStr))
			if dateErr != nil {
				lParseErrors = append(lParseErrors, fmt.Errorf("invalid date format: %s", dateErr))
			}
		}

		lNewArticle = core.Article{
			Title:       aTitle,
			Description: aDescription,
			Link:        link,
			PublishDate: aPublicationDate,
		}
	})
	_, lCallErr := http.Get(link)
	lArticleParser.Visit(link)
	if lCallErr == nil {
		return lNewArticle, lParseErrors
	}
	return lNewArticle, []error{lCallErr}
}

//func ParseTitle(aTitle string) {
//
//}
