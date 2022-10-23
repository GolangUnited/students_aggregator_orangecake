package handlers

import (
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
		lErrorsList = append(lErrorsList, lErr)
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

func ParseAppliedGoArticle(link string) (core.Article, error) {
	var lNewArticle core.Article
	var lCallErr, lParseErr error
	var lArticleParser = colly.NewCollector()
	lArticleParser.OnHTML("head", func(e *colly.HTMLElement) {
		aName := e.ChildAttr("meta[property=\"og:title\"]", "content")
		aDescription := e.ChildAttr("meta[property=\"og:description\"]", "content")
		aUrl := e.ChildAttr("meta[property=\"og:url\"]", "content")
		aPublicationDateStr := e.ChildAttr("meta[property=\"article:published_time\"]", "content")
		aPublicationDate, err := core.ParseDate(time.RFC3339, strings.TrimSpace(aPublicationDateStr))
		if err != nil {
			lParseErr = err
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
	_, lCallErr = http.Get(link)
	lArticleParser.Visit(link)
	if lCallErr == nil {
		return lNewArticle, lParseErr
	}
	return lNewArticle, lCallErr
}

//func ParseTitle(aTitle string) {
//
//}
