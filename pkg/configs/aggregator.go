package configs

type AggregatorConfig struct {
	Handlers *HandlersConfig
	DBConfig *DBConfig
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
		Handlers: handlersConfig,
		DBConfig: dbConfig,
	}, nil
}
