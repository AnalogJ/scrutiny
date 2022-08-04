package handler

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func SaveSettings(c *gin.Context) {
	logger := c.MustGet("LOGGER").(*logrus.Entry)
	deviceRepo := c.MustGet("DEVICE_REPOSITORY").(database.DeviceRepo)

	var settings models.Settings
	err := c.BindJSON(&settings)
	if err != nil {
		logger.Errorln("Cannot parse updated settings", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	err = deviceRepo.SaveSettings(c, settings)
	if err != nil {
		logger.Errorln("An error occurred while saving settings", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"settings": settings,
	})
}
