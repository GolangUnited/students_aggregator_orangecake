package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/indikator/aggregator_orange_cake/pkg/core"
	_ "github.com/joho/godotenv/autoload"
)

const (
	SERVER_PORT         = "OC_SERVER_PORT"
	DB_CONNECTIONSTRING = "OC_DB_CONNECTION_STRING"
)

type ServiceConfig struct {
	DBConnectionString string
	ServerPort         int
}

func NewServiceConfig() (*ServiceConfig, error) {
	lPortStr := os.Getenv(SERVER_PORT)
	if len(lPortStr) == 0 {
		return nil, fmt.Errorf("%w (%s)", core.ErrEmptyEnvVariable, SERVER_PORT)
	}
	lPortValue, lErr := strconv.Atoi(lPortStr)
	if lErr != nil {
		return nil, fmt.Errorf("%w (%s=%s)", core.ErrInvalidConfigValue, SERVER_PORT, lPortStr)
	}

	lConnectionString := strings.TrimSpace(os.Getenv(DB_CONNECTIONSTRING))
	if len(lConnectionString) == 0 {
		return nil, fmt.Errorf("%w (%s)", core.ErrEmptyEnvVariable, DB_CONNECTIONSTRING)
	}

	return &ServiceConfig{
		DBConnectionString: lConnectionString,
		ServerPort:         lPortValue,
	}, nil
}
