package configs

type AggregatorConfig struct {
	HandlersConfig *HandlersConfig
	DBConfig       *DBConfig
}

func NewAggregatorConfig() (*AggregatorConfig, error) {
	handlersConfig, err := NewHandlersConfig()
	if err != nil {
		return nil, err
	}

	dbConfig, err := NewDBConfig()
	if err != nil {
		return nil, err
	}

	return &AggregatorConfig{
		HandlersConfig: handlersConfig,
		DBConfig:       dbConfig,
	}, nil
}
