package handler

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func HealthCheck(c *gin.Context) {
	logger := c.MustGet("LOGGER").(*logrus.Entry)
	deviceRepo := c.MustGet("DEVICE_REPOSITORY").(database.DeviceRepo)
	logger.Infof("Checking Influxdb & Sqlite health")

	//check sqlite and influxdb health
	err := deviceRepo.HealthCheck(c)
	if err != nil {
		logger.Errorln("An error occurred during healthcheck", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	//TODO:
	// check if the /web folder is populated.

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
