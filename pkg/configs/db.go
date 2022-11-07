package configs

import (
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"os"
)

const (
	DB_CONNECTIONSTRING = "OC_DB_CONNECTION_STRING"
	DB_HOST             = "OC_DB_HOST"
	DB_PORT             = "OC_DB_PORT"
	DB_USERNAME         = "OC_DB_USERNAME"
	DB_SSLMOE           = "OC_DB_SSL_MODE"
	DB_NAME             = "OC_DB_NAME"
	DB_PASSWORD         = "OC_DB_PASSWORD"
)

type DBConfig struct {
	ConnectionString string
	Host             string
	Port             string
	Username         string
	SSLMode          string
	DBName           string
	Password         string
}

func NewDBConfig() (*DBConfig, error) {

	lConfig := DBConfig{
		ConnectionString: os.Getenv(DB_CONNECTIONSTRING),
		Host:             os.Getenv(DB_HOST),
		Port:             os.Getenv(DB_PORT),
		Username:         os.Getenv(DB_USERNAME),
		SSLMode:          os.Getenv(DB_SSLMOE),
		DBName:           os.Getenv(DB_NAME),
		Password:         os.Getenv(DB_PASSWORD),
	}

	lEmptyConfig := DBConfig{}

	if lConfig == lEmptyConfig {
		return nil, core.ErrEmptyConfig
	}

	return &lConfig, nil
}
