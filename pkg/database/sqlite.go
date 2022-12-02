package database

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewSqliteConnection(aFileName string) (*gorm.DB, error) {
	db, lErr := gorm.Open(sqlite.Open(aFileName), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if lErr != nil {
		return nil, fmt.Errorf("can't open database: %w", lErr)
	}

	sqlDb, lErr := db.DB()

	if err := sqlDb.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}
	return db, nil
}
