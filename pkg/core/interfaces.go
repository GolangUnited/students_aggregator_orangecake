package core

import (
    "context"
)

type Handler interface {
    ParseArticles() error
}

type DBReader interface {
    ReadArticleById(aCtx context.Context, aLink string) (*Article, error)
    ReadAllArticles(aCtx context.Context) ([]Article, error)
}

type DBWriter interface {
    WriteArticle(aCtx context.Context, lArticle Article) error
}

type Logger interface {
    LogWarning()
    LogError()
    LogTrace()
    LogInfo()
    LogDebug()
}