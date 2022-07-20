package models

// Settings is made up of parsed SettingEntry objects retrieved from the database
//type Settings struct {
//	MetricsNotifyLevel            pkg.MetricsNotifyLevel            `json:"metrics.notify.level" mapstructure:"metrics.notify.level"`
//	MetricsStatusFilterAttributes pkg.MetricsStatusFilterAttributes `json:"metrics.status.filter_attributes" mapstructure:"metrics.status.filter_attributes"`
//	MetricsStatusThreshold        pkg.MetricsStatusThreshold        `json:"metrics.status.threshold" mapstructure:"metrics.status.threshold"`
//}

type Settings struct {
	Metrics struct {
		Notify struct {
			Level int `json:"level" mapstructure:"level"`
		} `json:"notify" mapstructure:"notify"`
		Status struct {
			FilterAttributes int `json:"filter_attributes" mapstructure:"filter_attributes"`
			Threshold        int `json:"threshold" mapstructure:"threshold"`
		} `json:"status" mapstructure:"status"`
	} `json:"metrics" mapstructure:"metrics"`
}
