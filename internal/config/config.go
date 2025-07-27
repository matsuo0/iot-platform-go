package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	defaultKeepAlive      = 60
	defaultConnectTimeout = 30
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	MQTT     MQTTConfig
	JWT      JWTConfig
	Logging  LoggingConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
	Host string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

// MQTTConfig holds MQTT configuration
type MQTTConfig struct {
	Broker         string
	ClientID       string
	Username       string
	Password       string
	KeepAlive      int
	ConnectTimeout int
	QoS            byte
	CleanSession   bool
	AutoReconnect  bool
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret     string
	Expiration string
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level string
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "localhost"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "iot_platform"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		MQTT: MQTTConfig{
			Broker:         getEnv("MQTT_BROKER", "tcp://localhost:1883"),
			ClientID:       getEnv("MQTT_CLIENT_ID", "iot-platform-server"),
			Username:       getEnv("MQTT_USERNAME", ""),
			Password:       getEnv("MQTT_PASSWORD", ""),
			KeepAlive:      getEnvAsInt("MQTT_KEEP_ALIVE", defaultKeepAlive),
			ConnectTimeout: getEnvAsInt("MQTT_CONNECT_TIMEOUT", defaultConnectTimeout),
			QoS:            getEnvAsByte("MQTT_QOS", 1),
			CleanSession:   getEnvAsBool("MQTT_CLEAN_SESSION", true),
			AutoReconnect:  getEnvAsBool("MQTT_AUTO_RECONNECT", true),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key-here"),
			Expiration: getEnv("JWT_EXPIRATION", "24h"),
		},
		Logging: LoggingConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool gets an environment variable as a boolean or returns a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getEnvAsByte gets an environment variable as a byte or returns a default value
func getEnvAsByte(key string, defaultValue byte) byte {
	if value := os.Getenv(key); value != "" {
		if byteValue, err := strconv.ParseUint(value, 10, 8); err == nil {
			return byte(byteValue)
		}
	}
	return defaultValue
}

// GetDatabaseURL returns the database connection string
func (c *Config) GetDatabaseURL() string {
	return "postgres://" + c.Database.User + ":" + c.Database.Password + "@" +
		c.Database.Host + ":" + c.Database.Port + "/" + c.Database.Name +
		"?sslmode=" + c.Database.SSLMode
}
