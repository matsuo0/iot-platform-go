package device

import (
	"testing"

	"iot-platform-go/internal/config"
	"iot-platform-go/internal/database"
	"iot-platform-go/pkg/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDatabase(t *testing.T) *database.Database {
	// テスト用の設定
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			Name:     "iot_platform_test",
			User:     "postgres",
			Password: "password",
			SSLMode:  "disable",
		},
	}

	db, err := database.New(cfg)
	require.NoError(t, err)

	// テスト用のテーブルをクリーンアップ
	_, err = db.Exec("DELETE FROM device_data")
	require.NoError(t, err)
	_, err = db.Exec("DELETE FROM devices")
	require.NoError(t, err)

	return db
}

func createTestDeviceRequest() *models.CreateDeviceRequest {
	return &models.CreateDeviceRequest{
		Name:     "Test Device",
		Type:     "temperature",
		Location: "Test Room",
		Metadata: `{"manufacturer":"Test Corp","model":"TEMP-001"}`,
	}
}

func createTestUpdateRequest() *models.UpdateDeviceRequest {
	return &models.UpdateDeviceRequest{
		Name:     "Updated Device",
		Type:     "humidity",
		Location: "Updated Room",
		Metadata: `{"manufacturer":"Updated Corp","model":"HUM-001"}`,
	}
}

func TestRepository_Create(t *testing.T) {
	t.Skip("Skipping repository test as it requires database setup")
	db := setupTestDatabase(t)
	defer db.Close()

	repo := NewRepository(db)

	tests := []struct {
		name    string
		request *models.CreateDeviceRequest
		wantErr bool
	}{
		{
			name:    "successful device creation",
			request: createTestDeviceRequest(),
			wantErr: false,
		},
		{
			name: "device with empty name",
			request: &models.CreateDeviceRequest{
				Name:     "",
				Type:     "temperature",
				Location: "Test Room",
			},
			wantErr: false, // データベースレベルではエラーにならない
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			device, err := repo.Create(tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, device)
			assert.NotEmpty(t, device.ID)
			assert.Equal(t, tt.request.Name, device.Name)
			assert.Equal(t, tt.request.Type, device.Type)
			assert.Equal(t, tt.request.Location, device.Location)
			assert.Equal(t, "offline", device.Status)
			assert.NotNil(t, device.CreatedAt)
			assert.NotNil(t, device.UpdatedAt)
		})
	}
}

func TestRepository_GetByID(t *testing.T) {
	t.Skip("Skipping repository test as it requires database setup")
	db := setupTestDatabase(t)
	defer db.Close()

	repo := NewRepository(db)

	// テスト用のデバイスを作成
	createReq := createTestDeviceRequest()
	createdDevice, err := repo.Create(createReq)
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "successful device retrieval",
			id:      createdDevice.ID,
			wantErr: false,
		},
		{
			name:    "device not found",
			id:      "non-existent-id",
			wantErr: true,
		},
		{
			name:    "empty device ID",
			id:      "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			device, err := repo.GetByID(tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, device)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, device)
			assert.Equal(t, tt.id, device.ID)
			assert.Equal(t, createReq.Name, device.Name)
			assert.Equal(t, createReq.Type, device.Type)
			assert.Equal(t, createReq.Location, device.Location)
		})
	}
}

func TestRepository_GetAll(t *testing.T) {
	t.Skip("Skipping repository test as it requires database setup")
	db := setupTestDatabase(t)
	defer db.Close()

	repo := NewRepository(db)

	// 複数のテスト用デバイスを作成
	devices := []*models.CreateDeviceRequest{
		{Name: "Device 1", Type: "temperature", Location: "Room 1"},
		{Name: "Device 2", Type: "humidity", Location: "Room 2"},
		{Name: "Device 3", Type: "pressure", Location: "Room 3"},
	}

	for _, deviceReq := range devices {
		_, err := repo.Create(deviceReq)
		require.NoError(t, err)
	}

	t.Run("successful devices retrieval", func(t *testing.T) {
		retrievedDevices, err := repo.GetAll()

		assert.NoError(t, err)
		assert.NotNil(t, retrievedDevices)
		assert.Len(t, retrievedDevices, 3)

		// デバイスが作成日時の降順で取得されることを確認
		for i := 0; i < len(retrievedDevices)-1; i++ {
			assert.True(t, retrievedDevices[i].CreatedAt.After(retrievedDevices[i+1].CreatedAt) ||
				retrievedDevices[i].CreatedAt.Equal(retrievedDevices[i+1].CreatedAt))
		}
	})
}

func TestRepository_Update(t *testing.T) {
	t.Skip("Skipping repository test as it requires database setup")
	db := setupTestDatabase(t)
	defer db.Close()

	repo := NewRepository(db)

	// テスト用のデバイスを作成
	createReq := createTestDeviceRequest()
	createdDevice, err := repo.Create(createReq)
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		request *models.UpdateDeviceRequest
		wantErr bool
	}{
		{
			name:    "successful device update",
			id:      createdDevice.ID,
			request: createTestUpdateRequest(),
			wantErr: false,
		},
		{
			name:    "device not found",
			id:      "non-existent-id",
			request: createTestUpdateRequest(),
			wantErr: true,
		},
		{
			name: "partial update",
			id:   createdDevice.ID,
			request: &models.UpdateDeviceRequest{
				Name: "Partially Updated Device",
				// Type と Location は更新しない
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedDevice, err := repo.Update(tt.id, tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, updatedDevice)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, updatedDevice)
			assert.Equal(t, tt.id, updatedDevice.ID)

			// 更新されたフィールドを確認
			if tt.request.Name != "" {
				assert.Equal(t, tt.request.Name, updatedDevice.Name)
			}
			if tt.request.Type != "" {
				assert.Equal(t, tt.request.Type, updatedDevice.Type)
			}
			if tt.request.Location != "" {
				assert.Equal(t, tt.request.Location, updatedDevice.Location)
			}
			if tt.request.Metadata != "" {
				assert.Equal(t, tt.request.Metadata, updatedDevice.Metadata)
			}

			// UpdatedAtが更新されていることを確認
			assert.True(t, updatedDevice.UpdatedAt.After(createdDevice.UpdatedAt) ||
				updatedDevice.UpdatedAt.Equal(createdDevice.UpdatedAt))
		})
	}
}

func TestRepository_Delete(t *testing.T) {
	t.Skip("Skipping repository test as it requires database setup")
	db := setupTestDatabase(t)
	defer db.Close()

	repo := NewRepository(db)

	// テスト用のデバイスを作成
	createReq := createTestDeviceRequest()
	createdDevice, err := repo.Create(createReq)
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "successful device deletion",
			id:      createdDevice.ID,
			wantErr: false,
		},
		{
			name:    "device not found",
			id:      "non-existent-id",
			wantErr: true,
		},
		{
			name:    "empty device ID",
			id:      "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// デバイスが実際に削除されたことを確認
			_, err = repo.GetByID(tt.id)
			assert.Error(t, err)
		})
	}
}

func TestRepository_UpdateStatus(t *testing.T) {
	t.Skip("Skipping repository test as it requires database setup")
	db := setupTestDatabase(t)
	defer db.Close()

	repo := NewRepository(db)

	// テスト用のデバイスを作成
	createReq := createTestDeviceRequest()
	createdDevice, err := repo.Create(createReq)
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		status  string
		wantErr bool
	}{
		{
			name:    "successful status update to online",
			id:      createdDevice.ID,
			status:  "online",
			wantErr: false,
		},
		{
			name:    "successful status update to offline",
			id:      createdDevice.ID,
			status:  "offline",
			wantErr: false,
		},
		{
			name:    "device not found",
			id:      "non-existent-id",
			status:  "online",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.UpdateStatus(tt.id, tt.status)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// ステータスが実際に更新されたことを確認
			device, err := repo.GetByID(tt.id)
			assert.NoError(t, err)
			assert.Equal(t, tt.status, device.Status)
		})
	}
}

func TestRepository_Integration(t *testing.T) {
	t.Skip("Skipping repository test as it requires database setup")
	db := setupTestDatabase(t)
	defer db.Close()

	repo := NewRepository(db)

	t.Run("full CRUD operations", func(t *testing.T) {
		// Create
		createReq := createTestDeviceRequest()
		device, err := repo.Create(createReq)
		assert.NoError(t, err)
		assert.NotNil(t, device)

		// Read
		retrievedDevice, err := repo.GetByID(device.ID)
		assert.NoError(t, err)
		assert.Equal(t, device.ID, retrievedDevice.ID)

		// Update
		updateReq := createTestUpdateRequest()
		updatedDevice, err := repo.Update(device.ID, updateReq)
		assert.NoError(t, err)
		assert.Equal(t, updateReq.Name, updatedDevice.Name)

		// Update Status
		err = repo.UpdateStatus(device.ID, "online")
		assert.NoError(t, err)

		// Verify status update
		statusDevice, err := repo.GetByID(device.ID)
		assert.NoError(t, err)
		assert.Equal(t, "online", statusDevice.Status)

		// Delete
		err = repo.Delete(device.ID)
		assert.NoError(t, err)

		// Verify deletion
		_, err = repo.GetByID(device.ID)
		assert.Error(t, err)
	})
}

func TestRepository_DataValidation(t *testing.T) {
	t.Skip("Skipping repository test as it requires database setup")
	db := setupTestDatabase(t)
	defer db.Close()

	repo := NewRepository(db)

	t.Run("device with special characters", func(t *testing.T) {
		createReq := &models.CreateDeviceRequest{
			Name:     "Device with 特殊文字 & Symbols!@#$%",
			Type:     "temperature",
			Location: "Room with 日本語",
			Metadata: `{"special":"value with 特殊文字","number":123.45,"boolean":true}`,
		}

		device, err := repo.Create(createReq)
		assert.NoError(t, err)
		assert.Equal(t, createReq.Name, device.Name)
		assert.Equal(t, createReq.Location, device.Location)
		assert.Equal(t, createReq.Metadata, device.Metadata)
	})

	t.Run("device with very long name", func(t *testing.T) {
		longName := string(make([]byte, 1000)) // 1000文字の文字列
		for i := range longName {
			longName = longName[:i] + "a" + longName[i+1:]
		}

		createReq := &models.CreateDeviceRequest{
			Name:     longName,
			Type:     "temperature",
			Location: "Test Room",
		}

		device, err := repo.Create(createReq)
		assert.NoError(t, err)
		assert.Equal(t, createReq.Name, device.Name)
	})
}
