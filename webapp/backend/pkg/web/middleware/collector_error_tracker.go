package middleware

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/notify"
	"github.com/gin-gonic/gin"
)

func CollectorErrorTrackerMiddleware(tracker *notify.CollectorErrorTracker) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("COLLECTOR_ERROR_TRACKER", tracker)
		c.Next()
	}
}
