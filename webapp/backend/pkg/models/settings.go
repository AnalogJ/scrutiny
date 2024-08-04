package models

// Settings is made up of parsed SettingEntry objects retrieved from the database
//type Settings struct {
//	MetricsNotifyLevel            pkg.MetricsNotifyLevel            `json:"metrics.notify.level" mapstructure:"metrics.notify.level"`
//	MetricsStatusFilterAttributes pkg.MetricsStatusFilterAttributes `json:"metrics.status.filter_attributes" mapstructure:"metrics.status.filter_attributes"`
//	MetricsStatusThreshold        pkg.MetricsStatusThreshold        `json:"metrics.status.threshold" mapstructure:"metrics.status.threshold"`
//}

type Settings struct {
	Theme              string `json:"theme" mapstructure:"theme"`
	Layout             string `json:"layout" mapstructure:"layout"`
	DashboardDisplay   string `json:"dashboard_display" mapstructure:"dashboard_display"`
	DashboardSort      string `json:"dashboard_sort" mapstructure:"dashboard_sort"`
	TemperatureUnit    string `json:"temperature_unit" mapstructure:"temperature_unit"`
	FileSizeSIUnits    bool   `json:"file_size_si_units" mapstructure:"file_size_si_units"`
	LineStroke         string `json:"line_stroke" mapstructure:"line_stroke"`
	PoweredOnHoursUnit string `json:"powered_on_hours_unit" mapstructure:"powered_on_hours_unit"`

	Collector struct {
		RetrieveSCTHistory bool `json:"retrieve_sct_temperature_history" mapstructure:"retrieve_sct_temperature_history"`
	} `json:"collector" mapstructure:"collector"`

	Metrics struct {
		NotifyLevel            int  `json:"notify_level" mapstructure:"notify_level"`
		StatusFilterAttributes int  `json:"status_filter_attributes" mapstructure:"status_filter_attributes"`
		StatusThreshold        int  `json:"status_threshold" mapstructure:"status_threshold"`
		RepeatNotifications    bool `json:"repeat_notifications" mapstructure:"repeat_notifications"`
	} `json:"metrics" mapstructure:"metrics"`
}
