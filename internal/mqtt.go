package internal

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"

	cfg "monitoring/mqtt/config"
)

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	slog.Info("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	slog.Info(fmt.Sprintf("Connections lost: %v", err))
}

func RunMqtt(cfg cfg.Config) {
	opts := mqtt.NewClientOptions()
	opts.SetConnectTimeout(time.Second * 5)
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", cfg.Mqtt.Host, cfg.Mqtt.Port))
	opts.SetClientID(uuid.NewString())
	opts.SetUsername(cfg.Mqtt.User)
	opts.SetPassword(cfg.Mqtt.Password)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	monitors := make([]*Monitor, 0, len(cfg.Monitoring.Topics))
	tickers := make([]*time.Ticker, 0, len(cfg.Monitoring.Topics))

	for _, topic := range cfg.Monitoring.Topics {
		monitor := &Monitor{Count: 0, Topic: topic}
		monitors = append(monitors, monitor)
		sub(client, monitor)

		ticker := time.NewTicker(topic.Interval)
		defer ticker.Stop()
		tickers = append(tickers, ticker)

		go func(m *Monitor, t *time.Ticker) {
			for range t.C {
				m.CheckCount()
			}
		}(monitor, ticker)
	}
	// https://gobyexample.com/signals
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	slog.Info("Received termination signal...")
	client.Disconnect(250)
	slog.Info("Disconnected")
}

func sub(client mqtt.Client, m *Monitor) {
	token := client.Subscribe(m.Topic.Name, 1, func(c mqtt.Client, msg mqtt.Message) {
		m.mutex.Lock()
		m.Count++
		m.mutex.Unlock()
		slog.Info(fmt.Sprintf("Message: %s | Topic: %s | Count: %d\n", msg.Payload(), msg.Topic(), m.Count))
	})
	token.Wait()
	if token.Error() != nil {
		slog.Error(fmt.Sprintf("%v", token.Error()))
		os.Exit(1)
	}
	slog.Info(fmt.Sprintf("Subscribed to topic: %s\n", m.Topic.Name))
}

type Monitor struct {
	Count int
	Topic cfg.TopicConfig
	mutex sync.RWMutex
}

func (m *Monitor) CheckCount() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.Count < m.Topic.MinCount {
		slog.Warn(fmt.Sprintf("Not enough messages | %v\n", m.Topic.Name))
	}
	m.Count = 0
}
