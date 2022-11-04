package configs

type ServiceConfig struct {
	DBConfig *DBConfig
	Server   *Server
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
		Server:   serverConfig,
		DBConfig: dbConfig}, nil
}
