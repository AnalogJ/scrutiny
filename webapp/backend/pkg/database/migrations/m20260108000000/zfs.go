package m20260108000000

import (
	"time"
)

// ZFSPoolStatus represents the health status of a ZFS pool
type ZFSPoolStatus string

// ZFSScrubState represents the state of a ZFS scrub operation
type ZFSScrubState string

// ZFSVdevType represents the type of a ZFS vdev
type ZFSVdevType string

// ZFSPool represents a ZFS storage pool for migration
type ZFSPool struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	GUID   string `gorm:"primary_key"`
	Name   string
	HostID string

	Status ZFSPoolStatus
	Health string

	Size            int64
	Allocated       int64
	Free            int64
	Fragmentation   int
	CapacityPercent float64

	Ashift int

	ScrubState           ZFSScrubState
	ScrubStartTime       *time.Time
	ScrubEndTime         *time.Time
	ScrubScannedBytes    int64
	ScrubIssuedBytes     int64
	ScrubTotalBytes      int64
	ScrubErrorsCount     int64
	ScrubPercentComplete float64

	TotalReadErrors     int64
	TotalWriteErrors    int64
	TotalChecksumErrors int64

	Label    string
	Archived bool
	Muted    bool
}

// ZFSVdev represents a virtual device in a ZFS pool for migration
type ZFSVdev struct {
	ID       uint `gorm:"primary_key;autoIncrement"`
	PoolGUID string `gorm:"index;not null"`
	ParentID *uint  `gorm:"index"`

	Name   string
	Type   ZFSVdevType
	Status ZFSPoolStatus
	GUID   string
	Path   string

	ReadErrors     int64
	WriteErrors    int64
	ChecksumErrors int64

	Size      int64
	Allocated int64
}
