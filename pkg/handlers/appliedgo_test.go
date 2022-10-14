package handlers

import (
	"github.com/gocolly/colly"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestGetArticlesListAppliedGo(t *testing.T) {
	lTransport := &http.Transport{}
	lTransport.RegisterProtocol("file", http.NewFileTransport(http.Dir("./test_data/")))
	var lReceivedLinksList []string

	lArticleCollectorTester := colly.NewCollector()
	lArticleCollectorTester.WithTransport(lTransport)
	lArticleCollectorTester.OnHTML("header[class=\"article-header\"] > a", func(e *colly.HTMLElement) {
		lReceivedLinksList = append(lReceivedLinksList, e.Attr("href"))
	})
	lArticleCollectorTester.Visit("file://./AppliedGoMain.htm")

	lExpectedLinksList := []string{
		"https://appliedgo.net/rich/",
		"https://appliedgo.net/generictree/",
		"https://appliedgo.net/instantgo/",
		"https://appliedgo.net/mantil/",
		"https://appliedgo.net/auxin/",
	}
	assert.ElementsMatch(t, lExpectedLinksList, lReceivedLinksList)
}

func TestArticleScraping(t *testing.T) {
	var lReceivedData core.Article
	lTransport := &http.Transport{}
	lTransport.RegisterProtocol("file", http.NewFileTransport(http.Dir("./test_data/")))

	lArticleParserTester := colly.NewCollector()
	lArticleParserTester.WithTransport(lTransport)
	lArticleParserTester.OnHTML("head", func(e *colly.HTMLElement) {
		aName := e.ChildAttr("meta[property=\"og:title\"]", "content")
		aDescription := e.ChildAttr("meta[name=\"description\"]", "content")
		aUrl := e.ChildAttr("meta[property=\"og:url\"]", "content")
		aPublicationDateStr := e.ChildAttr("meta[property=\"article:published_time\"]", "content")
		aPublicationDate := parseDateRFC3339(strings.TrimSpace(aPublicationDateStr))
		lNewArticle := core.Article{
			Title:       aName,
			Description: aDescription,
			Link:        aUrl,
			PublishDate: aPublicationDate,
		}
		lReceivedData = lNewArticle
	})
	lArticleParserTester.Visit("file://./AppliedGoArticle.htm")
	lExpectedData := core.Article{
		Title:       "How I used Go to make my radio auto-switch to AUX-IN when a Raspi plays music - Applied Go",
		Author:      "",
		Link:        "https://appliedgo.net/auxin/",
		PublishDate: time.Date(2022, time.August, 20, 0, 0, 0, 0, time.UTC),
		Description: "How Go code detects music output on a Raspberry and switches a 3sixty radio to AUX-IN via Frontier Silicon API",
	}
	assert.Equal(t, lExpectedData, lReceivedData)

}
