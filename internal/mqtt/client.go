package mqtt

import (
	"fmt"
	"log"
	"time"

	"iot-platform-go/internal/config"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Client represents an MQTT client
type Client struct {
	client   mqtt.Client
	config   *config.MQTTConfig
	handlers map[string]MessageHandler
}

// MessageHandler is a function type for handling MQTT messages
type MessageHandler func(topic string, payload []byte)

// NewClient creates a new MQTT client
func NewClient(cfg *config.MQTTConfig) *Client {
	return &Client{
		config:   cfg,
		handlers: make(map[string]MessageHandler),
	}
}

// Connect establishes a connection to the MQTT broker
func (c *Client) Connect() error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(c.config.Broker)
	opts.SetClientID(c.config.ClientID)
	opts.SetKeepAlive(time.Duration(c.config.KeepAlive) * time.Second)
	opts.SetConnectTimeout(time.Duration(c.config.ConnectTimeout) * time.Second)
	opts.SetCleanSession(c.config.CleanSession)
	opts.SetAutoReconnect(c.config.AutoReconnect)
	opts.SetDefaultPublishHandler(c.defaultMessageHandler)

	// Set credentials if provided
	if c.config.Username != "" {
		opts.SetUsername(c.config.Username)
		opts.SetPassword(c.config.Password)
	}

	// Create client
	c.client = mqtt.NewClient(opts)

	// Connect to broker
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %v", token.Error())
	}

	log.Printf("Connected to MQTT broker: %s", c.config.Broker)
	return nil
}

// Disconnect closes the MQTT connection
func (c *Client) Disconnect() {
	if c.client != nil && c.client.IsConnected() {
		c.client.Disconnect(250) // 250ms timeout
		log.Println("Disconnected from MQTT broker")
	}
}

// Subscribe subscribes to a topic
func (c *Client) Subscribe(topic string, handler MessageHandler) error {
	if !c.client.IsConnected() {
		return fmt.Errorf("MQTT client is not connected")
	}

	// Store handler
	c.handlers[topic] = handler

	// Subscribe to topic
	token := c.client.Subscribe(topic, c.config.QoS, func(client mqtt.Client, msg mqtt.Message) {
		if handler, exists := c.handlers[msg.Topic()]; exists {
			handler(msg.Topic(), msg.Payload())
		}
	})

	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %v", topic, token.Error())
	}

	log.Printf("Subscribed to topic: %s", topic)
	return nil
}

// Unsubscribe unsubscribes from a topic
func (c *Client) Unsubscribe(topic string) error {
	if !c.client.IsConnected() {
		return fmt.Errorf("MQTT client is not connected")
	}

	token := c.client.Unsubscribe(topic)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to unsubscribe from topic %s: %v", topic, token.Error())
	}

	// Remove handler
	delete(c.handlers, topic)

	log.Printf("Unsubscribed from topic: %s", topic)
	return nil
}

// Publish publishes a message to a topic
func (c *Client) Publish(topic string, payload interface{}) error {
	if !c.client.IsConnected() {
		return fmt.Errorf("MQTT client is not connected")
	}

	token := c.client.Publish(topic, c.config.QoS, false, payload)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish to topic %s: %v", topic, token.Error())
	}

	log.Printf("Published message to topic: %s", topic)
	return nil
}

// IsConnected returns true if the client is connected
func (c *Client) IsConnected() bool {
	return c.client != nil && c.client.IsConnected()
}

// defaultMessageHandler handles messages that don't have a specific handler
func (c *Client) defaultMessageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message on topic %s: %s", msg.Topic(), string(msg.Payload()))
} 