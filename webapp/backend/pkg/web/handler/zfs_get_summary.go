package handler

import (
	"net/http"

	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetZFSPoolsSummary returns a summary of all ZFS pools for the dashboard
func GetZFSPoolsSummary(c *gin.Context) {
	deviceRepo := c.MustGet("DEVICE_REPOSITORY").(database.DeviceRepo)
	logger := c.MustGet("LOGGER").(*logrus.Entry)

	summary, err := deviceRepo.GetZFSPoolsSummary(c)
	if err != nil {
		logger.Errorln("An error occurred while getting ZFS pools summary", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"pools": summary,
		},
	})
}
