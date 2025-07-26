package handler

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func RegisterZfsPools(c *gin.Context) {
	logger := c.MustGet("LOGGER").(logrus.FieldLogger)
	deviceRepo := c.MustGet("DEVICE_REPOSITORY").(database.DeviceRepo)

	var pools models.ZfsPoolWrapper
	err := c.ShouldBindJSON(&pools)
	if err != nil {
		logger.Errorln("Cannot parse ZFS pool registration request", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "errors": []string{err.Error()}})
		return
	}

	logger.Infof("Registering %d ZFS pools", len(pools.Data))

	// Register the pools in the database
	err = deviceRepo.RegisterZfsPools(c, pools.Data)
	if err != nil {
		logger.Errorln("An error occurred while registering ZFS pools", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "errors": []string{err.Error()}})
		return
	}

	// Return the registered pools
	c.JSON(http.StatusOK, models.ZfsPoolWrapper{
		Success: true,
		Data:    pools.Data,
	})
}