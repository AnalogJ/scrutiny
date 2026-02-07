package handler

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg/metrics"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// GetMetrics handles Prometheus metrics endpoint
func GetMetrics(c *gin.Context) {
	logger := c.MustGet("LOGGER").(*logrus.Entry)
	collector, exists := c.MustGet("METRICS_COLLECTOR").(*metrics.Collector)

	if !exists || collector == nil {
		logger.Errorln("Metrics collector not found in context")
		c.String(500, "Metrics collector not initialized")
		return
	}

	handler := promhttp.HandlerFor(collector.GetRegistry(), promhttp.HandlerOpts{})
	handler.ServeHTTP(c.Writer, c.Request)
}
