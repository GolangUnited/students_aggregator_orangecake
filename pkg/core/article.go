package core

import "time"

type Article struct {
	Caption     string
	Author      string
	Link        string
	PublishDate time.Time
	Description string
}
