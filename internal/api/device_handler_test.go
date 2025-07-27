package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"iot-platform-go/internal/device"
	"iot-platform-go/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func createTestDevice() *models.Device {
	return &models.Device{
		ID:        uuid.New().String(),
		Name:      "Test Device",
		Type:      "temperature",
		Location:  "Test Room",
		Status:    "online",
		LastSeen:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata:  `{"manufacturer":"Test Corp","model":"TEMP-001"}`,
	}
}

func TestCreateDevice(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(*device.MockRepository)
		expectedStatus int
		expectedError  string
	}{
		{
			name:        "successful device creation",
			requestBody: `{"name":"Test Device","type":"temperature","location":"Test Room","metadata":"{\"manufacturer\":\"Test Corp\"}"}`,
			mockSetup: func(mock *device.MockRepository) {
				mock.SetCreateFunc(func(req *models.CreateDeviceRequest) (*models.Device, error) {
					return &models.Device{
						ID:       "test-id",
						Name:     req.Name,
						Type:     req.Type,
						Location: req.Location,
						Metadata: req.Metadata,
					}, nil
				})
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid JSON",
			requestBody:    `{"name":"Test Device","type":"temperature"`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name:        "repository error",
			requestBody: `{"name":"Test Device","type":"temperature","location":"Test Room"}`,
			mockSetup: func(mock *device.MockRepository) {
				mock.SetCreateFunc(func(req *models.CreateDeviceRequest) (*models.Device, error) {
					return nil, assert.AnError
				})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to create device",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := device.NewMockRepository()
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			handler := NewDeviceHandler(mockRepo)
			router := setupTestRouter()
			router.POST("/devices", handler.CreateDevice)

			// Create request
			req := httptest.NewRequest("POST", "/devices", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			} else {
				var device models.Device
				err := json.Unmarshal(w.Body.Bytes(), &device)
				assert.NoError(t, err)
				assert.NotEmpty(t, device.ID)
			}
		})
	}
}

func TestGetDevice(t *testing.T) {
	tests := []struct {
		name           string
		deviceID       string
		mockSetup      func(*device.MockRepository)
		expectedStatus int
		expectedError  string
	}{
		{
			name:     "successful device retrieval",
			deviceID: "test-id",
			mockSetup: func(mock *device.MockRepository) {
				mock.SetGetByIDFunc(func(id string) (*models.Device, error) {
					return createTestDevice(), nil
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing device ID",
			deviceID:       "",
			expectedStatus: http.StatusNotFound, // Ginのルーティングにより404になる
			expectedError:  "404 page not found",
		},
		{
			name:     "device not found",
			deviceID: "non-existent-id",
			mockSetup: func(mock *device.MockRepository) {
				mock.SetGetByIDFunc(func(id string) (*models.Device, error) {
					return nil, assert.AnError
				})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to get device",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := device.NewMockRepository()
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			handler := NewDeviceHandler(mockRepo)
			router := setupTestRouter()
			router.GET("/devices/:id", handler.GetDevice)

			// Create request
			url := "/devices/" + tt.deviceID
			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				if tt.expectedStatus == http.StatusNotFound {
					// 404の場合はHTMLが返されるため、文字列で確認
					assert.Contains(t, w.Body.String(), tt.expectedError)
				} else {
					var response map[string]interface{}
					err := json.Unmarshal(w.Body.Bytes(), &response)
					assert.NoError(t, err)
					assert.Contains(t, response["error"], tt.expectedError)
				}
			} else {
				var device models.Device
				err := json.Unmarshal(w.Body.Bytes(), &device)
				assert.NoError(t, err)
				assert.NotEmpty(t, device.ID)
			}
		})
	}
}

func TestGetAllDevices(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*device.MockRepository)
		expectedStatus int
		expectedCount  int
		expectedError  string
	}{
		{
			name: "successful devices retrieval",
			mockSetup: func(mock *device.MockRepository) {
				mock.SetGetAllFunc(func() ([]*models.Device, error) {
					return []*models.Device{
						createTestDevice(),
						createTestDevice(),
					}, nil
				})
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name: "empty devices list",
			mockSetup: func(mock *device.MockRepository) {
				mock.SetGetAllFunc(func() ([]*models.Device, error) {
					return []*models.Device{}, nil
				})
			},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name: "repository error",
			mockSetup: func(mock *device.MockRepository) {
				mock.SetGetAllFunc(func() ([]*models.Device, error) {
					return nil, assert.AnError
				})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to get devices",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := device.NewMockRepository()
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			handler := NewDeviceHandler(mockRepo)
			router := setupTestRouter()
			router.GET("/devices", handler.GetAllDevices)

			// Create request
			req := httptest.NewRequest("GET", "/devices", nil)
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			} else {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, float64(tt.expectedCount), response["count"])
			}
		})
	}
}

func TestUpdateDevice(t *testing.T) {
	tests := []struct {
		name           string
		deviceID       string
		requestBody    string
		mockSetup      func(*device.MockRepository)
		expectedStatus int
		expectedError  string
	}{
		{
			name:        "successful device update",
			deviceID:    "test-id",
			requestBody: `{"name":"Updated Device","location":"Updated Room"}`,
			mockSetup: func(mock *device.MockRepository) {
				mock.SetUpdateFunc(func(id string, req *models.UpdateDeviceRequest) (*models.Device, error) {
					return &models.Device{
						ID:       id,
						Name:     req.Name,
						Location: req.Location,
					}, nil
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing device ID",
			deviceID:       "",
			requestBody:    `{"name":"Updated Device"}`,
			expectedStatus: http.StatusNotFound, // Ginのルーティングにより404になる
			expectedError:  "404 page not found",
		},
		{
			name:           "invalid JSON",
			deviceID:       "test-id",
			requestBody:    `{"name":"Updated Device"`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name:        "device not found",
			deviceID:    "non-existent-id",
			requestBody: `{"name":"Updated Device"}`,
			mockSetup: func(mock *device.MockRepository) {
				mock.SetUpdateFunc(func(id string, req *models.UpdateDeviceRequest) (*models.Device, error) {
					return nil, assert.AnError
				})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to update device",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := device.NewMockRepository()
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			handler := NewDeviceHandler(mockRepo)
			router := setupTestRouter()
			router.PUT("/devices/:id", handler.UpdateDevice)

			// Create request
			url := "/devices/" + tt.deviceID
			req := httptest.NewRequest("PUT", url, strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				if tt.expectedStatus == http.StatusNotFound {
					// 404の場合はHTMLが返されるため、文字列で確認
					assert.Contains(t, w.Body.String(), tt.expectedError)
				} else {
					var response map[string]interface{}
					err := json.Unmarshal(w.Body.Bytes(), &response)
					assert.NoError(t, err)
					assert.Contains(t, response["error"], tt.expectedError)
				}
			} else {
				var device models.Device
				err := json.Unmarshal(w.Body.Bytes(), &device)
				assert.NoError(t, err)
				assert.NotEmpty(t, device.ID)
			}
		})
	}
}

func TestDeleteDevice(t *testing.T) {
	tests := []struct {
		name           string
		deviceID       string
		mockSetup      func(*device.MockRepository)
		expectedStatus int
		expectedError  string
	}{
		{
			name:     "successful device deletion",
			deviceID: "test-id",
			mockSetup: func(mock *device.MockRepository) {
				mock.SetDeleteFunc(func(id string) error {
					return nil
				})
			},
			expectedStatus: http.StatusOK, // 実装では200を返している
		},
		{
			name:           "missing device ID",
			deviceID:       "",
			expectedStatus: http.StatusNotFound, // Ginのルーティングにより404になる
			expectedError:  "404 page not found",
		},
		{
			name:     "device not found",
			deviceID: "non-existent-id",
			mockSetup: func(mock *device.MockRepository) {
				mock.SetDeleteFunc(func(id string) error {
					return assert.AnError
				})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to delete device",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := device.NewMockRepository()
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			handler := NewDeviceHandler(mockRepo)
			router := setupTestRouter()
			router.DELETE("/devices/:id", handler.DeleteDevice)

			// Create request
			url := "/devices/" + tt.deviceID
			req := httptest.NewRequest("DELETE", url, nil)
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				if tt.expectedStatus == http.StatusNotFound {
					// 404の場合はHTMLが返されるため、文字列で確認
					assert.Contains(t, w.Body.String(), tt.expectedError)
				} else {
					var response map[string]interface{}
					err := json.Unmarshal(w.Body.Bytes(), &response)
					assert.NoError(t, err)
					assert.Contains(t, response["error"], tt.expectedError)
				}
			}
		})
	}
}

func TestGetDeviceStatus(t *testing.T) {
	tests := []struct {
		name           string
		deviceID       string
		mockSetup      func(*device.MockRepository)
		expectedStatus int
		expectedError  string
	}{
		{
			name:     "successful status retrieval",
			deviceID: "test-id",
			mockSetup: func(mock *device.MockRepository) {
				mock.SetGetByIDFunc(func(id string) (*models.Device, error) {
					return &models.Device{
						ID:     id,
						Status: "online",
					}, nil
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing device ID",
			deviceID:       "",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Device ID is required",
		},
		{
			name:     "device not found",
			deviceID: "non-existent-id",
			mockSetup: func(mock *device.MockRepository) {
				mock.SetGetByIDFunc(func(id string) (*models.Device, error) {
					return nil, assert.AnError
				})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to get device status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := device.NewMockRepository()
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			handler := NewDeviceHandler(mockRepo)
			router := setupTestRouter()
			router.GET("/devices/:id/status", handler.GetDeviceStatus)

			// Create request
			url := "/devices/" + tt.deviceID + "/status"
			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			} else {
				var status models.DeviceStatus
				err := json.Unmarshal(w.Body.Bytes(), &status)
				assert.NoError(t, err)
				assert.NotEmpty(t, status.Status)
			}
		})
	}
}
