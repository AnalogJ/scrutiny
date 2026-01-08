package handler

import (
	"net/http"

	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ArchiveZFSPool archives a ZFS pool (hides it from the dashboard)
func ArchiveZFSPool(c *gin.Context) {
	deviceRepo := c.MustGet("DEVICE_REPOSITORY").(database.DeviceRepo)
	logger := c.MustGet("LOGGER").(*logrus.Entry)

	guid := c.Param("guid")

	err := deviceRepo.UpdateZFSPoolArchived(c, guid, true)
	if err != nil {
		logger.Errorln("An error occurred while archiving ZFS pool", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// UnarchiveZFSPool unarchives a ZFS pool (shows it on the dashboard)
func UnarchiveZFSPool(c *gin.Context) {
	deviceRepo := c.MustGet("DEVICE_REPOSITORY").(database.DeviceRepo)
	logger := c.MustGet("LOGGER").(*logrus.Entry)

	guid := c.Param("guid")

	err := deviceRepo.UpdateZFSPoolArchived(c, guid, false)
	if err != nil {
		logger.Errorln("An error occurred while unarchiving ZFS pool", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// MuteZFSPool mutes notifications for a ZFS pool
func MuteZFSPool(c *gin.Context) {
	deviceRepo := c.MustGet("DEVICE_REPOSITORY").(database.DeviceRepo)
	logger := c.MustGet("LOGGER").(*logrus.Entry)

	guid := c.Param("guid")

	err := deviceRepo.UpdateZFSPoolMuted(c, guid, true)
	if err != nil {
		logger.Errorln("An error occurred while muting ZFS pool", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// UnmuteZFSPool unmutes notifications for a ZFS pool
func UnmuteZFSPool(c *gin.Context) {
	deviceRepo := c.MustGet("DEVICE_REPOSITORY").(database.DeviceRepo)
	logger := c.MustGet("LOGGER").(*logrus.Entry)

	guid := c.Param("guid")

	err := deviceRepo.UpdateZFSPoolMuted(c, guid, false)
	if err != nil {
		logger.Errorln("An error occurred while unmuting ZFS pool", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// UpdateZFSPoolLabel updates the custom label for a ZFS pool
func UpdateZFSPoolLabel(c *gin.Context) {
	deviceRepo := c.MustGet("DEVICE_REPOSITORY").(database.DeviceRepo)
	logger := c.MustGet("LOGGER").(*logrus.Entry)

	guid := c.Param("guid")

	var payload struct {
		Label string `json:"label"`
	}
	if err := c.BindJSON(&payload); err != nil {
		logger.Errorln("Cannot parse label payload", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false})
		return
	}

	err := deviceRepo.UpdateZFSPoolLabel(c, guid, payload.Label)
	if err != nil {
		logger.Errorln("An error occurred while updating ZFS pool label", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DeleteZFSPool deletes a ZFS pool from tracking
func DeleteZFSPool(c *gin.Context) {
	deviceRepo := c.MustGet("DEVICE_REPOSITORY").(database.DeviceRepo)
	logger := c.MustGet("LOGGER").(*logrus.Entry)

	guid := c.Param("guid")

	err := deviceRepo.DeleteZFSPool(c, guid)
	if err != nil {
		logger.Errorln("An error occurred while deleting ZFS pool", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
