package web

import (
	"fmt"
	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	"github.com/analogj/scrutiny/webapp/backend/pkg/metadata"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	dbModels "github.com/analogj/scrutiny/webapp/backend/pkg/models/db"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type AppEngine struct {
	Config config.Interface
}

func (ae *AppEngine) Start() error {
	r := gin.Default()

	r.Use(database.DatabaseHandler(ae.Config.GetString("web.database.location")))

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
			})
		})

		//TODO: notifications
		api.GET("/devices", GetDevicesHandler)
		api.GET("/summary", GetDevicesSummary)
		api.POST("/device/:wwn/smart", UploadDeviceSmartData)
		api.POST("/device/:wwn/selftest", UploadDeviceSelfTestData)

		api.GET("/device/:wwn/details", GetDeviceDetails)
	}

	//Static request routing
	r.StaticFS("/web", http.Dir(ae.Config.GetString("web.src.frontend.path")))

	//redirect base url to /web
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/web")
	})

	//catch-all, serve index page.
	r.NoRoute(func(c *gin.Context) {
		c.File(fmt.Sprintf("%s/index.html", ae.Config.GetString("web.src.frontend.path")))
	})

	return r.Run(fmt.Sprintf("%s:%s", ae.Config.GetString("web.listen.host"), ae.Config.GetString("web.listen.port"))) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// Get all active disks for processing by collectors
func GetDevicesHandler(c *gin.Context) {
	storageDevices, err := RetrieveStorageDevices()

	db := c.MustGet("DB").(*gorm.DB)
	for _, dev := range storageDevices {
		//insert devices into DB if not already there.
		db.Where(dbModels.Device{WWN: dev.WWN}).FirstOrCreate(&dev)
	}

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    storageDevices,
		})
	}
}

func UploadDeviceSmartData(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)

	var collectorSmartData collector.SmartInfo
	err := c.BindJSON(&collectorSmartData)
	if err != nil {
		//TODO: cannot parse smart data
		log.Error("Cannot parse SMART data")
		c.JSON(http.StatusOK, gin.H{"success": false})

	}

	//update the device information if necessary
	var device dbModels.Device
	db.Where("wwn = ?", c.Param("wwn")).First(&device)
	device.UpdateFromCollectorSmartInfo(collectorSmartData)
	db.Model(&device).Updates(device)

	// insert smart info
	deviceSmartData := dbModels.Smart{}
	err = deviceSmartData.FromCollectorSmartInfo(c.Param("wwn"), collectorSmartData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}
	db.Create(&deviceSmartData)

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func UploadDeviceSelfTestData(c *gin.Context) {

}

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

func GetDevicesSummary(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)

	devices := []dbModels.Device{}

	//OLD: cant seem to figure out how to get the latest SmartResults for each Device, so instead
	// we're going to assume that results were retrieved at the same time, so we'll just get the last x number of results
	//var devicesCount int
	//db.Table("devices").Count(&devicesCount)

	//We need the last x (for now all) Smart objects for each Device, so that we can graph Temperature
	//We also need the last
	db.Debug().
		Preload("SmartResults", func(db *gorm.DB) *gorm.DB {
			return db.Order("smarts.created_at DESC") //OLD: .Limit(devicesCount)
		}).
		//Preload("SmartResults").
		// Preload("SmartResults.SmartAttributes").
		Find(&devices)

	//for _, dev := range devices {
	//	log.Printf("===== device: %s\n", dev.WWN)
	//	log.Print(len(dev.SmartResults))
	//}
	//a, _ := json.Marshal(devices) //get json byte array
	//n := len(a)   //Find the length of the byte array
	//s := string(a[:n]) //convert to string
	//log.Print(s) //write to response

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    devices,
	})
	//c.Data(http.StatusOK, "application/json", a)
}
