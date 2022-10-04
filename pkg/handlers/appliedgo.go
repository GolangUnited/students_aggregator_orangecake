package handlers

import (
	"github.com/gocolly/colly/v2"
	"strings"
	"sync"
	"time"
)

type Article struct {
	Caption     string
	PublishDate time.Time
	Description string
	Link        string
	Author      string
}

// fmt.Printf("Error: %s\n\n", lErr.Error())
// Returns the slice of json-marshalled structs "Article{ Name, Description, Date, Link string }"
// containing data about articles from appliedgo.net
func parseAppliedGo() []Article {
	var lArticlesMutex sync.Mutex
	var lArticlesList []Article

	// This parser is used to scrap the data from every article's personal page
	lArticleParser := colly.NewCollector(colly.Async(true))
	lArticleParser.OnHTML("head", func(e *colly.HTMLElement) {
		aName := e.ChildAttr("meta[property=\"og:title\"]", "content")
		aDescr := e.ChildAttr("meta[name=\"description\"]", "content")
		aUrl := e.ChildAttr("meta[property=\"og:url\"]", "content")
		aPublicationDateStr := e.ChildAttr("meta[property=\"article:published_time\"]", "content")
		aPublicationDate := parseDate(strings.TrimSpace(aPublicationDateStr))
		lNewArticle := Article{
			Caption:     aName,
			Description: aDescr,
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
		lArticleParser.Visit(e.Attr("href"))
	})

	lArticleCollector.Visit("https://appliedgo.net/")
	lArticleCollector.Wait()
	lArticleParser.Wait()
	return lArticlesList
}

func parseDate(dateStr string) time.Time {
	lDate, lErr := time.Parse(time.RFC3339, dateStr)
	if lErr != nil {
		lDate = time.Now()
	}

	return lDate.UTC()
}
