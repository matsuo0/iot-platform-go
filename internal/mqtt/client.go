package mqtt

import (
	"fmt"
	"log"
	"strings"
	"time"

	"iot-platform-go/internal/config"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	connectRetryInterval   = 5 * time.Second
	disconnectTimeout      = 250 // milliseconds
	connectionWaitTime     = 100 * time.Millisecond
	connectionWaitAttempts = 10
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
	opts.SetCleanSession(false) // Changed from c.config.CleanSession to false
	opts.SetAutoReconnect(c.config.AutoReconnect)
	opts.SetDefaultPublishHandler(c.defaultMessageHandler)

	// Add connection stability settings
	opts.SetMaxReconnectInterval(1 * time.Minute)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(connectRetryInterval)
	opts.SetOrderMatters(false)
	opts.SetResumeSubs(true)

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
		c.client.Disconnect(disconnectTimeout) // 250ms timeout
		log.Println("Disconnected from MQTT broker")
	}
}

// Subscribe subscribes to a topic
func (c *Client) Subscribe(topic string, handler MessageHandler) error {
	// Wait for connection to be established
	for i := 0; i < connectionWaitAttempts; i++ {
		if c.client.IsConnected() {
			break
		}
		time.Sleep(connectionWaitTime)
	}

	if !c.client.IsConnected() {
		return fmt.Errorf("MQTT client is not connected after waiting")
	}

	// Store handler
	c.handlers[topic] = handler

	// Subscribe to topic
	token := c.client.Subscribe(topic, c.config.QoS, func(client mqtt.Client, msg mqtt.Message) {
		// Find the appropriate handler for this topic
		// First try exact match
		if handler, exists := c.handlers[msg.Topic()]; exists {
			handler(msg.Topic(), msg.Payload())
			return
		}
		
		// Then try wildcard matches
		for pattern, handler := range c.handlers {
			if c.topicMatches(pattern, msg.Topic()) {
				handler(msg.Topic(), msg.Payload())
				return
			}
		}
		
		// If no handler found, use default handler
		c.defaultMessageHandler(client, msg)
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

// topicMatches checks if a topic matches a pattern (supports + and # wildcards)
func (c *Client) topicMatches(pattern, topic string) bool {
	// Simple wildcard matching implementation
	// This is a basic implementation - for production use a more robust MQTT topic matcher
	
	// Split both pattern and topic by '/'
	patternParts := strings.Split(pattern, "/")
	topicParts := strings.Split(topic, "/")
	
	// Handle # wildcard (matches everything after this point)
	if len(patternParts) > 0 && patternParts[len(patternParts)-1] == "#" {
		// Remove the # from pattern
		patternParts = patternParts[:len(patternParts)-1]
		// Check if topic starts with the pattern (excluding #)
		if len(topicParts) >= len(patternParts) {
			for i, part := range patternParts {
				if i >= len(topicParts) {
					return false
				}
				if part != "+" && part != topicParts[i] {
					return false
				}
			}
			return true
		}
		return false
	}
	
	// Handle + wildcard and exact matching
	if len(patternParts) != len(topicParts) {
		return false
	}
	
	for i, patternPart := range patternParts {
		if patternPart != "+" && patternPart != topicParts[i] {
			return false
		}
	}
	
	return true
}
