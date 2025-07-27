package main

import (
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
)

func main() {
	// Create log file
	logFile, err := os.OpenFile("mqtt-receiver.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, filePermission)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Set log output to both file and console
	log.SetOutput(os.Stdout)

	// Load configuration
	cfg := config.Load()

	// Create MQTT client
	mqttConfig := cfg.MQTT
	mqttConfig.CleanSession = false
	mqttConfig.ClientID = "mqtt-receiver-" + time.Now().Format("20060102150405")
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

	log.Printf("✅ RECEIVER Connected to MQTT broker: %s", cfg.MQTT.Broker)

	// Function to log message to file
	logToFile := func(message string) {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		logEntry := fmt.Sprintf("[%s] %s\n", timestamp, message)
		if _, err := logFile.WriteString(logEntry); err != nil {
			log.Printf("Failed to write to log file: %v", err)
		}
	}

	// Subscribe to exact topics (no wildcard for testing)
	err = client.Subscribe("devices/device001/data", func(topic string, payload []byte) {
		message := fmt.Sprintf("📡 RECEIVED DEVICE DATA from %s: %s", topic, string(payload))
		log.Print(message)
		logToFile(message)
	})
	if err != nil {
		logFile.Close()
		log.Fatalf("Failed to subscribe to device001/data: %v", err)
	}

	err = client.Subscribe("devices/device001/status", func(topic string, payload []byte) {
		message := fmt.Sprintf("📡 RECEIVED DEVICE STATUS from %s: %s", topic, string(payload))
		log.Print(message)
		logToFile(message)
	})
	if err != nil {
		logFile.Close()
		log.Fatalf("Failed to subscribe to device001/status: %v", err)
	}

	err = client.Subscribe("devices/device002/data", func(topic string, payload []byte) {
		message := fmt.Sprintf("📡 RECEIVED DEVICE DATA from %s: %s", topic, string(payload))
		log.Print(message)
		logToFile(message)
	})
	if err != nil {
		logFile.Close()
		log.Fatalf("Failed to subscribe to device002/data: %v", err)
	}

	log.Println("✅ RECEIVER Subscribed to topics:")
	log.Println("   - devices/device001/data")
	log.Println("   - devices/device001/status")
	log.Println("   - devices/device002/data")
	log.Println("")

	// Log startup message
	startupMessage := fmt.Sprintf("🚀 MQTT RECEIVER started at %s", time.Now().Format("2006-01-02 15:04:05"))
	log.Println(startupMessage)
	logToFile(startupMessage)

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("🔄 MQTT RECEIVER running...")
	log.Println("   - Waiting for messages...")
	log.Println("   - Logs saved to: mqtt-receiver.log")
	log.Println("   - Press Ctrl+C to exit")
	log.Println("")

	// Wait for signal
	<-sigChan

	// Log shutdown message
	shutdownMessage := fmt.Sprintf("🛑 MQTT RECEIVER stopped at %s", time.Now().Format("2006-01-02 15:04:05"))
	log.Println("")
	log.Println(shutdownMessage)
	logToFile(shutdownMessage)

	log.Println("🛑 Shutting down MQTT RECEIVER...")
	client.Disconnect()
	logFile.Close()
}
