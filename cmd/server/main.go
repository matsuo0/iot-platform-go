package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"iot-platform-go/internal/api"
	"iot-platform-go/internal/config"
	"iot-platform-go/internal/database"
	"iot-platform-go/internal/device"
	"iot-platform-go/internal/mqtt"

	"github.com/gin-gonic/gin"
)

// Global MQTT client for health check access
var mqttClient *mqtt.Client

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Initialize repositories
	deviceRepo := device.NewRepository(db)

	// Initialize handlers
	deviceHandler := api.NewDeviceHandler(deviceRepo)

	// Initialize MQTT client
	mqttConfig := cfg.MQTT
	mqttConfig.ClientID = "iot-platform-server-" + time.Now().Format("20060102150405")
	mqttClient = mqtt.NewClient(&mqttConfig)

	// Connect to MQTT broker
	log.Printf("Connecting to MQTT broker: %s", cfg.MQTT.Broker)
	if err := mqttClient.Connect(); err != nil {
		log.Printf("Failed to connect to MQTT broker: %v", err)
		log.Printf("Server will start without MQTT functionality")
	} else {
		log.Printf("✅ Successfully connected to MQTT broker")
		
		// Wait for connection to be established
		time.Sleep(2 * time.Second)
		
		if mqttClient.IsConnected() {
			log.Printf("✅ MQTT client is ready")
		} else {
			log.Printf("⚠️ MQTT client connection failed")
		}
	}

	// Setup Gin router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		mqttStatus := "disconnected"
		if mqttClient != nil && mqttClient.IsConnected() {
			mqttStatus = "connected"
		}
		
		c.JSON(http.StatusOK, gin.H{
			"status":     "ok",
			"message":    "IoT Platform is running",
			"mqtt_status": mqttStatus,
			"timestamp":  time.Now().Format(time.RFC3339),
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Device routes
		devices := api.Group("/devices")
		{
			devices.POST("", deviceHandler.CreateDevice)
			devices.GET("", deviceHandler.GetAllDevices)
			devices.GET("/:id", deviceHandler.GetDevice)
			devices.PUT("/:id", deviceHandler.UpdateDevice)
			devices.DELETE("/:id", deviceHandler.DeleteDevice)
			devices.GET("/:id/status", deviceHandler.GetDeviceStatus)
		}
	}

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
		log.Printf("Starting server on %s", addr)
		log.Printf("Health check: http://%s/health", addr)
		log.Printf("API documentation: http://%s/api", addr)

		if err := router.Run(addr); err != nil {
			log.Printf("Failed to start server: %v", err)
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("🛑 Shutting down IoT Platform...")
	
	// Disconnect MQTT client
	if mqttClient != nil && mqttClient.IsConnected() {
		mqttClient.Disconnect()
		log.Println("✅ MQTT client disconnected")
	}
	
	log.Println("✅ Server shutdown complete")
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
