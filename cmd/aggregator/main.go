package main

import (
	"fmt"
	_ "github.com/PuerkitoBio/goquery" // temporary import to init go.mod and go.sum and avoid compiler errors
	"github.com/indikator/aggregator_orange_cake/pkg/handlers"
	"github.com/indikator/aggregator_orange_cake/pkg/storage/sqlite"
)

const ConnectionString = "articles.db" // need take it from .env?

func main() {

	// scrape articles from dev.to
	h := handlers.NewDevtoHandler("https://dev.to/t/go")
	lArticles := h.Run()

	// connect to database or create it if not exist
	lStorage, lErr := sqlite.NewSqliteConnection(ConnectionString)
	if lErr != nil {
		fmt.Printf("error of new sqlite connection: %v", lErr)
	}

	// create new table article_dbs if it hasn't already been created
	lStorage.NewTable(lStorage.Db)

	// write articles from handler
	if lErr := lStorage.WriteArticles(lArticles); lErr != nil {
		fmt.Printf("error of write articles: %v", lErr)
	}
}
