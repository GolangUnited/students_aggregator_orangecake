package handlers

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"strings"
	"sync"
	"time"
)

// Returns the slice of json-marshalled structs "Article{ Name, Description, Date, Link string }"
// containing data about articles from appliedgo.net
func parseAppliedGo() []core.Article {
	var lArticlesMutex sync.Mutex
	var lArticlesList []core.Article

	// This parser is used to scrap the data from every article's personal page
	lArticleParser := colly.NewCollector(colly.Async(true))
	lArticleParser.OnHTML("head", func(e *colly.HTMLElement) {
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

		// Because of async usage of gocolly, I put a mutex here to prevent data race
		lArticlesMutex.Lock()
		lArticlesList = append(lArticlesList, lNewArticle)
		lArticlesMutex.Unlock()
	})

	// This parser is used to scrap all the articles from the page with list of them
	lArticleCollector := colly.NewCollector(colly.Async(true))
	lArticleCollector.OnHTML("header[class=\"article-header\"] > a", func(e *colly.HTMLElement) {
		// hasn't made error handling in func Visit yet, because don't know what to do with this mistakes
		lErr := lArticleParser.Visit(e.Attr("href"))
		if lErr != nil {
			fmt.Printf("Error: %s\n\n", lErr.Error())
		}
	})

	lErr := lArticleCollector.Visit("https://appliedgo.net/")
	if lErr != nil {
		fmt.Printf("Error: %s\n\n", lErr.Error())
	}
	lArticleCollector.Wait()
	lArticleParser.Wait()
	return lArticlesList
}

// Takes the string with datetime in 2006-01-02T15:04:05Z07:00 format
// and returns time.Time of this string in UTC format
func parseDateRFC3339(dateStr string) time.Time {
	lDate, lErr := time.Parse(time.RFC3339, dateStr)
	if lErr != nil {
		fmt.Printf("Error: %s\n\n", lErr.Error())
		lDate = time.Now()
	}

	return lDate.UTC()
}
