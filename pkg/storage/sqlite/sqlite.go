package sqlite

import (
	"fmt"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewSqliteConnection ... maybe conflict between gorm logger and custom logger, cause the same names
func NewSqliteConnection(aConnectionString string, aLogger core.Logger) (*SqliteStorage, error) {
	dbArticles, lErr := gorm.Open(sqlite.Open(aConnectionString), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if lErr != nil {
		return nil, fmt.Errorf("can't open database: %w", lErr)
	}

	lErr = dbArticles.AutoMigrate(&core.ArticleDB{})
	if lErr != nil {
		fmt.Printf("failed to migrate from Article struct: %v", lErr)
		return nil, lErr
	}

	return &SqliteStorage{db: dbArticles, logger: aLogger}, nil
}
