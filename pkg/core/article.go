package core

import "time"

type Article struct {
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Link        string    `json:"link"`
	PublishDate time.Time `json:"date"`
	Description string    `json:"desc"`
}

type ArticleDB struct {
	ID      uint    `gorm:"primary_key;auto_increment;not_null"`
	Article Article `gorm:"embedded"`
}
