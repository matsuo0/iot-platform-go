package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"iot-platform-go/internal/api"
	"iot-platform-go/internal/config"
	"iot-platform-go/internal/database"
	"iot-platform-go/internal/device"
	"iot-platform-go/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServer represents a test server instance
type TestServer struct {
	Router  *gin.Engine
	DB      *database.Database
	Repo    *device.Repository
	Handler *api.DeviceHandler
}

// NewTestServer creates a new test server
func NewTestServer(t *testing.T) *TestServer {
	// テスト用の設定
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			Name:     "iot_platform_test",
			User:     "postgres",
			Password: "password",
			SSLMode:  "disable",
		},
	}

	// データベース接続
	db, err := database.New(cfg)
	require.NoError(t, err)

	// リポジトリとハンドラーの作成
	repo := device.NewRepository(db)
	handler := api.NewDeviceHandler(repo)

	// ルーターの設定
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// APIルートの設定
	apiGroup := router.Group("/api")
	{
		devices := apiGroup.Group("/devices")
		{
			devices.POST("", handler.CreateDevice)
			devices.GET("", handler.GetAllDevices)
			devices.GET("/:id", handler.GetDevice)
			devices.PUT("/:id", handler.UpdateDevice)
			devices.DELETE("/:id", handler.DeleteDevice)
			devices.GET("/:id/status", handler.GetDeviceStatus)
		}
	}

	return &TestServer{
		Router:  router,
		DB:      db,
		Repo:    repo,
		Handler: handler,
	}
}

// Cleanup cleans up test data.
func (ts *TestServer) Cleanup() {
	// テストデータのクリーンアップ
	_, err := ts.DB.Exec("DELETE FROM device_data")
	if err != nil {
		// ログ出力のみ（テストクリーンアップなのでエラーは無視）
		_ = err
	}
	_, err = ts.DB.Exec("DELETE FROM devices")
	if err != nil {
		// ログ出力のみ（テストクリーンアップなのでエラーは無視）
		_ = err
	}
}

// Close closes the test server
func (ts *TestServer) Close() {
	ts.DB.Close()
}

// TestDeviceLifecycle tests the complete device lifecycle
func TestDeviceLifecycle(t *testing.T) {
	t.Skip("Skipping integration test as it requires database setup")
	server := NewTestServer(t)
	defer server.Close()
	defer server.Cleanup()

	t.Run("complete device lifecycle", func(t *testing.T) {
		// Step 1: Create a device
		createReq := &models.CreateDeviceRequest{
			Name:     "Integration Test Device",
			Type:     "temperature",
			Location: "Test Room",
			Metadata: `{"manufacturer":"Test Corp","model":"INT-001"}`,
		}

		createBody, _ := json.Marshal(createReq)
		req := httptest.NewRequest("POST", "/api/devices", bytes.NewBuffer(createBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		server.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var createdDevice models.Device
		err := json.Unmarshal(w.Body.Bytes(), &createdDevice)
		assert.NoError(t, err)
		assert.NotEmpty(t, createdDevice.ID)
		assert.Equal(t, createReq.Name, createdDevice.Name)
		assert.Equal(t, createReq.Type, createdDevice.Type)
		assert.Equal(t, createReq.Location, createdDevice.Location)
		assert.Equal(t, "offline", createdDevice.Status)

		deviceID := createdDevice.ID

		// Step 2: Get the created device
		req = httptest.NewRequest("GET", fmt.Sprintf("/api/devices/%s", deviceID), nil)
		w = httptest.NewRecorder()

		server.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var retrievedDevice models.Device
		err = json.Unmarshal(w.Body.Bytes(), &retrievedDevice)
		assert.NoError(t, err)
		assert.Equal(t, deviceID, retrievedDevice.ID)
		assert.Equal(t, createReq.Name, retrievedDevice.Name)

		// Step 3: Update the device
		updateReq := &models.UpdateDeviceRequest{
			Name:     "Updated Integration Test Device",
			Location: "Updated Test Room",
			Metadata: `{"manufacturer":"Updated Test Corp","model":"INT-002"}`,
		}

		updateBody, _ := json.Marshal(updateReq)
		req = httptest.NewRequest("PUT", fmt.Sprintf("/api/devices/%s", deviceID), bytes.NewBuffer(updateBody))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		server.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var updatedDevice models.Device
		err = json.Unmarshal(w.Body.Bytes(), &updatedDevice)
		assert.NoError(t, err)
		assert.Equal(t, deviceID, updatedDevice.ID)
		assert.Equal(t, updateReq.Name, updatedDevice.Name)
		assert.Equal(t, updateReq.Location, updatedDevice.Location)

		// Step 4: Get device status
		req = httptest.NewRequest("GET", fmt.Sprintf("/api/devices/%s/status", deviceID), nil)
		w = httptest.NewRecorder()

		server.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var status models.DeviceStatus
		err = json.Unmarshal(w.Body.Bytes(), &status)
		assert.NoError(t, err)
		assert.Equal(t, deviceID, status.DeviceID)
		assert.NotEmpty(t, status.Status)

		// Step 5: Get all devices
		req = httptest.NewRequest("GET", "/api/devices", nil)
		w = httptest.NewRecorder()

		server.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var devicesResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &devicesResponse)
		assert.NoError(t, err)
		assert.Contains(t, devicesResponse, "devices")
		assert.Contains(t, devicesResponse, "count")
		assert.Equal(t, float64(1), devicesResponse["count"])

		// Step 6: Delete the device
		req = httptest.NewRequest("DELETE", fmt.Sprintf("/api/devices/%s", deviceID), nil)
		w = httptest.NewRecorder()

		server.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Step 7: Verify device is deleted
		req = httptest.NewRequest("GET", fmt.Sprintf("/api/devices/%s", deviceID), nil)
		w = httptest.NewRecorder()

		server.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code) // または404、実装による
	})
}

// TestMultipleDevices tests operations with multiple devices
func TestMultipleDevices(t *testing.T) {
	t.Skip("Skipping integration test as it requires database setup")
	server := NewTestServer(t)
	defer server.Close()
	defer server.Cleanup()

	t.Run("multiple devices operations", func(t *testing.T) {
		// Create multiple devices
		deviceNames := []string{"Device 1", "Device 2", "Device 3"}
		deviceIDs := make([]string, len(deviceNames))

		for i, name := range deviceNames {
			createReq := &models.CreateDeviceRequest{
				Name:     name,
				Type:     "temperature",
				Location: fmt.Sprintf("Room %d", i+1),
			}

			createBody, _ := json.Marshal(createReq)
			req := httptest.NewRequest("POST", "/api/devices", bytes.NewBuffer(createBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.Router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusCreated, w.Code)

			var device models.Device
			err := json.Unmarshal(w.Body.Bytes(), &device)
			assert.NoError(t, err)
			deviceIDs[i] = device.ID
		}

		// Get all devices
		req := httptest.NewRequest("GET", "/api/devices", nil)
		w := httptest.NewRecorder()

		server.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var devicesResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &devicesResponse)
		assert.NoError(t, err)
		assert.Equal(t, float64(3), devicesResponse["count"])

		// Verify each device can be retrieved individually
		for _, deviceID := range deviceIDs {
			req = httptest.NewRequest("GET", fmt.Sprintf("/api/devices/%s", deviceID), nil)
			w = httptest.NewRecorder()

			server.Router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var device models.Device
			err := json.Unmarshal(w.Body.Bytes(), &device)
			assert.NoError(t, err)
			assert.Equal(t, deviceID, device.ID)
		}
	})
}

// TestErrorHandling tests error scenarios
func TestErrorHandling(t *testing.T) {
	t.Skip("Skipping integration test as it requires database setup")
	server := NewTestServer(t)
	defer server.Close()
	defer server.Cleanup()

	t.Run("invalid JSON in create request", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/devices", bytes.NewBufferString(`{"name":"Test Device"`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		server.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Invalid request body")
	})

	t.Run("get non-existent device", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/devices/non-existent-id", nil)
		w := httptest.NewRecorder()

		server.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code) // または404、実装による

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Device not found")
	})

	t.Run("update non-existent device", func(t *testing.T) {
		updateReq := &models.UpdateDeviceRequest{
			Name: "Updated Device",
		}

		updateBody, _ := json.Marshal(updateReq)
		req := httptest.NewRequest("PUT", "/api/devices/non-existent-id", bytes.NewBuffer(updateBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		server.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code) // または404、実装による

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Device not found")
	})

	t.Run("delete non-existent device", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/devices/non-existent-id", nil)
		w := httptest.NewRecorder()

		server.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code) // または404、実装による

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Device not found")
	})
}

// TestDataValidation tests data validation scenarios
func TestDataValidation(t *testing.T) {
	t.Skip("Skipping integration test as it requires database setup")
	server := NewTestServer(t)
	defer server.Close()
	defer server.Cleanup()

	t.Run("device with special characters", func(t *testing.T) {
		createReq := &models.CreateDeviceRequest{
			Name:     "Device with 特殊文字 & Symbols!@#$%",
			Type:     "temperature",
			Location: "Room with 日本語",
			Metadata: `{"special":"value with 特殊文字","number":123.45,"boolean":true}`,
		}

		createBody, _ := json.Marshal(createReq)
		req := httptest.NewRequest("POST", "/api/devices", bytes.NewBuffer(createBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		server.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var device models.Device
		err := json.Unmarshal(w.Body.Bytes(), &device)
		assert.NoError(t, err)
		assert.Equal(t, createReq.Name, device.Name)
		assert.Equal(t, createReq.Location, device.Location)
		assert.Equal(t, createReq.Metadata, device.Metadata)
	})

	t.Run("device with empty name", func(t *testing.T) {
		createReq := &models.CreateDeviceRequest{
			Name:     "",
			Type:     "temperature",
			Location: "Test Room",
		}

		createBody, _ := json.Marshal(createReq)
		req := httptest.NewRequest("POST", "/api/devices", bytes.NewBuffer(createBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		server.Router.ServeHTTP(w, req)

		// 実装によってはバリデーションエラーになる可能性がある
		// 現在の実装では成功する
		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

// TestConcurrentOperations tests concurrent operations
func TestConcurrentOperations(t *testing.T) {
	t.Skip("Skipping integration test as it requires database setup")
	server := NewTestServer(t)
	defer server.Close()
	defer server.Cleanup()

	t.Run("concurrent device creation", func(t *testing.T) {
		const numDevices = 10
		done := make(chan bool, numDevices)

		for i := 0; i < numDevices; i++ {
			go func(index int) {
				createReq := &models.CreateDeviceRequest{
					Name:     fmt.Sprintf("Concurrent Device %d", index),
					Type:     "temperature",
					Location: fmt.Sprintf("Room %d", index),
				}

				createBody, _ := json.Marshal(createReq)
				req := httptest.NewRequest("POST", "/api/devices", bytes.NewBuffer(createBody))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				server.Router.ServeHTTP(w, req)

				assert.Equal(t, http.StatusCreated, w.Code)
				done <- true
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < numDevices; i++ {
			<-done
		}

		// Verify all devices were created
		req := httptest.NewRequest("GET", "/api/devices", nil)
		w := httptest.NewRecorder()

		server.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var devicesResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &devicesResponse)
		assert.NoError(t, err)
		assert.Equal(t, float64(numDevices), devicesResponse["count"])
	})
}

// TestPerformance tests basic performance characteristics
func TestPerformance(t *testing.T) {
	t.Skip("Skipping integration test as it requires database setup")
	server := NewTestServer(t)
	defer server.Close()
	defer server.Cleanup()

	t.Run("bulk device creation performance", func(t *testing.T) {
		const numDevices = 100
		start := time.Now()

		for i := 0; i < numDevices; i++ {
			createReq := &models.CreateDeviceRequest{
				Name:     fmt.Sprintf("Performance Device %d", i),
				Type:     "temperature",
				Location: fmt.Sprintf("Room %d", i),
			}

			createBody, _ := json.Marshal(createReq)
			req := httptest.NewRequest("POST", "/api/devices", bytes.NewBuffer(createBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.Router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusCreated, w.Code)
		}

		duration := time.Since(start)
		t.Logf("Created %d devices in %v", numDevices, duration)

		// Verify all devices were created
		req := httptest.NewRequest("GET", "/api/devices", nil)
		w := httptest.NewRecorder()

		server.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var devicesResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &devicesResponse)
		assert.NoError(t, err)
		assert.Equal(t, float64(numDevices), devicesResponse["count"])
	})
}
