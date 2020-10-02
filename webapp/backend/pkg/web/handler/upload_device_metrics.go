package handler

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	dbModels "github.com/analogj/scrutiny/webapp/backend/pkg/models/db"
	"github.com/analogj/scrutiny/webapp/backend/pkg/notify"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"net/http"
)

func UploadDeviceMetrics(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	logger := c.MustGet("LOGGER").(logrus.FieldLogger)
	appConfig := c.MustGet("CONFIG").(config.Interface)

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

	//check for error
	if deviceSmartData.SmartStatus == dbModels.SmartStatusFailed {
		//send notifications
		testNotify := notify.Notify{
			Config: appConfig,
			Payload: notify.Payload{
				FailureType:  notify.NotifyFailureTypeSmartFailure,
				DeviceName:   device.DeviceName,
				DeviceType:   device.DeviceProtocol,
				DeviceSerial: device.SerialNumber,
				Test:         false,
			},
		}
		_ = testNotify.Send() //we ignore error message when sending notifications.
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
