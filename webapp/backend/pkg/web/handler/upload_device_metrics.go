package handler

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	dbModels "github.com/analogj/scrutiny/webapp/backend/pkg/models/db"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func UploadDeviceMetrics(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)

	var collectorSmartData collector.SmartInfo
	err := c.BindJSON(&collectorSmartData)
	if err != nil {
		//TODO: cannot parse smart data
		log.Error("Cannot parse SMART data")
		c.JSON(http.StatusOK, gin.H{"success": false})

	}

	//update the device information if necessary
	var device dbModels.Device
	db.Where("wwn = ?", c.Param("wwn")).First(&device)
	device.UpdateFromCollectorSmartInfo(collectorSmartData)
	db.Model(&device).Updates(device)

	// insert smart info
	deviceSmartData := dbModels.Smart{}
	err = deviceSmartData.FromCollectorSmartInfo(c.Param("wwn"), collectorSmartData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}
	db.Create(&deviceSmartData)

	c.JSON(http.StatusOK, gin.H{"success": true})
}
