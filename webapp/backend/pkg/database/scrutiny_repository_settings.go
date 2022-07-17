package database

import (
	"context"
	"fmt"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
)

func (sr *scrutinyRepository) GetSettings(ctx context.Context) (*models.Settings, error) {
	settingsEntries := []models.SettingEntry{}
	if err := sr.gormClient.WithContext(ctx).Find(&settingsEntries).Error; err != nil {
		return nil, fmt.Errorf("Could not get settings from DB: %v", err)
	}

	settings := models.Settings{}
	settings.PopulateFromSettingEntries(settingsEntries)

	return &settings, nil
}
func (sr *scrutinyRepository) SaveSettings(ctx context.Context, settings models.Settings) error {

	//get current settings
	settingsEntries := []models.SettingEntry{}
	if err := sr.gormClient.WithContext(ctx).Find(&settingsEntries).Error; err != nil {
		return fmt.Errorf("Could not get settings from DB: %v", err)
	}

	// override with values from settings object
	settingsEntries = settings.UpdateSettingEntries(settingsEntries)

	// store in database.
	return sr.gormClient.Updates(&settingsEntries).Error
}
