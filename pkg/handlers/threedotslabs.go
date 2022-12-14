package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
)

const THREE_DOTS_LABS_URL = "https://threedots.tech/"

type ThreeDotsLabsHandler struct {
	url string
	log core.Logger
}

func NewThreeDotsLabsHandler(aUrl string, aLog core.Logger) ThreeDotsLabsHandler {
	return ThreeDotsLabsHandler{
		url: aUrl,
		log: aLog,
	}
}

type threeDotsLabsParser struct {
	Article  core.Article
	Warnings []core.Warning
}

func newThreeDotsLabsParser() threeDotsLabsParser {
	return threeDotsLabsParser{
		Article:  core.Article{},
		Warnings: make([]core.Warning, 0),
	}
}

func (p *threeDotsLabsParser) addWarning(aWarning core.Warning) {
	p.Warnings = append(p.Warnings, aWarning)
}

func (p *threeDotsLabsParser) addWarningf(aFormat string, aArgs ...any) {
	p.addWarning(core.Warning(fmt.Sprintf(aFormat, aArgs...)))
}

func (p *threeDotsLabsParser) parseAuthorAndDateHeader(aNode *goquery.Selection) {
	// Author and PublishDate are optional fields, so we add warning if we cannot import them
	//  - Author can be empty
	//  - PublishDate filled with the default value (current UTC date) by core.ParseDate
	p.Article.Author = ""
	p.Article.PublishDate = core.NormalizeDate(time.Now())

	if len(aNode.Nodes) > 0 {
		lHeader := strings.TrimPrefix(strings.TrimSpace(aNode.Text()), "@")
		if len(lHeader) > 0 {
			lAuthorAndDate := strings.Split(lHeader, "·")
			if len(lAuthorAndDate) == 2 {
				lAuthor := strings.TrimSpace(lAuthorAndDate[0])
				lDateStr := strings.TrimSpace(lAuthorAndDate[1])

				if len(lAuthor) > 0 {
					p.Article.Author = lAuthor
				} else {
					p.addWarning("article Author is empty")
				}

				lDate, lError := core.ParseDate("Jan _2, 2006", lDateStr)
				if lError == nil {
					p.Article.PublishDate = lDate
				} else {
					p.addWarningf("cannot parse article date '%s'. %s", lDateStr, lError.Error())
				}
			} else {
				p.addWarning("invalid article Header format")
			}
		} else {
			p.addWarning("article Header is empty")
		}
	} else {
		p.addWarning("article Header node not found")
	}
}

func (p *threeDotsLabsParser) parseDescription(aNode *goquery.Selection) {
	// lDescription is an optional field, so we add warning if we cannot import it
	lDescription := ""
	if aNode != nil {
		lDescription = strings.TrimSpace(aNode.Text())
		if len(lDescription) <= 0 {
			p.addWarning("article Description is empty")
		}
	}

	p.Article.Description = lDescription
}

func (p *threeDotsLabsParser) parseTitleAndLink(aNode *goquery.Selection) error {
	// Title and Link are both required

	if len(aNode.Nodes) == 0 {
		return core.RequiredFieldError{ErrorType: core.ErrNodeNotFound, Field: core.LinkFieldName}
	}

	lTitle := strings.TrimSpace(aNode.Text())
	if len(lTitle) <= 0 {
		return core.RequiredFieldError{ErrorType: core.ErrFieldIsEmpty, Field: core.TitleFieldName}
	}

	lUrl, lExists := aNode.Attr("href")
	if !lExists {
		return core.RequiredFieldError{ErrorType: core.ErrAttributeNotExists, Field: core.LinkFieldName}
	}

	lLink := strings.TrimSpace(lUrl)
	if len(lLink) <= 0 {
		return core.RequiredFieldError{ErrorType: core.ErrFieldIsEmpty, Field: core.LinkFieldName}
	}

	p.Article.Title = lTitle
	p.Article.Link = lLink
	return nil
}

func (p *threeDotsLabsParser) parseArticle(aNode *goquery.Selection) error {
	p.Article = core.Article{}

	if lErr := p.parseTitleAndLink(aNode.Find("a.post-link").First()); lErr != nil {
		return lErr
	}

	p.parseAuthorAndDateHeader(aNode.Find("p.post-meta").First())
	p.parseDescription(aNode.Find("p.post-summary").First())

	return nil
}

func (h ThreeDotsLabsHandler) ParseHtml(aHtmlReader io.Reader) (aArticle []core.Article, aWarnings []core.Warning, aError error) {
	lHtml, lError := goquery.NewDocumentFromReader(aHtmlReader)
	if lError != nil {
		return nil, nil, lError
	}

	lArticles := make([]core.Article, 0)
	lWarnings := make([]core.Warning, 0)

	lHtml.Find("article.post-entry").Each(func(aIndex int, aNode *goquery.Selection) {
		lParser := newThreeDotsLabsParser()
		lErr := lParser.parseArticle(aNode)
		if lErr == nil {
			lArticles = append(lArticles, lParser.Article)
			if len(lParser.Warnings) > 0 {
				for i, lWarning := range lParser.Warnings {
					strWarning := fmt.Sprintf("Warning[%d,%d]: %s", aIndex, i, lWarning)
					h.log.Info(strWarning)
					lWarnings = append(lWarnings, core.Warning(strWarning))
				}
			}
		} else {
			strError := fmt.Sprintf("Error[%d]: %s", aIndex, lErr.Error())
			h.log.Warn(strError)
			lWarnings = append(lWarnings, core.Warning(strError))
		}
	})

	return lArticles, lWarnings, nil
}

func (h ThreeDotsLabsHandler) ParseArticles() (aArticles []core.Article, aWarnings []core.Warning, aError error) {

	lResponse, lError := http.Get(h.url)
	if lError != nil {
		return nil, nil, lError
	}
	defer lResponse.Body.Close()

	if lResponse.StatusCode != 200 {

		lMsg := core.ResponseError{Status: lResponse.Status, Code: lResponse.StatusCode}
		return nil, nil, lMsg
	}

	return h.ParseHtml(lResponse.Body)
}
