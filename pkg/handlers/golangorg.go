package handlers

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	GOLANG_ORG_URL = "https://tip.golang.org"
)

type GolangOrgHandler struct {
	url      string
	log      core.Logger
	articles []core.Article
	warnings []core.Warning
}

func NewGolangOrgHandler(aUrl string, aLogger core.Logger) GolangOrgHandler {
	return GolangOrgHandler{
		url:      aUrl,
		log:      aLogger,
		articles: make([]core.Article, 0),
		warnings: make([]core.Warning, 0),
	}
}

type golangOrgParser struct {
	article  core.Article
	warnings []core.Warning
}

func newGolangOrgParser() golangOrgParser {
	return golangOrgParser{
		article:  core.Article{},
		warnings: make([]core.Warning, 0),
	}
}

func (p *golangOrgParser) parseTitleLink(aSelection *goquery.Selection) error {
	// TODO: replace errors
	if aSelection.Nodes == nil {
		return core.RequiredFieldError{ErrorType: core.ErrNodeNotFound, Field: core.TitleFieldName}
	}

	lTitle := aSelection.Find("a[href]").Text()
	if len(lTitle) == 0 {
		return core.RequiredFieldError{ErrorType: core.ErrFieldIsEmpty, Field: core.TitleFieldName}
	}

	lUrl, lOk := aSelection.Find("a").Attr("href")
	if !lOk {
		return core.RequiredFieldError{ErrorType: core.ErrAttributeNotExists, Field: core.LinkFieldName}
	}

	lLink := strings.TrimSpace(lUrl)
	if len(lLink) == 0 {
		return core.RequiredFieldError{ErrorType: core.ErrFieldIsEmpty, Field: core.LinkFieldName}
	}

	lLink, lErr := resolveGolangOrgURL(GOLANG_ORG_URL, lLink)
	if lErr != nil {
		return errors.New("cannot resolve url")
	}

	p.article.Title = lTitle
	p.article.Link = lLink

	return nil
}

func (g *golangOrgParser) parseAuthor(aSelection *goquery.Selection) {
	lAuthor := ""
	// TODO: replace warnings
	if len(aSelection.Nodes) != 0 {
		lAuthor = strings.TrimSpace(aSelection.Find("span.author").Text())
		if len(lAuthor) == 0 {
			g.addWarning("article's author is empty")
		} else {
			g.article.Author = lAuthor
		}
	} else {
		g.addWarning("article's author node not found")
	}
}

func (g *golangOrgParser) parseDate(aSelection *goquery.Selection) {
	// TODO: replace warnings
	lDate := core.NormalizeDate(time.Now())
	var lErr error
	if aSelection.Nodes != nil {
		lDateStr := aSelection.Find("span.date").Text()
		lDateStr = strings.TrimSpace(lDateStr)
		lDate, lErr = core.ParseDate("_2 January 2006", lDateStr)
		if lErr != nil {
			g.addWarning(core.Warning(fmt.Sprintf("cannot parse article date '%s'. %s", lDateStr, lErr.Error())))
		}
	} else {
		g.addWarning("article date node not found")
	}

	g.article.PublishDate = lDate
}

func (g *golangOrgParser) parseDescription(aSelection *goquery.Selection) {
	// TODO: replace warnings
	g.article.Description = ""
	lDesc := aSelection.Next()
	if len(lDesc.Nodes) == 1 {
		lClass, lOk := lDesc.Attr("class")
		if lOk && lClass == "blogsummary" {
			g.article.Description = strings.TrimSpace(lDesc.Text())
		} else {
			g.addWarning("article description node not found")
		}
	}
	if g.article.Description == "" {
		g.addWarning("article description is empty")
	}
}

func (g GolangOrgHandler) ParseArticles() (aArticles []core.Article, aWarnings []core.Warning, aError error) {

	lResp, lErr := http.Get(g.url)
	if lErr != nil {
		// TODO: write error to log fmt.Errorf("http request returns an error: %w", lErr)
		return nil, nil, lErr
	}

	defer lResp.Body.Close()

	if lResp.StatusCode != 200 {
		lErr := fmt.Sprintf("Status code: %d %s", lResp.StatusCode, lResp.Status)
		return nil, nil, errors.New(lErr)
	}

	return g.GolangOrgScraper(lResp.Body)
}

func (p *golangOrgParser) parseArticle(aSelection *goquery.Selection) error {
	p.article = core.Article{}

	if lErr := p.parseTitleLink(aSelection); lErr != nil {
		return lErr
	}

	p.parseDate(aSelection)
	p.parseAuthor(aSelection)
	p.parseDescription(aSelection)

	return nil
}

// GolangOrgScraper takes data from tip.golang.com/blog/all and converts it into a structure of json.
func (g *GolangOrgHandler) GolangOrgScraper(aHtmlReader io.Reader) ([]core.Article, []core.Warning, error) {

	lArticles := make([]core.Article, 0)
	lWarnings := make([]core.Warning, 0)

	lDoc, lErr := goquery.NewDocumentFromReader(aHtmlReader)
	if lErr != nil {
		//TODO: write error to log fmt.Errorf("goquery.NewDocumentFromReader returns an error: %w", lErr)
		return nil, nil, lErr
	}

	// doc.Find("p.blogtitle").Each(func(aIndex int, aSelection *goquery.Selection) {
	lDoc.Find("p.blogtitle").Each(func(aIndex int, aSelection *goquery.Selection) {
		lParser := newGolangOrgParser()
		lErr := lParser.parseArticle(aSelection)
		if lErr == nil {
			lArticles = append(lArticles, lParser.article)
			if len(lParser.warnings) > 0 {
				for i, lWarning := range lParser.warnings {
					strWarning := fmt.Sprintf("Warning[%d,%d]: %s", aIndex, i, lWarning)
					g.log.Info(strWarning)
					lWarnings = append(lWarnings, core.Warning(strWarning))
				}
			}
		} else {
			strError := fmt.Sprintf("Error[%d]: %s", aIndex, lErr.Error())
			g.log.Warn(strError)
			lWarnings = append(lWarnings, core.Warning(strError))
		}
	})

	return lArticles, lWarnings, nil

}

func (g *golangOrgParser) addWarning(aWarning core.Warning) {
	g.warnings = append(g.warnings, aWarning)
}

func resolveGolangOrgURL(aBaseUrl, aPath string) (string, error) {
	lUrl, lErr := url.Parse(aPath)
	if lErr != nil {
		return "", lErr
	}

	lBase, lErr := url.Parse(aBaseUrl)
	if lErr != nil {
		return "", lErr
	}

	return lBase.ResolveReference(lUrl).String(), nil
}
