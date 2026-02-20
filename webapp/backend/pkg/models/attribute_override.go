package models

// AttributeOverride defines optional per-attribute handling rules.
// These are host-level overrides that allow you to ignore or force a status
// for specific SMART attributes without modifying collector output.
// Config path: smart.attribute_overrides
// Example:
// smart:
//
//	attribute_overrides:
//	  - protocol: ATA
//	    attribute_id: "187"
//	    wwn: "0x5000c5002df89099"
//	    action: "ignore"
//	  - protocol: NVMe
//	    attribute_id: "media_errors"
//	    action: "force_status"
//	    status: "passed"
//
// Supported actions: "ignore", "force_status" (status: passed|warn|failed)
// All matching fields are optional except protocol/action/attribute_id.
// WWN match is used when provided; otherwise override is applied to all devices.
// SerialNumber is accepted for future expansion but currently unused.
type AttributeOverride struct {
	Protocol     string `json:"protocol" mapstructure:"protocol"`
	AttributeId  string `json:"attribute_id" mapstructure:"attribute_id"`
	WWN          string `json:"wwn,omitempty" mapstructure:"wwn"`
	SerialNumber string `json:"serial_number,omitempty" mapstructure:"serial_number"`
	Action       string `json:"action" mapstructure:"action"`
	Status       string `json:"status,omitempty" mapstructure:"status"`
	// Optional numeric thresholds. If set, the attribute will be marked warn/failed
	// when its value exceeds the specified thresholds (using raw values for ATA, and
	// current Value for NVMe/SCSI).
	WarnAbove int64 `json:"warn_above" mapstructure:"warn_above"`
	FailAbove int64 `json:"fail_above" mapstructure:"fail_above"`

	// Internal flags to distinguish between an explicitly configured zero and an unset value.
	WarnAboveSet bool `json:"-" mapstructure:"-"`
	FailAboveSet bool `json:"-" mapstructure:"-"`
}
