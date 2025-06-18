package internal

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"

	cfg "monitoring/mqtt/config"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Message: %s | Topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v", err)
}

// Connects to MQTT and subscribes to topics.
// Timeout 5 seconds.
func RunMqtt(cfg cfg.Config) {
	opts := mqtt.NewClientOptions()
	opts.SetConnectTimeout(time.Second * 5)
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", cfg.Mqtt.Host, cfg.Mqtt.Port))
	opts.SetClientID(uuid.NewString())
	opts.SetUsername(cfg.Mqtt.User)
	opts.SetPassword(cfg.Mqtt.Password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	topics := cfg.Monitoring.GetTopics()
	log.Println("Subscribing to topics: ", topics)

	for _, topic := range topics {
		sub(client, topic)
	}

	// https://gobyexample.com/signals
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	log.Println("Received termination signal...")
	client.Disconnect(250)
	log.Println("Disconnected")
}

func sub(client mqtt.Client, topic string) {
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	if token.Error() != nil {
		log.Fatal(token.Error())
	}
	log.Printf("Subscribed to topic: %s\n", topic)
}
