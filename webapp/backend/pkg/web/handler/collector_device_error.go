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

func CollectorDeviceError(c *gin.Context) {
	logger := c.MustGet("LOGGER").(*logrus.Entry)
	appConfig := c.MustGet("CONFIG").(config.Interface)
	deviceRepo := c.MustGet("DEVICE_REPOSITORY").(database.DeviceRepo)

	scrutinyUuid, err := uuid.FromString(c.Param("scrutiny_uuid"))
	if err != nil {
		logger.Errorln("Invalid scrutiny uuid", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	var errorPayload collector.CollectorError
	err = c.BindJSON(&errorPayload)
	if err != nil {
		logger.Errorln("Cannot parse collector device error payload", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	logger.Errorf("Collector device error reported for %s: %s", scrutinyUuid, errorPayload.Error)

	// Try to look up the device in the database for full info
	device, err := deviceRepo.GetDeviceDetails(c, scrutinyUuid)
	if err != nil {
		// Device not found — build a minimal device with what we have
		logger.Warnf("Could not find device %s in database, sending notification with limited info", scrutinyUuid)
		device = models.Device{
			HostId:     errorPayload.HostId,
			DeviceName: fmt.Sprintf("unknown device (%s)", scrutinyUuid),
		}
	}

	liveNotify := notify.NewError(
		logger,
		appConfig,
		device,
		notify.NotifyFailureTypeCollectorDeviceError,
		errorPayload.Error,
	)
	_ = liveNotify.Send()

	c.JSON(http.StatusOK, gin.H{"success": true})
}
