package measurements

import (
	"fmt"
	"strings"

	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/thresholds"
)

// SmartAtaDeviceStatAttribute represents an ATA Device Statistics attribute
// from GP Log 0x04. Unlike standard SMART attributes which use numeric IDs,
// device statistics use string-based IDs in the format "devstat_<page>_<offset>".
type SmartAtaDeviceStatAttribute struct {
	AttributeId string `json:"attribute_id"` // e.g., "devstat_7_8"
	Value       int64  `json:"value"`
	Threshold   int64  `json:"thresh"`

	TransformedValue int64               `json:"transformed_value"`
	Status           pkg.AttributeStatus `json:"status"`
	StatusReason     string              `json:"status_reason,omitempty"`
	FailureRate      float64             `json:"failure_rate,omitempty"`
}

func (sa *SmartAtaDeviceStatAttribute) GetTransformedValue() int64 {
	return sa.TransformedValue
}

func (sa *SmartAtaDeviceStatAttribute) GetStatus() pkg.AttributeStatus {
	return sa.Status
}

func (sa *SmartAtaDeviceStatAttribute) Flatten() map[string]interface{} {
	return map[string]interface{}{
		fmt.Sprintf("attr.%s.attribute_id", sa.AttributeId): sa.AttributeId,
		fmt.Sprintf("attr.%s.value", sa.AttributeId):        sa.Value,
		fmt.Sprintf("attr.%s.thresh", sa.AttributeId):       sa.Threshold,

		// Generated Data
		fmt.Sprintf("attr.%s.transformed_value", sa.AttributeId): sa.TransformedValue,
		fmt.Sprintf("attr.%s.status", sa.AttributeId):            int64(sa.Status),
		fmt.Sprintf("attr.%s.status_reason", sa.AttributeId):     sa.StatusReason,
		fmt.Sprintf("attr.%s.failure_rate", sa.AttributeId):      sa.FailureRate,
	}
}

func (sa *SmartAtaDeviceStatAttribute) Inflate(key string, val interface{}) {
	if val == nil {
		return
	}

	keyParts := strings.Split(key, ".")

	switch keyParts[2] {
	case "attribute_id":
		sa.AttributeId = val.(string)
	case "value":
		sa.Value = val.(int64)
	case "thresh":
		sa.Threshold = val.(int64)

	// Generated
	case "transformed_value":
		sa.TransformedValue = val.(int64)
	case "status":
		sa.Status = pkg.AttributeStatus(val.(int64))
	case "status_reason":
		sa.StatusReason = val.(string)
	case "failure_rate":
		sa.FailureRate = val.(float64)
	}
}

// PopulateAttributeStatus sets the status based on device statistics metadata.
// Chainable.
func (sa *SmartAtaDeviceStatAttribute) PopulateAttributeStatus() *SmartAtaDeviceStatAttribute {
	// Look up metadata for this device statistic
	if metadata, ok := thresholds.AtaDeviceStatsMetadata[sa.AttributeId]; ok {
		// For percentage-based metrics like devstat_7_8 (Percentage Used Endurance Indicator),
		// check if value exceeds threshold based on ideal direction
		if metadata.Critical {
			// Default threshold for percentage used is 100 (device end of life)
			threshold := int64(100)
			if sa.Threshold > 0 {
				threshold = sa.Threshold
			}

			if metadata.Ideal == thresholds.ObservedThresholdIdealLow && sa.Value >= threshold {
				sa.Status = pkg.AttributeStatusSet(sa.Status, pkg.AttributeStatusFailedScrutiny)
				sa.StatusReason = "Device statistic exceeds recommended threshold"
			}
		}
	}

	// Set transformed value to the raw value for device statistics
	sa.TransformedValue = sa.Value

	return sa
}
