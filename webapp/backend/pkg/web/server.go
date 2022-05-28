package web

import (
	"fmt"
	"github.com/analogj/go-util/utils"
	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	"github.com/analogj/scrutiny/webapp/backend/pkg/errors"
	"github.com/analogj/scrutiny/webapp/backend/pkg/web/handler"
	"github.com/analogj/scrutiny/webapp/backend/pkg/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type AppEngine struct {
	Config config.Interface
}

func (ae *AppEngine) Setup(logger logrus.FieldLogger) *gin.Engine {
	r := gin.New()

	r.Use(middleware.LoggerMiddleware(logger))
	r.Use(middleware.RepositoryMiddleware(ae.Config, logger))
	r.Use(middleware.ConfigMiddleware(ae.Config))
	r.Use(gin.Recovery())

	basePath := ae.Config.GetString("web.listen.basepath")
	logger.Debugf("basepath: %s", basePath)

	base := r.Group(basePath)
	{
		api := base.Group("/api")
		{
			api.GET("/health", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"success": true,
				})
			})
			api.POST("/health/notify", handler.SendTestNotification) //check if notifications are configured correctly

			api.POST("/devices/register", handler.RegisterDevices)         //used by Collector to register new devices and retrieve filtered list
			api.GET("/summary", handler.GetDevicesSummary)                 //used by Dashboard
			api.GET("/summary/temp", handler.GetDevicesSummaryTempHistory) //used by Dashboard (Temperature history dropdown)
			api.POST("/device/:wwn/smart", handler.UploadDeviceMetrics)    //used by Collector to upload data
			api.POST("/device/:wwn/selftest", handler.UploadDeviceSelfTests)
			api.GET("/device/:wwn/details", handler.GetDeviceDetails) //used by Details
			api.DELETE("/device/:wwn", handler.DeleteDevice)          //used by UI to delete device

		}
	}

	//Static request routing
	base.StaticFS("/web", http.Dir(ae.Config.GetString("web.src.frontend.path")))

	//redirect base url to /web
	base.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, basePath+"/web")
	})

	//catch-all, serve index page.
	r.NoRoute(func(c *gin.Context) {
		c.File(fmt.Sprintf("%s/index.html", ae.Config.GetString("web.src.frontend.path")))
	})
	return r
}

func (ae *AppEngine) Start() error {
	//set the gin mode
	gin.SetMode(gin.ReleaseMode)
	if strings.ToLower(ae.Config.GetString("log.level")) == "debug" {
		gin.SetMode(gin.DebugMode)
	}

	logger := logrus.New()
	//set default log level
	logLevel, err := logrus.ParseLevel(ae.Config.GetString("log.level"))
	if err != nil {
		return err
	}
	logger.SetLevel(logLevel)
	//set the log file if present
	if len(ae.Config.GetString("log.file")) != 0 {
		logFile, err := os.OpenFile(ae.Config.GetString("log.file"), os.O_CREATE|os.O_WRONLY, 0644)
		defer logFile.Close()
		if err != nil {
			logrus.Errorf("Failed to open log file %s for output: %s", ae.Config.GetString("log.file"), err)
			return err
		}

		//configure the logrus default
		logger.SetOutput(io.MultiWriter(os.Stderr, logFile))
	}

	//check if the database parent directory exists, fail here rather than in a handler.
	if !utils.FileExists(filepath.Dir(ae.Config.GetString("web.database.location"))) {
		return errors.ConfigValidationError(fmt.Sprintf(
			"Database parent directory does not exist. Please check path (%s)",
			filepath.Dir(ae.Config.GetString("web.database.location"))))
	}

	r := ae.Setup(logger)

	return r.Run(fmt.Sprintf("%s:%s", ae.Config.GetString("web.listen.host"), ae.Config.GetString("web.listen.port")))
}
