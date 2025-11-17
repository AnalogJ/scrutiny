package metrics

import (
	"strconv"
	"strings"

	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
)

// SanitizeMetricName converts a string to a valid Prometheus metric name
// Example: converts "attr.5.raw_value" to "scrutiny_smart_attr_5_raw_value"
func SanitizeMetricName(name string) string {
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ToLower(name)

	if strings.HasPrefix(name, "attr_") {
		return "scrutiny_smart_" + name
	}
	return name
}

// TryParseFloat attempts to convert any type to float64
// Supports: int, int64, float32, float64, string, hexadecimal strings
func TryParseFloat(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, true
		}
		// Try parsing hexadecimal
		if strings.HasPrefix(v, "0x") || strings.HasPrefix(v, "0X") {
			if i, err := strconv.ParseInt(v, 0, 64); err == nil {
				return float64(i), true
			}
		}
	}
	return 0, false
}

// SelectLatestSmartResult selects the latest SMART result from a list (by timestamp)
func SelectLatestSmartResult(smartResults []measurements.Smart) *measurements.Smart {
	if len(smartResults) == 0 {
		return nil
	}

	latest := &smartResults[0]
	for i := 1; i < len(smartResults); i++ {
		if smartResults[i].Date.After(latest.Date) {
			latest = &smartResults[i]
		}
	}
	return latest
}
