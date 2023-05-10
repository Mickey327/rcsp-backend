package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	instance *Config
	once     sync.Once
)

type Config struct {
	Port     int    `env:"DB_PORT" env-default:"5432"`
	Host     string `env:"DB_HOST" env-default:"database"`
	Name     string `env:"DB_NAME" env-default:"gametrade"`
	User     string `env:"DB_USER" env-default:"postgres"`
	Password string `env:"DB_PASSWORD"`
}

func (c *Config) GenerateConnectPath() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Name)
}

func GetConfig() *Config {
	once.Do(func() {
		log.Println("gather database config")
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
