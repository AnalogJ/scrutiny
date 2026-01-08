package models

import (
	"time"
)

// ZFSPoolStatus represents the health status of a ZFS pool
type ZFSPoolStatus string

const (
	ZFSPoolStatusOnline   ZFSPoolStatus = "ONLINE"
	ZFSPoolStatusDegraded ZFSPoolStatus = "DEGRADED"
	ZFSPoolStatusFaulted  ZFSPoolStatus = "FAULTED"
	ZFSPoolStatusOffline  ZFSPoolStatus = "OFFLINE"
	ZFSPoolStatusRemoved  ZFSPoolStatus = "REMOVED"
	ZFSPoolStatusUnavail  ZFSPoolStatus = "UNAVAIL"
)

// ZFSScrubState represents the state of a ZFS scrub operation
type ZFSScrubState string

const (
	ZFSScrubStateNone     ZFSScrubState = "none"
	ZFSScrubStateScanning ZFSScrubState = "scanning"
	ZFSScrubStateFinished ZFSScrubState = "finished"
	ZFSScrubStateCanceled ZFSScrubState = "canceled"
)

// ZFSPoolWrapper wraps the response for ZFS pool API calls
type ZFSPoolWrapper struct {
	Success bool      `json:"success"`
	Errors  []error   `json:"errors,omitempty"`
	Data    []ZFSPool `json:"data"`
}

// ZFSPool represents a ZFS storage pool
type ZFSPool struct {
	// GORM attributes
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	// Pool identifier (GUID) - primary key
	GUID   string `json:"guid" gorm:"primary_key"`
	Name   string `json:"name"`
	HostID string `json:"host_id"`

	// Pool status
	Status ZFSPoolStatus `json:"status"`
	Health string        `json:"health"`

	// Capacity metrics (in bytes)
	Size            int64   `json:"size"`
	Allocated       int64   `json:"allocated"`
	Free            int64   `json:"free"`
	Fragmentation   int     `json:"fragmentation"`
	CapacityPercent float64 `json:"capacity_percent"`

	// Pool configuration
	Ashift int `json:"ashift"`

	// Scrub status (stored as JSON in SQLite)
	ScrubState          ZFSScrubState `json:"scrub_state"`
	ScrubStartTime      *time.Time    `json:"scrub_start_time,omitempty"`
	ScrubEndTime        *time.Time    `json:"scrub_end_time,omitempty"`
	ScrubScannedBytes   int64         `json:"scrub_scanned_bytes"`
	ScrubIssuedBytes    int64         `json:"scrub_issued_bytes"`
	ScrubTotalBytes     int64         `json:"scrub_total_bytes"`
	ScrubErrorsCount    int64         `json:"scrub_errors_count"`
	ScrubPercentComplete float64      `json:"scrub_percent_complete"`

	// Aggregate error counts (sum of all vdevs)
	TotalReadErrors     int64 `json:"total_read_errors"`
	TotalWriteErrors    int64 `json:"total_write_errors"`
	TotalChecksumErrors int64 `json:"total_checksum_errors"`

	// User provided metadata
	Label string `json:"label,omitempty"`

	// Pool management flags
	Archived bool `json:"archived"`
	Muted    bool `json:"muted"`

	// Vdevs associated with this pool (loaded via preload, not stored in this table)
	Vdevs []ZFSVdev `json:"vdevs,omitempty" gorm:"foreignKey:PoolGUID;references:GUID"`
}

// IsHealthy returns true if the pool status is ONLINE
func (p *ZFSPool) IsHealthy() bool {
	return p.Status == ZFSPoolStatusOnline
}

// IsDegraded returns true if the pool status is DEGRADED
func (p *ZFSPool) IsDegraded() bool {
	return p.Status == ZFSPoolStatusDegraded
}

// IsFaulted returns true if the pool status is FAULTED
func (p *ZFSPool) IsFaulted() bool {
	return p.Status == ZFSPoolStatusFaulted
}

// HasErrors returns true if the pool has any read, write, or checksum errors
func (p *ZFSPool) HasErrors() bool {
	return p.TotalReadErrors > 0 || p.TotalWriteErrors > 0 || p.TotalChecksumErrors > 0
}

// IsScrubbing returns true if a scrub is currently in progress
func (p *ZFSPool) IsScrubbing() bool {
	return p.ScrubState == ZFSScrubStateScanning
}
