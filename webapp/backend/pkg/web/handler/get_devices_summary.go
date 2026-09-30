package handler

import (
	"net/http"

	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetDevicesSummary(c *gin.Context) {
	logger := c.MustGet("LOGGER").(*logrus.Entry)
	appConfig := c.MustGet("CONFIG").(config.Interface)
	deviceRepo := c.MustGet("DEVICE_REPOSITORY").(database.DeviceRepo)

	summary, err := deviceRepo.GetSummary(c)
	if err != nil {
		logger.Errorln("An error occurred while retrieving device summary", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	friendlyNames := config.GetFriendlyNames(appConfig)
	for scrutinyUUID, deviceSummary := range summary {
		deviceSummary.Device.FriendlyName = friendlyNames[scrutinyUUID]
	}

	//this must match DeviceSummaryWrapper (webapp/backend/pkg/models/device_summary.go)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"summary": summary,
			//"temperature": tem
		},
	})
}
