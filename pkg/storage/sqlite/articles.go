package sqlite

import (
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"gorm.io/gorm"
	"log"
	"time"
)

const (
	LEN_OF_TITLE       = 200
	LEN_OF_AUTHOR      = 200
	LEN_OF_DESCRIPTION = 6000
)

type SqliteStorage struct {
	db *gorm.DB
}

// newArticleDB ...
func newArticleDB(aArticle *core.Article) *core.ArticleDB {
	return &core.ArticleDB{
		Article: core.Article{
			Title:       aArticle.Title,
			Author:      aArticle.Author,
			Link:        aArticle.Link,
			PublishDate: aArticle.PublishDate,
			Description: aArticle.Description,
		},
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
// TODO "&" is need it, like &len(&aArticle.Title)?
func validation(aArticle *core.Article) {
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

// WriteArticle adds manually one article
func (s *SqliteStorage) WriteArticle(aArticle core.Article) error {
	var lRowsAffected int64 = 0

	lErrTrans := s.db.Transaction(func(tx *gorm.DB) error {
		if countOfRecordsFound := tx.First(
			&core.ArticleDB{}, "link = ?", aArticle.Link).RowsAffected; countOfRecordsFound == 0 {

			validation(&aArticle)

			lResult := tx.Create(newArticleDB(&aArticle))
			if lResult.Error != nil {
				log.Printf("write article returns an error: %v", lResult.Error)
				return lResult.Error
			}
			lRowsAffected = lResult.RowsAffected
			log.Printf("wrote %d article", lRowsAffected)
			return nil
		} else {
			log.Printf("wrote %d articles", lRowsAffected)
		}
		return nil
	})

	if lErrTrans != nil {
		log.Printf("error of transaction: %v", lErrTrans)
	}

	return nil
}

// WriteArticles adds slice of articles after scraping data
func (s *SqliteStorage) WriteArticles(aArticles []core.Article) error {
	lCountOfWritingArticles := 0
	// unhandled error in Transaction
	lErrTrans := s.db.Transaction(func(tx *gorm.DB) error {
		for _, lArticle := range aArticles {
			//Validation length of fields
			validation(&lArticle)

			//TODO change 1680 on 168 after check
			if lArticle.PublishDate.After(core.NormalizeDate(time.Now()).Add(-1680 * time.Hour)) {

				if countOfRecordsFound := tx.First(
					&core.ArticleDB{},
					"link = ?", lArticle.Link).RowsAffected; countOfRecordsFound == 0 {

					lResult := tx.Create(newArticleDB(&lArticle))

					if lResult.Error != nil {
						log.Printf("can't write articles: %v", lResult.Error)
						return lResult.Error
					}
					lCountOfWritingArticles += int(lResult.RowsAffected)
				} else {
					log.Printf("article with link %s already exsist", lArticle.Link)
				}
			}
		}
		return nil
	})

	if lErrTrans != nil {
		log.Printf("error of transaction when write articles: %v", lErrTrans)
		return lErrTrans
	}

	log.Printf("wrote %d articles", lCountOfWritingArticles)

	return nil
}

// UpdateArticle updates record with id = aID
func (s *SqliteStorage) UpdateArticle(aID uint, aArticle core.Article) error {
	var lArticle core.ArticleDB

	lErrTrans := s.db.Transaction(func(tx *gorm.DB) error {

		lExistId := tx.First(&lArticle, aID).RowsAffected

		if lExistId == 1 {
			validation(&aArticle)

			lResult := tx.Model(&core.ArticleDB{}).Where("id = ?", aID).
				Updates(newArticleDB(&aArticle))
			if lResult.Error != nil {
				return lResult.Error
			}

			log.Printf("update %d articles with id #%d", lResult.RowsAffected, aID)
		} else if lExistId == 0 {
			log.Printf("nothing to update, id #%d isn't exist", aID)
		}
		return nil
	})

	if lErrTrans != nil {
		log.Printf("error of transaction when update article with id #%d: %v", aID, lErrTrans)
		return lErrTrans
	}

	return nil
}

// ReadArticleByID returns record with id = aID
func (s *SqliteStorage) ReadArticleByID(aID uint) (*core.ArticleDB, error) {
	var lArticle core.ArticleDB

	lResult := s.db.Where("id = ?", aID).First(&lArticle)

	lErr := lResult.Error
	if lErr != nil {
		log.Printf("error happens in row with id = %d: %v", aID, lErr)
		return nil, lErr
	}

	return &lArticle, nil
}

// ReadArticlesByDateRange returns records, that were written between aMin and aMax time frames
func (s *SqliteStorage) ReadArticlesByDateRange(aMin, aMax time.Time) ([]core.ArticleDB, error) {
	lArticles := make([]core.ArticleDB, 0)

	lResult := s.db.Where("publish_date BETWEEN ? AND ?", aMin, aMax).Find(&lArticles)
	lErr := lResult.Error
	if lErr != nil {
		log.Printf("error happens in rows between dates %s and %s: %v", aMin, aMax, lErr)
		return nil, lErr
	}
	log.Printf("%d rows was found.", lResult.RowsAffected)
	return lArticles, nil
}
