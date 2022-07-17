package models

import "github.com/analogj/scrutiny/webapp/backend/pkg"

// Settings is made up of parsed SettingEntry objects retrieved from the database
type Settings struct {
	MetricsNotifyLevel            pkg.MetricsNotifyLevel            `json:"metrics_notify_level"`
	MetricsStatusFilterAttributes pkg.MetricsStatusFilterAttributes `json:"metrics_status_filter_attributes"`
	MetricsStatusThreshold        pkg.MetricsStatusThreshold        `json:"metrics_status_threshold"`
}

func (s *Settings) PopulateFromSettingEntries(entries []SettingEntry) {
	for _, entry := range entries {
		if entry.SettingKeyName == "metrics.notify.level" {
			s.MetricsNotifyLevel = pkg.MetricsNotifyLevel(entry.SettingValueNumeric)
		} else if entry.SettingKeyName == "metrics.status.filter_attributes" {
			s.MetricsStatusFilterAttributes = pkg.MetricsStatusFilterAttributes(entry.SettingValueNumeric)
		} else if entry.SettingKeyName == "metrics.status.threshold" {
			s.MetricsStatusThreshold = pkg.MetricsStatusThreshold(entry.SettingValueNumeric)
		}
	}
}

func (s *Settings) UpdateSettingEntries(entries []SettingEntry) []SettingEntry {
	for _, entry := range entries {
		if entry.SettingKeyName == "metrics.notify.level" {
			entry.SettingValueNumeric = int64(s.MetricsNotifyLevel)
		} else if entry.SettingKeyName == "metrics.status.filter_attributes" {
			entry.SettingValueNumeric = int64(s.MetricsStatusFilterAttributes)
		} else if entry.SettingKeyName == "metrics.status.threshold" {
			entry.SettingValueNumeric = int64(s.MetricsStatusThreshold)
		}
	}
	return entries
}
