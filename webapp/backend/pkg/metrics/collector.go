package metrics

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
	metricsModels "github.com/analogj/scrutiny/webapp/backend/pkg/models/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/sirupsen/logrus"
)

// Collector manages Prometheus metrics for all devices
type Collector struct {
	mu       sync.RWMutex
	devices  map[string]*metricsModels.DeviceMetricsData // key: wwn
	registry *prometheus.Registry
	logger   *logrus.Entry
}

// NewCollector creates a new metrics collector
func NewCollector(logger *logrus.Entry) *Collector {
	mc := &Collector{
		devices:  make(map[string]*metricsModels.DeviceMetricsData),
		registry: prometheus.NewRegistry(),
		logger:   logger,
	}

	// Register Go runtime metrics (memory, GC, goroutines, etc.)
	mc.registry.MustRegister(collectors.NewGoCollector())

	// Register custom device metrics collector
	mc.registry.MustRegister(mc)
	return mc
}

// UpdateDeviceMetrics updates device metrics (called from UploadDeviceMetrics)
func (mc *Collector) UpdateDeviceMetrics(wwn string, device models.Device, smartData measurements.Smart) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.devices[wwn] = &metricsModels.DeviceMetricsData{
		Device:    device,
		SmartData: smartData,
		UpdatedAt: time.Now(),
	}
	mc.logger.Debugf("Updated metrics for device %s", wwn)
}

// LoadInitialData loads initial data from database (called at startup)
func (mc *Collector) LoadInitialData(deviceRepo database.DeviceRepo, ctx context.Context) error {
	start := time.Now()
	mc.logger.Info("Loading initial metrics data from database...")

	// Get device summary
	summary, err := deviceRepo.GetSummary(ctx)
	if err != nil {
		return fmt.Errorf("failed to load device summary: %w", err)
	}

	// Concurrently fetch latest SMART data for each device
	smartDataMap := make(map[string][]measurements.Smart)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for wwn := range summary {
		wg.Add(1)
		go func(w string) {
			defer wg.Done()
			smarts, err := deviceRepo.GetSmartAttributeHistory(ctx, w, "forever", 1, 0, nil)
			if err == nil && len(smarts) > 0 {
				mu.Lock()
				smartDataMap[w] = smarts
				mu.Unlock()
			}
		}(wwn)
	}

	wg.Wait()

	// Load into memory
	mc.mu.Lock()
	defer mc.mu.Unlock()

	for wwn, deviceSummary := range summary {
		if smartResults, ok := smartDataMap[wwn]; ok && len(smartResults) > 0 {
			mc.devices[wwn] = &metricsModels.DeviceMetricsData{
				Device:    deviceSummary.Device,
				SmartData: smartResults[0],
				UpdatedAt: time.Now(),
			}
		}
	}

	mc.logger.Infof("Loaded metrics for %d devices in %v", len(mc.devices), time.Since(start))
	return nil
}

// GetRegistry returns the Prometheus registry
func (mc *Collector) GetRegistry() *prometheus.Registry {
	return mc.registry
}

// Describe implements prometheus.Collector interface
func (mc *Collector) Describe(ch chan<- *prometheus.Desc) {
	// Dynamic metrics, no need to pre-describe
}

// Collect implements prometheus.Collector interface
func (mc *Collector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	mc.collectDeviceInfo(ch)
	mc.collectDeviceCapacity(ch)
	mc.collectDeviceStatus(ch)
	mc.collectSmartAttributes(ch)
	mc.collectSummaryMetrics(ch)
	mc.collectStatistics(ch)

	mc.logger.Debugf("Metrics collected in %v for %d devices", time.Since(start), len(mc.devices))
}

// collectDeviceInfo generates device information metrics
func (mc *Collector) collectDeviceInfo(ch chan<- prometheus.Metric) {
	for wwn, data := range mc.devices {
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc("scrutiny_device_info", "Device information",
				[]string{"wwn", "device_name", "model_name", "serial_number",
					"firmware", "protocol", "host_id", "form_factor"}, nil),
			prometheus.GaugeValue, 1,
			wwn, data.Device.DeviceName, data.Device.ModelName,
			data.Device.SerialNumber, data.Device.Firmware,
			data.Device.DeviceProtocol, data.Device.HostId, data.Device.FormFactor,
		)
	}
}

// collectDeviceCapacity generates device capacity metrics
func (mc *Collector) collectDeviceCapacity(ch chan<- prometheus.Metric) {
	for wwn, data := range mc.devices {
		if data.Device.Capacity > 0 {
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc("scrutiny_device_capacity_bytes", "Device capacity in bytes",
					[]string{"wwn", "device_name", "model_name", "protocol", "host_id"}, nil),
				prometheus.GaugeValue, float64(data.Device.Capacity),
				wwn, data.Device.DeviceName, data.Device.ModelName,
				data.Device.DeviceProtocol, data.Device.HostId,
			)
		}
	}
}

// collectDeviceStatus generates device status metrics
func (mc *Collector) collectDeviceStatus(ch chan<- prometheus.Metric) {
	for wwn, data := range mc.devices {
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc("scrutiny_device_status", "Device status (0=passed, 1=failed)",
				[]string{"wwn", "device_name", "model_name", "protocol", "host_id"}, nil),
			prometheus.GaugeValue, float64(data.Device.DeviceStatus),
			wwn, data.Device.DeviceName, data.Device.ModelName,
			data.Device.DeviceProtocol, data.Device.HostId,
		)
	}
}

// collectSmartAttributes generates SMART attribute metrics
func (mc *Collector) collectSmartAttributes(ch chan<- prometheus.Metric) {
	for wwn, data := range mc.devices {
		baseLabels := []string{wwn, data.Device.DeviceName, data.Device.ModelName,
			data.Device.DeviceProtocol, data.Device.HostId}

		for attrID, attr := range data.SmartData.Attributes {
			attrLabels := append(baseLabels, attrID)
			flattenedAttrs := attr.Flatten()

			for key, value := range flattenedAttrs {
				metricName := SanitizeMetricName(key)
				if floatVal, ok := TryParseFloat(value); ok {
					ch <- prometheus.MustNewConstMetric(
						prometheus.NewDesc(metricName, fmt.Sprintf("SMART attribute %s", key),
							[]string{"wwn", "device_name", "model_name", "protocol", "host_id", "attribute_id"}, nil),
						prometheus.GaugeValue, floatVal, attrLabels...,
					)
				}
			}
		}
	}
}

// collectSummaryMetrics generates summary metrics
func (mc *Collector) collectSummaryMetrics(ch chan<- prometheus.Metric) {
	for wwn, data := range mc.devices {
		labels := []string{wwn, data.Device.DeviceName, data.Device.ModelName,
			data.Device.DeviceProtocol, data.Device.HostId}

		if data.SmartData.Temp > 0 {
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc("scrutiny_smart_temperature_celsius",
					"Device temperature in Celsius",
					[]string{"wwn", "device_name", "model_name", "protocol", "host_id"}, nil),
				prometheus.GaugeValue, float64(data.SmartData.Temp), labels...,
			)
		}

		if data.SmartData.PowerOnHours > 0 {
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc("scrutiny_smart_power_on_hours", "Device power on hours",
					[]string{"wwn", "device_name", "model_name", "protocol", "host_id"}, nil),
				prometheus.GaugeValue, float64(data.SmartData.PowerOnHours), labels...,
			)
		}

		if data.SmartData.PowerCycleCount > 0 {
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc("scrutiny_smart_power_cycle_count", "Device power cycle count",
					[]string{"wwn", "device_name", "model_name", "protocol", "host_id"}, nil),
				prometheus.GaugeValue, float64(data.SmartData.PowerCycleCount), labels...,
			)
		}

		timestampMs := float64(data.SmartData.Date.Unix() * 1000)
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc("scrutiny_smart_collector_timestamp",
				"Timestamp of last data collection",
				[]string{"wwn", "device_name", "model_name", "protocol", "host_id"}, nil),
			prometheus.GaugeValue, timestampMs, labels...,
		)
	}
}

// collectStatistics generates statistics metrics
func (mc *Collector) collectStatistics(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc("scrutiny_devices_total", "Total number of monitored devices", nil, nil),
		prometheus.GaugeValue, float64(len(mc.devices)),
	)

	protocolCount := make(map[string]int)
	for _, data := range mc.devices {
		protocol := data.Device.DeviceProtocol
		if protocol == "" {
			protocol = "unknown"
		}
		protocolCount[protocol]++
	}

	for protocol, count := range protocolCount {
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc("scrutiny_devices_by_protocol", "Number of devices by protocol",
				[]string{"protocol"}, nil),
			prometheus.GaugeValue, float64(count), protocol,
		)
	}
}
