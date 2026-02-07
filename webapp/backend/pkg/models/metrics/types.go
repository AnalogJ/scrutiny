package metrics

import (
	"time"

	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
)

// DeviceMetricsData stores metrics data for a single device
type DeviceMetricsData struct {
	Device    models.Device      `json:"device"`
	SmartData measurements.Smart `json:"smart_data"`
	UpdatedAt time.Time          `json:"updated_at"`
}
