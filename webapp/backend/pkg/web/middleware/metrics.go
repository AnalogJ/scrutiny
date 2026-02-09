package middleware

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/metrics"
	"github.com/gin-gonic/gin"
)

// MetricsMiddleware injects metrics collector into gin context
func MetricsMiddleware(collector *metrics.Collector) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("METRICS_COLLECTOR", collector)
		c.Next()
	}
}
