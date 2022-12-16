package sqlite

import (
	"fmt"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

func NewSqliteConnection(aConnectionString string) (*SqliteStorage, error) {
	dbArticles, lErr := gorm.Open(sqlite.Open(aConnectionString), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if lErr != nil {
		return nil, fmt.Errorf("can't open database: %w", lErr)
	}

	lErr = dbArticles.AutoMigrate(&core.ArticleDB{})
	if lErr != nil {
		log.Printf("failed to migrate from Article struct: %v", lErr)
		return nil, lErr
	}

	return &SqliteStorage{db: dbArticles}, nil
}
