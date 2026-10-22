package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SetDeviceFriendlyName(c *gin.Context) {
	logger := c.MustGet("LOGGER").(*logrus.Entry)
	deviceRepo := c.MustGet("DEVICE_REPOSITORY").(database.DeviceRepo)

	scrutiny_uuid, err := uuid.FromString(c.Param("scrutiny_uuid"))
	if err != nil {
		logger.Errorln("Invalid scrutiny uuid", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false})
		return
	}

	var body struct {
		FriendlyName string `json:"friendly_name"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		logger.Errorln("Invalid friendly name request body", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false})
		return
	}
	friendlyName := strings.TrimSpace(body.FriendlyName)

	if _, err := deviceRepo.GetDeviceDetails(c, scrutiny_uuid); errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"success": false})
		return
	} else if err != nil {
		logger.Errorln("An error occurred while retrieving device", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	if err := deviceRepo.UpdateDeviceFriendlyName(c, scrutiny_uuid, friendlyName); err != nil {
		logger.Errorln("An error occurred while updating device friendly name", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"friendly_name": friendlyName}})
}
