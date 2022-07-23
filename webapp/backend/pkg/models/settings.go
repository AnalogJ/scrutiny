package models

// Settings is made up of parsed SettingEntry objects retrieved from the database
//type Settings struct {
//	MetricsNotifyLevel            pkg.MetricsNotifyLevel            `json:"metrics.notify.level" mapstructure:"metrics.notify.level"`
//	MetricsStatusFilterAttributes pkg.MetricsStatusFilterAttributes `json:"metrics.status.filter_attributes" mapstructure:"metrics.status.filter_attributes"`
//	MetricsStatusThreshold        pkg.MetricsStatusThreshold        `json:"metrics.status.threshold" mapstructure:"metrics.status.threshold"`
//}

type Settings struct {
	Theme            string `json:"theme" mapstructure:"theme"`
	Layout           string `json:"layout" mapstructure:"layout"`
	DashboardDisplay string `json:"dashboardDisplay" mapstructure:"dashboardDisplay"`
	DashboardSort    string `json:"dashboardSort" mapstructure:"dashboardSort"`
	TemperatureUnit  string `json:"temperatureUnit" mapstructure:"temperatureUnit"`

	Metrics struct {
		NotifyLevel            int `json:"notifyLevel" mapstructure:"notifyLevel"`
		StatusFilterAttributes int `json:"statusFilterAttributes" mapstructure:"statusFilterAttributes"`
		StatusThreshold        int `json:"statusThreshold" mapstructure:"statusThreshold"`
	} `json:"metrics" mapstructure:"metrics"`
}
