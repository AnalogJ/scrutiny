package middleware

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	"github.com/gin-gonic/gin"
)

func ConfigMiddleware(appConfig config.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("CONFIG", appConfig)
		c.Next()
	}
}
