package main

import (
	"fmt"
	"os"

	_ "github.com/PuerkitoBio/goquery" // temporary import to init go.mod and go.sum and avoid compiler errors
	"github.com/indikator/aggregator_orange_cake/pkg/core"
)

func main() {
	var logger core.Logger = core.NewZeroLogger(os.Stdout)
	logger.Info("Starting the aggregator's work.")

	var lArticles []core.Article

	for i, lArticle := range lArticles {
		lArticleDescr := fmt.Sprintf("Article %d: %s\n", i, lArticle.Title)
		lArticleDescr += fmt.Sprintf("  Author: %s\n", lArticle.Author)
		lArticleDescr += fmt.Sprintf("  Date: %s\n", lArticle.PublishDate.Format("Jan _2, 2006"))
		lArticleDescr += fmt.Sprintf("  URL: %s\n", lArticle.Link)
		lArticleDescr += fmt.Sprintf("  Description:\n    %s\n\n", lArticle.Description)

		logger.Info(lArticleDescr)
	}

	logger.Info("\n\n")
	logger.Info(fmt.Sprintf("  %d Articles detected.\n", len(lArticles)))
}
