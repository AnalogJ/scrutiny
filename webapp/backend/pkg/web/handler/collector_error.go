package handler

import (
	"fmt"
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

	notifyCollectorError(c, logger, appConfig, deviceRepo, errorPayload)

	c.JSON(http.StatusOK, models.DeviceWrapper{Success: true})
}

// notifyCollectorError logs a collector failure and, unless the user has opted out, notifies them.
func notifyCollectorError(c *gin.Context, logger *logrus.Entry, appConfig config.Interface, deviceRepo database.DeviceRepo, errorPayload collector.CollectorError) {
	logger.Errorf("Collector error reported for %q: %s", errorPayload.DeviceName, errorPayload.Error)

	if !appConfig.GetBool(fmt.Sprintf("%s.metrics.notify_collector_errors", config.DB_USER_SETTINGS_SUBKEY)) {
		logger.Debugln("Collector error notifications are disabled, skipping notification")
		return
	}

	//a device that is broken rather than failing reports the same error on every collector run, so
	//only notify again once the error changes, unless the user asked for repeat notifications.
	tracker := c.MustGet("COLLECTOR_ERROR_TRACKER").(*notify.CollectorErrorTracker)
	if !tracker.ShouldNotifyCollectorError(
		errorPayload.HostId,
		errorPayload.DeviceName,
		errorPayload.Error,
		appConfig.GetBool(fmt.Sprintf("%s.metrics.repeat_notifications", config.DB_USER_SETTINGS_SUBKEY)),
	) {
		logger.Debugf("Collector error for %q is unchanged since the last notification, skipping", errorPayload.DeviceName)
		return
	}

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
}
