package handler

import (
	dbModels "github.com/analogj/scrutiny/webapp/backend/pkg/models/db"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// filter devices that are detected by various collectors.
func RegisterDevices(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)

	var collectorDeviceWrapper dbModels.DeviceWrapper
	err := c.BindJSON(&collectorDeviceWrapper)
	if err != nil {
		log.Error("Cannot parse detected devices")
		c.JSON(http.StatusOK, gin.H{"success": false})
	}

	//TODO: filter devices here (remove excludes, force includes)

	for _, dev := range collectorDeviceWrapper.Data {
		//insert devices into DB if not already there.
		db.Where(dbModels.Device{WWN: dev.WWN}).FirstOrCreate(&dev)
	}

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
		})
	} else {
		c.JSON(http.StatusOK, dbModels.DeviceWrapper{
			Success: true,
			Data:    collectorDeviceWrapper.Data,
		})
	}
}
