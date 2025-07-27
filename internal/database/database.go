package database

import (
	"database/sql"
	"fmt"
	"log"

	"iot-platform-go/internal/config"

	_ "github.com/lib/pq"
)

// Database represents the database connection
type Database struct {
	*sql.DB
}

// New creates a new database connection
func New(cfg *config.Config) (*Database, error) {
	db, err := sql.Open("postgres", cfg.GetDatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &Database{db}
	
	// Initialize tables
	if err := database.initTables(); err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	log.Println("Database connection established successfully")
	return database, nil
}

// initTables creates the necessary tables if they don't exist
func (db *Database) initTables() error {
	// Create devices table
	createDevicesTable := `
	CREATE TABLE IF NOT EXISTS devices (
		id VARCHAR(255) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		type VARCHAR(100) NOT NULL,
		location VARCHAR(255) NOT NULL,
		status VARCHAR(50) DEFAULT 'offline',
		last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		metadata TEXT
	);`

	// Create device_data table
	createDeviceDataTable := `
	CREATE TABLE IF NOT EXISTS device_data (
		id VARCHAR(255) PRIMARY KEY,
		device_id VARCHAR(255) NOT NULL,
		timestamp TIMESTAMP NOT NULL,
		data JSONB NOT NULL,
		FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
	);`

	// Create indexes
	createIndexes := `
	CREATE INDEX IF NOT EXISTS idx_device_data_device_id ON device_data(device_id);
	CREATE INDEX IF NOT EXISTS idx_device_data_timestamp ON device_data(timestamp);
	CREATE INDEX IF NOT EXISTS idx_devices_status ON devices(status);
	CREATE INDEX IF NOT EXISTS idx_devices_type ON devices(type);
	`

	if _, err := db.Exec(createDevicesTable); err != nil {
		return fmt.Errorf("failed to create devices table: %w", err)
	}

	if _, err := db.Exec(createDeviceDataTable); err != nil {
		return fmt.Errorf("failed to create device_data table: %w", err)
	}

	if _, err := db.Exec(createIndexes); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	log.Println("Database tables initialized successfully")
	return nil
}

// Close closes the database connection
func (db *Database) Close() error {
	return db.DB.Close()
} 