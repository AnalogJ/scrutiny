package handler

import (
	dbModels "github.com/analogj/scrutiny/webapp/backend/pkg/models/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/sirupsen/logrus"
	"net/http"
)

// register devices that are detected by various collectors.
// This function is run everytime a collector is about to start a run. It can be used to update device data.
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

	errs := []error{}
	for _, dev := range collectorDeviceWrapper.Data {
		//insert devices into DB (and update specified columns if device is already registered)
		// update device fields that may change: (DeviceType, HostID)
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "wwn"}},
			DoUpdates: clause.AssignmentColumns([]string{"host_id", "device_name"}),
		}).Create(&dev).Error; err != nil {

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
