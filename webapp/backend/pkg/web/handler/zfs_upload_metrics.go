package handler

import (
	"net/http"

	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// UploadZFSPoolMetrics receives ZFS pool metrics from the collector and saves them
func UploadZFSPoolMetrics(c *gin.Context) {
	deviceRepo := c.MustGet("DEVICE_REPOSITORY").(database.DeviceRepo)
	logger := c.MustGet("LOGGER").(*logrus.Entry)

	guid := c.Param("guid")

	var pool models.ZFSPool
	err := c.BindJSON(&pool)
	if err != nil {
		logger.Errorln("Cannot parse ZFS pool metrics", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	// Ensure the GUID matches the URL parameter
	pool.GUID = guid

	// Update the pool in the database
	if err := deviceRepo.RegisterZFSPool(c, pool); err != nil {
		logger.Errorln("An error occurred while updating ZFS pool", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	// Save metrics to InfluxDB
	if err := deviceRepo.SaveZFSPoolMetrics(c, pool); err != nil {
		logger.Errorln("An error occurred while saving ZFS pool metrics", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
