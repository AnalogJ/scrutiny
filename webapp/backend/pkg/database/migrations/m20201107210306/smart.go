package m20201107210306

import (
	"gorm.io/gorm"
	"time"
)

type Smart struct {
	gorm.Model

	DeviceWWN string `json:"device_wwn"`
	Device    Device `json:"-" gorm:"foreignkey:DeviceWWN"` // use DeviceWWN as foreign key

	TestDate    time.Time `json:"date"`
	SmartStatus string    `json:"smart_status"` // SmartStatusPassed or SmartStatusFailed

	//Metrics
	Temp            int64 `json:"temp"`
	PowerOnHours    int64 `json:"power_on_hours"`
	PowerCycleCount int64 `json:"power_cycle_count"`

	AtaAttributes  []SmartAtaAttribute  `json:"ata_attributes" gorm:"foreignkey:SmartId"`
	NvmeAttributes []SmartNvmeAttribute `json:"nvme_attributes" gorm:"foreignkey:SmartId"`
	ScsiAttributes []SmartScsiAttribute `json:"scsi_attributes" gorm:"foreignkey:SmartId"`
}
