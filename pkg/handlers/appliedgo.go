package appliedgoParser

import (
	"encoding/json"
	"github.com/gocolly/colly/v2"
	"strings"
	"sync"
)

type Article struct {
	Name        string
	Date        string
	Description string
	Link        string
}

// Returns the slice of json-marshalled structs "Article{ Name, Description, Date, Link string }"
// containing data about articles from appliedgo.net
func parseAppliedGo() [][]byte {
	var listMutex sync.Mutex
	var articlesList [][]byte

	// This parser is used to scrap the data from every article's personal page
	articleParser := colly.NewCollector()
	articleParser.OnHTML("head", func(e *colly.HTMLElement) {
		name := e.ChildAttr("meta[property=\"og:title\"]", "content")
		descr := e.ChildAttr("meta[name=\"description\"]", "content")
		url := e.ChildAttr("meta[property=\"og:url\"]", "content")
		publicationTimeStr := e.ChildAttr("meta[property=\"article:published_time\"]", "content")
		newArticle := Article{
			Name:        name,
			Description: descr,
			Link:        url,
			Date:        strings.TrimSpace(publicationTimeStr),
		}
		// will add handling the error after getting known how we will handle mistakes
		newArticleJSON, _ := json.Marshal(&newArticle)

		// In case of async usage of gocolly, I put a mutex here to prevent data race
		listMutex.Lock()
		articlesList = append(articlesList, newArticleJSON)
		listMutex.Unlock()
	})

	// This parser is used to scrap all the articles from the page with list of them
	articleCollector := colly.NewCollector()
	articleCollector.OnHTML("header[class=\"article-header\"] > a", func(e *colly.HTMLElement) {
		// hasn't made error handling in func Visit yet, because don't know what to do with this mistakes
		articleParser.Visit(e.Attr("href"))
	})

	articleCollector.Visit("https://appliedgo.net/")
	return articlesList
}
