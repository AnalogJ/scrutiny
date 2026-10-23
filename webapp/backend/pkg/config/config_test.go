package config

import (
	"strings"
	"testing"
	"time"

	"github.com/analogj/scrutiny/webapp/backend/pkg/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func loadTestConfig(t *testing.T, yaml string) (*configuration, error) {
	testConfig := &configuration{}
	require.NoError(t, testConfig.Init())
	require.NoError(t, testConfig.MergeConfig(strings.NewReader(yaml)))
	return testConfig, testConfig.ValidateConfig()
}

func Test_GetDevices(t *testing.T) {
	testConfig, err := loadTestConfig(t, `
devices:
  - scrutiny_uuid: 9A4D34B5-B2EE-51EF-8506-90EEA09BE417
    friendly_name: Backup Disk
    notifications:
      on_missed_upload: false
      on_smartctl_error: false
      upload_period: 1h30m
      disabled: true
  - scrutiny_uuid: ecfaaf20-d1f6-558b-b33a-3e8db19a6c2c
`)
	require.NoError(t, err)

	devices, err := testConfig.GetDevices()
	require.NoError(t, err)
	require.Len(t, devices, 2)
	require.Equal(t, uuid.FromStringOrNil("9a4d34b5-b2ee-51ef-8506-90eea09be417"), devices[0].ScrutinyUUID)
	require.Equal(t, "Backup Disk", devices[0].FriendlyName)
	require.False(t, *devices[0].Notifications.OnMissedUpload)
	require.False(t, *devices[0].Notifications.OnSmartctlError)
	require.Equal(t, 90*time.Minute, devices[0].Notifications.UploadPeriod)
	require.Equal(t, uuid.FromStringOrNil("ecfaaf20-d1f6-558b-b33a-3e8db19a6c2c"), devices[1].ScrutinyUUID)
	require.Empty(t, devices[1].FriendlyName)
	require.Nil(t, devices[1].Notifications.OnMissedUpload)
	require.Nil(t, devices[1].Notifications.OnSmartctlError)
	require.Zero(t, devices[1].Notifications.UploadPeriod)
	require.True(t, devices[0].Notifications.Disabled)
	require.False(t, devices[1].Notifications.Disabled)

	found := testConfig.GetDeviceConfig(uuid.FromStringOrNil("9a4d34b5-b2ee-51ef-8506-90eea09be417"))
	require.NotNil(t, found)
	require.Equal(t, "Backup Disk", found.FriendlyName)
	require.True(t, found.Notifications.Disabled)
	require.Nil(t, testConfig.GetDeviceConfig(uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001")))

	require.False(t, *found.Notifications.OnSmartctlError)
	other := testConfig.GetDeviceConfig(uuid.FromStringOrNil("ecfaaf20-d1f6-558b-b33a-3e8db19a6c2c"))
	require.NotNil(t, other)
	require.Nil(t, other.Notifications.OnSmartctlError)
}

func Test_NotifyDefaults(t *testing.T) {
	testConfig, err := loadTestConfig(t, ``)
	require.NoError(t, err)
	require.True(t, testConfig.GetBool("notify.on_smartctl_error"))
}

func Test_ValidateConfig_InvalidDevices(t *testing.T) {
	for name, yaml := range map[string]string{
		"invalid uuid":    `devices: [{scrutiny_uuid: not-a-uuid}]`,
		"missing uuid":    `devices: [{notifications: {upload_period: 1h}}]`,
		"invalid period":  `devices: [{scrutiny_uuid: ecfaaf20-d1f6-558b-b33a-3e8db19a6c2c, notifications: {upload_period: daily}}]`,
		"negative period": `devices: [{scrutiny_uuid: ecfaaf20-d1f6-558b-b33a-3e8db19a6c2c, notifications: {upload_period: -1h}}]`,
	} {
		t.Run(name, func(t *testing.T) {
			_, err := loadTestConfig(t, yaml)
			require.IsType(t, errors.ConfigValidationError(""), err)
		})
	}
}

func Test_MergeConfigMap(t *testing.T) {
	//setup
	testConfig := configuration{
		Viper: viper.New(),
	}
	testConfig.Set("user.dashboard_display", "hello")
	testConfig.SetDefault("user.layout", "hello")

	mergeSettings := map[string]interface{}{
		"user": map[string]interface{}{
			"dashboard_display": "dashboard_display",
			"layout":            "layout",
		},
	}
	//test
	err := testConfig.MergeConfigMap(mergeSettings)

	//verify
	require.NoError(t, err)

	// if using Set, the MergeConfigMap functionality will not override
	// if using SetDefault, the MergeConfigMap will override correctly
	require.Equal(t, "hello", testConfig.GetString("user.dashboard_display"))
	require.Equal(t, "layout", testConfig.GetString("user.layout"))

}
