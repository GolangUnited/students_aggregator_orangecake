package configs

import "os"

type DBConfig struct {
	Host     string
	Port     string
	Username string
	DBName   string
	SSLMode  string
	Password string
}

func NewDBConfig() *DBConfig {

	return &DBConfig{
		Host:     os.Getenv("OC_DB_HOST"),
		Port:     os.Getenv("OC_DB_PORT"),
		Username: os.Getenv("OC_DB_USERNAME"),
		SSLMode:  os.Getenv("OC_DB_SSLMODE"),
		DBName:   os.Getenv("OC_DB_NAME"),
		Password: os.Getenv("OC_DB_PASSWORD"),
	}

}
