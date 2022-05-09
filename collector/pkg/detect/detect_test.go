package detect_test

import (
	mock_shell "github.com/analogj/scrutiny/collector/pkg/common/shell/mock"
	mock_config "github.com/analogj/scrutiny/collector/pkg/config/mock"
	"github.com/analogj/scrutiny/collector/pkg/detect"
	"github.com/analogj/scrutiny/collector/pkg/models"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func TestDetect_SmartctlScan(t *testing.T) {
	//setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("host.id").AnyTimes().Return("")
	fakeConfig.EXPECT().GetScanOverrides().AnyTimes().Return([]models.ScanOverride{})

	fakeShell := mock_shell.NewMockInterface(mockCtrl)
	testScanResults, err := ioutil.ReadFile("testdata/smartctl_scan_simple.json")
	fakeShell.EXPECT().Command(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(string(testScanResults), err)

	d := detect.Detect{
		Logger: logrus.WithFields(logrus.Fields{}),
		Shell:  fakeShell,
		Config: fakeConfig,
	}

	//test
	scannedDevices, err := d.SmartctlScan()

	//assert
	require.NoError(t, err)
	require.Equal(t, 7, len(scannedDevices))
	require.Equal(t, "scsi", scannedDevices[0].DeviceType)
}

func TestDetect_SmartctlScan_Megaraid(t *testing.T) {
	//setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("host.id").AnyTimes().Return("")
	fakeConfig.EXPECT().GetScanOverrides().AnyTimes().Return([]models.ScanOverride{})

	fakeShell := mock_shell.NewMockInterface(mockCtrl)
	testScanResults, err := ioutil.ReadFile("testdata/smartctl_scan_megaraid.json")
	fakeShell.EXPECT().Command(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(string(testScanResults), err)

	d := detect.Detect{
		Logger: logrus.WithFields(logrus.Fields{}),
		Shell:  fakeShell,
		Config: fakeConfig,
	}

	//test
	scannedDevices, err := d.SmartctlScan()

	//assert
	require.NoError(t, err)
	require.Equal(t, 2, len(scannedDevices))
	require.Equal(t, []models.Device{
		models.Device{DeviceName: "bus/0", DeviceType: "megaraid,0"},
		models.Device{DeviceName: "bus/0", DeviceType: "megaraid,1"},
	}, scannedDevices)
}

func TestDetect_SmartctlScan_Nvme(t *testing.T) {
	//setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("host.id").AnyTimes().Return("")
	fakeConfig.EXPECT().GetScanOverrides().AnyTimes().Return([]models.ScanOverride{})

	fakeShell := mock_shell.NewMockInterface(mockCtrl)
	testScanResults, err := ioutil.ReadFile("testdata/smartctl_scan_nvme.json")
	fakeShell.EXPECT().Command(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(string(testScanResults), err)

	d := detect.Detect{
		Logger: logrus.WithFields(logrus.Fields{}),
		Shell:  fakeShell,
		Config: fakeConfig,
	}

	//test
	scannedDevices, err := d.SmartctlScan()

	//assert
	require.NoError(t, err)
	require.Equal(t, 1, len(scannedDevices))
	require.Equal(t, []models.Device{
		models.Device{DeviceName: "nvme0", DeviceType: "nvme"},
	}, scannedDevices)
}

func TestDetect_TransformDetectedDevices_Empty(t *testing.T) {
	//setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("host.id").AnyTimes().Return("")
	fakeConfig.EXPECT().GetScanOverrides().AnyTimes().Return([]models.ScanOverride{})
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

	//test
	transformedDevices := d.TransformDetectedDevices(detectedDevices)

	//assert
	require.Equal(t, "sda", transformedDevices[0].DeviceName)
	require.Equal(t, "scsi", transformedDevices[0].DeviceType)
}

func TestDetect_TransformDetectedDevices_Ignore(t *testing.T) {
	//setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("host.id").AnyTimes().Return("")
	fakeConfig.EXPECT().GetScanOverrides().AnyTimes().Return([]models.ScanOverride{{Device: "/dev/sda", DeviceType: nil, Ignore: true}})
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

	//test
	transformedDevices := d.TransformDetectedDevices(detectedDevices)

	//assert
	require.Equal(t, []models.Device{}, transformedDevices)
}

func TestDetect_TransformDetectedDevices_Raid(t *testing.T) {
	//setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("host.id").AnyTimes().Return("")
	fakeConfig.EXPECT().GetScanOverrides().AnyTimes().Return([]models.ScanOverride{
		{
			Device:     "/dev/bus/0",
			DeviceType: []string{"megaraid,14", "megaraid,15", "megaraid,18", "megaraid,19", "megaraid,20", "megaraid,21"},
			Ignore:     false,
		},
		{
			Device:     "/dev/twa0",
			DeviceType: []string{"3ware,0", "3ware,1", "3ware,2", "3ware,3", "3ware,4", "3ware,5"},
			Ignore:     false,
		}})
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

	//test
	transformedDevices := d.TransformDetectedDevices(detectedDevices)

	//assert
	require.Equal(t, 12, len(transformedDevices))
}

func TestDetect_TransformDetectedDevices_Simple(t *testing.T) {
	//setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetString("host.id").AnyTimes().Return("")
	fakeConfig.EXPECT().GetScanOverrides().AnyTimes().Return([]models.ScanOverride{{Device: "/dev/sda", DeviceType: []string{"sat+megaraid"}}})
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

	//test
	transformedDevices := d.TransformDetectedDevices(detectedDevices)

	//assert
	require.Equal(t, 1, len(transformedDevices))
	require.Equal(t, "sat+megaraid", transformedDevices[0].DeviceType)
}
