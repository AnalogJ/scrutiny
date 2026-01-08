package handler

import (
	"net/http"

	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetZFSPoolDetails returns detailed information about a specific ZFS pool
func GetZFSPoolDetails(c *gin.Context) {
	deviceRepo := c.MustGet("DEVICE_REPOSITORY").(database.DeviceRepo)
	logger := c.MustGet("LOGGER").(*logrus.Entry)

	guid := c.Param("guid")

	// Get the pool details with vdev hierarchy
	pool, err := deviceRepo.GetZFSPoolDetails(c, guid)
	if err != nil {
		logger.Errorln("An error occurred while getting ZFS pool details", err)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Pool not found",
		})
		return
	}

	// Get metrics history (default to week)
	durationKey := c.DefaultQuery("duration_key", "week")
	metricsHistory, err := deviceRepo.GetZFSPoolMetricsHistory(c, guid, durationKey)
	if err != nil {
		logger.Warnln("Could not get ZFS pool metrics history", err)
		// Continue without metrics history
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"pool":            pool,
			"metrics_history": metricsHistory,
		},
	})
}
