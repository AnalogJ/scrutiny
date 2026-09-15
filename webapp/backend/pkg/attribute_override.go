package pkg

import "strings"

// AttributeOverride replaces Scrutiny's own analysis for a single SMART attribute on a
// single device, letting a user re-tune an attribute whose observed failure rates do not
// match their hardware.
//
// Thresholds are pointers so that an explicit 0 is distinguishable from "not set" -- for
// most attributes 0 is a meaningful threshold.
//
// A threshold is crossed when the value is strictly worse than it -- matching the
// existing NVMe/SCSI threshold convention. Whether "worse" means larger or smaller depends
// on the attribute's Ideal (low/high) in the thresholds metadata, so the direction is
// resolved per attribute rather than assumed. Without that, an override on an ideal-high
// attribute such as NVMe available_spare would compare backwards.
type AttributeOverride struct {
	Protocol      string `mapstructure:"protocol"`
	AttributeID   string `mapstructure:"attribute_id"`
	WarnThreshold *int64 `mapstructure:"warn_threshold"`
	FailThreshold *int64 `mapstructure:"fail_threshold"`
}

// AttributeOverrideSet is the set of overrides configured for one device.
type AttributeOverrideSet []AttributeOverride

// Find returns the first override matching the given protocol and attribute id, or nil.
func (s AttributeOverrideSet) Find(protocol string, attributeID string) *AttributeOverride {
	for ndx := range s {
		if !strings.EqualFold(s[ndx].Protocol, protocol) {
			continue
		}
		if s[ndx].AttributeID != attributeID {
			continue
		}
		return &s[ndx]
	}
	return nil
}
