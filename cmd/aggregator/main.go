package main

import (
	"fmt"
	_ "github.com/PuerkitoBio/goquery" // temporary import to init go.mod and go.sum and avoid compiler errors
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"github.com/indikator/aggregator_orange_cake/pkg/handlers"
	"log"
)

type ScrapperConstruct func(string, *core.Logger) interface{}

var ScrappersMap = map[string]ScrapperConstruct{
	"devto": func(url string, logger *core.Logger) interface{} {
		return handlers.NewDevtoHandler(url /*TODO add logger to constructor*/)
	},
	"bitfield": func(url string, logger *core.Logger) interface{} {
		return handlers.NewBitfieldScrapper(url /*TODO add logger to constructor*/)
	},
	"threedotslabs": func(url string, logger *core.Logger) interface{} {
		return handlers.NewThreeDotsLabsHandler(url /*TODO add logger to constructor*/)
	},
	"hashnode": func(url string, logger *core.Logger) interface{} {
		return handlers.NewHashnodeScraper(&log.Logger{} /*TODO change log.logger to core.logger*/, url)
	},
	"appliedgo": func(url string, logger *core.Logger) interface{} {
		return handlers.NewAppliedGoParser( /*TODO add url to constructor*/ /*TODO add logger to constructor*/ )
	},
	"golangorg": func(url string, logger *core.Logger) interface{} {
		return handlers.NewGolangOrgHandler(url /*TODO add logger to constructor*/)
	},
}

func main() {
	lConfig, err := NewAggregatorConfig()
	if err != nil {
		fmt.Println(err)
	}

	scrappers := CreateScrappers(lConfig, nil)

	articles, warnings := GetArticles(scrappers, nil)

	for _, article := range articles {
		fmt.Printf("Node %s:\n", article.Title)
		fmt.Printf("  Author: %s\n", article.Author)
		fmt.Printf("  Date: %s\n", article.PublishDate.Format("Jan _2, 2006"))
		fmt.Printf("  URL: %s\n", article.Link)
		fmt.Printf("  Description:\n    %s\n\n", article.Description)
	}

	fmt.Println("----WARNINGS----")
	fmt.Println(warnings)
}

func CreateScrappers(aConfig *AggregatorConfig, lLogger *core.Logger) []core.ArticleScraper {
	lScrappers := make([]core.ArticleScraper, 0)
	for _, value := range aConfig.Handlers {
		scrapper := ScrappersMap[value.Handler](value.URL, nil /*TODO pass logger*/)
		result, ok := scrapper.(core.ArticleScraper)
		if !ok {
			//TODO log
			fmt.Println(fmt.Errorf("%s, %w", value.Handler, core.ErrUnableCastToInterface))
			continue
		}
		lScrappers = append(lScrappers, result)
	}

	return lScrappers
}

func GetArticles(aScrappers []core.ArticleScraper, logger *core.Logger) ([]core.Article, []core.Warning) {
	lArticles := make([]core.Article, 0)
	lWarnings := make([]core.Warning, 0)
	for _, scrapper := range aScrappers {
		articles, warnings, err := scrapper.ParseArticles()
		if err != nil {
			//TODO log
			fmt.Println(err)
			continue
		}
		lArticles = append(lArticles, articles...)
		lWarnings = append(lWarnings, warnings...)
	}

	return lArticles, lWarnings
}
