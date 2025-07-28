package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"iot-platform-go/internal/config"
	"iot-platform-go/internal/mqtt"
)

const (
	filePermission     = 0644
	connectionWaitTime = 2 * time.Second
	messageInterval    = 10 * time.Second
	baseTemperature    = 20.0
	baseHumidity       = 50.0
	humidityMultiplier = 2.0
)

func main() {
	// Create log file in the same directory as the executable
	logFile, err := os.OpenFile("mqtt-test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, filePermission)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Set log output to both file and console
	log.SetOutput(os.Stdout)

	// Function to log message to file
	logToFile := func(message string) {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		logEntry := fmt.Sprintf("[%s] %s\n", timestamp, message)
		if _, err := logFile.WriteString(logEntry); err != nil {
			log.Printf("Failed to write to log file: %v", err)
		}
	}

	// Load configuration
	cfg := config.Load()

	// Create MQTT client
	mqttConfig := cfg.MQTT
	mqttConfig.CleanSession = false
	mqttConfig.ClientID = "mqtt-test-" + time.Now().Format("20060102150405")
	client := mqtt.NewClient(&mqttConfig)

	// Connect to MQTT broker
	if err := client.Connect(); err != nil {
		log.Printf("Failed to connect to MQTT broker: %v", err)
		logFile.Close()
		os.Exit(1)
	}

	// Wait for connection
	time.Sleep(connectionWaitTime)

	if !client.IsConnected() {
		logFile.Close()
		client.Disconnect()
		log.Fatal("MQTT client is not connected")
	}

	log.Printf("✅ Connected to MQTT broker: %s", cfg.MQTT.Broker)

	// Subscribe to device topics
	err = client.Subscribe("devices/+/data", func(topic string, payload []byte) {
		message := fmt.Sprintf("📡 RECEIVED DEVICE DATA from %s: %s", topic, string(payload))
		log.Print(message)
		logToFile(message)
	})
	if err != nil {
		logFile.Close()
		log.Fatalf("Failed to subscribe to device data: %v", err)
	}

	err = client.Subscribe("devices/+/status", func(topic string, payload []byte) {
		message := fmt.Sprintf("📡 RECEIVED DEVICE STATUS from %s: %s", topic, string(payload))
		log.Print(message)
		logToFile(message)
	})
	if err != nil {
		logFile.Close()
		log.Fatalf("Failed to subscribe to device status: %v", err)
	}

	// Also subscribe to specific topics for testing
	err = client.Subscribe("devices/test-device/data", func(topic string, payload []byte) {
		message := fmt.Sprintf("📡 RECEIVED TEST DEVICE DATA from %s: %s", topic, string(payload))
		log.Print(message)
		logToFile(message)
	})
	if err != nil {
		logFile.Close()
		log.Fatalf("Failed to subscribe to test-device/data: %v", err)
	}

	err = client.Subscribe("devices/test-device/status", func(topic string, payload []byte) {
		message := fmt.Sprintf("📡 RECEIVED TEST DEVICE STATUS from %s: %s", topic, string(payload))
		log.Print(message)
		logToFile(message)
	})
	if err != nil {
		logFile.Close()
		log.Fatalf("Failed to subscribe to test-device/status: %v", err)
	}

	log.Println("✅ Subscribed to topics:")
	log.Println("   - devices/+/data (wildcard)")
	log.Println("   - devices/+/status (wildcard)")
	log.Println("   - devices/test-device/data (specific)")
	log.Println("   - devices/test-device/status (specific)")
	log.Println("")

	// Log startup message
	startupMessage := fmt.Sprintf("🚀 MQTT TEST CLIENT started at %s", time.Now().Format("2006-01-02 15:04:05"))
	log.Println(startupMessage)
	logToFile(startupMessage)

	// Send test messages every 10 seconds
	go func() {
		ticker := time.NewTicker(messageInterval)
		defer ticker.Stop()

		messageCount := 0
		for range ticker.C {
			messageCount++

			// Test different device IDs to verify wildcard matching
			deviceIDs := []string{"test-device", "sensor-001", "camera-002", "thermostat-003"}
			deviceID := deviceIDs[messageCount%len(deviceIDs)]

			// Send device data
			deviceData := map[string]interface{}{
				"device_id":   deviceID,
				"temperature": baseTemperature + float64(messageCount),
				"humidity":    baseHumidity + float64(messageCount)*humidityMultiplier,
				"timestamp":   time.Now().Format(time.RFC3339),
			}

			deviceDataJSON, _ := json.Marshal(deviceData)
			if err := client.Publish("devices/"+deviceID+"/data", string(deviceDataJSON)); err != nil {
				errorMsg := fmt.Sprintf("❌ Failed to publish device data: %v", err)
				log.Print(errorMsg)
				logToFile(errorMsg)
			} else {
				sentMsg := fmt.Sprintf("📤 SENT DEVICE DATA: %s", string(deviceDataJSON))
				log.Print(sentMsg)
				logToFile(sentMsg)
			}

			// Send device status
			deviceStatus := map[string]interface{}{
				"device_id": deviceID,
				"status":    "online",
				"last_seen": time.Now().Format(time.RFC3339),
			}

			deviceStatusJSON, _ := json.Marshal(deviceStatus)
			if err := client.Publish("devices/"+deviceID+"/status", string(deviceStatusJSON)); err != nil {
				errorMsg := fmt.Sprintf("❌ Failed to publish device status: %v", err)
				log.Print(errorMsg)
				logToFile(errorMsg)
			} else {
				sentMsg := fmt.Sprintf("📤 SENT DEVICE STATUS: %s", string(deviceStatusJSON))
				log.Print(sentMsg)
				logToFile(sentMsg)
			}

			log.Println("")
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("🔄 MQTT Test Client running...")
	log.Println("   - Receiving messages from devices/+/data and devices/+/status")
	log.Println("   - Sending test messages every 10 seconds")
	log.Println("   - Logs saved to: mqtt-test.log")
	log.Println("   - Press Ctrl+C to exit")
	log.Println("")

	// Wait for signal
	<-sigChan

	// Log shutdown message
	shutdownMessage := fmt.Sprintf("🛑 MQTT TEST CLIENT stopped at %s", time.Now().Format("2006-01-02 15:04:05"))
	log.Println("")
	log.Println(shutdownMessage)
	logToFile(shutdownMessage)

	log.Println("🛑 Shutting down MQTT Test Client...")
	client.Disconnect()
	logFile.Close()
}
