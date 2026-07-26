package handler

import (
	"net/http"

	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// register devices that are detected by various collectors.
// This function is run everytime a collector is about to start a run. It can be used to update device metadata.
func RegisterDevices(c *gin.Context) {
	deviceRepo := c.MustGet("DEVICE_REPOSITORY").(database.DeviceRepo)
	logger := c.MustGet("LOGGER").(*logrus.Entry)

	var collectorDeviceWrapper models.DeviceWrapper
	err := c.BindJSON(&collectorDeviceWrapper)
	if err != nil {
		logger.Errorln("Cannot parse detected devices", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	// Ignore any device without a Scrutiny UUID. This should never happen...
	detectedStorageDevices := make([]models.Device, 0, len(collectorDeviceWrapper.Data))
	for _, dev := range collectorDeviceWrapper.Data {
		if dev.ScrutinyUUID.IsNil() {
			logger.Errorf("Device %s has no scrutiny UUID; skipping registration (no data association possible).", dev.DeviceName)
			continue
		}
		detectedStorageDevices = append(detectedStorageDevices, dev)
	}

	errs := []error{}
	for _, dev := range detectedStorageDevices {
		//insert devices into DB (and update specified columns if device is already registered)
		// update device fields that may change: (DeviceType, HostID)
		if err := deviceRepo.RegisterDevice(c, dev); err != nil {
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
		c.JSON(http.StatusOK, models.DeviceWrapper{
			Success: true,
			Data:    detectedStorageDevices,
		})
		return
	}
}
