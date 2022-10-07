package core

import "time"

//scrappin url
type Scrapper interface {
	GetArticles() ([]Article, error)
}

//reader for data store
type DSReader interface {
	FindLastArticles(endDate time.Time) ([]Article, error)
	FindArticle(Article) (bool, error)
}

//writer for data store
type DSWriter interface {
	WriteArticles([]Article) error
}
