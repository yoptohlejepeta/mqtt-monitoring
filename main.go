package main

// import src "monitoring/mqtt/internal"
import (
	cfg "monitoring/mqtt/config"
	src "monitoring/mqtt/internal"
)

func main() {
	config := cfg.Get_config()

	src.RunMqtt(config)
}
