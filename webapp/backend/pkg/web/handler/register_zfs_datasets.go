package handler

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func RegisterZfsDatasets(c *gin.Context) {
	logger := c.MustGet("LOGGER").(logrus.FieldLogger)
	deviceRepo := c.MustGet("DEVICE_REPOSITORY").(database.DeviceRepo)

	var datasets models.ZfsDatasetWrapper
	err := c.ShouldBindJSON(&datasets)
	if err != nil {
		logger.Errorln("Cannot parse ZFS dataset registration request", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "errors": []string{err.Error()}})
		return
	}

	logger.Infof("Registering %d ZFS datasets", len(datasets.Data))

	// Register the datasets in the database
	err = deviceRepo.RegisterZfsDatasets(c, datasets.Data)
	if err != nil {
		logger.Errorln("An error occurred while registering ZFS datasets", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "errors": []string{err.Error()}})
		return
	}

	// Return the registered datasets
	c.JSON(http.StatusOK, models.ZfsDatasetWrapper{
		Success: true,
		Data:    datasets.Data,
	})
}