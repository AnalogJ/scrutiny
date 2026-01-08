package handler

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/analogj/go-util/utils"
	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func HealthCheck(c *gin.Context) {
	logger := c.MustGet("LOGGER").(*logrus.Entry)
	deviceRepo := c.MustGet("DEVICE_REPOSITORY").(database.DeviceRepo)
	appConfig := c.MustGet("CONFIG").(config.Interface)
	logger.Infof("Checking Influxdb & Sqlite health")

	//check sqlite and influxdb health
	err := deviceRepo.HealthCheck(c)
	if err != nil {
		logger.Errorln("An error occurred during healthcheck", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	// check if the /web folder is populated with expected frontend files
	frontendPath := appConfig.GetString("web.src.frontend.path")
	indexPath := filepath.Join(frontendPath, "index.html")
	if !utils.FileExists(indexPath) {
		errMsg := fmt.Sprintf("Frontend files not found. Expected index.html at: %s", indexPath)
		logger.Errorln(errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
