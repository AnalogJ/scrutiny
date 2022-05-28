package config_test

import (
	"github.com/analogj/scrutiny/collector/pkg/config"
	"github.com/analogj/scrutiny/collector/pkg/models"
	"github.com/stretchr/testify/require"
	"path"
	"testing"
)

func TestConfiguration_InvalidConfigPath(t *testing.T) {
	t.Parallel()

	//setup
	testConfig, _ := config.Create()

	//test
	err := testConfig.ReadConfig("does_not_exist.yaml")

	//assert
	require.Error(t, err, "should return an error")
}

func TestConfiguration_GetScanOverrides_Simple(t *testing.T) {
	t.Parallel()

	//setup
	testConfig, _ := config.Create()

	//test
	err := testConfig.ReadConfig(path.Join("testdata", "simple_device.yaml"))
	require.NoError(t, err, "should correctly load simple device config")
	scanOverrides := testConfig.GetDeviceOverrides()

	//assert
	require.Equal(t, []models.ScanOverride{{Device: "/dev/sda", DeviceType: []string{"sat"}, Ignore: false}}, scanOverrides)
}

func TestConfiguration_GetScanOverrides_Ignore(t *testing.T) {
	t.Parallel()

	//setup
	testConfig, _ := config.Create()

	//test
	err := testConfig.ReadConfig(path.Join("testdata", "ignore_device.yaml"))
	require.NoError(t, err, "should correctly load ignore device config")
	scanOverrides := testConfig.GetDeviceOverrides()

	//assert
	require.Equal(t, []models.ScanOverride{{Device: "/dev/sda", DeviceType: nil, Ignore: true}}, scanOverrides)
}

func TestConfiguration_GetScanOverrides_Raid(t *testing.T) {
	t.Parallel()

	//setup
	testConfig, _ := config.Create()

	//test
	err := testConfig.ReadConfig(path.Join("testdata", "raid_device.yaml"))
	require.NoError(t, err, "should correctly load ignore device config")
	scanOverrides := testConfig.GetDeviceOverrides()

	//assert
	require.Equal(t, []models.ScanOverride{
		{
			Device:     "/dev/bus/0",
			DeviceType: []string{"megaraid,14", "megaraid,15", "megaraid,18", "megaraid,19", "megaraid,20", "megaraid,21"},
			Ignore:     false,
		},
		{
			Device:     "/dev/twa0",
			DeviceType: []string{"3ware,0", "3ware,1", "3ware,2", "3ware,3", "3ware,4", "3ware,5"},
			Ignore:     false,
		}}, scanOverrides)
}

func TestConfiguration_InvalidCommands_MissingJson(t *testing.T) {
	t.Parallel()

	//setup
	testConfig, _ := config.Create()

	//test
	err := testConfig.ReadConfig(path.Join("testdata", "invalid_commands_missing_json.yaml"))
	require.EqualError(t, err, `ConfigValidationError: "configuration key 'commands.metrics_scan_args' is missing '--json' flag"`, "should throw an error because json flag is missing")
}

func TestConfiguration_InvalidCommands_IncludesDevice(t *testing.T) {
	t.Parallel()

	//setup
	testConfig, _ := config.Create()

	//test
	err := testConfig.ReadConfig(path.Join("testdata", "invalid_commands_includes_device.yaml"))
	require.EqualError(t, err, `ConfigValidationError: "configuration key 'commands.metrics_info_args' must not contain '--device' or '-d' flag, configuration key 'commands.metrics_smart_args' must not contain '--device' or '-d' flag"`, "should throw an error because device flags detected")
}

func TestConfiguration_OverrideCommands(t *testing.T) {
	t.Parallel()

	//setup
	testConfig, _ := config.Create()

	//test
	err := testConfig.ReadConfig(path.Join("testdata", "override_commands.yaml"))
	require.NoError(t, err, "should not throw an error")
	require.Equal(t, "--xall --json -T permissive", testConfig.GetString("commands.metrics_smart_args"))
}

func TestConfiguration_OverrideDeviceCommands_MetricsInfoArgs(t *testing.T) {
	t.Parallel()

	//setup
	testConfig, _ := config.Create()

	//test
	err := testConfig.ReadConfig(path.Join("testdata", "override_device_commands.yaml"))
	require.NoError(t, err, "should correctly override device command")

	//assert
	require.Equal(t, "--info --json -T permissive", testConfig.GetCommandMetricsInfoArgs("/dev/sda"))
	require.Equal(t, "--info --json", testConfig.GetCommandMetricsInfoArgs("/dev/sdb"))
	//require.Equal(t, []models.ScanOverride{{Device: "/dev/sda", DeviceType: nil, Commands: {MetricsInfoArgs: "--info --json -T "}}}, scanOverrides)
}
