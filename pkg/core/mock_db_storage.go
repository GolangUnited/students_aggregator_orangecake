package core

import (
	"fmt"
	"time"
)

type MockDBStorage struct {
	ReadArticleByIdFunc         func(aId uint) (ArticleDB, error)
	ReadArticlesByDateRangeFunc func(aMin, aMax time.Time) ([]ArticleDB, error)
	WriteArticleFunc            func(aArticle Article) error
	WriteArticlesFunc           func(aArticles []Article) error
	UpdateArticleFunc           func(aId uint, aArticle Article) error
}

const errorTemplate = "%w (MockDBStorage.%s)"

func (s MockDBStorage) ReadArticleByID(aId uint) (ArticleDB, error) {
	if s.ReadArticleByIdFunc != nil {
		return s.ReadArticleByIdFunc(aId)
	}

	return ArticleDB{}, fmt.Errorf(errorTemplate, ErrNotImplemented, "ReadArticleByID")
}

func (s MockDBStorage) ReadArticlesByDateRange(aMin, aMax time.Time) ([]ArticleDB, error) {
	if s.ReadArticlesByDateRangeFunc != nil {
		return s.ReadArticlesByDateRangeFunc(aMin, aMax)
	}

	return nil, fmt.Errorf(errorTemplate, ErrNotImplemented, "ReadArticlesByDateRange")
}

func (s MockDBStorage) WriteArticle(aArticle Article) error {
	if s.WriteArticleFunc != nil {
		return s.WriteArticleFunc(aArticle)
	}

	return fmt.Errorf(errorTemplate, ErrNotImplemented, "WriteArticle")
}

func (s MockDBStorage) WriteArticles(aArticles []Article) error {
	if s.WriteArticlesFunc != nil {
		return s.WriteArticlesFunc(aArticles)
	}

	return fmt.Errorf(errorTemplate, ErrNotImplemented, "WriteArticles")
}

func (s MockDBStorage) UpdateArticle(aId uint, aArticle Article) error {
	if s.UpdateArticleFunc != nil {
		return s.UpdateArticleFunc(aId, aArticle)
	}

	return fmt.Errorf(errorTemplate, ErrNotImplemented, "UpdateArticle")
}
