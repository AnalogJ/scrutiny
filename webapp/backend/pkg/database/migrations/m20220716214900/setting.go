package m20220716214900

import (
	"gorm.io/gorm"
)

type Setting struct {
	//GORM attributes, see: http://gorm.io/docs/conventions.html
	gorm.Model

	SettingKeyName        string `json:"setting_key_name"`
	SettingKeyDescription string `json:"setting_key_description"`
	SettingDataType       string `json:"setting_data_type"`

	SettingValueNumeric int    `json:"setting_value_numeric"`
	SettingValueString  string `json:"setting_value_string"`
	SettingValueBool    bool   `json:"setting_value_bool"`
}
