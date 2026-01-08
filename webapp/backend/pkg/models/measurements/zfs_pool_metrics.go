package measurements

import (
	"time"
)

// ZFSPoolMetrics represents time-series metrics for a ZFS pool stored in InfluxDB
type ZFSPoolMetrics struct {
	Date     time.Time `json:"date"`
	PoolGUID string    `json:"pool_guid"` // tag
	PoolName string    `json:"pool_name"` // tag

	// Capacity metrics (fields)
	Size            int64   `json:"size"`
	Allocated       int64   `json:"allocated"`
	Free            int64   `json:"free"`
	CapacityPercent float64 `json:"capacity_percent"`
	Fragmentation   int     `json:"fragmentation"`

	// Health status (field - stored as string)
	Status string `json:"status"`

	// Error counts (fields)
	ReadErrors     int64 `json:"read_errors"`
	WriteErrors    int64 `json:"write_errors"`
	ChecksumErrors int64 `json:"checksum_errors"`

	// Scrub metrics (fields)
	ScrubState   string  `json:"scrub_state"`
	ScrubPercent float64 `json:"scrub_percent"`
	ScrubErrors  int64   `json:"scrub_errors"`
}

// Flatten converts the ZFSPoolMetrics struct to tags and fields for InfluxDB
func (m *ZFSPoolMetrics) Flatten() (tags map[string]string, fields map[string]interface{}) {
	tags = map[string]string{
		"pool_guid": m.PoolGUID,
		"pool_name": m.PoolName,
	}

	fields = map[string]interface{}{
		"size":             m.Size,
		"allocated":        m.Allocated,
		"free":             m.Free,
		"capacity_percent": m.CapacityPercent,
		"fragmentation":    m.Fragmentation,
		"status":           m.Status,
		"read_errors":      m.ReadErrors,
		"write_errors":     m.WriteErrors,
		"checksum_errors":  m.ChecksumErrors,
		"scrub_state":      m.ScrubState,
		"scrub_percent":    m.ScrubPercent,
		"scrub_errors":     m.ScrubErrors,
	}

	return tags, fields
}

// NewZFSPoolMetricsFromInfluxDB creates a ZFSPoolMetrics from InfluxDB query result
func NewZFSPoolMetricsFromInfluxDB(attrs map[string]interface{}) (*ZFSPoolMetrics, error) {
	m := ZFSPoolMetrics{
		Date:     attrs["_time"].(time.Time),
		PoolGUID: attrs["pool_guid"].(string),
		PoolName: attrs["pool_name"].(string),
	}

	// Parse optional fields
	if val, ok := attrs["size"]; ok && val != nil {
		m.Size = val.(int64)
	}
	if val, ok := attrs["allocated"]; ok && val != nil {
		m.Allocated = val.(int64)
	}
	if val, ok := attrs["free"]; ok && val != nil {
		m.Free = val.(int64)
	}
	if val, ok := attrs["capacity_percent"]; ok && val != nil {
		m.CapacityPercent = val.(float64)
	}
	if val, ok := attrs["fragmentation"]; ok && val != nil {
		m.Fragmentation = int(val.(int64))
	}
	if val, ok := attrs["status"]; ok && val != nil {
		m.Status = val.(string)
	}
	if val, ok := attrs["read_errors"]; ok && val != nil {
		m.ReadErrors = val.(int64)
	}
	if val, ok := attrs["write_errors"]; ok && val != nil {
		m.WriteErrors = val.(int64)
	}
	if val, ok := attrs["checksum_errors"]; ok && val != nil {
		m.ChecksumErrors = val.(int64)
	}
	if val, ok := attrs["scrub_state"]; ok && val != nil {
		m.ScrubState = val.(string)
	}
	if val, ok := attrs["scrub_percent"]; ok && val != nil {
		m.ScrubPercent = val.(float64)
	}
	if val, ok := attrs["scrub_errors"]; ok && val != nil {
		m.ScrubErrors = val.(int64)
	}

	return &m, nil
}

// ZFSPoolCapacityHistory represents a simplified capacity data point for charts
type ZFSPoolCapacityHistory struct {
	Date            time.Time `json:"date"`
	CapacityPercent float64   `json:"capacity_percent"`
}
