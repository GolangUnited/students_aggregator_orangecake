package configs

type DBConfig struct {
	Host     string
	Port     int
	Username string
	DBName   string
	SSLMode  string
	Password string
}

func NewDBConfig() (*DBConfig, error) {

	return &DBConfig{}, nil
}
