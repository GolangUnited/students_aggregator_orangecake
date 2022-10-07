package handlers

import (
	"fmt"
	"time"

	"github.com/indikator/aggregator_orange_cake/pkg/core"

	"github.com/gocolly/colly"
)

type hashnodeScraper struct {
	Articles []core.Article
}

//create Hashnode scrapper struct
func NewHashnodeScraper() *hashnodeScraper {
	return &hashnodeScraper{
		Articles: []core.Article{},
	}
}

//srappin url https://hashnode.com/n/go
func (h *hashnodeScraper) ScrapUrl() error {

	lUrl := "https://hashnode.com/n/go"
	lDomains := []string{"www.hashnode.com", "hashnode.com"}

	lC := colly.NewCollector(colly.AllowedDomains(lDomains...))

	lC.OnHTML("div.css-4gdbui", func(el *colly.HTMLElement) {

		lArticle := core.Article{}

		lDOM := el.DOM

		lTitle := lDOM.Find("div.css-1wg9be8 div.css-16fbhyp h1.css-1j1qyv3 a.css-4zleql")
		lArticle.Title = lTitle.Text()

		lLink, _ := lTitle.Attr("href")
		lArticle.Link = lLink

		lDescription := lDOM.Find("div.css-1wg9be8 div.css-16fbhyp p.css-1072ocs a.css-4zleql")
		lArticle.Description = lDescription.Text()

		lAuthor := lDOM.Find("div.css-dxz0om div.css-tel74u div.css-2wkyxu div.css-1ajtyzd a.css-c3r4j7")
		lArticle.Author = lAuthor.Text()

		lDate := lDOM.Find("div.css-dxz0om div.css-tel74u div.css-2wkyxu div.css-1n08q4e a.css-1u6dh35")

		lDateTime, err := time.Parse("Jan _2, 2006", lDate.Text())
		if err != nil {
			fmt.Println(err)
			lDateTime = time.Now()
		}
		lArticle.PublishDate = lDateTime

		h.Articles = append(h.Articles, lArticle)

	})

	err := lC.Visit(lUrl)
	if err != nil {
		return fmt.Errorf("visit error %w", err)
	}

	return nil
}
