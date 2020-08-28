package handler

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/metadata"
	dbModels "github.com/analogj/scrutiny/webapp/backend/pkg/models/db"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func GetDeviceDetails(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	device := dbModels.Device{}

	db.Debug().
		Preload("SmartResults", func(db *gorm.DB) *gorm.DB {
			return db.Order("smarts.created_at DESC").Limit(40)
		}).
		Preload("SmartResults.AtaAttributes").
		Preload("SmartResults.NvmeAttributes").
		Preload("SmartResults.ScsiAttributes").
		Where("wwn = ?", c.Param("wwn")).
		First(&device)

	device.SquashHistory()
	device.ApplyMetadataRules()

	var deviceMetadata interface{}
	if device.IsAta() {
		deviceMetadata = metadata.AtaMetadata
	} else if device.IsNvme() {
		deviceMetadata = metadata.NmveMetadata
	} else if device.IsScsi() {
		deviceMetadata = metadata.ScsiMetadata
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": device, "metadata": deviceMetadata})
}
