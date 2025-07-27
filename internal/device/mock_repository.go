package device

import (
	"fmt"
	"iot-platform-go/pkg/models"
	"time"
)

// MockRepository is a mock implementation of the device repository for testing
type MockRepository struct {
	devices          map[string]*models.Device
	createFunc       func(req *models.CreateDeviceRequest) (*models.Device, error)
	getByIDFunc      func(id string) (*models.Device, error)
	getAllFunc       func() ([]*models.Device, error)
	updateFunc       func(id string, req *models.UpdateDeviceRequest) (*models.Device, error)
	deleteFunc       func(id string) error
	updateStatusFunc func(id string, status string) error
}

// NewMockRepository creates a new mock repository
func NewMockRepository() *MockRepository {
	return &MockRepository{
		devices: make(map[string]*models.Device),
	}
}

// Create creates a new device
func (m *MockRepository) Create(req *models.CreateDeviceRequest) (*models.Device, error) {
	if m.createFunc != nil {
		return m.createFunc(req)
	}

	device := &models.Device{
		ID:        "mock-device-id",
		Name:      req.Name,
		Type:      req.Type,
		Location:  req.Location,
		Status:    "offline",
		LastSeen:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata:  req.Metadata,
	}

	m.devices[device.ID] = device
	return device, nil
}

// GetByID retrieves a device by ID
func (m *MockRepository) GetByID(id string) (*models.Device, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(id)
	}

	device, exists := m.devices[id]
	if !exists {
		return nil, fmt.Errorf("device not found")
	}

	return device, nil
}

// GetAll retrieves all devices
func (m *MockRepository) GetAll() ([]*models.Device, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc()
	}

	var devices []*models.Device
	for _, device := range m.devices {
		devices = append(devices, device)
	}

	return devices, nil
}

// Update updates a device
func (m *MockRepository) Update(id string, req *models.UpdateDeviceRequest) (*models.Device, error) {
	if m.updateFunc != nil {
		return m.updateFunc(id, req)
	}

	device, exists := m.devices[id]
	if !exists {
		return nil, fmt.Errorf("device not found")
	}

	if req.Name != "" {
		device.Name = req.Name
	}
	if req.Type != "" {
		device.Type = req.Type
	}
	if req.Location != "" {
		device.Location = req.Location
	}
	if req.Metadata != "" {
		device.Metadata = req.Metadata
	}

	device.UpdatedAt = time.Now()
	m.devices[id] = device

	return device, nil
}

// Delete deletes a device
func (m *MockRepository) Delete(id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}

	if _, exists := m.devices[id]; !exists {
		return fmt.Errorf("device not found")
	}

	delete(m.devices, id)
	return nil
}

// UpdateStatus updates device status
func (m *MockRepository) UpdateStatus(id string, status string) error {
	if m.updateStatusFunc != nil {
		return m.updateStatusFunc(id, status)
	}

	device, exists := m.devices[id]
	if !exists {
		return fmt.Errorf("device not found")
	}

	device.Status = status
	device.LastSeen = time.Now()
	device.UpdatedAt = time.Now()
	m.devices[id] = device

	return nil
}

// SetCreateFunc sets a custom create function for testing
func (m *MockRepository) SetCreateFunc(fn func(req *models.CreateDeviceRequest) (*models.Device, error)) {
	m.createFunc = fn
}

// SetGetByIDFunc sets a custom get by ID function for testing
func (m *MockRepository) SetGetByIDFunc(fn func(id string) (*models.Device, error)) {
	m.getByIDFunc = fn
}

// SetGetAllFunc sets a custom get all function for testing
func (m *MockRepository) SetGetAllFunc(fn func() ([]*models.Device, error)) {
	m.getAllFunc = fn
}

// SetUpdateFunc sets a custom update function for testing
func (m *MockRepository) SetUpdateFunc(fn func(id string, req *models.UpdateDeviceRequest) (*models.Device, error)) {
	m.updateFunc = fn
}

// SetDeleteFunc sets a custom delete function for testing
func (m *MockRepository) SetDeleteFunc(fn func(id string) error) {
	m.deleteFunc = fn
}

// SetUpdateStatusFunc sets a custom update status function for testing
func (m *MockRepository) SetUpdateStatusFunc(fn func(id string, status string) error) {
	m.updateStatusFunc = fn
}

// AddDevice adds a device to the mock repository for testing
func (m *MockRepository) AddDevice(device *models.Device) {
	m.devices[device.ID] = device
}

// Clear clears all devices from the mock repository
func (m *MockRepository) Clear() {
	m.devices = make(map[string]*models.Device)
}
