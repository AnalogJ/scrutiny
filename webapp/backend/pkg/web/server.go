package web

import (
	"fmt"
	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	"github.com/analogj/scrutiny/webapp/backend/pkg/web/handler"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AppEngine struct {
	Config config.Interface
}

func (ae *AppEngine) Setup() *gin.Engine {
	r := gin.Default()

	r.Use(database.DatabaseHandler(ae.Config.GetString("web.database.location")))

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
			})
		})

		api.POST("/devices/register", handler.RegisterDevices)
		api.GET("/summary", handler.GetDevicesSummary)
		api.POST("/device/:wwn/smart", handler.UploadDeviceMetrics)
		api.POST("/device/:wwn/selftest", handler.UploadDeviceSelfTests)

		api.GET("/device/:wwn/details", handler.GetDeviceDetails)
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
	return r
}

func (ae *AppEngine) Start() error {
	r := ae.Setup()

	return r.Run(fmt.Sprintf("%s:%s", ae.Config.GetString("web.listen.host"), ae.Config.GetString("web.listen.port")))
}
