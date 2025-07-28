package mqtt

import (
	"fmt"
	"os"
	"testing"
	"time"

	"iot-platform-go/internal/config"
)

func TestNewClient(t *testing.T) {
	cfg := &config.MQTTConfig{
		Broker:         "tcp://localhost:1883",
		ClientID:       "test-client",
		KeepAlive:      60,
		ConnectTimeout: 30,
		QoS:            1,
		CleanSession:   true,
		AutoReconnect:  true,
	}

	client := NewClient(cfg)
	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	if client.config != cfg {
		t.Error("Expected config to be set correctly")
	}

	if len(client.handlers) != 0 {
		t.Error("Expected handlers map to be empty")
	}
}

func TestClientConnection(t *testing.T) {
	// Skip this test in CI/CD environment
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping MQTT connection test in CI environment")
	}

	// Test connection to running MQTT broker
	cfg := &config.MQTTConfig{
		Broker:         "tcp://localhost:1883",
		ClientID:       "test-client-" + time.Now().Format("20060102150405"),
		KeepAlive:      60,
		ConnectTimeout: 30,
		QoS:            1,
		CleanSession:   true,
		AutoReconnect:  true,
	}

	client := NewClient(cfg)

	// Test connection with timeout
	connectChan := make(chan error, 1)
	go func() {
		connectChan <- client.Connect()
	}()

	select {
	case err := <-connectChan:
		if err != nil {
			t.Skipf("Skipping test - MQTT broker not available: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Skip("Skipping test - MQTT broker connection timeout")
	}

	// Test connection status
	if !client.IsConnected() {
		t.Error("Expected client to be connected")
	}

	// Test disconnect
	client.Disconnect()
	if client.IsConnected() {
		t.Error("Expected client to be disconnected")
	}
}

func TestMessageHandler(t *testing.T) {
	cfg := &config.MQTTConfig{
		Broker:         "tcp://localhost:1883",
		ClientID:       "test-client",
		KeepAlive:      60,
		ConnectTimeout: 30,
		QoS:            1,
		CleanSession:   true,
		AutoReconnect:  true,
	}

	client := NewClient(cfg)

	// Test message handler registration
	messageReceived := false
	handler := func(topic string, payload []byte) {
		messageReceived = true
	}

	// This would normally be called after connection
	client.handlers["test/topic"] = handler

	if len(client.handlers) != 1 {
		t.Error("Expected handler to be registered")
	}

	// Test handler execution
	client.handlers["test/topic"]("test/topic", []byte("test message"))
	if !messageReceived {
		t.Error("Expected message handler to be called")
	}
}

func TestMessagePublishSubscribe(t *testing.T) {
	// Skip this test in CI/CD environment
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping MQTT publish/subscribe test in CI environment")
	}

	// Create two clients for testing publish/subscribe
	publisher := NewClient(&config.MQTTConfig{
		Broker:         "tcp://localhost:1883",
		ClientID:       "test-publisher-" + time.Now().Format("20060102150405"),
		KeepAlive:      60,
		ConnectTimeout: 30,
		QoS:            1,
		CleanSession:   true,
		AutoReconnect:  true,
	})

	subscriber := NewClient(&config.MQTTConfig{
		Broker:         "tcp://localhost:1883",
		ClientID:       "test-subscriber-" + time.Now().Format("20060102150405"),
		KeepAlive:      60,
		ConnectTimeout: 30,
		QoS:            1,
		CleanSession:   true,
		AutoReconnect:  true,
	})

	// Connect both clients with timeout
	connectChan := make(chan error, 2)
	go func() {
		connectChan <- publisher.Connect()
	}()
	go func() {
		connectChan <- subscriber.Connect()
	}()

	// Wait for both connections with timeout
	select {
	case err := <-connectChan:
		if err != nil {
			t.Skipf("Skipping test - MQTT broker not available: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Skip("Skipping test - MQTT broker connection timeout")
	}

	select {
	case err := <-connectChan:
		if err != nil {
			t.Skipf("Skipping test - MQTT broker not available: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Skip("Skipping test - MQTT broker connection timeout")
	}

	defer publisher.Disconnect()
	defer subscriber.Disconnect()

	// Test topic
	topic := "test/message/" + time.Now().Format("20060102150405")
	expectedMessage := "Hello MQTT!"

	// Channel to receive message
	messageReceived := make(chan string, 1)

	// Subscribe to topic
	err := subscriber.Subscribe(topic, func(topic string, payload []byte) {
		messageReceived <- string(payload)
	})
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// Wait a bit for subscription to be established
	time.Sleep(100 * time.Millisecond)

	// Publish message
	err = publisher.Publish(topic, expectedMessage)
	if err != nil {
		t.Fatalf("Failed to publish: %v", err)
	}

	// Wait for message to be received
	select {
	case receivedMessage := <-messageReceived:
		if receivedMessage != expectedMessage {
			t.Errorf("Expected message '%s', got '%s'", expectedMessage, receivedMessage)
		}
	case <-time.After(5 * time.Second):
		t.Error("Timeout waiting for message")
	}
}

func TestMultipleSubscribers(t *testing.T) {
	// Skip this test in CI/CD environment
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping MQTT multiple subscribers test in CI environment")
	}

	// Create multiple subscribers
	subscribers := make([]*Client, 3)
	receivedMessages := make([]chan string, 3)

	// Cleanup function to disconnect all subscribers
	cleanup := func() {
		for _, subscriber := range subscribers {
			if subscriber != nil {
				subscriber.Disconnect()
			}
		}
	}
	defer cleanup()

	for i := 0; i < 3; i++ {
		subscribers[i] = NewClient(&config.MQTTConfig{
			Broker:         "tcp://localhost:1883",
			ClientID:       fmt.Sprintf("test-subscriber-%d-%s", i, time.Now().Format("20060102150405")),
			KeepAlive:      60,
			ConnectTimeout: 30,
			QoS:            1,
			CleanSession:   true,
			AutoReconnect:  true,
		})

		receivedMessages[i] = make(chan string, 1)

		// Connect with timeout
		connectChan := make(chan error, 1)
		go func(client *Client) {
			connectChan <- client.Connect()
		}(subscribers[i])

		select {
		case err := <-connectChan:
			if err != nil {
				t.Skipf("Skipping test - MQTT broker not available: %v", err)
			}
		case <-time.After(10 * time.Second):
			t.Skip("Skipping test - MQTT broker connection timeout")
		}
	}

	// Create publisher
	publisher := NewClient(&config.MQTTConfig{
		Broker:         "tcp://localhost:1883",
		ClientID:       "test-publisher-multi-" + time.Now().Format("20060102150405"),
		KeepAlive:      60,
		ConnectTimeout: 30,
		QoS:            1,
		CleanSession:   true,
		AutoReconnect:  true,
	})

	// Connect publisher with timeout
	connectChan := make(chan error, 1)
	go func() {
		connectChan <- publisher.Connect()
	}()

	select {
	case err := <-connectChan:
		if err != nil {
			t.Skipf("Skipping test - MQTT broker not available: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Skip("Skipping test - MQTT broker connection timeout")
	}

	defer publisher.Disconnect()

	// Test topic
	topic := "test/multi/" + time.Now().Format("20060102150405")
	expectedMessage := "Hello Multiple Subscribers!"

	// Subscribe all subscribers to the same topic
	for i, subscriber := range subscribers {
		err := subscriber.Subscribe(topic, func(topic string, payload []byte) {
			receivedMessages[i] <- string(payload)
		})
		if err != nil {
			t.Fatalf("Failed to subscribe subscriber %d: %v", i, err)
		}
	}

	// Wait a bit for subscriptions to be established
	time.Sleep(100 * time.Millisecond)

	// Publish message
	err := publisher.Publish(topic, expectedMessage)
	if err != nil {
		t.Fatalf("Failed to publish: %v", err)
	}

	// Wait for all subscribers to receive the message
	for i, messageChan := range receivedMessages {
		select {
		case receivedMessage := <-messageChan:
			if receivedMessage != expectedMessage {
				t.Errorf("Subscriber %d: Expected message '%s', got '%s'", i, expectedMessage, receivedMessage)
			}
		case <-time.After(5 * time.Second):
			t.Errorf("Subscriber %d: Timeout waiting for message", i)
		}
	}
}
