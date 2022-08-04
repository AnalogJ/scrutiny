package database

import (
	"context"
	"fmt"
	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/mitchellh/mapstructure"
	"strings"
)

// LoadSettings will retrieve settings from the database, store them in the AppConfig object, and return a Settings struct
func (sr *scrutinyRepository) LoadSettings(ctx context.Context) (*models.Settings, error) {
	settingsEntries := []models.SettingEntry{}
	if err := sr.gormClient.WithContext(ctx).Find(&settingsEntries).Error; err != nil {
		return nil, fmt.Errorf("Could not get settings from DB: %v", err)
	}

	// store retrieved settings in the AppConfig obj
	for _, settingsEntry := range settingsEntries {
		configKey := fmt.Sprintf("%s.%s", config.DB_USER_SETTINGS_SUBKEY, settingsEntry.SettingKeyName)

		if settingsEntry.SettingDataType == "numeric" {
			sr.appConfig.SetDefault(configKey, settingsEntry.SettingValueNumeric)
		} else if settingsEntry.SettingDataType == "string" {
			sr.appConfig.SetDefault(configKey, settingsEntry.SettingValueString)
		} else if settingsEntry.SettingDataType == "bool" {
			sr.appConfig.SetDefault(configKey, settingsEntry.SettingValueBool)
		}
	}

	// unmarshal the dbsetting object data to a settings object.
	var settings models.Settings
	err := sr.appConfig.UnmarshalKey(config.DB_USER_SETTINGS_SUBKEY, &settings)
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

// testing
// curl -d '{"metrics": { "notify_level": 5, "status_filter_attributes": 5, "status_threshold": 5 }}' -H "Content-Type: application/json" -X POST http://localhost:9090/api/settings
// SaveSettings will update settings in AppConfig object, then save the settings to the database.
func (sr *scrutinyRepository) SaveSettings(ctx context.Context, settings models.Settings) error {
	//save the entries to the appconfig
	settingsMap := &map[string]interface{}{}
	err := mapstructure.Decode(settings, &settingsMap)
	if err != nil {
		return err
	}
	settingsWrapperMap := map[string]interface{}{}
	settingsWrapperMap[config.DB_USER_SETTINGS_SUBKEY] = *settingsMap
	err = sr.appConfig.MergeConfigMap(settingsWrapperMap)
	if err != nil {
		return err
	}
	sr.logger.Debugf("after merge settings: %v", sr.appConfig.AllSettings())
	//retrieve current settings from the database
	settingsEntries := []models.SettingEntry{}
	if err := sr.gormClient.WithContext(ctx).Find(&settingsEntries).Error; err != nil {
		return fmt.Errorf("Could not get settings from DB: %v", err)
	}

	//update settingsEntries
	for ndx, settingsEntry := range settingsEntries {
		configKey := fmt.Sprintf("%s.%s", config.DB_USER_SETTINGS_SUBKEY, strings.ToLower(settingsEntry.SettingKeyName))

		if settingsEntry.SettingDataType == "numeric" {
			settingsEntries[ndx].SettingValueNumeric = sr.appConfig.GetInt(configKey)
		} else if settingsEntry.SettingDataType == "string" {
			settingsEntries[ndx].SettingValueString = sr.appConfig.GetString(configKey)
		} else if settingsEntry.SettingDataType == "bool" {
			settingsEntries[ndx].SettingValueBool = sr.appConfig.GetBool(configKey)
		}

		// store in database.
		//TODO: this should be `sr.gormClient.Updates(&settingsEntries).Error`
		err := sr.gormClient.Model(&models.SettingEntry{}).Where([]uint{settingsEntry.ID}).Select("setting_value_numeric", "setting_value_string", "setting_value_bool").Updates(settingsEntries[ndx]).Error
		if err != nil {
			return err
		}

	}
	return nil
}
