package handler

import (
	"net/http"

	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/analogj/scrutiny/webapp/backend/pkg/notify"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CollectorScanError(c *gin.Context) {
	logger := c.MustGet("LOGGER").(*logrus.Entry)
	appConfig := c.MustGet("CONFIG").(config.Interface)

	var errorPayload collector.CollectorError
	err := c.BindJSON(&errorPayload)
	if err != nil {
		logger.Errorln("Cannot parse collector scan error payload", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	logger.Errorf("Collector scan error reported: %s (host: %s)", errorPayload.Error, errorPayload.HostId)

	device := models.Device{
		HostId:     errorPayload.HostId,
		DeviceName: "collector scan",
	}

	liveNotify := notify.NewError(
		logger,
		appConfig,
		device,
		notify.NotifyFailureTypeCollectorScanError,
		errorPayload.Error,
	)
	_ = liveNotify.Send()

	c.JSON(http.StatusOK, gin.H{"success": true})
}
