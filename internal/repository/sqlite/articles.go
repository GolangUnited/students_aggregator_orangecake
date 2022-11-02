package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db}
}

func (s *Storage) WriteArticle(aCtx context.Context, lArticle core.Article) error {
	lStatement := "INSERT INTO articles (title, author, link, pablish_date, description) VALUES (?, ?, ?, ?, ?)"

	if _, lErr := s.db.ExecContext(aCtx, lStatement, lArticle.Title, lArticle.Author, lArticle.Link, lArticle.PublishDate, lArticle.Description); lErr != nil {
		return fmt.Errorf("can't write article: %w", lErr)
	}

	return nil
}

func (s *Storage) ReadArticleById(aCtx context.Context, aId int) (core.Article, error) {
	var lArticle core.Article //нужен конструктор для создания пустого статьи ИЛИ aArticle в аргумент?
	lStatement := "SELECT title, author, link, pablish_date, description FROM articles WHERE id=?"

	lErr := s.db.QueryRowContext(aCtx, lStatement, aId).Scan(&lArticle.Title, &lArticle.Author, &lArticle.Link, &lArticle.PublishDate, &lArticle.Description)
	if lErr == sql.ErrNoRows {
		return lArticle, fmt.Errorf("article #%d not found: %w", aId, lErr)
	}
	if lErr != nil {
		return core.Article{}, fmt.Errorf("can't read article #%d: %w", aId, lErr)
	}

	return lArticle, nil
}

func (s *Storage) ReadAllArticles(aCtx context.Context) ([]core.Article, error) {
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
