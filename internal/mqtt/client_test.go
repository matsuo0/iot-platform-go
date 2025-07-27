package mqtt

import (
	"testing"
	"time"

	"iot-platform-go/internal/config"
)

func TestNewClient(t *testing.T) {
	cfg := &config.MQTTConfig{
		Broker:        "tcp://localhost:1883",
		ClientID:      "test-client",
		KeepAlive:     60,
		ConnectTimeout: 30,
		QoS:           1,
		CleanSession:  true,
		AutoReconnect: true,
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
	// Skip if MQTT broker is not running
	t.Skip("Skipping connection test - requires running MQTT broker")

	cfg := &config.MQTTConfig{
		Broker:        "tcp://localhost:1883",
		ClientID:      "test-client-" + time.Now().Format("20060102150405"),
		KeepAlive:     60,
		ConnectTimeout: 30,
		QoS:           1,
		CleanSession:  true,
		AutoReconnect: true,
	}

	client := NewClient(cfg)

	// Test connection
	err := client.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
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
		Broker:        "tcp://localhost:1883",
		ClientID:      "test-client",
		KeepAlive:     60,
		ConnectTimeout: 30,
		QoS:           1,
		CleanSession:  true,
		AutoReconnect: true,
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