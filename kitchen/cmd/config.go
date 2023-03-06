package main

import (
	"github.com/caarlos0/env"
)

type Config struct {
	RabbitMQUser     string `env:"RABBITMQ_USER,required"`
	RabbitMQPassword string `env:"RABBITMQ_PASSWORD,required"`
	RabbitMQAddress  string `env:"RABBITMQ_ADDRESS,required"`
	RabbitMQPort     string `env:"RABBITMQ_PORT,required"`
	PostgresUser     string `env:"POSTGRES_USER,required"`
	PostgresPassword string `env:"POSTGRES_PASSWORD,required"`
	PostgresAddr     string `env:"POSTGRES_ADDR,required"`
	PostgresPort     string `env:"POSTGRES_PORT,required"`
	PostgresDBName   string `env:"POSTGRES_DB,required"`
}

func ParseConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
