package handlers

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
)

type HashnodeScraper struct {
	articles []core.Article
	warnings []core.Warning
	url      string
	log      core.Logger
}

const (
	HASHNODE_URL  = "https://hashnode.com/n/go"
	ARTICLE_CLASS = "div.css-4gdbui"
	TITLE_PATH    = "div.css-1wg9be8 div.css-16fbhyp h1.css-1j1qyv3 a.css-4zleql"
	DESCR_PATH    = "div.css-1wg9be8 div.css-16fbhyp p.css-1072ocs a.css-4zleql"
	AUTHOR_PATH   = "div.css-dxz0om div.css-tel74u div.css-2wkyxu div.css-1ajtyzd a.css-c3r4j7"
	DATE_PATH     = "div.css-dxz0om div.css-tel74u div.css-2wkyxu div.css-1n08q4e a.css-1u6dh35"
)

// NewHashnodeScraper create Hashnode scrapper struct for "https://hashnode.com/n/go"
func NewHashnodeScraper(aUrl string, logger core.Logger) *HashnodeScraper {
	return &HashnodeScraper{
		articles: []core.Article{},
		warnings: []core.Warning{},
		url:      aUrl,
		log:      logger,
	}
}

// ScrapUrl scrapping url
func (h *HashnodeScraper) scrapUrl() error {

	lC := colly.NewCollector()

	lC.OnHTML(ARTICLE_CLASS, h.elementSearch)

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

// ElementSearch colly searching func
func (h *HashnodeScraper) elementSearch(el *colly.HTMLElement) {

	lArticle := core.Article{}
	lDOM := el.DOM

	lTitle := lDOM.Find(TITLE_PATH)

	//link inside Title field
	if lTitle.Nodes == nil {
		strError := core.RequiredFieldError{ErrorType: core.ErrNodeNotFound, Field: core.TitleFieldName}.Error()
		h.log.Warn(strError)
		h.warnings = append(h.warnings, core.Warning(strError))
		return
	}

	lArticle.Title = lTitle.Text()
	lArticle.Link, _ = lTitle.Attr("href")

	lDescription := lDOM.Find(DESCR_PATH)
	if lDescription.Nodes == nil {
		strWarning := fmt.Sprintf("For article %s, %s field Description not found", lArticle.Title, lArticle.Link)
		h.log.Info(strWarning)
		h.warnings = append(h.warnings, core.Warning(strWarning))
	} else {
		lArticle.Description = lDescription.Text()
	}

	lAuthor := lDOM.Find(AUTHOR_PATH)
	if lAuthor.Nodes == nil {
		strWarning := fmt.Sprintf("For article %s, %s field Author not found", lArticle.Title, lArticle.Link)
		h.log.Info(strWarning)
		h.warnings = append(h.warnings, core.Warning(strWarning))
	} else {
		lArticle.Author = lAuthor.Text()
	}

	lDate := lDOM.Find(DATE_PATH)

	var lDateString string
	if lDate.Nodes == nil {
		strWarning := fmt.Sprintf("For article %s, %s field Data not found", lArticle.Title, lArticle.Link)
		h.log.Info(strWarning)
		h.warnings = append(h.warnings, core.Warning(strWarning))
	} else {
		lDateString = lDate.Text()
	}

	lPubDate, err := core.ParseDate("Jan _2, 2006", lDateString)
	if err != nil {
		strWarning := fmt.Sprintf("For article %s, %s DateErr: %s ", lArticle.Title, lArticle.Link, err.Error())
		h.log.Info(strWarning)
		h.warnings = append(h.warnings, core.Warning(strWarning))
	}
	lArticle.PublishDate = lPubDate

	if lArticle.Title == "" {
		strError := core.RequiredFieldError{ErrorType: core.ErrFieldIsEmpty, Field: core.TitleFieldName}.Error()
		h.log.Warn(strError)
		h.warnings = append(h.warnings, core.Warning(strError))
		return
	}

	if lArticle.Link == "" {
		strError := core.RequiredFieldError{ErrorType: core.ErrFieldIsEmpty, Field: core.LinkFieldName}.Error()
		h.log.Warn(strError)
		h.warnings = append(h.warnings, core.Warning(strError))
		return
	}

	h.articles = append(h.articles, lArticle)
}

// ParseArticles get articles from scrapper
func (h *HashnodeScraper) ParseArticles() ([]core.Article, []core.Warning, error) {
	err := h.scrapUrl()
	if err != nil {
		return nil, h.warnings, err
	}

	return h.articles, h.warnings, nil
}
