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
		articleDescr := fmt.Sprintf("Node %d: %s\n"+
			"  Author: %s\n"+
			"  Date: %s\n"+
			"  URL: %s\n"+
			"  Description:\n    %s\n\n",
			i, lArticle.Title,
			lArticle.Author,
			lArticle.PublishDate.Format("Jan _2, 2006"),
			lArticle.Link,
			lArticle.Description)
		logger.Info(articleDescr)
	}

	logger.Info("")
	logger.Info("")
	logger.Info(fmt.Sprintf("  %d Articles detected.", len(lArticles)))
	logger.Info("")
}
