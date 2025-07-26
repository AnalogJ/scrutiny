package models

import (
	"time"
)

type ZfsDatasetWrapper struct {
	Success bool         `json:"success"`
	Errors  []error      `json:"errors"`
	Data    []ZfsDataset `json:"data"`
}

type ZfsDataset struct {
	//GORM attributes
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	ID       uint   `json:"id" gorm:"primary_key;autoIncrement"`
	Name     string `json:"name" gorm:"uniqueIndex:idx_dataset_host"`
	HostId   string `json:"host_id" gorm:"uniqueIndex:idx_dataset_host"`

	// Dataset information
	Type      string `json:"type"`       // FILESYSTEM, VOLUME, SNAPSHOT
	Pool      string `json:"pool"`       // Parent pool name
	CreateTxg string `json:"createtxg"`  // Creation transaction group

	// Space information
	Used       string `json:"used"`       // Space used by dataset
	Available  string `json:"available"`  // Space available
	Referenced string `json:"referenced"` // Space referenced by dataset
	Mountpoint string `json:"mountpoint"` // Mount point path
}