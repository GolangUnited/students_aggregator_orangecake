package handlers

import (
	"fmt"
	"log"

	"github.com/indikator/aggregator_orange_cake/pkg/core"

	"github.com/gocolly/colly"
)

type hashnodeScraper struct {
	Articles      []core.Article
	URL           string
	Log           *log.Logger
	Err           error
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
func NewHashnodeScraper(log *log.Logger) *hashnodeScraper {
	return &hashnodeScraper{
		Articles:      []core.Article{},
		URL:           HASHNODE_URL,
		Log:           log,
		ArticlesFound: 0,
	}
}

// TODO: errors will be replaced
// srappin url
func (h *hashnodeScraper) ScrapUrl() error {

	lC := colly.NewCollector()

	lC.OnHTML(ARTICLE_CLASS, h.ElementSearch)

	err := lC.Visit(h.URL)
	if err != nil {
		return fmt.Errorf("visit error %w", err)
	}

	lC.Wait()

	if h.Err != nil {
		return h.Err
	}

	if h.ArticlesFound == 0 {
		return fmt.Errorf("unable to find articles")
	}

	if len(h.Articles) == 0 {
		return fmt.Errorf("no correct articles")
	}

	return nil
}

// TODO: errors will be replaced
// colly searching func
func (h *hashnodeScraper) ElementSearch(el *colly.HTMLElement) {

	//articles counter
	h.ArticlesFound++

	//Handler already have critical error. Skip further crawl
	if h.Err != nil {
		return
	}

	lArticle := core.Article{}
	lDOM := el.DOM

	lTitle := lDOM.Find(TITLE_PATH)
	if lTitle.Nodes == nil {
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

	lDescription := lDOM.Find(DESCR_PATH)
	if lDescription.Nodes == nil {
		h.Err = fmt.Errorf("unable to find required field")
		h.Log.Println("unable to find required field")
		return
	}

	lArticle.Description = lDescription.Text()

	lAuthor := lDOM.Find(AUTHOR_PATH)
	if lAuthor.Nodes == nil {
		h.Err = fmt.Errorf("unable to find required field")
		h.Log.Println("unable to find required field")
		return
	}

	lArticle.Author = lAuthor.Text()

	lDate := lDOM.Find(DATE_PATH)
	if lDate.Nodes == nil {
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
