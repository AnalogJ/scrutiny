package database

import (
	"context"

	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
)

type DeviceRepo interface {
	Close() error
	HealthCheck(ctx context.Context) error

	RegisterDevice(ctx context.Context, dev models.Device) error
	GetDevices(ctx context.Context) ([]models.Device, error)
	UpdateDevice(ctx context.Context, wwn string, collectorSmartData collector.SmartInfo) (models.Device, error)
	UpdateDeviceStatus(ctx context.Context, wwn string, status pkg.DeviceStatus) (models.Device, error)
	GetDeviceDetails(ctx context.Context, wwn string) (models.Device, error)
	DeleteDevice(ctx context.Context, wwn string) error

	SaveSmartAttributes(ctx context.Context, wwn string, collectorSmartData collector.SmartInfo) (measurements.Smart, error)
	GetSmartAttributeHistory(ctx context.Context, wwn string, durationKey string, n int, offset int, attributes []string) ([]measurements.Smart, error)

	SaveSmartTemperature(ctx context.Context, wwn string, deviceProtocol string, collectorSmartData collector.SmartInfo) error

	GetSummary(ctx context.Context) (map[string]*models.DeviceSummary, error)
	GetSmartTemperatureHistory(ctx context.Context, durationKey string) (map[string][]measurements.SmartTemperature, error)

	LoadSettings(ctx context.Context) (*models.Settings, error)
	SaveSettings(ctx context.Context, settings models.Settings) error
}
