package handlers

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
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
	err := lArticleCollectorTester.Visit("file://./AppliedGoMain.htm")
	if err != nil {
		fmt.Print(err)
	}

	lExpectedLinksList := []string{
		"https://appliedgo.net/rich/",
		"https://appliedgo.net/generictree/",
		"https://appliedgo.net/instantgo/",
		"https://appliedgo.net/mantil/",
		"https://appliedgo.net/auxin/",
	}
	assert.ElementsMatch(t, lExpectedLinksList, lReceivedLinksList)
}
