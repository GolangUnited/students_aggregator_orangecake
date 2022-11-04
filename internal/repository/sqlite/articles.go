package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
)

type SQLiteStorage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *SQLiteStorage {
	return &SQLiteStorage{db}
}

func (s *SQLiteStorage) WriteArticle(aCtx context.Context, lArticle core.Article) error {
	lStatement := "INSERT INTO articles (title, author, link, publish_date, description) VALUES (?, ?, ?, ?, ?)"

	if _, lErr := s.db.ExecContext(aCtx, lStatement, lArticle.Title, lArticle.Author, lArticle.Link, lArticle.PublishDate, lArticle.Description); lErr != nil {
		return fmt.Errorf("can't write article: %w", lErr)
	}

	return nil
}

func (s *SQLiteStorage) ReadArticleById(aCtx context.Context, aLink string) (*core.Article, error) {
	var lArticle core.Article //нужен конструктор для создания пустого статьи ИЛИ aArticle в аргумент?
	lStatement := "SELECT title, author, link, publish_date, description FROM articles WHERE id=?"

	lErr := s.db.QueryRowContext(aCtx, lStatement, aLink).Scan(&lArticle.Title, &lArticle.Author, &lArticle.Link, &lArticle.PublishDate, &lArticle.Description)
	if lErr == sql.ErrNoRows {
		return nil, fmt.Errorf("article with link #%s not found: %w", aLink, lErr)
	}
	if lErr != nil {
		return nil, fmt.Errorf("can't read article with link #%s: %w", aLink, lErr)
	}

	return &lArticle, nil
}

func (s *SQLiteStorage) ReadAllArticles(aCtx context.Context) ([]core.Article, error) {
	lRows, lErr := s.db.QueryContext(aCtx, "SELECT title, author, link, publish_date, description FROM articles")
	if lErr != nil {
		return nil, lErr
	}

	lArticles := make([]core.Article, 0)
	for lRows.Next() {
		var lArticle core.Article
		if err := lRows.Scan(&lArticle.Title, &lArticle.Author, &lArticle.Link, &lArticle.PublishDate, &lArticle.Description); err != nil {
			return nil, err
		}

		lArticles = append(lArticles, lArticle)
	}

	return lArticles, lRows.Err() //обернуть ошибку
}
