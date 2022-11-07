package configs

import (
	"os"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
)

type Server struct {
	Port string
}

func NewServerConfig() (*Server, error) {
	lConfig := Server{
		Port: os.Getenv("OC_SERVER_PORT"),
	}

	lEmptyConfig := Server{}

	if lConfig == lEmptyConfig {
		return nil, core.ErrEmptyConfig
	}

	return &lConfig, nil
}
