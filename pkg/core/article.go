package core

import "time"

type Article struct {
	Title       string
	Author      string
	Link        string
	PublishDate time.Time
	Description string
}

type ArticleDB struct {
	ID          uint `gorm:"primary_key;auto_increment;not_null"`
	Title       string
	Author      string
	Link        string
	PublishDate time.Time
	Description string
}
