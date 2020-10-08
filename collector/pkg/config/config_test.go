package config_test

import (
	"github.com/analogj/scrutiny/collector/pkg/config"
	"github.com/analogj/scrutiny/collector/pkg/models"
	"github.com/stretchr/testify/require"
	"path"
	"testing"
)

func TestConfiguration_GetScanOverrides_Simple(t *testing.T) {
	t.Parallel()

	//setup
	testConfig, _ := config.Create()

	//test
	err := testConfig.ReadConfig(path.Join("testdata", "simple_device.yaml"))
	require.NoError(t, err, "should correctly load simple device config")
	scanOverrides := testConfig.GetScanOverrides()

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
	scanOverrides := testConfig.GetScanOverrides()

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
	scanOverrides := testConfig.GetScanOverrides()

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
