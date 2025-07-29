package api

import (
	"net/http"
	"strconv"
	"time"

	"iot-platform-go/internal/influxdb"

	"github.com/gin-gonic/gin"
)

const (
	// InfluxDB API limits
	InfluxDBDefaultLimit = 100
	InfluxDBMaxLimit     = 1000
)

// InfluxDBHandler handles InfluxDB-related API endpoints
type InfluxDBHandler struct {
	influxClient *influxdb.Client
}

// NewInfluxDBHandler creates a new InfluxDB handler
func NewInfluxDBHandler(influxClient *influxdb.Client) *InfluxDBHandler {
	return &InfluxDBHandler{
		influxClient: influxClient,
	}
}

// GetDeviceDataFromInfluxDB gets device data from InfluxDB
func (h *InfluxDBHandler) GetDeviceDataFromInfluxDB(c *gin.Context) {
	if h.influxClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "InfluxDB not available"})
		return
	}

	deviceID := c.Param("id")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device ID is required"})
		return
	}

	// Get query parameters
	dataType := c.Query("type")
	limitStr := c.DefaultQuery("limit", strconv.Itoa(InfluxDBDefaultLimit))
	startStr := c.Query("start")
	endStr := c.Query("end")

	// Parse limit
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = InfluxDBDefaultLimit
	}
	if limit > InfluxDBMaxLimit {
		limit = InfluxDBMaxLimit
	}

	// Parse time range
	end := time.Now()
	start := end.Add(-24 * time.Hour) // Default to last 24 hours

	if startStr != "" {
		if parsed, err := time.Parse(time.RFC3339, startStr); err == nil {
			start = parsed
		}
	}

	if endStr != "" {
		if parsed, err := time.Parse(time.RFC3339, endStr); err == nil {
			end = parsed
		}
	}

	// Query data from InfluxDB
	data, err := h.influxClient.QueryDeviceData(deviceID, dataType, start, end, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query data from InfluxDB"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"device_id": deviceID,
		"data":      data,
		"count":     len(data),
		"limit":     limit,
		"start":     start.Format(time.RFC3339),
		"end":       end.Format(time.RFC3339),
		"source":    "influxdb",
	})
}

// GetLatestDeviceDataFromInfluxDB gets the latest data point for a device from InfluxDB
func (h *InfluxDBHandler) GetLatestDeviceDataFromInfluxDB(c *gin.Context) {
	if h.influxClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "InfluxDB not available"})
		return
	}

	deviceID := c.Param("id")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device ID is required"})
		return
	}

	dataType := c.Query("type")

	// Query latest data from InfluxDB
	data, err := h.influxClient.GetLatestDeviceData(deviceID, dataType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No data found for device"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"device_id":   deviceID,
		"latest_data": data,
		"source":      "influxdb",
	})
}
