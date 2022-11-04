package database

import (
	"database/sql"
	"fmt"
)

type ConnectionInfo struct {
	FileName string
}

func NewSqliteConnection(info ConnectionInfo) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", info.FileName)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	return db, nil
}
