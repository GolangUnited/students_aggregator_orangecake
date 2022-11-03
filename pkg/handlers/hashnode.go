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
	badArticle    bool
	articlesFound int
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
func NewHashnodeScraper(log *log.Logger, aUrl string) *HashnodeScraper {
	return &HashnodeScraper{
		Articles:      []core.Article{},
		URL:           aUrl,
		Log:           log,
		articlesFound: 0,
	}
}

// TODO: log will be replaced
// srappin url
func (h *HashnodeScraper) ScrapUrl() error {

	lC := colly.NewCollector()

	lC.OnHTML(ARTICLE_CLASS, h.ElementSearch)

	err := lC.Visit(h.URL)
	if err != nil {
		return fmt.Errorf("%s: %w", core.ErrUrlVisit.Error(), err)
	}

	lC.Wait()

	if h.articlesFound == 0 {
		return core.ErrArticlesNotFound
	}

	if len(h.Articles) == 0 {
		return core.ErrNoArticles
	}
	return nil
}

// TODO: log will be replaced
// colly searching func
func (h *HashnodeScraper) ElementSearch(el *colly.HTMLElement) {

	h.articlesFound++
	h.badArticle = false

	lArticle := core.Article{}
	lDOM := el.DOM

	lTitle := lDOM.Find(TITLE_PATH)

	if lTitle.Nodes == nil {
		h.badArticle = true
		lArticle.Title = "unable to find field"
	} else {
		lArticle.Title = lTitle.Text()
		if lTitle.Text() == "" {
			h.badArticle = true
			lArticle.Title = "critical field is empty"
		}

		lLink, _ := lTitle.Attr("href")
		lArticle.Link = lLink
		if lLink == "" {
			h.badArticle = true
			lArticle.Link = "critical field is empty"
		}
	}

	lDescription := lDOM.Find(DESCR_PATH)
	if lDescription.Nodes == nil {
		h.badArticle = true
		lArticle.Description = "unable to find field"
	} else {
		lArticle.Description = lDescription.Text()
	}

	lAuthor := lDOM.Find(AUTHOR_PATH)
	if lAuthor.Nodes == nil {
		h.badArticle = true
		lArticle.Author = "unable to find field"
	} else {
		lArticle.Author = lAuthor.Text()
	}

	lDate := lDOM.Find(DATE_PATH)
	if lDate.Nodes == nil {
		h.badArticle = true
		h.Log.Printf("For article %s, %s field Data not found", lArticle.Title, lArticle.Link)
		lArticle.PublishDate = core.NormalizeDate(time.Now())
	} else {
		lPubDate, err := core.ParseDate("Jan _2, 2006", lDate.Text())
		if err != nil {
			h.Log.Printf("For article %s, %s DataErr: %s ", lArticle.Title, lArticle.Link, err.Error())
		}
		lArticle.PublishDate = lPubDate
	}

	if h.badArticle {
		h.Log.Printf("Bad Article: %#v", lArticle)
		return
	}

	h.Articles = append(h.Articles, lArticle)
}
