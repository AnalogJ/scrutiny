package handler

import (
	dbModels "github.com/analogj/scrutiny/webapp/backend/pkg/models/db"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"net/http"
)

// filter devices that are detected by various collectors.
func RegisterDevices(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	logger := c.MustGet("LOGGER").(logrus.FieldLogger)

	var collectorDeviceWrapper dbModels.DeviceWrapper
	err := c.BindJSON(&collectorDeviceWrapper)
	if err != nil {
		logger.Errorln("Cannot parse detected devices", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	//TODO: filter devices here (remove excludes, force includes)
	errs := []error{}
	for _, dev := range collectorDeviceWrapper.Data {
		//insert devices into DB if not already there.
		if err := db.Where(dbModels.Device{WWN: dev.WWN}).FirstOrCreate(&dev).Error; err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		logger.Errorln("An error occurred while registering devices", errs)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
		})
		return
	} else {
		c.JSON(http.StatusOK, dbModels.DeviceWrapper{
			Success: true,
			Data:    collectorDeviceWrapper.Data,
		})
		return
	}
}
