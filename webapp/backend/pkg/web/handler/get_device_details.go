package handler

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	"github.com/analogj/scrutiny/webapp/backend/pkg/thresholds"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func GetDeviceDetails(c *gin.Context) {
	logger := c.MustGet("LOGGER").(logrus.FieldLogger)
	deviceRepo := c.MustGet("DEVICE_REPOSITORY").(database.DeviceRepo)

	device, err := deviceRepo.GetDeviceDetails(c, c.Param("wwn"))
	if err != nil {
		logger.Errorln("An error occurred while retrieving device details", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	smartResults, err := deviceRepo.GetSmartAttributeHistory(c, c.Param("wwn"), "", nil)

	var deviceMetadata interface{}
	if device.IsAta() {
		deviceMetadata = thresholds.AtaMetadata
	} else if device.IsNvme() {
		deviceMetadata = thresholds.NmveMetadata
	} else if device.IsScsi() {
		deviceMetadata = thresholds.ScsiMetadata
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": map[string]interface{}{"device": device, "smart_results": smartResults}, "metadata": deviceMetadata})
}
