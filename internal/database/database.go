package database

import (
	"database/sql"
	"fmt"
	"log"

	"iot-platform-go/internal/config"

	_ "github.com/lib/pq"
)

// Database represents the database connection.
type Database struct {
	*sql.DB
}

// New creates a new database connection.
func New(cfg *config.Config) (*Database, error) {
	dsn := cfg.GetDatabaseURL()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &Database{DB: db}

	// Initialize tables
	if err := database.initTables(); err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	return database, nil
}

// initTables creates the necessary tables if they don't exist.
func (d *Database) initTables() error {
	// Create devices table
	createDevicesTable := `
		CREATE TABLE IF NOT EXISTS devices (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			type VARCHAR(100) NOT NULL,
			location VARCHAR(255),
			status VARCHAR(50) DEFAULT 'offline',
			metadata TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			last_seen TIMESTAMP
		)
	`

	_, err := d.Exec(createDevicesTable)
	if err != nil {
		return fmt.Errorf("failed to create devices table: %w", err)
	}

	// Create device_data table
	createDeviceDataTable := `
		CREATE TABLE IF NOT EXISTS device_data (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			device_id UUID NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			data_type VARCHAR(100) NOT NULL,
			value REAL NOT NULL,
			unit VARCHAR(50),
			metadata TEXT
		)
	`

	_, err = d.Exec(createDeviceDataTable)
	if err != nil {
		return fmt.Errorf("failed to create device_data table: %w", err)
	}

	// Create indexes
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_devices_status ON devices(status)",
		"CREATE INDEX IF NOT EXISTS idx_devices_type ON devices(type)",
		"CREATE INDEX IF NOT EXISTS idx_device_data_device_id ON device_data(device_id)",
		"CREATE INDEX IF NOT EXISTS idx_device_data_timestamp ON device_data(timestamp)",
		"CREATE INDEX IF NOT EXISTS idx_device_data_type ON device_data(data_type)",
	}

	for _, index := range indexes {
		_, err := d.Exec(index)
		if err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	log.Println("Database tables initialized successfully")
	return nil
}

// Close closes the database connection.
func (d *Database) Close() error {
	return d.DB.Close()
}
