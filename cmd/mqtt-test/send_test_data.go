package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"iot-platform-go/internal/config"
	"iot-platform-go/internal/mqtt"
)

// DeviceDataMessage represents device data structure
type DeviceDataMessage struct {
	DeviceID  string                 `json:"device_id"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// DeviceStatusMessage represents device status structure
type DeviceStatusMessage struct {
	DeviceID string                 `json:"device_id"`
	Status   string                 `json:"status"`
	LastSeen string                 `json:"last_seen"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

const (
	// Test data generation constants
	dataSendInterval   = 5 * time.Second
	statusSendInterval = 3 // Send status every 3 data batches
	temperatureBase    = 20.0
	temperatureRange   = 10.0
	humidityBase       = 40.0
	humidityRange      = 30.0
	pressureBase       = 1000.0
	pressureRange      = 50.0
	voltageBase        = 3.0
	voltageRange       = 0.5
	batteryBase        = 80
	batteryRange       = 20
	signalBase         = 70
	signalRange        = 30
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create MQTT client
	mqttConfig := cfg.MQTT
	mqttConfig.ClientID = "test-sender-" + time.Now().Format("20060102150405")
	client := mqtt.NewClient(&mqttConfig)

	// Connect to MQTT broker
	log.Printf("Connecting to MQTT broker: %s", mqttConfig.Broker)
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", err)
	}
	defer client.Disconnect()

	log.Println("✅ Connected to MQTT broker")

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start sending test data
	go sendTestData(client)

	// Wait for shutdown signal
	<-sigChan
	log.Println("🛑 Shutting down test sender...")
}

func sendTestData(client *mqtt.Client) {
	// Use the created device ID
	deviceIDs := []string{
		"0a0e35e6-eeba-49ea-a02f-444a722fabe1", // Test Temperature Sensor
	}

	statuses := []string{"online", "offline", "error", "maintenance"}

	ticker := time.NewTicker(dataSendInterval) // Send data every 5 seconds
	defer ticker.Stop()

	counter := 0
	for range ticker.C {
		counter++

		// Send device data
		for _, deviceID := range deviceIDs {
			// Generate random sensor data
			data := map[string]interface{}{
				"temperature": temperatureBase + rand.Float64()*temperatureRange, // 20-30°C
				"humidity":    humidityBase + rand.Float64()*humidityRange,       // 40-70%
				"pressure":    pressureBase + rand.Float64()*pressureRange,       // 1000-1050 hPa
				"voltage":     voltageBase + rand.Float64()*voltageRange,         // 3.0-3.5V
			}

			deviceData := DeviceDataMessage{
				DeviceID:  deviceID,
				Timestamp: time.Now().Format(time.RFC3339),
				Data:      data,
				Metadata: map[string]interface{}{
					"sequence": counter,
					"quality":  "good",
				},
			}

			payload, err := json.Marshal(deviceData)
			if err != nil {
				log.Printf("❌ Failed to marshal device data: %v", err)
				continue
			}

			topic := fmt.Sprintf("devices/%s/data", deviceID)
			if err := client.Publish(topic, payload); err != nil {
				log.Printf("❌ Failed to publish device data: %v", err)
			} else {
				log.Printf("📤 Sent device data to %s", topic)
			}
		}

		// Send device status (less frequently)
		if counter%statusSendInterval == 0 { // Every 15 seconds
			for _, deviceID := range deviceIDs {
				status := statuses[rand.Intn(len(statuses))]

				deviceStatus := DeviceStatusMessage{
					DeviceID: deviceID,
					Status:   status,
					LastSeen: time.Now().Format(time.RFC3339),
					Metadata: map[string]interface{}{
						"battery": batteryBase + rand.Intn(batteryRange), // 80-100%
						"signal":  signalBase + rand.Intn(signalRange),   // 70-100%
					},
				}

				payload, err := json.Marshal(deviceStatus)
				if err != nil {
					log.Printf("❌ Failed to marshal device status: %v", err)
					continue
				}

				topic := fmt.Sprintf("devices/%s/status", deviceID)
				if err := client.Publish(topic, payload); err != nil {
					log.Printf("❌ Failed to publish device status: %v", err)
				} else {
					log.Printf("📤 Sent device status to %s: %s", topic, status)
				}
			}
		}

		log.Printf("📊 Sent test data batch #%d", counter)
	}
}
