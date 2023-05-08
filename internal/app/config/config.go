package config

import "sync"

const (
	configPath = "/configs/app-config.yaml"
)

var (
	instance *Config
	once     sync.Once
)

type Config struct {
	Port string `yaml:"port" env:"PORT" env-default:"8080"`
	Host string `yaml:"host" env:"HOST" env-default:"localhost"`
}
