package sqlite

import (
	"fmt"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"gorm.io/gorm"
	"time"
)

const (
	LEN_OF_TITLE       = 300
	LEN_OF_AUTHOR      = 200
	LEN_OF_DESCRIPTION = 6000
)

type SqliteStorage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) *SqliteStorage {
	return &SqliteStorage{db}
}

// NewTable create table articles with field like struct core.ArticleDB
func (s *SqliteStorage) NewTable(db *gorm.DB) {
	lErr := db.AutoMigrate(&core.ArticleDB{})
	if lErr != nil {
		panic("failed to migrate from Article struct")
	}
}

func cut(aText string, aLimit int) string {
	lRunes := []rune(aText)
	if len(lRunes) >= aLimit {
		return string(lRunes[:aLimit])
	}
	return aText
}

func (s *SqliteStorage) WriteArticles(aArticles []core.Article) error {
	for _, lArticle := range aArticles {
		//Validation length of fields
		if len(lArticle.Title) > LEN_OF_TITLE {
			lArticle.Title = cut(lArticle.Title, 300)
		}
		if len(lArticle.Author) > LEN_OF_AUTHOR {
			lArticle.Author = cut(lArticle.Author, 200)
		}
		if len(lArticle.Description) > LEN_OF_DESCRIPTION {
			lArticle.Description = cut(lArticle.Description, 6000)
		}

		s.db.Create(&core.ArticleDB{
			Title:       lArticle.Title,
			Author:      lArticle.Author,
			Link:        lArticle.Link,
			PublishDate: lArticle.PublishDate,
			Description: lArticle.Description})
	}

	return nil
}

func (s *SqliteStorage) ReadArticleByID(aID uint) (*core.ArticleDB, error) {
	var lArticle core.ArticleDB

	lResult := s.db.Where("id = ?", aID).First(&lArticle)

	lErr := lResult.Error
	if lErr != nil {
		fmt.Errorf("Error happens in row with id = %d: %w", aID, lErr)
		return nil, lErr
	}

	return &lArticle, nil
}

func (s *SqliteStorage) ReadArticlesByDateRange(aMin, aMax time.Time) ([]core.ArticleDB, error) {
	lArticles := make([]core.ArticleDB, 0)

	lResult := s.db.Where("publish_date BETWEEN ? AND ?", aMin, aMax).Find(&lArticles)
	lErr := lResult.Error
	if lErr != nil {
		fmt.Errorf("Error happens in rows between dates %s and %s: %w", aMin, aMax, lErr)
		return nil, lErr
	}
	fmt.Printf("%d rows was found.", lResult.RowsAffected)
	return lArticles, nil
}

func (s *SqliteStorage) UpdateArticles(aArticles []core.Article) error {
	var lLastWriteDate time.Time
	var lArticleForDate core.ArticleDB

	s.db.Raw("SELECT * FROM article_dbs ORDER BY publish_date DESC LIMIT 1").Last(&lArticleForDate)
	//TODO Wrap error
	lLastWriteDate = lArticleForDate.PublishDate

	for _, lArticle := range aArticles {
		if lArticle.PublishDate.After(lLastWriteDate) {
			s.db.Create(&core.ArticleDB{
				Title:       lArticle.Title,
				Author:      lArticle.Author,
				Link:        lArticle.Link,
				PublishDate: lArticle.PublishDate,
				Description: lArticle.Description})
		}
	}

	/*//TODO Write func validate()
	//Validation length of fields
	if len(lArticle.Title) > 300 {
		lArticle.Title = cut(lArticle.Title, 300)
	}
	if len(lArticle.Author) > 200 {
		lArticle.Author = cut(lArticle.Author, 200)
	}
	if len(lArticle.Description) > 6000 {
		lArticle.Description = cut(lArticle.Description, 6000)
	}*/

	return nil

	//TODO Проверить стыковочные записи
	//TODO ВНИМАНИЕ! Изолировать в одной транзакции чтение из базы (когда дату для стыковки определяю) и запись в базу.
}

func (s *SqliteStorage) AddOneArticle(aArticle *core.ArticleDB) error {
	return nil
	//Добавления Артикла вручную
}
