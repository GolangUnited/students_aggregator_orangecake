package core

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var ErrMockImplementation = errors.New("mock implementation")

func TestReadArticleByIdDef(t *testing.T) {
	lDBReader := (DBReader)(&MockDBStorage{})
	_, lErr := lDBReader.ReadArticleByID(1)

	assert.ErrorIs(t, lErr, ErrNotImplemented)
}

func TestReadArticleByIdMock(t *testing.T) {
	lStorage := &MockDBStorage{}
	lStorage.ReadArticleByIdFunc = func(aId uint) (ArticleDB, error) {
		return ArticleDB{}, ErrMockImplementation
	}

	lDBReader := (DBReader)(lStorage)
	_, lErr := lDBReader.ReadArticleByID(1)

	assert.ErrorIs(t, lErr, ErrMockImplementation)
}

func TestReadArticlesByDateRangeDef(t *testing.T) {
	lDBReader := (DBReader)(&MockDBStorage{})
	_, lErr := lDBReader.ReadArticlesByDateRange(time.Now(), time.Now())

	assert.ErrorIs(t, lErr, ErrNotImplemented)
}

func TestReadArticlesByDateRangeMock(t *testing.T) {
	lStorage := &MockDBStorage{}
	lStorage.ReadArticlesByDateRangeFunc = func(aMin, aMax time.Time) ([]ArticleDB, error) {
		return nil, ErrMockImplementation
	}

	lDBReader := (DBReader)(lStorage)
	_, lErr := lDBReader.ReadArticlesByDateRange(time.Now(), time.Now())

	assert.ErrorIs(t, lErr, ErrMockImplementation)
}

func TestWriteArticleDef(t *testing.T) {
	lDBWriter := (DBWriter)(&MockDBStorage{})
	lErr := lDBWriter.WriteArticle(Article{})

	assert.ErrorIs(t, lErr, ErrNotImplemented)
}

func TestWriteArticleMock(t *testing.T) {
	lStorage := &MockDBStorage{}
	lStorage.WriteArticleFunc = func(aArticle Article) error {
		return ErrMockImplementation
	}

	lDBWriter := (DBWriter)(lStorage)
	lErr := lDBWriter.WriteArticle(Article{})

	assert.ErrorIs(t, lErr, ErrMockImplementation)
}

func TestWriteArticlesDef(t *testing.T) {
	lDBWriter := (DBWriter)(&MockDBStorage{})
	lErr := lDBWriter.WriteArticles(nil)

	assert.ErrorIs(t, lErr, ErrNotImplemented)
}

func TestWriteArticlesMock(t *testing.T) {
	lStorage := &MockDBStorage{}
	lStorage.WriteArticlesFunc = func(aArticles []Article) error {
		return ErrMockImplementation
	}

	lDBWriter := (DBWriter)(lStorage)
	lErr := lDBWriter.WriteArticles(nil)

	assert.ErrorIs(t, lErr, ErrMockImplementation)
}

func TestUpdateArticleDef(t *testing.T) {
	lDBWriter := (DBWriter)(&MockDBStorage{})
	lErr := lDBWriter.UpdateArticle(1, Article{})

	assert.ErrorIs(t, lErr, ErrNotImplemented)
}

func TestUpdateArticleMock(t *testing.T) {
	lStorage := &MockDBStorage{}
	lStorage.UpdateArticleFunc = func(aId uint, aArticle Article) error {
		return ErrMockImplementation
	}

	lDBWriter := (DBWriter)(lStorage)
	lErr := lDBWriter.UpdateArticle(1, Article{})

	assert.ErrorIs(t, lErr, ErrMockImplementation)
}
