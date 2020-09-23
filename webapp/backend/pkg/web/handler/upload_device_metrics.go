package handler

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	dbModels "github.com/analogj/scrutiny/webapp/backend/pkg/models/db"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"net/http"
)

func UploadDeviceMetrics(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	logger := c.MustGet("LOGGER").(logrus.FieldLogger)

	var collectorSmartData collector.SmartInfo
	err := c.BindJSON(&collectorSmartData)
	if err != nil {
		logger.Errorln("Cannot parse SMART data", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	//update the device information if necessary
	var device dbModels.Device
	db.Where("wwn = ?", c.Param("wwn")).First(&device)
	device.UpdateFromCollectorSmartInfo(collectorSmartData)
	if err := db.Model(&device).Updates(device).Error; err != nil {
		logger.Errorln("An error occurred while updating device data from smartctl metrics", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	// insert smart info
	deviceSmartData := dbModels.Smart{}
	err = deviceSmartData.FromCollectorSmartInfo(c.Param("wwn"), collectorSmartData)
	if err != nil {
		logger.Errorln("Could not process SMART metrics", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}
	if err := db.Create(&deviceSmartData).Error; err != nil {
		logger.Errorln("An error occurred while saving smartctl metrics", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
