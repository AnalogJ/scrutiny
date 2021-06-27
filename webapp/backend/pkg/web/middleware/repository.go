package middleware

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RepositoryMiddleware(appConfig config.Interface, globalLogger logrus.FieldLogger) gin.HandlerFunc {

	deviceRepo, err := database.NewScrutinyRepository(appConfig, globalLogger)
	if err != nil {
		panic(err)
	}

	//TODO: determine where we can call defer deviceRepo.Close()
	return func(c *gin.Context) {
		c.Set("DEVICE_REPOSITORY", deviceRepo)
		c.Next()
	}
}
