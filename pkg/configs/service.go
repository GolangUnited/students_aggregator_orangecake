package configs

type ServiceConfig struct {
	DBConfig *DBConfig
	Server   *Server
}

func NewServiceConfig() (*ServiceConfig, error) {
	dbConfig := NewDBConfig()
	serverConfig := NewServerConfig()

	return &ServiceConfig{
		Server:   serverConfig,
		DBConfig: dbConfig,
	}, nil
}
