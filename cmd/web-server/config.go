package main

import (
	"fmt"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	_ "github.com/joho/godotenv/autoload"
	"os"
)

const (
	SERVER_PORT         = "OC_SERVER_PORT"
	DB_CONNECTIONSTRING = "OC_DB_CONNECTION_STRING"
	ERROR_MESSAGE       = "new service config error"
)

type ServiceConfig struct {
	DBConnectionString string
	ServerPort         string
}

func NewServiceConfig() (*ServiceConfig, error) {
	lPort := os.Getenv(SERVER_PORT)
	if len(lPort) == 0 {
		return nil, fmt.Errorf("%s: %s: %w", ERROR_MESSAGE, SERVER_PORT, core.ErrEmptyEnvVariable)
	}

	lConnectionString := os.Getenv(DB_CONNECTIONSTRING)
	if len(lConnectionString) == 0 {
		return nil, fmt.Errorf("%s: %s: %w", ERROR_MESSAGE, DB_CONNECTIONSTRING, core.ErrEmptyEnvVariable)
	}

	return &ServiceConfig{
		DBConnectionString: lPort,
		ServerPort:         lConnectionString,
	}, nil
}
