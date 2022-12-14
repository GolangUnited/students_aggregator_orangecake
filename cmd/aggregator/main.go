package main

import (
	"fmt"
	"os"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"github.com/indikator/aggregator_orange_cake/pkg/handlers"
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
		return nil //handlers.NewHashnodeScraper(&log.Logger{} /*TODO change log.logger to core.logger*/, url)
	},
	"appliedgo": func(url string, logger *core.Logger) interface{} {
		return handlers.NewAppliedGoScraper(url, *logger)
	},
	"golangorg": func(url string, logger *core.Logger) interface{} {
		return handlers.NewGolangOrgHandler(url /*TODO add logger to constructor*/)
	},
}

func main() {
	var logger core.Logger = core.NewZeroLogger(os.Stdout)
	logger.Info("Starting the aggregator's work.")

	lConfig, err := NewAggregatorConfig()
	if err != nil {
		fmt.Println(err)
	}

	scrappers := CreateScrappers(lConfig, nil)

	lArticles, lWarnings := GetArticles(scrappers, nil)

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
	fmt.Println("----WARNINGS----")
	fmt.Println(lWarnings)
}

func CreateScrappers(aConfig *AggregatorConfig, aLogger *core.Logger) []core.ArticleScraper {
	lScrappers := make([]core.ArticleScraper, 0)
	for _, value := range aConfig.Handlers {
		scrapper := ScrappersMap[value.Handler](value.URL, aLogger)
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
