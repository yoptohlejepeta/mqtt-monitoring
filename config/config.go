package config

import (
	"fmt"

	env "github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Host     string `env:"MQTT_HOST"`
	Port     int    `env:"MQTT_PORT"`
	User     string `env:"MQTT_USERNAME"`
	Password string `env:"MQTT_PASSWORD"`
	Topic    string `env:"MQTT_TOPIC"`
}

func Get_config() (Config, error) {
	var cfg Config

	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg)

	return cfg, nil
}
