package main

import (
	"fmt"

	_ "github.com/PuerkitoBio/goquery" // temporary import to init go.mod and go.sum and avoid compiler errors
	"github.com/indikator/aggregator_orange_cake/pkg/core"
)

func main() {
	var lArticles []core.Article

	for i, lArticle := range lArticles {
		fmt.Printf("Node %d: %s\n", i, lArticle.Title)
		fmt.Printf("  Author: %s\n", lArticle.Author)
		fmt.Printf("  Date: %s\n", lArticle.PublishDate.Format("Jan _2, 2006"))
		fmt.Printf("  URL: %s\n", lArticle.Link)
		fmt.Printf("  Description:\n    %s\n\n", lArticle.Description)
	}

	fmt.Println()
	fmt.Println()
	fmt.Printf("  %d Articles detected.", len(lArticles))
	fmt.Println()
}
