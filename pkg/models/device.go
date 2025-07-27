package models

import (
	"time"
)

// Device represents an IoT device
type Device struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Type        string    `json:"type" db:"type"`
	Location    string    `json:"location" db:"location"`
	Status      string    `json:"status" db:"status"` // online, offline, error
	LastSeen    time.Time `json:"last_seen" db:"last_seen"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Metadata    string    `json:"metadata" db:"metadata"` // JSON string for additional data
}

// DeviceData represents sensor data from a device
type DeviceData struct {
	ID        string                 `json:"id" db:"id"`
	DeviceID  string                 `json:"device_id" db:"device_id"`
	Timestamp time.Time              `json:"timestamp" db:"timestamp"`
	Data      map[string]interface{} `json:"data" db:"data"`
}

// CreateDeviceRequest represents the request to create a new device
type CreateDeviceRequest struct {
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Location string `json:"location" binding:"required"`
	Metadata string `json:"metadata"`
}

// UpdateDeviceRequest represents the request to update a device
type UpdateDeviceRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Location string `json:"location"`
	Status   string `json:"status"`
	Metadata string `json:"metadata"`
}

// DeviceStatus represents the current status of a device
type DeviceStatus struct {
	DeviceID  string    `json:"device_id"`
	Status    string    `json:"status"`
	LastSeen  time.Time `json:"last_seen"`
	Data      map[string]interface{} `json:"data,omitempty"`
} 