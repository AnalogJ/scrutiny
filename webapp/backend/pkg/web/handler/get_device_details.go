package handler

import (
	"fmt"
	"net/http"

	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	"github.com/analogj/scrutiny/webapp/backend/pkg/thresholds"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetDeviceDetails(c *gin.Context) {
	logger := c.MustGet("LOGGER").(*logrus.Entry)
	deviceRepo := c.MustGet("DEVICE_REPOSITORY").(database.DeviceRepo)

	device, err := deviceRepo.GetDeviceDetails(c, c.Param("wwn"))
	if err != nil {
		logger.Errorln("An error occurred while retrieving device details", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	durationKey, exists := c.GetQuery("duration_key")
	if !exists {
		durationKey = "forever"
	}

	smartResults, err := deviceRepo.GetSmartAttributeHistory(c, c.Param("wwn"), durationKey, 0, 0, nil)
	if err != nil {
		logger.Errorln("An error occurred while retrieving device smart results", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	var deviceMetadata interface{}
	if device.IsAta() {
		// Merge standard ATA SMART attribute metadata with ATA Device Statistics metadata
		// Device statistics (like devstat_7_8 for Percentage Used Endurance Indicator)
		// are critical for enterprise SSD monitoring
		mergedMetadata := make(map[string]interface{})
		for k, v := range thresholds.AtaMetadata {
			mergedMetadata[fmt.Sprintf("%d", k)] = v
		}
		for k, v := range thresholds.AtaDeviceStatsMetadata {
			mergedMetadata[k] = v
		}
		deviceMetadata = mergedMetadata
	} else if device.IsNvme() {
		deviceMetadata = thresholds.NmveMetadata
	} else if device.IsScsi() {
		deviceMetadata = thresholds.ScsiMetadata
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": map[string]interface{}{"device": device, "smart_results": smartResults}, "metadata": deviceMetadata})
}
