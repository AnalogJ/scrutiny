package models

import (
	"time"
)

type ZfsPoolWrapper struct {
	Success bool      `json:"success"`
	Errors  []error   `json:"errors"`
	Data    []ZfsPool `json:"data"`
}

type ZfsPool struct {
	//GORM attributes
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	PoolGuid string `json:"pool_guid" gorm:"primary_key"`
	Name     string `json:"name"`
	HostId   string `json:"host_id"`

	// Pool status information
	State       string `json:"state"`        // ONLINE, DEGRADED, FAULTED, OFFLINE, etc.
	Txg         string `json:"txg"`          // Transaction group
	SpaVersion  string `json:"spa_version"`  // Storage pool version
	ZplVersion  string `json:"zpl_version"`  // ZFS POSIX layer version
	Status      string `json:"status"`       // Human-readable status message
	Action      string `json:"action"`       // Recommended action
	ErrorCount  string `json:"error_count"` // Total error count

	// Space information (from zpool status)
	AllocSpace string `json:"alloc_space"` // Allocated space
	TotalSpace string `json:"total_space"` // Total space
	DefSpace   string `json:"def_space"`   // Deferred space

	// Pool properties (from zpool list)
	Size            string `json:"size,omitempty"`             // Pool size
	Allocated       string `json:"allocated,omitempty"`        // Allocated space
	Free            string `json:"free,omitempty"`             // Free space
	Fragmentation   string `json:"fragmentation,omitempty"`    // Fragmentation percentage
	CapacityPercent string `json:"capacity_percent,omitempty"` // Capacity percentage
	Dedupratio      string `json:"dedupratio,omitempty"`       // Deduplication ratio

	// Error counters
	ReadErrors     string `json:"read_errors"`
	WriteErrors    string `json:"write_errors"`
	ChecksumErrors string `json:"checksum_errors"`

	// Scan information (scrub/resilver)
	ScanFunction   string `json:"scan_function,omitempty"`    // SCRUB, RESILVER
	ScanState      string `json:"scan_state,omitempty"`       // SCANNING, FINISHED, CANCELED
	ScanStartTime  string `json:"scan_start_time,omitempty"`  // Start time
	ScanEndTime    string `json:"scan_end_time,omitempty"`    // End time
	ScanToExamine  string `json:"scan_to_examine,omitempty"`  // Bytes to examine
	ScanExamined   string `json:"scan_examined,omitempty"`    // Bytes examined
	ScanProcessed  string `json:"scan_processed,omitempty"`   // Bytes processed
	ScanErrors     string `json:"scan_errors,omitempty"`      // Scan errors
	ScanIssued     string `json:"scan_issued,omitempty"`      // Bytes issued

	// Related vdevs
	Vdevs []ZfsVdev `json:"vdevs" gorm:"foreignKey:PoolGuid;references:PoolGuid"`
}

type ZfsVdev struct {
	//GORM attributes
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	ID       uint   `json:"id" gorm:"primary_key;autoIncrement"`
	PoolGuid string `json:"pool_guid"`
	ParentId *uint  `json:"parent_id,omitempty"` // For hierarchical vdevs
	Guid     string `json:"guid"`
	Name     string `json:"name"`

	// Vdev information
	VdevType string `json:"vdev_type"` // root, mirror, raidz, raidz2, raidz3, disk, file, etc.
	Class    string `json:"class"`     // normal, spare, cache, etc.
	State    string `json:"state"`     // ONLINE, DEGRADED, FAULTED, OFFLINE, etc.

	// Space information
	AllocSpace  string `json:"alloc_space,omitempty"`  // Allocated space
	TotalSpace  string `json:"total_space,omitempty"`  // Total space
	DefSpace    string `json:"def_space,omitempty"`    // Deferred space
	RepDevSize  string `json:"rep_dev_size,omitempty"` // Reported device size
	PhysSpace   string `json:"phys_space,omitempty"`   // Physical space

	// Error counters
	ReadErrors     string `json:"read_errors"`
	WriteErrors    string `json:"write_errors"`
	ChecksumErrors string `json:"checksum_errors"`
	SlowIos        string `json:"slow_ios,omitempty"`

	// Device path information (for leaf devices)
	Path     string `json:"path,omitempty"`      // Device path
	PhysPath string `json:"phys_path,omitempty"` // Physical path
	DevId    string `json:"devid,omitempty"`     // Device ID

	// Scan information
	ScanProcessed string `json:"scan_processed,omitempty"` // Bytes processed during scan

	// Child vdevs for hierarchical structure
	Children []ZfsVdev `json:"children,omitempty" gorm:"foreignKey:ParentId;references:ID"`
}