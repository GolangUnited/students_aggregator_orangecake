package configs

import (
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"os"
)

const (
	DB_CONNECTIONSTRING = "OC_DB_CONNECTION_STRING"
)

type DBConfig struct {
	ConnectionString string
}

func NewDBConfig() (*DBConfig, error) {

	lConfig := DBConfig{
		ConnectionString: os.Getenv(DB_CONNECTIONSTRING),
	}

	lEmptyConfig := DBConfig{}

	if lConfig == lEmptyConfig {
		return nil, core.ErrEmptyConfig
	}

	return &lConfig, nil
}
