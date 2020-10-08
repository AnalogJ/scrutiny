package handler

import (
	dbModels "github.com/analogj/scrutiny/webapp/backend/pkg/models/db"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

func GetDevicesSummary(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	logger := c.MustGet("LOGGER").(logrus.FieldLogger)
	devices := []dbModels.Device{}

	//We need the last x (for now all) Smart objects for each Device, so that we can graph Temperature
	//We also need the last
	if err := db.Preload("SmartResults", func(db *gorm.DB) *gorm.DB {
		return db.Order("smarts.created_at DESC") //OLD: .Limit(devicesCount)
	}).
		Find(&devices).Error; err != nil {
		logger.Errorln("Could not get device summary from DB", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    devices,
	})
}
