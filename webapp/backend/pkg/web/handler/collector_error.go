package handler

import (
	"net/http"

	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/analogj/scrutiny/webapp/backend/pkg/notify"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/sirupsen/logrus"
)

// CollectorError is used by the collector to report that it could not gather data.
func CollectorError(c *gin.Context) {
	logger := c.MustGet("LOGGER").(*logrus.Entry)
	appConfig := c.MustGet("CONFIG").(config.Interface)
	deviceRepo := c.MustGet("DEVICE_REPOSITORY").(database.DeviceRepo)

	var errorPayload collector.CollectorError
	if err := c.BindJSON(&errorPayload); err != nil {
		logger.Errorln("Cannot parse collector error", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false})
		return
	}

	logger.Errorf("Collector error reported for %q: %s", errorPayload.DeviceName, errorPayload.Error)

	device := models.Device{
		HostId:     errorPayload.HostId,
		DeviceName: errorPayload.DeviceName,
		DeviceType: errorPayload.DeviceType,
	}
	// a device that failed `smartctl --info` was never registered, so details are only available
	// when the collector managed to send a UUID
	if scrutinyUUID, err := uuid.FromString(errorPayload.ScrutinyUUID); err == nil && !scrutinyUUID.IsNil() {
		if storedDevice, err := deviceRepo.GetDeviceDetails(c, scrutinyUUID); err == nil {
			device = storedDevice
		} else {
			logger.Warnf("Could not look up device %s: %v", scrutinyUUID, err)
		}
	}

	collectorNotify := notify.NewCollectorError(logger, appConfig, device, errorPayload.Error)
	_ = collectorNotify.Send() //we ignore error message when sending notifications.

	c.JSON(http.StatusOK, models.DeviceWrapper{Success: true})
}
