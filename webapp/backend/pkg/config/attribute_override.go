package config

import (
	"fmt"
	"strings"

	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/errors"
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
		if override.WarnThreshold == nil && override.FailThreshold == nil {
			return errors.ConfigValidationError(fmt.Sprintf(
				"%s must set `warn_threshold` and/or `fail_threshold`, otherwise it would have no effect", location))
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
