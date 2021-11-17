package handler

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func GetDevicesSummaryTempHistory(c *gin.Context) {
	logger := c.MustGet("LOGGER").(logrus.FieldLogger)
	deviceRepo := c.MustGet("DEVICE_REPOSITORY").(database.DeviceRepo)

	durationKey, exists := c.GetQuery("duration_key")
	if !exists {
		durationKey = "week"
	}

	tempHistory, err := deviceRepo.GetSmartTemperatureHistory(c, durationKey)
	if err != nil {
		logger.Errorln("An error occurred while retrieving summary/temp history", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"temp_history": tempHistory,
		},
	})
}
