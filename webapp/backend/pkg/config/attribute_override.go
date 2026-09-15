package config

import (
	"fmt"
	"strings"

	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/errors"
)

// deviceOverrideConfig mirrors the `devices` block of scrutiny.yaml:
//
//	devices:
//	  - scrutiny_uuid: 106089ed-2273-54e0-b498-ec4bdfc8ca6c
//	    attribute_overrides:
//	      - protocol: ATA
//	        attribute_id: "199"
//	        warn_threshold: 100
//	        fail_threshold: 200
type deviceOverrideConfig struct {
	ScrutinyUUID       string                   `mapstructure:"scrutiny_uuid"`
	AttributeOverrides pkg.AttributeOverrideSet `mapstructure:"attribute_overrides"`
}

// DeviceAttributeOverrides returns the configured attribute overrides, keyed by lowercased
// scrutiny_uuid. Devices with no overrides are omitted.
func DeviceAttributeOverrides(c Interface) (map[string]pkg.AttributeOverrideSet, error) {
	devices, err := parseDeviceOverrides(c)
	if err != nil {
		return nil, err
	}

	overrides := map[string]pkg.AttributeOverrideSet{}
	for _, device := range devices {
		if len(device.AttributeOverrides) == 0 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(device.ScrutinyUUID))
		overrides[key] = append(overrides[key], device.AttributeOverrides...)
	}
	return overrides, nil
}

func parseDeviceOverrides(c Interface) ([]deviceOverrideConfig, error) {
	if !c.IsSet("devices") {
		return nil, nil
	}
	devices := []deviceOverrideConfig{}
	if err := c.UnmarshalKey("devices", &devices); err != nil {
		return nil, errors.ConfigValidationError(fmt.Sprintf("`devices` configuration could not be parsed: %v", err))
	}
	return devices, nil
}

// validateDeviceOverrides rejects overrides that could never take effect. A silently ignored
// override is indistinguishable from a working one, so these fail loudly at startup instead.
func validateDeviceOverrides(c Interface) error {
	devices, err := parseDeviceOverrides(c)
	if err != nil {
		return err
	}

	for ndx, device := range devices {
		if len(strings.TrimSpace(device.ScrutinyUUID)) == 0 {
			return errors.ConfigValidationError(fmt.Sprintf("`devices[%d]` is missing a `scrutiny_uuid`", ndx))
		}
		for oNdx, override := range device.AttributeOverrides {
			location := fmt.Sprintf("`devices[%d].attribute_overrides[%d]` (device %s)", ndx, oNdx, device.ScrutinyUUID)

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
