package handler

import (
	dbModels "github.com/analogj/scrutiny/webapp/backend/pkg/models/db"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func GetDevicesSummary(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	devices := []dbModels.Device{}

	//We need the last x (for now all) Smart objects for each Device, so that we can graph Temperature
	//We also need the last
	db.Debug().
		Preload("SmartResults", func(db *gorm.DB) *gorm.DB {
			return db.Order("smarts.created_at DESC") //OLD: .Limit(devicesCount)
		}).
		Find(&devices)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    devices,
	})
}
