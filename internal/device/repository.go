package device

import (
	"database/sql"
	"fmt"
	"time"

	"iot-platform-go/internal/database"
	"iot-platform-go/pkg/models"

	"github.com/google/uuid"
)

// RepositoryInterface defines the interface for device repository operations
type RepositoryInterface interface {
	Create(req *models.CreateDeviceRequest) (*models.Device, error)
	GetByID(id string) (*models.Device, error)
	GetAll() ([]*models.Device, error)
	Update(id string, req *models.UpdateDeviceRequest) (*models.Device, error)
	Delete(id string) error
	UpdateStatus(id string, status string) error
}

// Repository handles database operations for devices
type Repository struct {
	db *database.Database
}

// NewRepository creates a new device repository
func NewRepository(db *database.Database) *Repository {
	return &Repository{db: db}
}

// Create creates a new device
func (r *Repository) Create(req *models.CreateDeviceRequest) (*models.Device, error) {
	device := &models.Device{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Type:      req.Type,
		Location:  req.Location,
		Status:    "offline",
		LastSeen:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata:  req.Metadata,
	}

	query := `
		INSERT INTO devices (id, name, type, location, status, last_seen, created_at, updated_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.Exec(query, device.ID, device.Name, device.Type, device.Location,
		device.Status, device.LastSeen, device.CreatedAt, device.UpdatedAt, device.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to create device: %w", err)
	}

	return device, nil
}

// GetByID retrieves a device by ID
func (r *Repository) GetByID(id string) (*models.Device, error) {
	device := &models.Device{}
	query := `
		SELECT id, name, type, location, status, last_seen, created_at, updated_at, metadata
		FROM devices WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&device.ID, &device.Name, &device.Type, &device.Location,
		&device.Status, &device.LastSeen, &device.CreatedAt, &device.UpdatedAt, &device.Metadata)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("device not found")
		}
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	return device, nil
}

// GetAll retrieves all devices
func (r *Repository) GetAll() ([]*models.Device, error) {
	query := `
		SELECT id, name, type, location, status, metadata, created_at, updated_at, last_seen
		FROM devices
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query devices: %w", err)
	}
	defer rows.Close()

	var devices []*models.Device
	for rows.Next() {
		device := &models.Device{}
		err := rows.Scan(
			&device.ID,
			&device.Name,
			&device.Type,
			&device.Location,
			&device.Status,
			&device.Metadata,
			&device.CreatedAt,
			&device.UpdatedAt,
			&device.LastSeen,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan device: %w", err)
		}
		devices = append(devices, device)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return devices, nil
}

// Update updates a device
func (r *Repository) Update(id string, req *models.UpdateDeviceRequest) (*models.Device, error) {
	device, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		device.Name = req.Name
	}
	if req.Type != "" {
		device.Type = req.Type
	}
	if req.Location != "" {
		device.Location = req.Location
	}
	if req.Status != "" {
		device.Status = req.Status
	}
	if req.Metadata != "" {
		device.Metadata = req.Metadata
	}

	device.UpdatedAt = time.Now()

	query := `
		UPDATE devices 
		SET name = $1, type = $2, location = $3, status = $4, metadata = $5, updated_at = $6
		WHERE id = $7
	`

	_, err = r.db.Exec(query, device.Name, device.Type, device.Location,
		device.Status, device.Metadata, device.UpdatedAt, device.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update device: %w", err)
	}

	return device, nil
}

// Delete deletes a device
func (r *Repository) Delete(id string) error {
	query := `DELETE FROM devices WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete device: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("device not found")
	}

	return nil
}

// UpdateStatus updates the status and last seen time of a device
func (r *Repository) UpdateStatus(id string, status string) error {
	query := `
		UPDATE devices 
		SET status = $1, last_seen = $2, updated_at = $3
		WHERE id = $4
	`

	_, err := r.db.Exec(query, status, time.Now(), time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update device status: %w", err)
	}

	return nil
}
