package handler

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/metadata"
	dbModels "github.com/analogj/scrutiny/webapp/backend/pkg/models/db"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

func GetDeviceDetails(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	logger := c.MustGet("LOGGER").(logrus.FieldLogger)
	device := dbModels.Device{}

	if err := db.Preload("SmartResults", func(db *gorm.DB) *gorm.DB {
		return db.Order("smarts.created_at DESC").Limit(40)
	}).
		Preload("SmartResults.AtaAttributes").
		Preload("SmartResults.NvmeAttributes").
		Preload("SmartResults.ScsiAttributes").
		Where("wwn = ?", c.Param("wwn")).
		First(&device).Error; err != nil {

		logger.Errorln("An error occurred while retrieving device details", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	if err := device.SquashHistory(); err != nil {
		logger.Errorln("An error occurred while squashing device history", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	if err := device.ApplyMetadataRules(); err != nil {
		logger.Errorln("An error occurred while applying scrutiny thresholds & rules", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

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
