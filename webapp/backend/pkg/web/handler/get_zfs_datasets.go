package handler

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func GetZfsDatasets(c *gin.Context) {
	logger := c.MustGet("LOGGER").(logrus.FieldLogger)
	deviceRepo := c.MustGet("DEVICE_REPOSITORY").(database.DeviceRepo)

	// Get optional pool filter
	poolName := c.Query("pool")
	// Get optional host filter
	hostId := c.Query("host_id")

	var datasets []models.ZfsDataset
	var err error

	if poolName != "" {
		datasets, err = deviceRepo.GetZfsDatasetsByPool(c, poolName)
	} else if hostId != "" {
		datasets, err = deviceRepo.GetZfsDatasetsByHost(c, hostId)
	} else {
		datasets, err = deviceRepo.GetZfsDatasets(c)
	}

	if err != nil {
		logger.Errorln("An error occurred while retrieving ZFS datasets", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, models.ZfsDatasetWrapper{
		Success: true,
		Data:    datasets,
	})
}