package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"iot-platform-go/internal/api"
	"iot-platform-go/internal/config"
	"iot-platform-go/internal/database"
	"iot-platform-go/internal/device"
	"iot-platform-go/internal/mqtt"

	"github.com/gin-gonic/gin"
)

// Device data structure for MQTT messages
type DeviceDataMessage struct {
	DeviceID   string                 `json:"device_id"`
	Timestamp  string                 `json:"timestamp"`
	Data       map[string]interface{} `json:"data"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// Device status structure for MQTT messages
type DeviceStatusMessage struct {
	DeviceID  string `json:"device_id"`
	Status    string `json:"status"`
	LastSeen  string `json:"last_seen"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// Application holds all dependencies
type Application struct {
	config     *config.Config
	db         *database.Database
	deviceRepo *device.Repository
	mqttClient *mqtt.Client
	router     *gin.Engine
	server     *http.Server
}

// NewApplication creates a new application instance
func NewApplication(cfg *config.Config) (*Application, error) {
	// Initialize database
	db, err := database.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	// Initialize repositories
	deviceRepo := device.NewRepository(db)

	// Initialize MQTT client
	mqttConfig := cfg.MQTT
	mqttConfig.CleanSession = false
	mqttConfig.ClientID = "iot-platform-server-" + time.Now().Format("20060102150405")
	mqttClient := mqtt.NewClient(&mqttConfig)

	// Setup Gin router
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	app := &Application{
		config:     cfg,
		db:         db,
		deviceRepo: deviceRepo,
		mqttClient: mqttClient,
		router:     router,
	}

	// Setup routes
	app.setupRoutes()

	return app, nil
}

// setupRoutes configures all application routes
func (app *Application) setupRoutes() {
	// Health check endpoint
	app.router.GET("/health", app.healthCheckHandler)

	// API routes
	apiGroup := app.router.Group("/api")
	{
		// Device routes
		deviceHandler := api.NewDeviceHandler(app.deviceRepo)
		devices := apiGroup.Group("/devices")
		{
			devices.POST("", deviceHandler.CreateDevice)
			devices.GET("", deviceHandler.GetAllDevices)
			devices.GET("/:id", deviceHandler.GetDevice)
			devices.PUT("/:id", deviceHandler.UpdateDevice)
			devices.DELETE("/:id", deviceHandler.DeleteDevice)
			devices.GET("/:id/status", deviceHandler.GetDeviceStatus)
		}
	}
}

// healthCheckHandler handles health check requests
func (app *Application) healthCheckHandler(c *gin.Context) {
	mqttStatus := "disconnected"
	if app.mqttClient != nil && app.mqttClient.IsConnected() {
		mqttStatus = "connected"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "ok",
		"message":     "IoT Platform is running",
		"mqtt_status": mqttStatus,
		"timestamp":   time.Now().Format(time.RFC3339),
	})
}

// Start initializes and starts the application
func (app *Application) Start() error {
	// Connect to MQTT broker
	log.Printf("Connecting to MQTT broker: %s", app.config.MQTT.Broker)
	if err := app.mqttClient.Connect(); err != nil {
		log.Printf("Failed to connect to MQTT broker: %v", err)
		log.Printf("Server will start without MQTT functionality")
	} else {
		log.Printf("✅ Successfully connected to MQTT broker")

		// Wait for connection to be established
		time.Sleep(2 * time.Second)

		if app.mqttClient.IsConnected() {
			log.Printf("✅ MQTT client is ready")
			
			// Subscribe to MQTT topics
			if err := app.subscribeToMQTTTopics(); err != nil {
				log.Printf("⚠️ Failed to subscribe to MQTT topics: %v", err)
			} else {
				log.Printf("✅ Successfully subscribed to MQTT topics")
			}
		} else {
			log.Printf("⚠️ MQTT client connection failed")
		}
	}

	// Setup HTTP server
	addr := fmt.Sprintf("%s:%s", app.config.Server.Host, app.config.Server.Port)
	app.server = &http.Server{
		Addr:    addr,
		Handler: app.router,
	}

	log.Printf("Starting server on %s", addr)
	log.Printf("Health check: http://%s/health", addr)
	log.Printf("API documentation: http://%s/api", addr)

	return app.server.ListenAndServe()
}

// Stop gracefully shuts down the application
func (app *Application) Stop(ctx context.Context) error {
	log.Println("🛑 Shutting down IoT Platform...")

	// Disconnect MQTT client
	if app.mqttClient != nil && app.mqttClient.IsConnected() {
		app.mqttClient.Disconnect()
		log.Println("✅ MQTT client disconnected")
	}

	// Close database
	if app.db != nil {
		if err := app.db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}

	// Shutdown HTTP server
	if app.server != nil {
		if err := app.server.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down server: %v", err)
		}
	}

	log.Println("✅ Server shutdown complete")
	return nil
}

// subscribeToMQTTTopics subscribes to device data and status topics
func (app *Application) subscribeToMQTTTopics() error {
	// Subscribe to device data topics with wildcard
	if err := app.mqttClient.Subscribe("devices/+/data", app.handleDeviceData); err != nil {
		return fmt.Errorf("failed to subscribe to device data topics: %v", err)
	}

	// Subscribe to device status topics with wildcard
	if err := app.mqttClient.Subscribe("devices/+/status", app.handleDeviceStatus); err != nil {
		return fmt.Errorf("failed to subscribe to device status topics: %v", err)
	}

	// Subscribe to all device topics (optional - for debugging)
	if err := app.mqttClient.Subscribe("devices/#", app.handleAllDeviceMessages); err != nil {
		log.Printf("⚠️ Failed to subscribe to all device topics: %v", err)
	}

	log.Println("📡 Subscribed to MQTT topics:")
	log.Println("   - devices/+/data (device data)")
	log.Println("   - devices/+/status (device status)")
	log.Println("   - devices/# (all device messages - debug)")

	return nil
}

// handleDeviceData processes incoming device data messages
func (app *Application) handleDeviceData(topic string, payload []byte) {
	msg := fmt.Sprintf("📡 RECEIVED DEVICE DATA from %s: %s", topic, string(payload))
	log.Println(msg)
	logToFile(msg)

	// Parse the JSON payload
	var deviceData DeviceDataMessage
	if err := json.Unmarshal(payload, &deviceData); err != nil {
		log.Printf("❌ Failed to parse device data JSON: %v", err)
		log.Printf("   Raw payload: %s", string(payload))
		return
	}

	// Validate required fields
	if deviceData.DeviceID == "" {
		log.Printf("❌ Device data missing required field: device_id")
		return
	}

	if deviceData.Timestamp == "" {
		log.Printf("❌ Device data missing required field: timestamp")
		return
	}

	// Parse timestamp
	timestamp, err := time.Parse(time.RFC3339, deviceData.Timestamp)
	if err != nil {
		log.Printf("❌ Failed to parse timestamp '%s': %v", deviceData.Timestamp, err)
		return
	}

	// Log the received data
	log.Printf("✅ Processed device data:")
	log.Printf("   Device ID: %s", deviceData.DeviceID)
	log.Printf("   Timestamp: %s", timestamp.Format(time.RFC3339))
	log.Printf("   Data points: %d", len(deviceData.Data))

	// TODO: Save to database (will be implemented in next step)
	log.Printf("📊 Device data ready for database storage")
}

// handleDeviceStatus processes incoming device status messages
func (app *Application) handleDeviceStatus(topic string, payload []byte) {
	msg := fmt.Sprintf("📡 RECEIVED DEVICE STATUS from %s: %s", topic, string(payload))
	log.Println(msg)
	logToFile(msg)

	// Parse the JSON payload
	var deviceStatus DeviceStatusMessage
	if err := json.Unmarshal(payload, &deviceStatus); err != nil {
		log.Printf("❌ Failed to parse device status JSON: %v", err)
		log.Printf("   Raw payload: %s", string(payload))
		return
	}

	// Validate required fields
	if deviceStatus.DeviceID == "" {
		log.Printf("❌ Device status missing required field: device_id")
		return
	}

	if deviceStatus.Status == "" {
		log.Printf("❌ Device status missing required field: status")
		return
	}

	// Parse last seen timestamp if provided
	var lastSeen time.Time
	var err error
	if deviceStatus.LastSeen != "" {
		lastSeen, err = time.Parse(time.RFC3339, deviceStatus.LastSeen)
		if err != nil {
			log.Printf("❌ Failed to parse last_seen timestamp '%s': %v", deviceStatus.LastSeen, err)
			lastSeen = time.Now()
		}
	} else {
		lastSeen = time.Now()
	}

	// Log the received status
	log.Printf("✅ Processed device status:")
	log.Printf("   Device ID: %s", deviceStatus.DeviceID)
	log.Printf("   Status: %s", deviceStatus.Status)
	log.Printf("   Last Seen: %s", lastSeen.Format(time.RFC3339))

	// TODO: Update device status in database (will be implemented in next step)
	log.Printf("📊 Device status ready for database update")
}

// handleAllDeviceMessages processes all device messages for debugging
func (app *Application) handleAllDeviceMessages(topic string, payload []byte) {
	// Only log if it's not already handled by specific handlers
	if !strings.HasSuffix(topic, "/data") && !strings.HasSuffix(topic, "/status") {
		msg := fmt.Sprintf("📡 RECEIVED OTHER DEVICE MESSAGE from %s: %s", topic, string(payload))
		log.Println(msg)
		logToFile(msg)
	}
}

func main() {
	// Load configuration
	cfg := config.Load()

	// Create application
	app, err := NewApplication(cfg)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := app.Start(); err != nil && err != http.ErrServerClosed {
			log.Printf("Failed to start server: %v", err)
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal
	<-sigChan

	// Create context for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Stop application
	if err := app.Stop(ctx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// --- ファイル出力用の関数 ---
func logToFile(message string) {
	logFile, err := os.OpenFile("cmd/server/mqtt-received.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open mqtt-received.log: %v", err)
		return
	}
	defer logFile.Close()
	
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] %s\n", timestamp, message)
	if _, err := logFile.WriteString(logEntry); err != nil {
		log.Printf("Failed to write to mqtt-received.log: %v", err)
	}
}
