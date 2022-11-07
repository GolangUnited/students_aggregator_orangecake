package configs

import (
	"github.com/indikator/aggregator_orange_cake/pkg/core"
	"os"
)

const (
	SERVER_PORT = "OC_SERVER_PORT"
)

type Server struct {
	Port string
}

func NewServerConfig() (*Server, error) {
	lConfig := Server{
		Port: os.Getenv(SERVER_PORT),
	}

	lEmptyConfig := Server{}

	if lConfig == lEmptyConfig {
		return nil, core.ErrEmptyConfig
	}

	return &lConfig, nil
}
