package database

import (
	"context"

	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
)

// Create mock using:
// mockgen -source=webapp/backend/pkg/database/interface.go -destination=webapp/backend/pkg/database/mock/mock_database.go
type DeviceRepo interface {
	Close() error
	HealthCheck(ctx context.Context) error

	RegisterDevice(ctx context.Context, dev models.Device) error
	GetDevices(ctx context.Context) ([]models.Device, error)
	UpdateDevice(ctx context.Context, wwn string, collectorSmartData collector.SmartInfo) (models.Device, error)
	UpdateDeviceStatus(ctx context.Context, wwn string, status pkg.DeviceStatus) (models.Device, error)
	GetDeviceDetails(ctx context.Context, wwn string) (models.Device, error)
	UpdateDeviceArchived(ctx context.Context, wwn string, archived bool) error
	DeleteDevice(ctx context.Context, wwn string) error

	SaveSmartAttributes(ctx context.Context, wwn string, collectorSmartData collector.SmartInfo) (measurements.Smart, error)
	GetSmartAttributeHistory(ctx context.Context, wwn string, durationKey string, selectEntries int, selectEntriesOffset int, attributes []string) ([]measurements.Smart, error)

	SaveSmartTemperature(ctx context.Context, wwn string, deviceProtocol string, collectorSmartData collector.SmartInfo) error

	GetSummary(ctx context.Context) (map[string]*models.DeviceSummary, error)
	GetSmartTemperatureHistory(ctx context.Context, durationKey string) (map[string][]measurements.SmartTemperature, error)

	LoadSettings(ctx context.Context) (*models.Settings, error)
	SaveSettings(ctx context.Context, settings models.Settings) error

	// ZFS Pool methods
	RegisterZfsPools(ctx context.Context, pools []models.ZfsPool) error
	GetZfsPools(ctx context.Context) ([]models.ZfsPool, error)
	GetZfsPoolByGuid(ctx context.Context, poolGuid string) (models.ZfsPool, error)
	GetZfsPoolsByHost(ctx context.Context, hostId string) ([]models.ZfsPool, error)
	DeleteZfsPool(ctx context.Context, poolGuid string) error

	// ZFS Dataset methods
	RegisterZfsDatasets(ctx context.Context, datasets []models.ZfsDataset) error
	GetZfsDatasets(ctx context.Context) ([]models.ZfsDataset, error)
	GetZfsDatasetsByPool(ctx context.Context, poolName string) ([]models.ZfsDataset, error)
	GetZfsDatasetsByHost(ctx context.Context, hostId string) ([]models.ZfsDataset, error)
}
