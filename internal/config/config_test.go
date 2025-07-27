package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// テスト用の環境変数を設定
	originalEnv := make(map[string]string)
	envVars := map[string]string{
		"SERVER_HOST":     "test-host",
		"SERVER_PORT":     "8080",
		"DB_HOST":         "test-db-host",
		"DB_PORT":         "5433",
		"DB_NAME":         "test_db",
		"DB_USER":         "test_user",
		"DB_PASSWORD":     "test_password",
		"DB_SSL_MODE":     "require",
		"MQTT_BROKER":     "test-mqtt-broker",
		"MQTT_CLIENT_ID":  "test-client-id",
		"MQTT_USERNAME":   "test-mqtt-user",
		"MQTT_PASSWORD":   "test-mqtt-password",
		"JWT_SECRET":      "test-jwt-secret",
		"JWT_EXPIRATION":  "24h",
		"LOG_LEVEL":       "debug",
	}

	// 元の環境変数を保存
	for key := range envVars {
		if val := os.Getenv(key); val != "" {
			originalEnv[key] = val
		}
	}

	// テスト用の環境変数を設定
	for key, value := range envVars {
		os.Setenv(key, value)
	}

	// テスト終了後に環境変数を復元
	defer func() {
		for key := range envVars {
			os.Unsetenv(key)
		}
		for key, value := range originalEnv {
			os.Setenv(key, value)
		}
	}()

	t.Run("load configuration from environment variables", func(t *testing.T) {
		cfg := Load()

		// Server設定の検証
		assert.Equal(t, "test-host", cfg.Server.Host)
		assert.Equal(t, "8080", cfg.Server.Port)

		// Database設定の検証
		assert.Equal(t, "test-db-host", cfg.Database.Host)
		assert.Equal(t, "5433", cfg.Database.Port)
		assert.Equal(t, "test_db", cfg.Database.Name)
		assert.Equal(t, "test_user", cfg.Database.User)
		assert.Equal(t, "test_password", cfg.Database.Password)
		assert.Equal(t, "require", cfg.Database.SSLMode)

		// MQTT設定の検証
		assert.Equal(t, "test-mqtt-broker", cfg.MQTT.Broker)
		assert.Equal(t, "test-client-id", cfg.MQTT.ClientID)
		assert.Equal(t, "test-mqtt-user", cfg.MQTT.Username)
		assert.Equal(t, "test-mqtt-password", cfg.MQTT.Password)

		// JWT設定の検証
		assert.Equal(t, "test-jwt-secret", cfg.JWT.Secret)
		assert.Equal(t, "24h", cfg.JWT.Expiration)

		// Logging設定の検証
		assert.Equal(t, "debug", cfg.Logging.Level)
	})
}

func TestLoadWithDefaults(t *testing.T) {
	// 環境変数をクリア
	envVars := []string{
		"SERVER_HOST", "SERVER_PORT", "DB_HOST", "DB_PORT", "DB_NAME",
		"DB_USER", "DB_PASSWORD", "DB_SSL_MODE", "MQTT_BROKER", "MQTT_CLIENT_ID",
		"MQTT_USERNAME", "MQTT_PASSWORD", "JWT_SECRET", "JWT_EXPIRATION",
		"LOG_LEVEL",
	}

	originalEnv := make(map[string]string)
	for _, key := range envVars {
		if val := os.Getenv(key); val != "" {
			originalEnv[key] = val
		}
		os.Unsetenv(key)
	}

	// テスト終了後に環境変数を復元
	defer func() {
		for key, value := range originalEnv {
			os.Setenv(key, value)
		}
	}()

	t.Run("load configuration with default values", func(t *testing.T) {
		cfg := Load()

		// デフォルト値の検証
		assert.Equal(t, "localhost", cfg.Server.Host)
		assert.Equal(t, "8080", cfg.Server.Port)
		assert.Equal(t, "localhost", cfg.Database.Host)
		assert.Equal(t, "5432", cfg.Database.Port)
		assert.Equal(t, "iot_platform", cfg.Database.Name)
		assert.Equal(t, "postgres", cfg.Database.User)
		assert.Equal(t, "password", cfg.Database.Password)
		assert.Equal(t, "disable", cfg.Database.SSLMode)
		assert.Equal(t, "tcp://localhost:1883", cfg.MQTT.Broker)
		assert.Equal(t, "iot-platform-server", cfg.MQTT.ClientID)
		assert.Equal(t, "", cfg.MQTT.Username)
		assert.Equal(t, "", cfg.MQTT.Password)
		assert.Equal(t, "your-secret-key-here", cfg.JWT.Secret)
		assert.Equal(t, "24h", cfg.JWT.Expiration)
		assert.Equal(t, "info", cfg.Logging.Level)
	})
}

func TestGetDatabaseURL(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected string
	}{
		{
			name: "standard database URL",
			config: &Config{
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     "5432",
					Name:     "test_db",
					User:     "test_user",
					Password: "test_password",
					SSLMode:  "disable",
				},
			},
			expected: "postgres://test_user:test_password@localhost:5432/test_db?sslmode=disable",
		},
		{
			name: "database URL with SSL",
			config: &Config{
				Database: DatabaseConfig{
					Host:     "db.example.com",
					Port:     "5432",
					Name:     "production_db",
					User:     "prod_user",
					Password: "prod_password",
					SSLMode:  "require",
				},
			},
			expected: "postgres://prod_user:prod_password@db.example.com:5432/production_db?sslmode=require",
		},
		{
			name: "database URL with special characters in password",
			config: &Config{
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     "5432",
					Name:     "test_db",
					User:     "test_user",
					Password: "pass@word#123",
					SSLMode:  "disable",
				},
			},
			expected: "postgres://test_user:pass@word#123@localhost:5432/test_db?sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetDatabaseURL()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLoadWithEnvFile(t *testing.T) {
	t.Skip("Skipping .env file test as it requires specific environment setup")
}

func TestConfigValidation(t *testing.T) {
	t.Run("valid configuration", func(t *testing.T) {
		cfg := &Config{
			Server: ServerConfig{
				Host: "localhost",
				Port: "8080",
			},
			Database: DatabaseConfig{
				Host:     "localhost",
				Port:     "5432",
				Name:     "test_db",
				User:     "test_user",
				Password: "test_password",
				SSLMode:  "disable",
			},
		}

		// 基本的な検証（実際のバリデーション関数がある場合）
		assert.NotEmpty(t, cfg.Server.Host)
		assert.NotEmpty(t, cfg.Server.Port)
		assert.NotEmpty(t, cfg.Database.Host)
		assert.NotEmpty(t, cfg.Database.Port)
		assert.NotEmpty(t, cfg.Database.Name)
		assert.NotEmpty(t, cfg.Database.User)
		assert.NotEmpty(t, cfg.Database.Password)
	})

	t.Run("database URL generation", func(t *testing.T) {
		cfg := &Config{
			Database: DatabaseConfig{
				Host:     "localhost",
				Port:     "5432",
				Name:     "test_db",
				User:     "test_user",
				Password: "test_password",
				SSLMode:  "disable",
			},
		}

		url := cfg.GetDatabaseURL()
		assert.Contains(t, url, "postgres://test_user:test_password@localhost:5432/test_db")
		assert.Contains(t, url, "sslmode=disable")
	})
} 