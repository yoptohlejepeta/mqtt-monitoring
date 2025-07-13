package config

import (
	"fmt"
	"log"
	"os"
	"time"

	env "github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
	"gopkg.in/yaml.v3"
)

type TopicConfig struct {
	Name     string        `yaml:"name"`
	Hours    []int         `yaml:"hours"`
	Interval time.Duration `yaml:"interval"`
	MinCount int           `yaml:"min_count"`
}

type MonitoringConfig struct {
	Topics []TopicConfig `yml:"topics"`
}

type MqttConfig struct {
	Host     string `env:"MQTT_HOST"`
	Port     int    `env:"MQTT_PORT"`
	User     string `env:"MQTT_USERNAME"`
	Password string `env:"MQTT_PASSWORD"`
}

type Config struct {
	Mqtt       MqttConfig
	Monitoring MonitoringConfig
}

func Get_config() Config {
	var mqtt_cfg MqttConfig
	var monitoring_config MonitoringConfig

	cfg := Config{
		Mqtt:       *mqtt_cfg.ParseEnv(),
		Monitoring: *monitoring_config.ParseYml(),
	}
	return cfg
}

func (c *MqttConfig) ParseEnv() *MqttConfig {
	err := env.Parse(c)
	if err != nil {
		panic(err)
	}

	return c
}

func (c *MonitoringConfig) ParseYml() *MonitoringConfig {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get working directory.")
	}

	configPath := fmt.Sprint(cwd, "/config.yml")
	fileContent, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal("Failed to read config file. Path: ", configPath)
	}

	if err := yaml.Unmarshal(fileContent, c); err != nil {
		log.Fatal("Failed to parse config file. | ", err)
	}

	return c
}
