package sqlite

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewSqliteConnection(aConnectionString string) (*SqliteStorage, error) {
	dbArticles, lErr := gorm.Open(sqlite.Open(aConnectionString), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if lErr != nil {
		return nil, fmt.Errorf("can't open database: %w", lErr)
	}

	sqlDb, lErr := dbArticles.DB()

	if err := sqlDb.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	return &SqliteStorage{Db: dbArticles}, nil
}
