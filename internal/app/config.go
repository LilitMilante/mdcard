package app

import (
	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
)

type Config struct {
	Port string `env:"PORT"`

	Database DBConfig
}

type DBConfig struct {
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT"`
	Name     string `env:"DB_NAME"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
}

func NewConfig() (c Config, err error) {
	err = godotenv.Load(".env")
	if err != nil {
		return c, err
	}

	err = env.Parse(&c)
	if err != nil {
		return c, err
	}

	return c, nil
}
