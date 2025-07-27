package models

import "time"

// Device represents an IoT device.
type Device struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Location  string    `json:"location"`
	Status    string    `json:"status"`
	Metadata  string    `json:"metadata,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LastSeen  time.Time `json:"last_seen,omitempty"`
}

// DeviceData represents sensor data from a device.
type DeviceData struct {
	ID        string    `json:"id"`
	DeviceID  string    `json:"device_id"`
	Timestamp time.Time `json:"timestamp"`
	DataType  string    `json:"data_type"`
	Value     float64   `json:"value"`
	Unit      string    `json:"unit,omitempty"`
	Metadata  string    `json:"metadata,omitempty"`
}

// CreateDeviceRequest represents the request to create a new device.
type CreateDeviceRequest struct {
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Location string `json:"location"`
	Metadata string `json:"metadata,omitempty"`
}

// UpdateDeviceRequest represents the request to update a device.
type UpdateDeviceRequest struct {
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	Location string `json:"location,omitempty"`
	Status   string `json:"status,omitempty"`
	Metadata string `json:"metadata,omitempty"`
}

// DeviceStatus represents the current status of a device.
type DeviceStatus struct {
	DeviceID string    `json:"device_id"`
	Status   string    `json:"status"`
	LastSeen time.Time `json:"last_seen"`
}
