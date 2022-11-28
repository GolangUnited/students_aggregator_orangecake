package handlers

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
)

type HashnodeScraper struct {
	articles []core.Article
	url      string
	log      *log.Logger
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
		articles: []core.Article{},
		url:      aUrl,
		log:      log,
	}
}

// TODO: log will be replaced
// srappin url
func (h *HashnodeScraper) ScrapUrl() error {

	lC := colly.NewCollector()

	lC.OnHTML(ARTICLE_CLASS, h.ElementSearch)

	err := lC.Visit(h.url)
	if err != nil {
		return fmt.Errorf("%s: %w", core.ErrUrlVisit.Error(), err)
	}

	lC.Wait()

	if len(h.articles) == 0 {
		return core.ErrNoArticles
	}

	return nil
}

// TODO: log will be replaced
// colly searching func
func (h *HashnodeScraper) ElementSearch(el *colly.HTMLElement) {

	lArticle := core.Article{}
	lDOM := el.DOM

	lTitle := lDOM.Find(TITLE_PATH)

	//link inside Title field
	if lTitle.Nodes == nil {
		h.log.Println("cant find Title(link) field")
		return
	}

	lArticle.Title = lTitle.Text()
	lArticle.Link, _ = lTitle.Attr("href")

	lDescription := lDOM.Find(DESCR_PATH)
	if lDescription.Nodes == nil {
		h.log.Printf("For article %s, %s field Description not found", lArticle.Title, lArticle.Link)
	} else {
		lArticle.Description = lDescription.Text()
	}

	lAuthor := lDOM.Find(AUTHOR_PATH)
	if lAuthor.Nodes == nil {
		h.log.Printf("For article %s, %s field Author not found", lArticle.Title, lArticle.Link)
	} else {
		lArticle.Author = lAuthor.Text()
	}

	lDate := lDOM.Find(DATE_PATH)

	var lDateString string
	if lDate.Nodes == nil {
		h.log.Printf("For article %s, %s field Data not found", lArticle.Title, lArticle.Link)
	} else {
		lDateString = lDate.Text()
	}

	lPubDate, err := core.ParseDate("Jan _2, 2006", lDateString)
	if err != nil {
		h.log.Printf("For article %s, %s DataErr: %s ", lArticle.Title, lArticle.Link, err.Error())
	}
	lArticle.PublishDate = lPubDate

	if lArticle.Title == "" {
		h.log.Printf("For article %s Title field is empty", lArticle.Link)
		return
	}

	if lArticle.Link == "" {
		h.log.Printf("For article %s Link field is empty", lArticle.Title)
		return
	}

	h.articles = append(h.articles, lArticle)
}
