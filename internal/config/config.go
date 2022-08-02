package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	RedisHost string `env:"REDIS_HOST"`
	RedisPort string `env:"REDIS_PORT"`
	PortRPC   string `env:"RPC_POSITION_SERVICE_PORT"`
}

func GetConfig() (*Config, error) {
	config := Config{}
	err := env.Parse(&config)
	if err != nil {
		return nil, fmt.Errorf("config / GetConfig / err parse : %v", err)
	}
	return &config, nil
}
