package configs

import (
	_ "embed"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"gopkg.in/yaml.v3"
)

// if file not found error will be printed in std out
//
//go:embed handlers.config.example.yaml
var data []byte

type HandlersConfig struct {
	Handlers []struct {
		Handler string `yaml:"handler"`
		URL     string `yaml:"url"`
	} `yaml:"handlers"`
}

func NewHandlersConfig() (*HandlersConfig, error) {
	if len(data) == 0 {
		return nil, core.ErrConfigFileIsEmpty
	}
	var lConfig HandlersConfig
	err := yaml.Unmarshal(data, &lConfig)
	if err != nil {
		return nil, err
	}

	if lConfig.Handlers == nil {
		return nil, core.ErrNoHandlersInConfig
	}

	return &lConfig, nil
}
