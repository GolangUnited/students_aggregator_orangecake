package sqlite

import (
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"gorm.io/gorm"
	"log"
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

// NewTable create table articles with field like struct core.ArticleDB.
func (s *SqliteStorage) NewTable(db *gorm.DB) {
	lErr := db.AutoMigrate(&core.ArticleDB{})
	if lErr != nil {
		panic("failed to migrate from Article struct") // or log.Fatal?
	}
}

// cut decrease the length of field
func cut(aText string, aLimit int) string {
	lRunes := []rune(aText)
	if len(lRunes) >= aLimit {
		return string(lRunes[:aLimit])
	}
	return aText
}

// validation check length of fields and cut it, if need it.
func validation(aArticle core.Article) {
	if len(aArticle.Title) > LEN_OF_TITLE {
		aArticle.Title = cut(aArticle.Title, LEN_OF_TITLE)
	}

	if len(aArticle.Author) > LEN_OF_AUTHOR {
		aArticle.Author = cut(aArticle.Author, LEN_OF_AUTHOR)
	}

	if len(aArticle.Description) > LEN_OF_DESCRIPTION {
		aArticle.Description = cut(aArticle.Description, LEN_OF_DESCRIPTION)
	}
}

func (s *SqliteStorage) WriteArticle(aArticle core.Article) error {
	lResult := s.db.Create(&core.ArticleDB{
		Title:       aArticle.Title,
		Author:      aArticle.Author,
		Link:        aArticle.Link,
		PublishDate: aArticle.PublishDate,
		Description: aArticle.Description,
	})

	if lResult.Error != nil {
		log.Printf("write article returns an error: %w", lResult.Error)
		return lResult.Error
	}
	log.Printf("wrote %d article", lResult.RowsAffected)
	return nil
}

func (s *SqliteStorage) WriteArticles(aArticles []core.Article) error {
	lCountOfWritingArticles := 0

	for _, lArticle := range aArticles {
		//Validation length of fields
		validation(lArticle)

		//TODO Добавить условие про последнюю неделю, или отсчитать последние 10 артиклов из каждого хэндлера!
		lResult := s.db.Create(&core.ArticleDB{
			Title:       lArticle.Title,
			Author:      lArticle.Author,
			Link:        lArticle.Link,
			PublishDate: lArticle.PublishDate,
			Description: lArticle.Description})

		if lResult.Error != nil {
			log.Printf("can't write articles: %w", lResult.Error)
			return lResult.Error
		}
		lCountOfWritingArticles += int(lResult.RowsAffected)
	}

	if lCountOfWritingArticles != len(aArticles) {
		log.Println("Count of records and length of []Articles mismatch")
	}

	return nil
}

func (s *SqliteStorage) UpdateArticles(aArticles []core.Article) error {
	var lLastWriteDate time.Time
	var lArticleForDate core.ArticleDB
	//TODO Добавить условие про последнюю неделю
	s.db.Raw("SELECT * FROM article_dbs ORDER BY publish_date DESC LIMIT 1").Last(&lArticleForDate)
	//TODO Wrap error
	lLastWriteDate = lArticleForDate.PublishDate
	//TODO Validate
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

	return nil
}

func (s *SqliteStorage) ReadArticleByID(aID uint) (*core.ArticleDB, error) {
	var lArticle core.ArticleDB

	lResult := s.db.Where("id = ?", aID).First(&lArticle)

	lErr := lResult.Error
	if lErr != nil {
		log.Printf("error happens in row with id = %d: %w", aID, lErr)
		return nil, lErr
	}

	return &lArticle, nil
}

func (s *SqliteStorage) ReadArticlesByDateRange(aMin, aMax time.Time) ([]core.ArticleDB, error) {
	lArticles := make([]core.ArticleDB, 0)

	lResult := s.db.Where("publish_date BETWEEN ? AND ?", aMin, aMax).Find(&lArticles)
	lErr := lResult.Error
	if lErr != nil {
		log.Printf("error happens in rows between dates %s and %s: %w", aMin, aMax, lErr)
		return nil, lErr
	}
	log.Printf("%d rows was found.", lResult.RowsAffected)
	return lArticles, nil
}
