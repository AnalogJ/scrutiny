package handler

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/metadata"
	dbModels "github.com/analogj/scrutiny/webapp/backend/pkg/models/db"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func GetDeviceDetails(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	device := dbModels.Device{}

	db.Debug().
		Preload("SmartResults", func(db *gorm.DB) *gorm.DB {
			return db.Order("smarts.created_at DESC").Limit(40)
		}).
		Preload("SmartResults.SmartAttributes").
		Where("wwn = ?", c.Param("wwn")).
		First(&device)

	device.SquashHistory()
	device.ApplyMetadataRules()

	c.JSON(http.StatusOK, gin.H{"success": true, "data": device, "lookup": metadata.AtaSmartAttributes})
}
