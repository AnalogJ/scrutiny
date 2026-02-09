package detect_test

import (
	"os"
	"strings"
	"testing"

	mock_shell "github.com/analogj/scrutiny/collector/pkg/common/shell/mock"
	mock_config "github.com/analogj/scrutiny/collector/pkg/config/mock"
	"github.com/analogj/scrutiny/collector/pkg/detect"
	"github.com/analogj/scrutiny/collector/pkg/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDetect_SmartctlScan(t *testing.T) {
	// setup
	mockCtrl := gomock.NewController(t)
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("host.id").AnyTimes().Return("")
	fakeConfig.EXPECT().GetDeviceOverrides().AnyTimes().Return([]models.ScanOverride{})
	fakeConfig.EXPECT().GetString("commands.metrics_smartctl_bin").AnyTimes().Return("smartctl")
	fakeConfig.EXPECT().GetString("commands.metrics_scan_args").AnyTimes().Return("--scan --json")
	fakeConfig.EXPECT().IsAllowlistedDevice(gomock.Any()).AnyTimes().Return(true)

	fakeShell := mock_shell.NewMockInterface(mockCtrl)
	testScanResults, err := os.ReadFile("testdata/smartctl_scan_simple.json")
	fakeShell.EXPECT().Command(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(string(testScanResults), err)

	d := detect.Detect{
		Logger: logrus.WithFields(logrus.Fields{}),
		Shell:  fakeShell,
		Config: fakeConfig,
	}

	// test
	scannedDevices, err := d.SmartctlScan()

	// assert
	require.NoError(t, err)
	require.Equal(t, 7, len(scannedDevices))
	require.Equal(t, "scsi", scannedDevices[0].DeviceType)
}

func TestDetect_SmartctlScan_Megaraid(t *testing.T) {
	// setup
	mockCtrl := gomock.NewController(t)
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("host.id").AnyTimes().Return("")
	fakeConfig.EXPECT().GetDeviceOverrides().AnyTimes().Return([]models.ScanOverride{})
	fakeConfig.EXPECT().GetString("commands.metrics_smartctl_bin").AnyTimes().Return("smartctl")
	fakeConfig.EXPECT().GetString("commands.metrics_scan_args").AnyTimes().Return("--scan --json")
	fakeConfig.EXPECT().IsAllowlistedDevice(gomock.Any()).AnyTimes().Return(true)

	fakeShell := mock_shell.NewMockInterface(mockCtrl)
	testScanResults, err := os.ReadFile("testdata/smartctl_scan_megaraid.json")
	fakeShell.EXPECT().Command(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(string(testScanResults), err)

	d := detect.Detect{
		Logger: logrus.WithFields(logrus.Fields{}),
		Shell:  fakeShell,
		Config: fakeConfig,
	}

	// test
	scannedDevices, err := d.SmartctlScan()

	// assert
	require.NoError(t, err)
	require.Equal(t, 2, len(scannedDevices))
	require.Equal(t, []models.Device{
		{DeviceName: "bus/0", DeviceType: "megaraid,0"},
		{DeviceName: "bus/0", DeviceType: "megaraid,1"},
	}, scannedDevices)
}

func TestDetect_SmartctlScan_Nvme(t *testing.T) {
	// setup
	mockCtrl := gomock.NewController(t)
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("host.id").AnyTimes().Return("")
	fakeConfig.EXPECT().GetDeviceOverrides().AnyTimes().Return([]models.ScanOverride{})
	fakeConfig.EXPECT().GetString("commands.metrics_smartctl_bin").AnyTimes().Return("smartctl")
	fakeConfig.EXPECT().GetString("commands.metrics_scan_args").AnyTimes().Return("--scan --json")
	fakeConfig.EXPECT().IsAllowlistedDevice(gomock.Any()).AnyTimes().Return(true)

	fakeShell := mock_shell.NewMockInterface(mockCtrl)
	testScanResults, err := os.ReadFile("testdata/smartctl_scan_nvme.json")
	fakeShell.EXPECT().Command(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(string(testScanResults), err)

	d := detect.Detect{
		Logger: logrus.WithFields(logrus.Fields{}),
		Shell:  fakeShell,
		Config: fakeConfig,
	}

	// test
	scannedDevices, err := d.SmartctlScan()

	// assert
	require.NoError(t, err)
	require.Equal(t, 1, len(scannedDevices))
	require.Equal(t, []models.Device{
		{DeviceName: "nvme0", DeviceType: "nvme"},
	}, scannedDevices)
}

func TestDetect_TransformDetectedDevices_Empty(t *testing.T) {
	// setup
	mockCtrl := gomock.NewController(t)
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("host.id").AnyTimes().Return("")
	fakeConfig.EXPECT().GetDeviceOverrides().AnyTimes().Return([]models.ScanOverride{})
	fakeConfig.EXPECT().GetString("commands.metrics_smartctl_bin").AnyTimes().Return("smartctl")
	fakeConfig.EXPECT().GetString("commands.metrics_scan_args").AnyTimes().Return("--scan --json")
	fakeConfig.EXPECT().IsAllowlistedDevice(gomock.Any()).AnyTimes().Return(true)

	detectedDevices := models.Scan{
		Devices: []models.ScanDevice{
			{
				Name:     "/dev/sda",
				InfoName: "/dev/sda",
				Protocol: "scsi",
				Type:     "scsi",
			},
		},
	}

	d := detect.Detect{
		Config: fakeConfig,
	}

	// test
	transformedDevices := d.TransformDetectedDevices(detectedDevices)

	// assert
	require.Equal(t, "sda", transformedDevices[0].DeviceName)
	require.Equal(t, "scsi", transformedDevices[0].DeviceType)
}

func TestDetect_TransformDetectedDevices_Ignore(t *testing.T) {
	// setup
	mockCtrl := gomock.NewController(t)
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("host.id").AnyTimes().Return("")
	fakeConfig.EXPECT().GetDeviceOverrides().AnyTimes().Return([]models.ScanOverride{{Device: "/dev/sda", DeviceType: nil, Ignore: true}})
	fakeConfig.EXPECT().GetString("commands.metrics_smartctl_bin").AnyTimes().Return("smartctl")
	fakeConfig.EXPECT().GetString("commands.metrics_scan_args").AnyTimes().Return("--scan --json")
	fakeConfig.EXPECT().IsAllowlistedDevice(gomock.Any()).AnyTimes().Return(true)

	detectedDevices := models.Scan{
		Devices: []models.ScanDevice{
			{
				Name:     "/dev/sda",
				InfoName: "/dev/sda",
				Protocol: "scsi",
				Type:     "scsi",
			},
		},
	}

	d := detect.Detect{
		Config: fakeConfig,
	}

	// test
	transformedDevices := d.TransformDetectedDevices(detectedDevices)

	// assert
	require.Equal(t, []models.Device{}, transformedDevices)
}

func TestDetect_TransformDetectedDevices_Raid(t *testing.T) {
	// setup
	mockCtrl := gomock.NewController(t)
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("host.id").AnyTimes().Return("")
	fakeConfig.EXPECT().GetString("commands.metrics_smartctl_bin").AnyTimes().Return("smartctl")
	fakeConfig.EXPECT().GetString("commands.metrics_scan_args").AnyTimes().Return("--scan --json")
	fakeConfig.EXPECT().IsAllowlistedDevice(gomock.Any()).AnyTimes().Return(true)
	fakeConfig.EXPECT().GetDeviceOverrides().AnyTimes().Return([]models.ScanOverride{
		{
			Device:     "/dev/bus/0",
			DeviceType: []string{"megaraid,14", "megaraid,15", "megaraid,18", "megaraid,19", "megaraid,20", "megaraid,21"},
			Ignore:     false,
		},
		{
			Device:     "/dev/twa0",
			DeviceType: []string{"3ware,0", "3ware,1", "3ware,2", "3ware,3", "3ware,4", "3ware,5"},
			Ignore:     false,
		},
	})
	detectedDevices := models.Scan{
		Devices: []models.ScanDevice{
			{
				Name:     "/dev/bus/0",
				InfoName: "/dev/bus/0",
				Protocol: "scsi",
				Type:     "scsi",
			},
		},
	}

	d := detect.Detect{
		Config: fakeConfig,
	}

	// test
	transformedDevices := d.TransformDetectedDevices(detectedDevices)

	// assert
	require.Equal(t, 12, len(transformedDevices))
}

func TestDetect_TransformDetectedDevices_Simple(t *testing.T) {
	// setup
	mockCtrl := gomock.NewController(t)
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("host.id").AnyTimes().Return("")
	fakeConfig.EXPECT().GetString("commands.metrics_smartctl_bin").AnyTimes().Return("smartctl")
	fakeConfig.EXPECT().GetString("commands.metrics_scan_args").AnyTimes().Return("--scan --json")
	fakeConfig.EXPECT().GetDeviceOverrides().AnyTimes().Return([]models.ScanOverride{{Device: "/dev/sda", DeviceType: []string{"sat+megaraid"}}})
	fakeConfig.EXPECT().IsAllowlistedDevice(gomock.Any()).AnyTimes().Return(true)
	detectedDevices := models.Scan{
		Devices: []models.ScanDevice{
			{
				Name:     "/dev/sda",
				InfoName: "/dev/sda",
				Protocol: "ata",
				Type:     "ata",
			},
		},
	}

	d := detect.Detect{
		Config: fakeConfig,
	}

	// test
	transformedDevices := d.TransformDetectedDevices(detectedDevices)

	// assert
	require.Equal(t, 1, len(transformedDevices))
	require.Equal(t, "sat+megaraid", transformedDevices[0].DeviceType)
}

// test https://github.com/AnalogJ/scrutiny/issues/255#issuecomment-1164024126
func TestDetect_TransformDetectedDevices_WithoutDeviceTypeOverride(t *testing.T) {
	// setup
	mockCtrl := gomock.NewController(t)
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("host.id").AnyTimes().Return("")
	fakeConfig.EXPECT().GetString("commands.metrics_smartctl_bin").AnyTimes().Return("smartctl")
	fakeConfig.EXPECT().GetString("commands.metrics_scan_args").AnyTimes().Return("--scan --json")
	fakeConfig.EXPECT().GetDeviceOverrides().AnyTimes().Return([]models.ScanOverride{{Device: "/dev/sda"}})
	fakeConfig.EXPECT().IsAllowlistedDevice(gomock.Any()).AnyTimes().Return(true)
	detectedDevices := models.Scan{
		Devices: []models.ScanDevice{
			{
				Name:     "/dev/sda",
				InfoName: "/dev/sda",
				Protocol: "ata",
				Type:     "scsi",
			},
		},
	}

	d := detect.Detect{
		Config: fakeConfig,
	}

	// test
	transformedDevices := d.TransformDetectedDevices(detectedDevices)

	// assert
	require.Equal(t, 1, len(transformedDevices))
	require.Equal(t, "scsi", transformedDevices[0].DeviceType)
}

func TestDetect_TransformDetectedDevices_WhenDeviceNotDetected(t *testing.T) {
	// setup
	mockCtrl := gomock.NewController(t)
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("host.id").AnyTimes().Return("")
	fakeConfig.EXPECT().GetString("commands.metrics_smartctl_bin").AnyTimes().Return("smartctl")
	fakeConfig.EXPECT().GetString("commands.metrics_scan_args").AnyTimes().Return("--scan --json")
	fakeConfig.EXPECT().GetDeviceOverrides().AnyTimes().Return([]models.ScanOverride{{Device: "/dev/sda"}})
	detectedDevices := models.Scan{}

	d := detect.Detect{
		Config: fakeConfig,
	}

	// test
	transformedDevices := d.TransformDetectedDevices(detectedDevices)

	// assert
	require.Equal(t, 1, len(transformedDevices))
	require.Equal(t, "ata", transformedDevices[0].DeviceType)
}

func TestDetect_TransformDetectedDevices_AllowListFilters(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("host.id").AnyTimes().Return("")
	fakeConfig.EXPECT().GetString("commands.metrics_smartctl_bin").AnyTimes().Return("smartctl")
	fakeConfig.EXPECT().GetString("commands.metrics_scan_args").AnyTimes().Return("--scan --json")
	fakeConfig.EXPECT().GetDeviceOverrides().AnyTimes().Return([]models.ScanOverride{{Device: "/dev/sda", DeviceType: []string{"sat+megaraid"}}})
	fakeConfig.EXPECT().IsAllowlistedDevice("/dev/sda").Return(true)
	fakeConfig.EXPECT().IsAllowlistedDevice("/dev/sdb").Return(false)
	detectedDevices := models.Scan{
		Devices: []models.ScanDevice{
			{
				Name:     "/dev/sda",
				InfoName: "/dev/sda",
				Protocol: "ata",
				Type:     "ata",
			},
			{
				Name:     "/dev/sdb",
				InfoName: "/dev/sdb",
				Protocol: "ata",
				Type:     "ata",
			},
		},
	}

	d := detect.Detect{
		Config: fakeConfig,
	}

	// test
	transformedDevices := d.TransformDetectedDevices(detectedDevices)

	// assert
	require.Equal(t, 1, len(transformedDevices))
	require.Equal(t, "sda", transformedDevices[0].DeviceName)
}

func TestDetect_SmartCtlInfo(t *testing.T) {
	t.Run("should report nvme info", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		const (
			someArgs = "--info --json"

			// device info
			someDeviceName           = "some-device-name"
			someModelName            = "KCD61LUL3T84"
			someSerialNumber         = "61Q0A05UT7B8"
			someFirmware             = "8002"
			someDeviceProtocol       = "NVMe"
			someDeviceType           = "nvme"
			someCapacity       int64 = 3840755982336
		)

		fakeConfig := mock_config.NewMockInterface(ctrl)
		fakeConfig.EXPECT().
			GetCommandMetricsInfoArgs("/dev/" + someDeviceName).
			Return(someArgs)
		fakeConfig.EXPECT().
			GetString("commands.metrics_smartctl_bin").
			Return("smartctl")

		someLogger := logrus.WithFields(logrus.Fields{})

		smartctlInfoResults, err := os.ReadFile("testdata/smartctl_info_nvme.json")
		require.NoError(t, err)

		fakeShell := mock_shell.NewMockInterface(ctrl)
		fakeShell.EXPECT().
			Command(someLogger, "smartctl", append(strings.Split(someArgs, " "), "/dev/"+someDeviceName), "", gomock.Any()).
			Return(string(smartctlInfoResults), err)

		d := detect.Detect{
			Logger: someLogger,
			Shell:  fakeShell,
			Config: fakeConfig,
		}

		someDevice := &models.Device{
			WWN:        "some wwn",
			DeviceName: someDeviceName,
		}

		require.NoError(t, d.SmartCtlInfo(someDevice))

		assert.Equal(t, someDeviceName, someDevice.DeviceName)
		assert.Equal(t, someModelName, someDevice.ModelName)
		assert.Equal(t, someSerialNumber, someDevice.SerialNumber)
		assert.Equal(t, someFirmware, someDevice.Firmware)
		assert.Equal(t, someDeviceProtocol, someDevice.DeviceProtocol)
		assert.Equal(t, someDeviceType, someDevice.DeviceType)
		assert.Equal(t, someCapacity, someDevice.Capacity)
	})
}
