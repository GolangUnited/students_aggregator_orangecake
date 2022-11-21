package database

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSqliteConnection(aFileName string) (*sql.DB, error) {
	db, lErr := gorm.Open(sqlite.Open(aFileName), &gorm.Config{})
	if lErr != nil {
		return nil, fmt.Errorf("can't open database: %w", lErr)
	}

	sqlDb, lErr := db.DB()

	if err := sqlDb.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}
	return sqlDb, nil
}
