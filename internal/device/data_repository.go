package device

import (
	"database/sql"
	"fmt"
	"time"

	"iot-platform-go/internal/database"
	"iot-platform-go/pkg/models"
)

// DataRepositoryInterface defines the interface for device data repository operations
type DataRepositoryInterface interface {
	SaveData(data *models.DeviceData) error
	GetDeviceData(deviceID string, limit int) ([]*models.DeviceData, error)
	GetDeviceDataByType(deviceID string, dataType string, limit int) ([]*models.DeviceData, error)
	GetLatestData(deviceID string) (*models.DeviceData, error)
	DeleteOldData(deviceID string, olderThan time.Time) error
}

// DataRepository handles database operations for device data
type DataRepository struct {
	db *database.Database
}

// NewDataRepository creates a new device data repository
func NewDataRepository(db *database.Database) *DataRepository {
	return &DataRepository{db: db}
}

// SaveData saves device data to the database
func (r *DataRepository) SaveData(data *models.DeviceData) error {
	query := `
		INSERT INTO device_data (id, device_id, timestamp, data_type, value, unit, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(query, data.ID, data.DeviceID, data.Timestamp, data.DataType, data.Value, data.Unit, data.Metadata)
	if err != nil {
		return fmt.Errorf("failed to save device data: %w", err)
	}

	return nil
}

// GetDeviceData retrieves device data with limit
func (r *DataRepository) GetDeviceData(deviceID string, limit int) ([]*models.DeviceData, error) {
	query := `
		SELECT id, device_id, timestamp, data_type, value, unit, metadata
		FROM device_data 
		WHERE device_id = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`

	rows, err := r.db.Query(query, deviceID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query device data: %w", err)
	}
	defer rows.Close()

	var data []*models.DeviceData
	for rows.Next() {
		item := &models.DeviceData{}
		err := rows.Scan(
			&item.ID,
			&item.DeviceID,
			&item.Timestamp,
			&item.DataType,
			&item.Value,
			&item.Unit,
			&item.Metadata,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan device data: %w", err)
		}
		data = append(data, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return data, nil
}

// GetDeviceDataByType retrieves device data filtered by data type
func (r *DataRepository) GetDeviceDataByType(deviceID string, dataType string, limit int) ([]*models.DeviceData, error) {
	query := `
		SELECT id, device_id, timestamp, data_type, value, unit, metadata
		FROM device_data 
		WHERE device_id = $1 AND data_type = $2
		ORDER BY timestamp DESC
		LIMIT $3
	`

	rows, err := r.db.Query(query, deviceID, dataType, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query device data by type: %w", err)
	}
	defer rows.Close()

	var data []*models.DeviceData
	for rows.Next() {
		item := &models.DeviceData{}
		err := rows.Scan(
			&item.ID,
			&item.DeviceID,
			&item.Timestamp,
			&item.DataType,
			&item.Value,
			&item.Unit,
			&item.Metadata,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan device data: %w", err)
		}
		data = append(data, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return data, nil
}

// GetLatestData retrieves the most recent data for a device
func (r *DataRepository) GetLatestData(deviceID string) (*models.DeviceData, error) {
	query := `
		SELECT id, device_id, timestamp, data_type, value, unit, metadata
		FROM device_data 
		WHERE device_id = $1
		ORDER BY timestamp DESC
		LIMIT 1
	`

	data := &models.DeviceData{}
	err := r.db.QueryRow(query, deviceID).Scan(
		&data.ID,
		&data.DeviceID,
		&data.Timestamp,
		&data.DataType,
		&data.Value,
		&data.Unit,
		&data.Metadata,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no data found for device")
		}
		return nil, fmt.Errorf("failed to get latest device data: %w", err)
	}

	return data, nil
}

// DeleteOldData deletes device data older than the specified time
func (r *DataRepository) DeleteOldData(deviceID string, olderThan time.Time) error {
	query := `DELETE FROM device_data WHERE device_id = $1 AND timestamp < $2`

	result, err := r.db.Exec(query, deviceID, olderThan)
	if err != nil {
		return fmt.Errorf("failed to delete old device data: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	fmt.Printf("Deleted %d old data records for device %s", rowsAffected, deviceID)
	return nil
}
