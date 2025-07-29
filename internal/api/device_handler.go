package api

import (
	"net/http"
	"strconv"

	"iot-platform-go/internal/device"
	"iot-platform-go/pkg/models"

	"github.com/gin-gonic/gin"
)

const (
	// Error messages
	ErrDeviceNotFound = "device not found"

	// API limits
	DefaultLimit = 100
	MaxLimit     = 1000
)

// DeviceHandler handles HTTP requests for devices
type DeviceHandler struct {
	repo     device.RepositoryInterface
	dataRepo device.DataRepositoryInterface
}

// NewDeviceHandler creates a new device handler
func NewDeviceHandler(repo device.RepositoryInterface, dataRepo device.DataRepositoryInterface) *DeviceHandler {
	return &DeviceHandler{
		repo:     repo,
		dataRepo: dataRepo,
	}
}

// CreateDevice handles POST /api/devices
func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	var req models.CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	device, err := h.repo.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create device: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, device)
}

// GetDevice handles GET /api/devices/:id.
func (h *DeviceHandler) GetDevice(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device ID is required"})
		return
	}

	device, err := h.repo.GetByID(id)
	if err != nil {
		if err.Error() == ErrDeviceNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": ErrDeviceNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get device"})
		return
	}

	c.JSON(http.StatusOK, device)
}

// GetAllDevices handles GET /api/devices
func (h *DeviceHandler) GetAllDevices(c *gin.Context) {
	devices, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get devices: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"devices": devices,
		"count":   len(devices),
	})
}

// UpdateDevice handles PUT /api/devices/:id.
func (h *DeviceHandler) UpdateDevice(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device ID is required"})
		return
	}

	var req models.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	device, err := h.repo.Update(id, &req)
	if err != nil {
		if err.Error() == ErrDeviceNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": ErrDeviceNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update device: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, device)
}

// DeleteDevice handles DELETE /api/devices/:id.
func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device ID is required"})
		return
	}

	err := h.repo.Delete(id)
	if err != nil {
		if err.Error() == ErrDeviceNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": ErrDeviceNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete device: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device deleted successfully"})
}

// GetDeviceStatus handles GET /api/devices/:id/status.
func (h *DeviceHandler) GetDeviceStatus(c *gin.Context) {
	id := c.Param("id")
	device, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"device_id": device.ID,
		"status":    device.Status,
		"last_seen": device.LastSeen,
	})
}

// GetDeviceData gets the data for a device
func (h *DeviceHandler) GetDeviceData(c *gin.Context) {
	deviceID := c.Param("id")

	// Get limit from query parameter (default: 100)
	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit // Maximum limit
	}

	// Get data type filter from query parameter
	dataType := c.Query("type")

	var data []*models.DeviceData
	var dataErr error

	if dataType != "" {
		data, dataErr = h.dataRepo.GetDeviceDataByType(deviceID, dataType, limit)
	} else {
		data, dataErr = h.dataRepo.GetDeviceData(deviceID, limit)
	}

	if dataErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get device data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"device_id": deviceID,
		"data":      data,
		"count":     len(data),
		"limit":     limit,
	})
}

// GetLatestDeviceData gets the latest data for a device
func (h *DeviceHandler) GetLatestDeviceData(c *gin.Context) {
	deviceID := c.Param("id")

	data, err := h.dataRepo.GetLatestData(deviceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No data found for device"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"device_id":   deviceID,
		"latest_data": data,
	})
}
