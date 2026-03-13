package database

import (
	"context"

	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
	"github.com/gofrs/uuid/v5"
)

// Create mock using:
// mockgen -source=webapp/backend/pkg/database/interface.go -destination=webapp/backend/pkg/database/mock/mock_database.go
type DeviceRepo interface {
	Close() error
	HealthCheck(ctx context.Context) error

	RegisterDevice(ctx context.Context, dev models.Device) error
	GetDevices(ctx context.Context) ([]models.Device, error)
	UpdateDevice(ctx context.Context, scrutiny_uuid uuid.UUID, collectorSmartData collector.SmartInfo) (models.Device, error)
	UpdateDeviceStatus(ctx context.Context, scrutiny_uuid uuid.UUID, status pkg.DeviceStatus) (models.Device, error)
	GetDeviceDetails(ctx context.Context, scrutiny_uuid uuid.UUID) (models.Device, error)
	UpdateDeviceArchived(ctx context.Context, scrutiny_uuid uuid.UUID, archived bool) error
	DeleteDevice(ctx context.Context, scrutiny_uuid uuid.UUID) error

	SaveSmartAttributes(ctx context.Context, scrutiny_uuid uuid.UUID, collectorSmartData collector.SmartInfo) (measurements.Smart, error)
	GetSmartAttributeHistory(ctx context.Context, scrutiny_uuid uuid.UUID, durationKey string, selectEntries int, selectEntriesOffset int, attributes []string) ([]measurements.Smart, error)

	SaveSmartTemperature(ctx context.Context, scrutiny_uuid uuid.UUID, deviceProtocol string, collectorSmartData collector.SmartInfo, discardSCTTempHistory bool) error

	GetSummary(ctx context.Context) (map[uuid.UUID]*models.DeviceSummary, error)
	GetSmartTemperatureHistory(ctx context.Context, durationKey string) (map[uuid.UUID][]measurements.SmartTemperature, error)

	LoadSettings(ctx context.Context) (*models.Settings, error)
	SaveSettings(ctx context.Context, settings models.Settings) error
}
