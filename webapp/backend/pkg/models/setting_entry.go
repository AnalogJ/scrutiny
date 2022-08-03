package models

import (
	"gorm.io/gorm"
)

// SettingEntry matches a setting row in the database
type SettingEntry struct {
	//GORM attributes, see: http://gorm.io/docs/conventions.html
	gorm.Model

	SettingKeyName        string `json:"setting_key_name" gorm:"unique;not null"`
	SettingKeyDescription string `json:"setting_key_description"`
	SettingDataType       string `json:"setting_data_type"`

	SettingValueNumeric int    `json:"setting_value_numeric"`
	SettingValueString  string `json:"setting_value_string"`
	SettingValueBool    bool   `json:"setting_value_bool"`
}

func (s SettingEntry) TableName() string {
	return "settings"
}
