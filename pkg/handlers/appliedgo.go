package handlers

import (
	"github.com/gocolly/colly"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"strings"
	"time"
)

// ParseAppliedGo returns ...
// ToDo about args
func ParseAppliedGo(aStartLink string,
	aMainCollector *colly.Collector,
	aArticleCollector *colly.Collector) (aArticles []core.Article, aWarnings []error, aError error) {
	linksList, err := ParseAppliedGoMain(aStartLink, aMainCollector)

	var lArticlesList []core.Article
	var lErrorsList []error
	for _, value := range linksList {
		lArticle, lErr := ParseAppliedGoArticle(value, aArticleCollector)
		lArticlesList = append(lArticlesList, lArticle)
		lErrorsList = append(lErrorsList, lErr)
	}
	return lArticlesList, lErrorsList, err
}

func ParseAppliedGoMain(link string, aArticleCollector *colly.Collector) ([]string, error) {
	var lLinksList []string
	aArticleCollector.OnHTML("header[class=\"article-header\"] > a", func(e *colly.HTMLElement) {
		lLinksList = append(lLinksList, e.Attr("href"))
	})
	err := aArticleCollector.Visit(link)
	return lLinksList, err
}

func ParseAppliedGoArticle(link string, aArticleParser *colly.Collector) (core.Article, error) {
	var lNewArticle core.Article
	var callErr, parseErr error
	aArticleParser.OnHTML("head", func(e *colly.HTMLElement) {
		aName := e.ChildAttr("meta[property=\"og:title\"]", "content")
		aDescription := e.ChildAttr("meta[name=\"description\"]", "content")
		aUrl := e.ChildAttr("meta[property=\"og:url\"]", "content")
		aPublicationDateStr := e.ChildAttr("meta[property=\"article:published_time\"]", "content")
		aPublicationDate, lErr := core.ParseDate(time.RFC3339, strings.TrimSpace(aPublicationDateStr))
		if lErr != nil {
			parseErr = lErr
			//lWarnings = append(lWarnings, fmt.Sprintf("Warning[%d,%d]: %s", aIndex, i, lWarning))
			//lWarnings = append(lWarnings, fmt.Sprintf("Error[%d]: %s", aIndex, lErr.Error()))
		}
		lNewArticle = core.Article{
			Title:       aName,
			Description: aDescription,
			Link:        aUrl,
			PublishDate: aPublicationDate,
		}
	})
	callErr = aArticleParser.Visit(link)
	if callErr == nil {
		return lNewArticle, parseErr
	}
	return lNewArticle, callErr
}
