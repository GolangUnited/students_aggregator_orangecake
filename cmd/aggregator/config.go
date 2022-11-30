package main

import (
	"fmt"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	_ "github.com/joho/godotenv/autoload"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

const (
	HANDLERS_FILE       = "handlers.config.yaml"
	DB_CONNECTIONSTRING = "OC_DB_CONNECTION_STRING"
	ERROR_MESSAGE       = "new aggregator config error"
)

type AggregatorConfig struct {
	ConnectionString string
	Handlers         []struct {
		Handler string `yaml:"handler"`
		URL     string `yaml:"url"`
	} `yaml:"handlers"`
}

func NewAggregatorConfig() (*AggregatorConfig, error) {
	var lConfig AggregatorConfig

	lConnectionString := os.Getenv(DB_CONNECTIONSTRING)
	if len(lConnectionString) == 0 {
		return nil, fmt.Errorf("%s: %w", ERROR_MESSAGE, core.ErrEmptyEnvVariable)
	}

	lConfig.ConnectionString = lConnectionString

	lFile, err := os.Open(HANDLERS_FILE)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ERROR_MESSAGE, err)
	}

	lFileBuffer, err := io.ReadAll(lFile)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ERROR_MESSAGE, err)
	}

	err = yaml.Unmarshal(lFileBuffer, &lConfig)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ERROR_MESSAGE, err)
	}

	return &lConfig, nil
}
