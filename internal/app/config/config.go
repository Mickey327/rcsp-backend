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
	ApiPort            string `env:"API_PORT"`
	ApiHost            string `env:"API_HOST"`
	OuterClientAddress string `env:"OUTER_ADDRESS"`
	ClientHost         string `env:"CLIENT_HOST"`
	ClientPort         string `env:"CLIENT_PORT"`
	Email              string `env:"MAIL_EMAIL"`
	MailHost           string `env:"MAIL_HOST"`
	MailPort           int    `env:"MAIL_PORT"`
	MailPassword       string `env:"MAIL_PASSWORD"`
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
