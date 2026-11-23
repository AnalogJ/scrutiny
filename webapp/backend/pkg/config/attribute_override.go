package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/errors"
	"github.com/analogj/scrutiny/webapp/backend/pkg/thresholds"
)

// validateAttributeOverrides rejects overrides that could never take effect. A silently ignored
// override is indistinguishable from a working one, so these fail loudly at startup instead.
func validateAttributeOverrides(device DeviceConfig) error {
	for ndx, override := range device.AttributeOverrides {
		location := fmt.Sprintf("`attribute_overrides[%d]` of device %s", ndx, device.ScrutinyUUID)

		if !isSupportedProtocol(override.Protocol) {
			return errors.ConfigValidationError(fmt.Sprintf(
				"%s has an unsupported `protocol` (%q). Must be one of %s, %s, %s",
				location, override.Protocol, pkg.DeviceProtocolAta, pkg.DeviceProtocolNvme, pkg.DeviceProtocolScsi))
		}
		if len(strings.TrimSpace(override.AttributeID)) == 0 {
			return errors.ConfigValidationError(fmt.Sprintf("%s is missing an `attribute_id`", location))
		}
		if strings.EqualFold(override.Protocol, pkg.DeviceProtocolAta) {
			if _, err := strconv.Atoi(override.AttributeID); err != nil {
				return errors.ConfigValidationError(fmt.Sprintf("%s has an `attribute_id` (%q) that isn't an ATA attribute number", location, override.AttributeID))
			}
		}
		if override.WarnThreshold == nil && override.FailThreshold == nil {
			return errors.ConfigValidationError(fmt.Sprintf(
				"%s must set `warn_threshold` and/or `fail_threshold`, otherwise it would have no effect", location))
		}
		if override.WarnThreshold != nil && override.FailThreshold != nil {
			warn, fail := *override.WarnThreshold, *override.FailThreshold
			largerIsWorse := thresholds.LargerIsWorse(override.Protocol, override.AttributeID)
			if (largerIsWorse && warn > fail) || (!largerIsWorse && warn < fail) {
				return errors.ConfigValidationError(fmt.Sprintf("%s has a `warn_threshold` (%d) worse than its `fail_threshold` (%d), so it could never warn", location, warn, fail))
			}
		}
	}
	return nil
}

func isSupportedProtocol(protocol string) bool {
	for _, supported := range []string{pkg.DeviceProtocolAta, pkg.DeviceProtocolNvme, pkg.DeviceProtocolScsi} {
		if strings.EqualFold(protocol, supported) {
			return true
		}
	}
	return false
}
