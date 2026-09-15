package measurements

import (
	"fmt"

	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/thresholds"
)

// evaluateAttributeOverride resolves a user-configured override into an attribute status.
//
// It fully replaces Scrutiny's own analysis for the attribute: when an override matches, the
// observed-failure-rate buckets (ATA) or the recommended-threshold comparison (NVMe/SCSI) are
// not consulted at all. Without that, an attribute could never be brought back to "passed" --
// critical attributes fall through to a warning whenever no bucket matches.
//
// idealLow reports whether a larger value is worse for this attribute.
func evaluateAttributeOverride(value int64, idealLow bool, override *pkg.AttributeOverride) (pkg.AttributeStatus, string) {
	worseThan := func(threshold int64) bool {
		if idealLow {
			return value > threshold
		}
		return value < threshold
	}
	direction := "above"
	if !idealLow {
		direction = "below"
	}

	if override.FailThreshold != nil && worseThan(*override.FailThreshold) {
		return pkg.AttributeStatusFailedScrutiny | pkg.AttributeStatusFailedOverride,
			fmt.Sprintf("Attribute is %s the configured fail threshold (%d)", direction, *override.FailThreshold)
	}
	if override.WarnThreshold != nil && worseThan(*override.WarnThreshold) {
		return pkg.AttributeStatusWarningScrutiny | pkg.AttributeStatusWarningOverride,
			fmt.Sprintf("Attribute is %s the configured warn threshold (%d)", direction, *override.WarnThreshold)
	}
	return pkg.AttributeStatusPassedOverride, "Attribute is within the configured override thresholds"
}

// idealIsLow reports whether a larger value is worse. Attributes with no metadata, or with no
// stated ideal, are treated as ideal-low: nearly every SMART attribute is an error counter.
func idealIsLow(ideal string) bool {
	return ideal != thresholds.ObservedThresholdIdealHigh
}
