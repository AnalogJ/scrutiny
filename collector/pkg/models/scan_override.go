package models

type ScanOverride struct {
	Device     string   `mapstructure:"device"`
	DeviceType []string `mapstructure:"type"`
	Ignore     bool     `mapstructure:"ignore"`
	// nil means "not set for this device", so the top level setting applies
	NotifyOnSmartctlError *bool `mapstructure:"notify_on_smartctl_error"`
	Commands              struct {
		MetricsInfoArgs  string `mapstructure:"metrics_info_args"`
		MetricsSmartArgs string `mapstructure:"metrics_smart_args"`
	} `mapstructure:"commands"`
}
