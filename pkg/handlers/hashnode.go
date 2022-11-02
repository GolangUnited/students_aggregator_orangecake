package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/indikator/aggregator_orange_cake/pkg/core"

	"github.com/gocolly/colly"
)

type HashnodeScraper struct {
	Articles      []core.Article
	URL           string
	Log           *log.Logger
	BadArticle    bool
	ArticlesFound int
}

const (
	HASHNODE_URL  = "https://hashnode.com/n/go"
	ARTICLE_CLASS = "div.css-4gdbui"
	TITLE_PATH    = "div.css-1wg9be8 div.css-16fbhyp h1.css-1j1qyv3 a.css-4zleql"
	DESCR_PATH    = "div.css-1wg9be8 div.css-16fbhyp p.css-1072ocs a.css-4zleql"
	AUTHOR_PATH   = "div.css-dxz0om div.css-tel74u div.css-2wkyxu div.css-1ajtyzd a.css-c3r4j7"
	DATE_PATH     = "div.css-dxz0om div.css-tel74u div.css-2wkyxu div.css-1n08q4e a.css-1u6dh35"
)

// create Hashnode scrapper struct for "https://hashnode.com/n/go"
func NewHashnodeScraper(log *log.Logger) *HashnodeScraper {
	return &HashnodeScraper{
		Articles:      []core.Article{},
		URL:           HASHNODE_URL,
		Log:           log,
		ArticlesFound: 0,
	}
}

// TODO: errors and log messages will be replaced
// srappin url
func (h *HashnodeScraper) ScrapUrl() error {

	lC := colly.NewCollector()

	lC.OnHTML(ARTICLE_CLASS, h.ElementSearch)

	err := lC.Visit(h.URL)
	if err != nil {
		return fmt.Errorf("visit error %w", err)
	}

	lC.Wait()

	if h.ArticlesFound == 0 {
		return fmt.Errorf("unable to find articles")
	}

	if len(h.Articles) == 0 {
		return fmt.Errorf("no correct articles")
	}

	return nil
}

// TODO: errors and log messages will be replaced
// colly searching func
func (h *HashnodeScraper) ElementSearch(el *colly.HTMLElement) {

	h.ArticlesFound++
	h.BadArticle = false

	lArticle := core.Article{}
	lDOM := el.DOM

	lTitle := lDOM.Find(TITLE_PATH)

	if lTitle.Nodes == nil {
		h.BadArticle = true
		lArticle.Title = "unable to find field"
	} else {
		lArticle.Title = lTitle.Text()
		if lTitle.Text() == "" {
			h.BadArticle = true
			lArticle.Title = "critical field is empty"
		}

		lLink, _ := lTitle.Attr("href")
		lArticle.Link = lLink
		if lLink == "" {
			h.BadArticle = true
			lArticle.Link = "critical field is empty"
		}
	}

	lDescription := lDOM.Find(DESCR_PATH)
	if lDescription.Nodes == nil {
		h.BadArticle = true
		lArticle.Description = "unable to find field"
	} else {
		lArticle.Description = lDescription.Text()
	}

	lAuthor := lDOM.Find(AUTHOR_PATH)
	if lAuthor.Nodes == nil {
		h.BadArticle = true
		lArticle.Author = "unable to find field"
	} else {
		lArticle.Author = lAuthor.Text()
	}

	lDate := lDOM.Find(DATE_PATH)
	if lDate.Nodes == nil {
		h.BadArticle = true
		log.Println("unable to find field ")
		lArticle.PublishDate = core.NormalizeDate(time.Now())
	} else {
		lPubDate, err := core.ParseDate("Jan _2, 2006", lDate.Text())
		if err != nil {
			h.Log.Printf("For article %s, %s DataErr: %s ", lArticle.Title, lArticle.Link, err.Error())
		}
		lArticle.PublishDate = lPubDate
	}

	if h.BadArticle {
		log.Printf("Bad Article: %#v", lArticle)
		return
	}

	h.Articles = append(h.Articles, lArticle)
}
