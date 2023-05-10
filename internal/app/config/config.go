package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	instance *Config
	once     sync.Once
)

type Config struct {
	ApiPort    string `env:"API_PORT"`
	ApiHost    string `env:"API_HOST"`
	ClientHost string `env:"CLIENT_HOST"`
	ClientPort string `env:"CLIENT_PORT"`
}

func GetConfig() *Config {
	once.Do(func() {
		log.Println("gather app config")
		instance = &Config{}

		if err := cleanenv.ReadEnv(instance); err != nil {
			helpText := "Gametrade - the best gaming store"
			description, _ := cleanenv.GetDescription(instance, &helpText)
			log.Println(description)
			log.Fatal(err)
		}
	})
	return instance
}
