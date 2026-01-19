package web

import (
	"context"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/analogj/go-util/utils"
	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	"github.com/analogj/scrutiny/webapp/backend/pkg/errors"
	"github.com/analogj/scrutiny/webapp/backend/pkg/metrics"
	"github.com/analogj/scrutiny/webapp/backend/pkg/web/handler"
	"github.com/analogj/scrutiny/webapp/backend/pkg/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AppEngine struct {
	Config           config.Interface
	Logger           *logrus.Entry
	MetricsCollector *metrics.Collector
}

func (ae *AppEngine) Setup(logger *logrus.Entry) *gin.Engine {
	// Register additional MIME types for proper file serving
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".mjs", "application/javascript")
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".woff", "font/woff")
	mime.AddExtensionType(".woff2", "font/woff2")
	mime.AddExtensionType(".ttf", "font/ttf")
	mime.AddExtensionType(".eot", "application/vnd.ms-fontobject")
	mime.AddExtensionType(".otf", "font/otf")
	mime.AddExtensionType(".svg", "image/svg+xml")
	mime.AddExtensionType(".json", "application/json")
	
	r := gin.New()

	r.Use(middleware.LoggerMiddleware(logger))
	r.Use(middleware.RepositoryMiddleware(ae.Config, logger))
	r.Use(middleware.ConfigMiddleware(ae.Config))

	// Initialize metrics collector if enabled
	if ae.Config.GetBool("web.metrics.enabled") {
		if ae.MetricsCollector == nil {
			ae.MetricsCollector = metrics.NewCollector(logger)
		}
		r.Use(middleware.MetricsMiddleware(ae.MetricsCollector))
		logger.Info("Prometheus metrics endpoint enabled")
	} else {
		logger.Info("Prometheus metrics endpoint disabled")
	}
	r.Use(gin.Recovery())

	basePath := ae.Config.GetString("web.listen.basepath")
	logger.Debugf("basepath: %s", basePath)

	base := r.Group(basePath)
	{
		api := base.Group("/api")
		{
			api.GET("/health", handler.HealthCheck)
			api.HEAD("/health", handler.HealthCheck)
			api.POST("/health/notify", handler.SendTestNotification) //check if notifications are configured correctly

			api.POST("/devices/register", handler.RegisterDevices)         //used by Collector to register new devices and retrieve filtered list
			api.GET("/summary", handler.GetDevicesSummary)                 //used by Dashboard
			api.GET("/summary/temp", handler.GetDevicesSummaryTempHistory) //used by Dashboard (Temperature history dropdown)

			// Prometheus metrics endpoint (only registered if enabled)
			if ae.Config.GetBool("web.metrics.enabled") {
				api.GET("/metrics", handler.GetMetrics)
			}

			api.POST("/device/:wwn/smart", handler.UploadDeviceMetrics) //used by Collector to upload data
			api.POST("/device/:wwn/selftest", handler.UploadDeviceSelfTests)
			api.GET("/device/:wwn/details", handler.GetDeviceDetails)   //used by Details
			api.POST("/device/:wwn/archive", handler.ArchiveDevice)     //used by UI to archive device
			api.POST("/device/:wwn/unarchive", handler.UnarchiveDevice) //used by UI to unarchive device
			api.POST("/device/:wwn/mute", handler.MuteDevice)           //used by UI to mute device
			api.POST("/device/:wwn/unmute", handler.UnmuteDevice)       //used by UI to unmute device
			api.POST("/device/:wwn/label", handler.UpdateDeviceLabel)   //used by UI to set device label
			api.DELETE("/device/:wwn", handler.DeleteDevice)            //used by UI to delete device

			api.GET("/settings", handler.GetSettings)   //used to get settings
			api.POST("/settings", handler.SaveSettings) //used to save settings

			// ZFS Pool API endpoints
			zfs := api.Group("/zfs")
			{
				zfs.POST("/pools/register", handler.RegisterZFSPools)        //used by ZFS Collector to register pools
				zfs.GET("/summary", handler.GetZFSPoolsSummary)              //used by ZFS Dashboard
				zfs.POST("/pool/:guid/metrics", handler.UploadZFSPoolMetrics) //used by ZFS Collector to upload metrics
				zfs.GET("/pool/:guid/details", handler.GetZFSPoolDetails)    //used by ZFS Pool Details view
				zfs.POST("/pool/:guid/archive", handler.ArchiveZFSPool)      //used by UI to archive pool
				zfs.POST("/pool/:guid/unarchive", handler.UnarchiveZFSPool)  //used by UI to unarchive pool
				zfs.POST("/pool/:guid/mute", handler.MuteZFSPool)            //used by UI to mute pool
				zfs.POST("/pool/:guid/unmute", handler.UnmuteZFSPool)        //used by UI to unmute pool
				zfs.POST("/pool/:guid/label", handler.UpdateZFSPoolLabel)    //used by UI to set pool label
				zfs.DELETE("/pool/:guid", handler.DeleteZFSPool)             //used by UI to delete pool
			}
		}
	}

	//Static request routing
	// Determine the actual frontend path - check if browser/ subdirectory exists
	frontendPath := ae.Config.GetString("web.src.frontend.path")
	browserPath := filepath.Join(frontendPath, "browser")
	indexPath := filepath.Join(browserPath, "index.html")
	
	// Use browser subdirectory if it exists, otherwise use the configured path directly
	actualFrontendPath := frontendPath
	if utils.FileExists(indexPath) {
		actualFrontendPath = browserPath
		logger.Debugf("Serving frontend from browser subdirectory: %s", actualFrontendPath)
	} else {
		logger.Debugf("Serving frontend from configured path: %s", actualFrontendPath)
	}

	// Create file server - it will automatically use the MIME types registered globally above
	fileServer := http.FileServer(http.Dir(actualFrontendPath))
	
	// Serve static files with proper MIME types and SPA routing support
	base.GET("/web", func(c *gin.Context) {
		c.File(filepath.Join(actualFrontendPath, "index.html"))
	})
	
	base.GET("/web/*filepath", func(c *gin.Context) {
		file := c.Param("filepath")
		if file == "" || file == "/" {
			c.File(filepath.Join(actualFrontendPath, "index.html"))
			return
		}
		
		// Remove leading slash if present
		if strings.HasPrefix(file, "/") {
			file = file[1:]
		}
		
		// Check if file exists
		fullPath := filepath.Join(actualFrontendPath, file)
		if !utils.FileExists(fullPath) {
			// For SPA routing, serve index.html for non-existent files
			c.File(filepath.Join(actualFrontendPath, "index.html"))
			return
		}
		
		// Serve the file using the file server
		// MIME type will be automatically set based on registered types above
		c.Request.URL.Path = "/" + file
		fileServer.ServeHTTP(c.Writer, c.Request)
	})

	//redirect base url to /web
	base.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, basePath+"/web")
	})

	//catch-all, serve index page for any unmatched routes
	r.NoRoute(func(c *gin.Context) {
		c.File(filepath.Join(actualFrontendPath, "index.html"))
	})
	return r
}

func (ae *AppEngine) Start() error {
	//set the gin mode
	gin.SetMode(gin.ReleaseMode)
	if strings.ToLower(ae.Config.GetString("log.level")) == "debug" {
		gin.SetMode(gin.DebugMode)
	}

	//check if the database parent directory exists, fail here rather than in a handler.
	if !utils.FileExists(filepath.Dir(ae.Config.GetString("web.database.location"))) {
		return errors.ConfigValidationError(fmt.Sprintf(
			"Database parent directory does not exist. Please check path (%s)",
			filepath.Dir(ae.Config.GetString("web.database.location"))))
	}

	r := ae.Setup(ae.Logger)

	// Load initial metrics data asynchronously at startup (if metrics enabled)
	if ae.Config.GetBool("web.metrics.enabled") && ae.MetricsCollector != nil {
		go func() {
			deviceRepo, err := database.NewScrutinyRepository(ae.Config, ae.Logger)
			if err != nil {
				ae.Logger.Errorln("Failed to create repository for loading metrics:", err)
				return
			}
			defer deviceRepo.Close()

			if err := ae.MetricsCollector.LoadInitialData(deviceRepo, context.Background()); err != nil {
				ae.Logger.Errorln("Failed to load initial metrics data:", err)
			}
		}()
	}

	return r.Run(fmt.Sprintf("%s:%s", ae.Config.GetString("web.listen.host"), ae.Config.GetString("web.listen.port")))
}
