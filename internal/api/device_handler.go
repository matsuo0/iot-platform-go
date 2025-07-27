package api

import (
	"net/http"

	"iot-platform-go/internal/device"
	"iot-platform-go/pkg/models"

	"github.com/gin-gonic/gin"
)

const (
	// Error messages
	ErrDeviceNotFound = "device not found"
)

// DeviceHandler handles device-related HTTP requests.
type DeviceHandler struct {
	repo device.RepositoryInterface
}

// NewDeviceHandler creates a new device handler
func NewDeviceHandler(repo device.RepositoryInterface) *DeviceHandler {
	return &DeviceHandler{repo: repo}
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get device status: " + err.Error()})
		return
	}

	status := models.DeviceStatus{
		DeviceID: device.ID,
		Status:   device.Status,
		LastSeen: device.LastSeen,
	}

	c.JSON(http.StatusOK, status)
}
