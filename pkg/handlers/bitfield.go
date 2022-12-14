package handlers

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"net/url"
	"strings"
	"time"
)

const BITFIELD_URL = "https://bitfieldconsulting.com/golang"

type BitfieldHandler struct {
	url      string
	log      core.Logger
	articles []core.Article
	warnings []core.Warning
}

// NewBitfieldScrapper - return new bitfield scrapper struct
func NewBitfieldScrapper(aUrl string, aLog core.Logger) core.ArticleScraper {
	return &BitfieldHandler{
		url:      aUrl,
		log:      aLog,
		articles: make([]core.Article, 0),
		warnings: make([]core.Warning, 0),
	}
}

type bitfieldParser struct {
	log      core.Logger
	article  core.Article
	warnings []core.Warning
}

func newBitfieldParser(aLog core.Logger) bitfieldParser {
	return bitfieldParser{
		log:      aLog,
		article:  core.Article{},
		warnings: make([]core.Warning, 0),
	}
}

func (b *bitfieldParser) parseTitleLink(selection *goquery.Selection) error {
	if selection.Nodes == nil {
		return core.RequiredFieldError{ErrorType: core.ErrNodeNotFound, Field: core.TitleFieldName}
	}

	lTitle := strings.TrimSpace(selection.First().Clone().Children().Remove().End().Text())
	if len(lTitle) == 0 {
		return core.RequiredFieldError{ErrorType: core.ErrFieldIsEmpty, Field: core.TitleFieldName}
	}

	lUrl, lExists := selection.Attr("href")
	if !lExists {
		return core.RequiredFieldError{ErrorType: core.ErrAttributeNotExists, Field: core.LinkFieldName}
	}

	lLink := strings.TrimSpace(lUrl)
	if len(lLink) == 0 {
		return core.RequiredFieldError{ErrorType: core.ErrFieldIsEmpty, Field: core.LinkFieldName}
	}

	lLink, lErr := resolveURL(BITFIELD_URL, lLink)
	if lErr != nil {
		return errors.New("cannot resolve url")
	}

	b.article.Title = lTitle
	b.article.Link = lLink

	return nil
}

func (b *bitfieldParser) parseDescription(selection *goquery.Selection) {
	lDescription := ""
	// TODO: replace warnings
	if selection.Nodes != nil {
		lDescription = strings.TrimSpace(selection.Text())
		if len(lDescription) == 0 {
			b.addWarning("article description is empty")
		}
	} else {
		b.addWarning("article description node not found")
	}

	b.article.Description = lDescription
}

func (b *bitfieldParser) parseAuthor(selection *goquery.Selection) {
	lAuthor := ""
	// TODO: replace warnings
	if selection.Nodes != nil {
		lAuthor = strings.TrimSpace(selection.Text())
		if len(lAuthor) == 0 {
			b.addWarning("article author is empty")
		}
	} else {
		b.addWarning("article author node not found")
	}

	b.article.Author = lAuthor
}

func (b *bitfieldParser) parseDate(selection *goquery.Selection) {
	// TODO: replace warnings
	lDate := core.NormalizeDate(time.Now())
	var lErr error
	if selection.Nodes != nil {
		lDateStr, lExists := selection.Attr("datetime")
		if !lExists {
			b.addWarning("article date attribute not exists")
		} else {
			lDateStr = strings.TrimSpace(lDateStr)
			lDate, lErr = core.ParseDate("2006-01-02", lDateStr)
			if lErr != nil {
				b.addWarning(core.Warning(fmt.Sprintf("cannot parse article date '%s'. %s", lDateStr, lErr.Error())))
			}
		}
	} else {
		b.addWarning("article date node not found")
	}

	b.article.PublishDate = lDate
}

func (b *bitfieldParser) addWarning(aWarning core.Warning) {
	b.log.Info(string(aWarning))
	b.warnings = append(b.warnings, aWarning)
}

// articlesSearching - searching for article by specific selector
func (b *BitfieldHandler) articlesSearching() error {
	lCollector := colly.NewCollector()

	lCollector.OnHTML("article > div.entry-text", func(element *colly.HTMLElement) {
		lParser := newBitfieldParser(b.log)
		lErr := lParser.parseArticle(element)
		if lErr == nil {
			b.articles = append(b.articles, lParser.article)
			if len(lParser.warnings) > 0 {
				for i, warning := range lParser.warnings {
					b.warnings = append(b.warnings, core.Warning(fmt.Sprintf("Warning[%d,%d]: %s", element.Index, i, warning)))
				}
			}
		}
		if lErr != nil {
			b.log.Warn(lErr.Error())
			b.warnings = append(b.warnings, core.Warning(fmt.Sprintf("Error[%d]: %s", element.Index, lErr.Error())))
		}
	})

	lErr := lCollector.Visit(b.url)
	if lErr != nil {
		return lErr
	}

	return nil
}

// parseArticle - parse found article
func (b *bitfieldParser) parseArticle(element *colly.HTMLElement) error {
	lErr := b.parseTitleLink(element.DOM.Find("a.u-url"))
	if lErr != nil {
		return lErr
	}

	b.parseDate(element.DOM.Find("time.dt-published.date-callout"))
	b.parseDescription(element.DOM.Find("div.p-summary"))
	b.parseAuthor(element.DOM.Find("a.blog-author-name"))

	return nil
}

// ParseArticles - return array of articles
func (b *BitfieldHandler) ParseArticles() ([]core.Article, []core.Warning, error) {
	lErr := b.articlesSearching()
	if lErr != nil {
		return nil, nil, lErr
	}

	return b.articles, b.warnings, nil
}

func resolveURL(baseUrl, path string) (string, error) {
	u, err := url.Parse(path)
	if err != nil {
		return "", err
	}

	base, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}

	return base.ResolveReference(u).String(), nil
}
