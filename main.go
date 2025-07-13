package main

import (
	"fmt"
	"github.com/a-h/templ"
	"log"
	cfg "monitoring/mqtt/config"
	src "monitoring/mqtt/internal"
	"net/http"
)

func main() {
	config := cfg.Get_config()

	// Start HTTP server in a goroutine
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/", templ.Handler(Index(config.Monitoring.Topics)))
		fmt.Println("Starting HTTP server on :8000")
		if err := http.ListenAndServe(":8000", mux); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Run MQTT client
	src.RunMqtt(config)
}
