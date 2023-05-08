package config

import "sync"

const (
	configPath = "/configs/database-config.yaml"
)

var (
	instance *Config
	once     sync.Once
)

type Config struct {
	Port     string `env:"PORT" env-default:"5432"`
	Host     string `env:"HOST" env-default:"localhost"`
	Name     string `env:"NAME" env-default:"postgres"`
	User     string `env:"USER" env-default:"user"`
	Password string `env:"PASSWORD"`
}
