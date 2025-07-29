package influxdb

import (
	"context"
	"fmt"
	"log"
	"time"

	"iot-platform-go/internal/config"
	"iot-platform-go/pkg/models"

	"github.com/google/uuid"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

// Client represents an InfluxDB client
type Client struct {
	client   influxdb2.Client
	writeAPI api.WriteAPIBlocking
	queryAPI api.QueryAPI
	config   *config.InfluxDBConfig
}

// NewClient creates a new InfluxDB client
func NewClient(cfg *config.InfluxDBConfig) (*Client, error) {
	client := influxdb2.NewClient(cfg.URL, cfg.Token)

	// Test the connection
	_, err := client.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to InfluxDB: %w", err)
	}

	writeAPI := client.WriteAPIBlocking(cfg.Org, cfg.Bucket)
	queryAPI := client.QueryAPI(cfg.Org)

	log.Printf("✅ Connected to InfluxDB at %s", cfg.URL)

	return &Client{
		client:   client,
		writeAPI: writeAPI,
		queryAPI: queryAPI,
		config:   cfg,
	}, nil
}

// WriteDeviceData writes device data to InfluxDB
func (c *Client) WriteDeviceData(data *models.DeviceData) error {
	point := influxdb2.NewPoint(
		"device_data",
		map[string]string{
			"device_id": data.DeviceID,
			"data_type": data.DataType,
			"unit":      data.Unit,
		},
		map[string]interface{}{
			"value": data.Value,
		},
		data.Timestamp,
	)

	err := c.writeAPI.WritePoint(context.Background(), point)
	if err != nil {
		return fmt.Errorf("failed to write data point: %w", err)
	}

	return nil
}

// QueryDeviceData queries device data from InfluxDB
func (c *Client) QueryDeviceData(deviceID string, dataType string, start time.Time, end time.Time, limit int) (
	[]*models.DeviceData, error) {
	query := fmt.Sprintf(`
		from(bucket: %q)
			|> range(start: %s, stop: %s)
			|> filter(fn: (r) => r["_measurement"] == "device_data")
			|> filter(fn: (r) => r["device_id"] == %q)
	`, c.config.Bucket, start.Format(time.RFC3339), end.Format(time.RFC3339), deviceID)

	if dataType != "" {
		query += fmt.Sprintf(`|> filter(fn: (r) => r["data_type"] == %q)`, dataType)
	}

	query += fmt.Sprintf(`
		|> sort(columns: ["_time"])
		|> limit(n: %d)
	`, limit)

	result, err := c.queryAPI.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}
	defer result.Close()

	var dataPoints []*models.DeviceData
	for result.Next() {
		record := result.Record()

		// Parse the value using type switch
		var value float64
		switch v := record.Value().(type) {
		case float64:
			value = v
		case int64:
			value = float64(v)
		default:
			continue // Skip non-numeric values
		}

		// Get device_id from tags
		deviceID := ""
		if deviceIDVal, ok := record.ValueByKey("device_id").(string); ok {
			deviceID = deviceIDVal
		}

		// Get data_type from tags
		dataType := ""
		if dataTypeVal, ok := record.ValueByKey("data_type").(string); ok {
			dataType = dataTypeVal
		}

		// Get unit from tags
		unit := ""
		if unitVal, ok := record.ValueByKey("unit").(string); ok {
			unit = unitVal
		}

		dataPoint := &models.DeviceData{
			ID:        uuid.New().String(), // Generate new UUID for API response
			DeviceID:  deviceID,
			Timestamp: record.Time(),
			DataType:  dataType,
			Value:     value,
			Unit:      unit,
			Metadata:  "",
		}
		dataPoints = append(dataPoints, dataPoint)
	}

	return dataPoints, nil
}

// GetLatestDeviceData gets the latest data point for a device
func (c *Client) GetLatestDeviceData(deviceID string, dataType string) (*models.DeviceData, error) {
	end := time.Now()
	start := end.Add(-24 * time.Hour) // Last 24 hours

	query := fmt.Sprintf(`
		from(bucket: %q)
			|> range(start: %s, stop: %s)
			|> filter(fn: (r) => r["_measurement"] == "device_data")
			|> filter(fn: (r) => r["device_id"] == %q)
	`, c.config.Bucket, start.Format(time.RFC3339), end.Format(time.RFC3339), deviceID)

	if dataType != "" {
		query += fmt.Sprintf(`|> filter(fn: (r) => r["data_type"] == %q)`, dataType)
	}

	query += `
		|> sort(columns: ["_time"], desc: true)
		|> limit(n: 1)
	`

	result, err := c.queryAPI.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to query latest data: %w", err)
	}
	defer result.Close()

	if !result.Next() {
		return nil, fmt.Errorf("no data found for device %s", deviceID)
	}

	record := result.Record()

	// Parse the value using type switch
	var value float64
	switch v := record.Value().(type) {
	case float64:
		value = v
	case int64:
		value = float64(v)
	default:
		return nil, fmt.Errorf("invalid value type for device %s", deviceID)
	}

	// Get device_id from tags
	deviceIDFromRecord := ""
	if deviceIDVal, ok := record.ValueByKey("device_id").(string); ok {
		deviceIDFromRecord = deviceIDVal
	}

	// Get data_type from tags
	dataTypeFromRecord := ""
	if dataTypeVal, ok := record.ValueByKey("data_type").(string); ok {
		dataTypeFromRecord = dataTypeVal
	}

	// Get unit from tags
	unit := ""
	if unitVal, ok := record.ValueByKey("unit").(string); ok {
		unit = unitVal
	}

	return &models.DeviceData{
		ID:        uuid.New().String(), // Generate new UUID for API response
		DeviceID:  deviceIDFromRecord,
		Timestamp: record.Time(),
		DataType:  dataTypeFromRecord,
		Value:     value,
		Unit:      unit,
		Metadata:  "",
	}, nil
}

// Close closes the InfluxDB client
func (c *Client) Close() {
	c.client.Close()
}
