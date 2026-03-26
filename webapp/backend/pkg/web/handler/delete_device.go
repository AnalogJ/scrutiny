package handler

import (
	"net/http"

	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/sirupsen/logrus"
)

func DeleteDevice(c *gin.Context) {
	logger := c.MustGet("LOGGER").(*logrus.Entry)
	deviceRepo := c.MustGet("DEVICE_REPOSITORY").(database.DeviceRepo)
	scrutiny_uuid, err := uuid.FromString(c.Param("scrutiny_uuid"))
	if err != nil {
		logger.Errorln("Invalid scrutiny uuid", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}
	err = deviceRepo.DeleteDevice(c, scrutiny_uuid)
	if err != nil {
		logger.Errorln("An error occurred while deleting device", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
