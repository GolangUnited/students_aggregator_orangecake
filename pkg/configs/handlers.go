package configs

import (
	_ "embed"
	"gopkg.in/yaml.v3"
)

//go:embed config.handlers.yaml
var data []byte

type HandlersConfig struct {
	Handlers []struct {
		Handler string `yaml:"handler"`
		URL     string `yaml:"url"`
	} `yaml:"handlers"`
}

func NewHandlersConfig() (*HandlersConfig, error) {
	var lConfig HandlersConfig
	err := yaml.Unmarshal(data, &lConfig)
	if err != nil {
		return nil, err
	}

	return &lConfig, nil
}
