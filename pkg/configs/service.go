package configs

type ServiceConfig struct {
	DBConfig     *DBConfig
	ServerConfig *ServerConfig
}

func NewServiceConfig() (*ServiceConfig, error) {
	dbConfig, err := NewDBConfig()
	if err != nil {
		return nil, err
	}

	serverConfig, err := NewServerConfig()
	if err != nil {
		return nil, err
	}

	return &ServiceConfig{
		ServerConfig: serverConfig,
		DBConfig:     dbConfig,
	}, nil
}
