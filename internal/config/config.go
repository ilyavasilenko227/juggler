package config

import (
	"juggler/internal/utils/logger"

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

func Config() (cfg App, err error) {
	if err = env.Parse(&cfg); err != nil {
		logger.Zap.Errorf("error: %s reading config from env", err.Error())
		return cfg, err
	}
	return cfg, nil
}
