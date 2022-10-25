package handlers

import (
	"fmt"
	"log"

	"github.com/indikator/aggregator_orange_cake/pkg/core"

	"github.com/gocolly/colly"
)

type hashnodeScraper struct {
	Articles []core.Article
	URL      string
	Log      *log.Logger
	Err      error
}

const HASHNODE_URL = "https://hashnode.com/n/go"

// create Hashnode scrapper struct for "https://hashnode.com/n/go"
func NewHashnodeScraper(log *log.Logger) *hashnodeScraper {
	return &hashnodeScraper{
		Articles: []core.Article{},
		URL:      HASHNODE_URL,
		Log:      log,
	}
}

// srappin url
func (h *hashnodeScraper) ScrapUrl() error {

	lC := colly.NewCollector()

	lC.OnHTML("div.css-4gdbui", h.ElementSearch)

	err := lC.Visit(h.URL)
	if err != nil {
		return fmt.Errorf("visit error %w", err)
	}

	lC.Wait()

	if h.Err != nil {
		return h.Err
	}

	if len(h.Articles) == 0 {
		return fmt.Errorf("no correct articles")
	}

	return nil
}

func (h *hashnodeScraper) ElementSearch(el *colly.HTMLElement) {

	//Handler already have critical error. Skip further crawl
	if h.Err != nil {
		return
	}

	lArticle := core.Article{}
	lDOM := el.DOM

	lTitle := lDOM.Find("div.css-1wg9be8 div.css-16fbhyp h1.css-1j1qyv3 a.css-4zleql")
	if lTitle.Length() == 0 {
		h.Err = fmt.Errorf("unable to find required field")
		h.Log.Println("unable to find required field")
		return
	}

	lArticle.Title = lTitle.Text()
	if lTitle.Text() == "" {
		h.Log.Println("Critical field is empty")
		return
	}

	lLink, _ := lTitle.Attr("href")
	lArticle.Link = lLink
	if lLink == "" {
		h.Log.Println("Critical field is empty")
		return
	}

	lDescription := lDOM.Find("div.css-1wg9be8 div.css-16fbhyp p.css-1072ocs a.css-4zleql")
	if lDescription.Length() == 0 {
		h.Err = fmt.Errorf("unable to find required field")
		h.Log.Println("unable to find required field")
		return
	}

	lArticle.Description = lDescription.Text()

	lAuthor := lDOM.Find("div.css-dxz0om div.css-tel74u div.css-2wkyxu div.css-1ajtyzd a.css-c3r4j7")
	if lAuthor.Length() == 0 {
		h.Err = fmt.Errorf("unable to find required field")
		h.Log.Println("unable to find required field")
		return
	}

	lArticle.Author = lAuthor.Text()

	lDate := lDOM.Find("div.css-dxz0om div.css-tel74u div.css-2wkyxu div.css-1n08q4e a.css-1u6dh35")
	if lTitle.Length() == 0 {
		h.Err = fmt.Errorf("unable to find required field")
		h.Log.Println("unable to find required field")
		return
	}

	lPubDate, err := core.ParseDate("Jan _2, 2006", lDate.Text())
	if err != nil {
		h.Log.Printf("For article %s, %s DataErr: %s ", lArticle.Title, lArticle.Link, err.Error())
	}
	lArticle.PublishDate = lPubDate

	h.Articles = append(h.Articles, lArticle)

}
