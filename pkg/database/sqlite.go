package database

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSqliteConnection(aFileName string) (*gorm.DB, error) {
	db, lErr := gorm.Open(sqlite.Open(aFileName), &gorm.Config{})
	if lErr != nil {
		return nil, fmt.Errorf("can't open database: %w", lErr)
	}

	//TODO Ping in GORM
	/*	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}*/

	return db, nil
}
