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

// ZFSVdevType represents the type of a ZFS vdev
type ZFSVdevType string

const (
	ZFSVdevTypeDisk    ZFSVdevType = "disk"
	ZFSVdevTypeFile    ZFSVdevType = "file"
	ZFSVdevTypeMirror  ZFSVdevType = "mirror"
	ZFSVdevTypeRaidz1  ZFSVdevType = "raidz1"
	ZFSVdevTypeRaidz2  ZFSVdevType = "raidz2"
	ZFSVdevTypeRaidz3  ZFSVdevType = "raidz3"
	ZFSVdevTypeDraid1  ZFSVdevType = "draid1"
	ZFSVdevTypeDraid2  ZFSVdevType = "draid2"
	ZFSVdevTypeDraid3  ZFSVdevType = "draid3"
	ZFSVdevTypeSpare   ZFSVdevType = "spare"
	ZFSVdevTypeLog     ZFSVdevType = "log"
	ZFSVdevTypeCache   ZFSVdevType = "cache"
	ZFSVdevTypeSpecial ZFSVdevType = "special"
	ZFSVdevTypeDedup   ZFSVdevType = "dedup"
)

// ZFSPoolWrapper wraps the response for ZFS pool API calls
type ZFSPoolWrapper struct {
	Success bool      `json:"success"`
	Errors  []error   `json:"errors,omitempty"`
	Data    []ZFSPool `json:"data"`
}

// ZFSPool represents a ZFS storage pool
type ZFSPool struct {
	GUID   string `json:"guid"`
	Name   string `json:"name"`
	HostID string `json:"host_id"`

	Status ZFSPoolStatus `json:"status"`
	Health string        `json:"health"`

	Size            int64   `json:"size"`
	Allocated       int64   `json:"allocated"`
	Free            int64   `json:"free"`
	Fragmentation   int     `json:"fragmentation"`
	CapacityPercent float64 `json:"capacity_percent"`

	Ashift int `json:"ashift"`

	ScrubState           ZFSScrubState `json:"scrub_state"`
	ScrubStartTime       *time.Time    `json:"scrub_start_time,omitempty"`
	ScrubEndTime         *time.Time    `json:"scrub_end_time,omitempty"`
	ScrubScannedBytes    int64         `json:"scrub_scanned_bytes"`
	ScrubIssuedBytes     int64         `json:"scrub_issued_bytes"`
	ScrubTotalBytes      int64         `json:"scrub_total_bytes"`
	ScrubErrorsCount     int64         `json:"scrub_errors_count"`
	ScrubPercentComplete float64       `json:"scrub_percent_complete"`

	TotalReadErrors     int64 `json:"total_read_errors"`
	TotalWriteErrors    int64 `json:"total_write_errors"`
	TotalChecksumErrors int64 `json:"total_checksum_errors"`

	Vdevs []ZFSVdev `json:"vdevs,omitempty"`
}

// ZFSVdev represents a virtual device in a ZFS pool
type ZFSVdev struct {
	Name   string        `json:"name"`
	Type   ZFSVdevType   `json:"type"`
	Status ZFSPoolStatus `json:"status"`
	GUID   string        `json:"guid,omitempty"`
	Path   string        `json:"path,omitempty"`

	ReadErrors     int64 `json:"read_errors"`
	WriteErrors    int64 `json:"write_errors"`
	ChecksumErrors int64 `json:"checksum_errors"`

	Size      int64 `json:"size,omitempty"`
	Allocated int64 `json:"allocated,omitempty"`

	Children []ZFSVdev `json:"children,omitempty"`
}
