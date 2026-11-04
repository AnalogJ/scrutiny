package measurements

import (
	"fmt"

	"github.com/analogj/scrutiny/webapp/backend/pkg"
)

// applyAttributeOverride resolves a user-configured override into an attribute status.
//
// It fully replaces Scrutiny's own analysis for the attribute: when an override matches, the
// observed-failure-rate buckets (ATA) or the recommended-threshold comparison (NVMe/SCSI) are
// not consulted at all. Without that, an attribute could never be brought back to "passed" --
// critical attributes fall through to a warning whenever no bucket matches.
func applyAttributeOverride(status *pkg.AttributeStatus, statusReason *string, value int64, largerIsWorse bool, override *pkg.AttributeOverride) {
	worseThan := func(threshold int64) bool {
		if largerIsWorse {
			return value > threshold
		}
		return value < threshold
	}
	direction := "above"
	if !largerIsWorse {
		direction = "below"
	}

	overrideStatus := pkg.AttributeStatusPassedOverride
	reason := "Attribute is within the configured override thresholds"
	if override.FailThreshold != nil && worseThan(*override.FailThreshold) {
		overrideStatus = pkg.AttributeStatusFailedScrutiny | pkg.AttributeStatusFailedOverride
		reason = fmt.Sprintf("Attribute is %s the configured fail threshold (%d)", direction, *override.FailThreshold)
	} else if override.WarnThreshold != nil && worseThan(*override.WarnThreshold) {
		overrideStatus = pkg.AttributeStatusWarningScrutiny | pkg.AttributeStatusWarningOverride
		reason = fmt.Sprintf("Attribute is %s the configured warn threshold (%d)", direction, *override.WarnThreshold)
	}

	*status = pkg.AttributeStatusSet(*status, overrideStatus)
	if *statusReason != "" {
		*statusReason += "; "
	}
	*statusReason += reason
}
