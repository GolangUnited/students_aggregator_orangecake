package core

import "time"

type Article struct {
	Title       string
	Author      string
	Link        string
	PublishDate time.Time
	Description string
}
