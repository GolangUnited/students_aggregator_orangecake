package handlers

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"strings"
	"time"
)

const (
	TGO_DOMAIN     = "tip.golang.org"
	TGO_WWW_DOMAIN = "www.tip.golang.org"
	TGO_FULL_URL   = "https://tip.golang.org/blog/all"
)

// Scraping takes data from tip.golang.com/blog/all and converts it into a structure of json.
func Scraping(aURL string) ([]core.Article, error) {

	lArticles := make([]core.Article, 0, 0)
	lDescriptions := make([]string, 0, 0)

	//Create a new collector
	lC := colly.NewCollector()

	// Call the OnHTML method to take needed data from website
	lC.OnHTML("p.blogtitle", func(h *colly.HTMLElement) {
		lT, lErr := time.Parse("_2 January 2006", h.ChildText("span.date"))

		if lErr != nil {
			fmt.Println("date can't be formatted: ", lErr)
			return
		}

		article := core.Article{
			Link:        TGO_FULL_URL + h.ChildAttr("a", "href"),
			Title:       h.ChildText("a[href]"),
			PublishDate: lT.UTC(),
			Author:      h.ChildText("span.author"),
		}
		lArticles = append(lArticles, article)
	})

	// Call the OnHTML method secondly with another selector, because I couldn't use two different selectors in OnHTML func
	lC.OnHTML("p.blogsummary", func(h *colly.HTMLElement) {
		lDescription := strings.TrimSpace(h.Text)
		lDescriptions = append(lDescriptions, lDescription)
	})

	// Before making a request print "Visiting"
	lC.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on URL
	lErr := lC.Visit(aURL)
	if lErr != nil {
		return nil, lErr
	}

	// Add value to summary field into the Article struct
	for i := 0; i <= len(lArticles)-1; i++ {
		lArticles[i].Description = lDescriptions[i]
	}

	return lArticles, nil

}
