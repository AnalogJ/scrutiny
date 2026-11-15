package thresholds

import (
	"strconv"
	"strings"

	"github.com/analogj/scrutiny/webapp/backend/pkg"
)

// LargerIsWorse reports whether a larger comparable value is worse for an attribute. Attributes
// with no metadata, or with no stated ideal, default to true: nearly all are error counters.
func LargerIsWorse(protocol string, attributeID string) bool {
	switch {
	case strings.EqualFold(protocol, pkg.DeviceProtocolNvme):
		return NmveMetadata[attributeID].Ideal != ObservedThresholdIdealHigh
	case strings.EqualFold(protocol, pkg.DeviceProtocolScsi):
		return ScsiMetadata[attributeID].Ideal != ObservedThresholdIdealHigh
	default:
		id, _ := strconv.Atoi(attributeID)
		metadata := AtaMetadata[id]
		// normalized values fall as a drive degrades, even when Ideal describes a raw counter
		return metadata.DisplayType != AtaSmartAttributeDisplayTypeNormalized && metadata.Ideal != ObservedThresholdIdealHigh
	}
}
