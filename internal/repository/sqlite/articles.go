package sqlite

import (
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"gorm.io/gorm"
	"time"
)

type SqliteStorage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) *SqliteStorage {
	return &SqliteStorage{db}
}

// NewTable create table articles with field like struct core.ArticleDB
func NewTable(db *gorm.DB) {
	lErr := db.AutoMigrate(&core.ArticleDB{})
	if lErr != nil {
		panic("failed to migrate from Article struct")
	}
}

func (s *SqliteStorage) WriteArticles(lArticles []core.Article) error {
	for _, lArticle := range lArticles {
		s.db.Create(&core.ArticleDB{
			Title:       lArticle.Title,
			Author:      lArticle.Author,
			Link:        lArticle.Link,
			PublishDate: lArticle.PublishDate,
			Description: lArticle.Description})
	}

	return nil
}

func (s *SqliteStorage) ReadArticleByID(ID uint) (core.ArticleDB, error) {
	var lArticle core.ArticleDB
	//TODO Wrap errors
	s.db.Where("id = ?", ID).Find(&lArticle)

	return lArticle, nil
}

func (s *SqliteStorage) ReadArticlesByDateRange(aMin, aMax time.Time) ([]core.ArticleDB, error) {
	lArticles := make([]core.ArticleDB, 0)
	//TODO Wrap errors
	s.db.Where("publish_date BETWEEN ? AND ?", aMin, aMax).Find(&lArticles)

	return lArticles, nil
}

func (s *SqliteStorage) UpdateArticles(aArticles []core.ArticleDB) error {
	return nil
}

func (s *SqliteStorage) AddOneArticle(aArticle core.ArticleDB) error {
	return nil
}
