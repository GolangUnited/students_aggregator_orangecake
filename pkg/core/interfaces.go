package core

import (
    "context"
    "time"
)

type ArticleScraper interface {
    ParseArticles() error
    GetArticles() []Article
}

type DBReader interface {
    ReadArticleById(aCtx context.Context, aLink string) (*Article, error)
    ReadArticlesByDateRange(aCtx context.Context, aMin, aMax time.Time) ([]Article, error)
}

type DBWriter interface {
    WriteArticle(aCtx context.Context, aArticle Article) error
}

type Logger interface {
    Info(aMessage string, aValues ...interface{})
    Error(aMessage string, aValues ...interface{})
    Warn(aMessage string, aValues ...interface{})
    Debug(aMessage string, aValues ...interface{})
    Trace(aMessage string, aValues ...interface{})
}