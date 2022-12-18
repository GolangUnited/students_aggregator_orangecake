package main

import (
	"fmt"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"github.com/indikator/aggregator_orange_cake/pkg/handlers"
	"os"
)

type ScrapperConstruct func(string, core.Logger) core.ArticleScraper

var ScrappersMap = map[string]ScrapperConstruct{
	"devto": func(url string, logger core.Logger) core.ArticleScraper {
		return nil //handlers.NewDevtoHandler(url /*TODO add logger to constructor*/)
	},
	"bitfield": func(url string, logger core.Logger) core.ArticleScraper {
		return handlers.NewBitfieldScrapper(url, logger)
	},
	"threedotslabs": func(url string, logger core.Logger) core.ArticleScraper {
		return handlers.NewThreeDotsLabsHandler(url, logger)
	},
	"hashnode": func(url string, logger core.Logger) core.ArticleScraper {
		return handlers.NewHashnodeScraper(url, logger)
	},
	"appliedgo": func(url string, logger core.Logger) core.ArticleScraper {
		return handlers.NewAppliedGoScraper(url, logger)
	},
	"golangorg": func(url string, logger core.Logger) core.ArticleScraper {
		return handlers.NewGolangOrgHandler(url, logger)
	},
}

func main() {
	var logger core.Logger = core.NewZeroLogger(os.Stdout)
	logger.Info("Starting the aggregator's work.")

	lConfig, err := NewAggregatorConfig()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	scrappers := CreateScrappers(lConfig, logger)

	lArticles, _ := GetArticles(scrappers, logger)

	for i, lArticle := range lArticles {
		lArticleDescr := fmt.Sprintf("Article %d: %s\n", i, lArticle.Title)
		lArticleDescr += fmt.Sprintf("  Author: %s\n", lArticle.Author)
		lArticleDescr += fmt.Sprintf("  Date: %s\n", lArticle.PublishDate.Format("Jan _2, 2006"))
		lArticleDescr += fmt.Sprintf("  URL: %s\n", lArticle.Link)
		lArticleDescr += fmt.Sprintf("  Description:\n    %s\n\n", lArticle.Description)

		logger.Info(lArticleDescr)
	}

	logger.Info("  %d Articles detected.\n", len(lArticles))
}

func CreateScrappers(aConfig *AggregatorConfig, aLogger core.Logger) []core.ArticleScraper {
	lScrappers := make([]core.ArticleScraper, 0)
	for _, value := range aConfig.Handlers {
		scrapperConstruct, ok := ScrappersMap[value.Handler]
		if ok {
			scrapper := scrapperConstruct(value.URL, aLogger)
			if scrapper == nil {
				aLogger.Warn("handler: %s error: %s", value, "is nil")
				continue
			}
			lScrappers = append(lScrappers, scrapper)
		}
	}

	return lScrappers
}

func GetArticles(aScrappers []core.ArticleScraper, logger core.Logger) ([]core.Article, []core.Warning) {
	lArticles := make([]core.Article, 0)
	lWarnings := make([]core.Warning, 0)
	for i, scrapper := range aScrappers {
		articles, warnings, err := scrapper.ParseArticles()
		if err != nil {
			logger.Warn("[%d] scrapper error: %s", i, err.Error())
			continue
		}
		lArticles = append(lArticles, articles...)
		lWarnings = append(lWarnings, warnings...)
	}

	return lArticles, lWarnings
}
