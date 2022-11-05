package configs

import "os"

type Server struct {
	Port string
}

func NewServerConfig() *Server {
	port := os.Getenv("OC_SERVER_PORT")

	return &Server{Port: port}
}
