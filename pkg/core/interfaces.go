package core

import (
    "time"
)

type ArticleScraper interface {
    ParseArticles() ([]Article, []Warning, error)
}

type DBReader interface {
    ReadArticleByID(ID uint) (Article, error)
    ReadArticlesByDateRange(aMin, aMax time.Time) ([]Article, error)
}

type DBWriter interface {
    WriteArticle(aArticle Article) error
    WriteArticles(aArticles []Article) error
    UpdateArticle(aID uint, aArticle Article) error
}

type Logger interface {
    Info(aMessage string, aValues ...interface{})
    Error(aMessage string, aValues ...interface{})
    Warn(aMessage string, aValues ...interface{})
    Debug(aMessage string, aValues ...interface{})
    Trace(aMessage string, aValues ...interface{})
}