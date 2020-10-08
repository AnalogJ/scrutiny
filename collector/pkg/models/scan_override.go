package models

type ScanOverride struct {
	Device     string   `mapstructure:"device"`
	DeviceType []string `mapstructure:"type"`
	Ignore     bool     `mapstructure:"ignore"`
}
