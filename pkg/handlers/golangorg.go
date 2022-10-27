package handlers

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"net/http"
	"strings"
)

const (
	TGO_SURL = "https://tip.golang.org"
)

// GolangorgScraper takes data from tip.golang.com/blog/all and converts it into a structure of json.
func GolangorgScraper(aURL string) ([]core.Article, error) {

	lArticles := make([]core.Article, 0, 0)

	resp, lErr := http.Get(aURL)
	if lErr != nil {
		// TODO: write error to log
		fmt.Println(errors.New("http request returns an error: "), lErr)
	}

	defer resp.Body.Close()

	if resp.StatusCode > 400 {
		fmt.Println("Status code: ", resp.StatusCode)
	}

	//rewrite the code with goquery, without goColly
	doc, lErr := goquery.NewDocumentFromReader(resp.Body)
	if lErr != nil {
		// TODO: write error to log
		fmt.Println(errors.New("goquery.NewDocumentFromReader returns an error: "), lErr)
	}

	doc.Find("p.blogtitle, p.blogsummary").Each(func(i int, selection *goquery.Selection) {
		if i%2 == 0 {
			lLink, _ := selection.Find("a").Attr("href")
			lT, lErr := core.ParseDate("_2 January 2006", selection.Find("span.date").Text())
			if lErr != nil {
				// TODO: write error to log
				fmt.Println(errors.New("error of ParseDate function: "), lErr)
			}

			lArticle := core.Article{
				Title:       selection.Find("a[href]").Text(),
				Author:      selection.Find("span.author").Text(),
				PublishDate: lT,
				Link:        TGO_SURL + lLink,
				Description: strings.TrimSpace(selection.NextFiltered("p.blogsummary").Text()),
			}

			lArticles = append(lArticles, lArticle)
		}
	})

	return lArticles, nil

}
