package handler

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func GetZfsPools(c *gin.Context) {
	logger := c.MustGet("LOGGER").(logrus.FieldLogger)
	deviceRepo := c.MustGet("DEVICE_REPOSITORY").(database.DeviceRepo)

	// Get optional host filter
	hostId := c.Query("host_id")

	var pools []models.ZfsPool
	var err error

	if hostId != "" {
		pools, err = deviceRepo.GetZfsPoolsByHost(c, hostId)
	} else {
		pools, err = deviceRepo.GetZfsPools(c)
	}

	if err != nil {
		logger.Errorln("An error occurred while retrieving ZFS pools", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, models.ZfsPoolWrapper{
		Success: true,
		Data:    pools,
	})
}

func GetZfsPoolDetails(c *gin.Context) {
	logger := c.MustGet("LOGGER").(logrus.FieldLogger)
	deviceRepo := c.MustGet("DEVICE_REPOSITORY").(database.DeviceRepo)

	poolGuid := c.Param("poolGuid")
	if poolGuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "errors": []string{"Pool GUID is required"}})
		return
	}

	pool, err := deviceRepo.GetZfsPoolByGuid(c, poolGuid)
	if err != nil {
		logger.Errorln("An error occurred while retrieving ZFS pool details", err)
		c.JSON(http.StatusNotFound, gin.H{"success": false, "errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    pool,
	})
}