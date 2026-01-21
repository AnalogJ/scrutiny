package measurements_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/analogj/scrutiny/webapp/backend/pkg"
	mock_config "github.com/analogj/scrutiny/webapp/backend/pkg/config/mock"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestSmart_Flatten(t *testing.T) {
	//setup
	timeNow := time.Now()
	smart := measurements.Smart{
		Date:            timeNow,
		DeviceWWN:       "test-wwn",
		DeviceProtocol:  pkg.DeviceProtocolAta,
		Temp:            50,
		PowerOnHours:    10,
		PowerCycleCount: 10,
		LogicalBlockSize: 512,
		Attributes:      nil,
		Status:          0,
	}

	//test
	tags, fields := smart.Flatten()

	//assert
	require.Equal(t, map[string]string{"device_protocol": "ATA", "device_wwn": "test-wwn"}, tags)
	require.Equal(t, map[string]interface{}{"logical_block_size": int64(512), "power_cycle_count": int64(10), "power_on_hours": int64(10), "temp": int64(50)}, fields)
}

func TestSmart_Flatten_ATA(t *testing.T) {
	//setup
	timeNow := time.Now()
	smart := measurements.Smart{
		Date:            timeNow,
		DeviceWWN:       "test-wwn",
		DeviceProtocol:  pkg.DeviceProtocolAta,
		Temp:            50,
		PowerOnHours:    10,
		PowerCycleCount: 10,
		LogicalBlockSize: 512,
		Status:          0,
		Attributes: map[string]measurements.SmartAttribute{
			"1": &measurements.SmartAtaAttribute{
				AttributeId: 1,
				Value:       100,
				Threshold:   1,
				Worst:       100,
				RawValue:    0,
				RawString:   "0",
				WhenFailed:  "",
			},
			"2": &measurements.SmartAtaAttribute{
				AttributeId: 2,
				Value:       135,
				Threshold:   54,
				Worst:       135,
				RawValue:    108,
				RawString:   "108",
				WhenFailed:  "",
			},
		},
	}

	//test
	tags, fields := smart.Flatten()

	//assert
	require.Equal(t, map[string]string{"device_protocol": "ATA", "device_wwn": "test-wwn"}, tags)
	require.Equal(t, map[string]interface{}{
		"attr.1.attribute_id":      "1",
		"attr.1.failure_rate":      float64(0),
		"attr.1.raw_string":        "0",
		"attr.1.raw_value":         int64(0),
		"attr.1.status":            int64(0),
		"attr.1.status_reason":     "",
		"attr.1.thresh":            int64(1),
		"attr.1.transformed_value": int64(0),
		"attr.1.value":             int64(100),
		"attr.1.when_failed":       "",
		"attr.1.worst":             int64(100),

		"attr.2.attribute_id":      "2",
		"attr.2.failure_rate":      float64(0),
		"attr.2.raw_string":        "108",
		"attr.2.raw_value":         int64(108),
		"attr.2.status":            int64(0),
		"attr.2.status_reason":     "",
		"attr.2.thresh":            int64(54),
		"attr.2.transformed_value": int64(0),
		"attr.2.value":             int64(135),
		"attr.2.when_failed":       "",
		"attr.2.worst":             int64(135),

		"logical_block_size": int64(512),
		"power_cycle_count": int64(10),
		"power_on_hours":    int64(10),
		"temp":              int64(50),
	}, fields)
}

func TestSmart_Flatten_SCSI(t *testing.T) {
	//setup
	timeNow := time.Now()
	smart := measurements.Smart{
		Date:            timeNow,
		DeviceWWN:       "test-wwn",
		DeviceProtocol:  pkg.DeviceProtocolScsi,
		Temp:            50,
		PowerOnHours:    10,
		PowerCycleCount: 10,
		LogicalBlockSize: 512,
		Status:          0,
		Attributes: map[string]measurements.SmartAttribute{
			"read_errors_corrected_by_eccfast": &measurements.SmartScsiAttribute{
				AttributeId: "read_errors_corrected_by_eccfast",
				Value:       int64(300357663),
			},
		},
	}

	//test
	tags, fields := smart.Flatten()

	//assert
	require.Equal(t, map[string]string{"device_protocol": "SCSI", "device_wwn": "test-wwn"}, tags)
	require.Equal(t, map[string]interface{}{
		"attr.read_errors_corrected_by_eccfast.attribute_id":      "read_errors_corrected_by_eccfast",
		"attr.read_errors_corrected_by_eccfast.failure_rate":      float64(0),
		"attr.read_errors_corrected_by_eccfast.status":            int64(0),
		"attr.read_errors_corrected_by_eccfast.status_reason":     "",
		"attr.read_errors_corrected_by_eccfast.thresh":            int64(0),
		"attr.read_errors_corrected_by_eccfast.transformed_value": int64(0),
		"attr.read_errors_corrected_by_eccfast.value":             int64(300357663),
		"logical_block_size": int64(512),
		"power_cycle_count": int64(10),
		"power_on_hours":    int64(10),
		"temp":              int64(50)},
		fields)
}

func TestSmart_Flatten_NVMe(t *testing.T) {
	//setup
	timeNow := time.Now()
	smart := measurements.Smart{
		Date:            timeNow,
		DeviceWWN:       "test-wwn",
		DeviceProtocol:  pkg.DeviceProtocolNvme,
		Temp:            50,
		PowerOnHours:    10,
		PowerCycleCount: 10,
		LogicalBlockSize: 512,
		Status:          0,
		Attributes: map[string]measurements.SmartAttribute{
			"available_spare": &measurements.SmartNvmeAttribute{
				AttributeId: "available_spare",
				Value:       int64(100),
			},
		},
	}

	//test
	tags, fields := smart.Flatten()

	//assert
	require.Equal(t, map[string]string{"device_protocol": "NVMe", "device_wwn": "test-wwn"}, tags)
	require.Equal(t, map[string]interface{}{
		"attr.available_spare.attribute_id":      "available_spare",
		"attr.available_spare.failure_rate":      float64(0),
		"attr.available_spare.status":            int64(0),
		"attr.available_spare.status_reason":     "",
		"attr.available_spare.thresh":            int64(0),
		"attr.available_spare.transformed_value": int64(0),
		"attr.available_spare.value":             int64(100),
		"logical_block_size":                     int64(512),
		"power_cycle_count":                      int64(10),
		"power_on_hours":                         int64(10),
		"temp":                                   int64(50)}, fields)
}

func TestNewSmartFromInfluxDB_ATA(t *testing.T) {
	//setup
	timeNow := time.Now()
	attrs := map[string]interface{}{
		"_time":                    timeNow,
		"device_wwn":               "test-wwn",
		"device_protocol":          pkg.DeviceProtocolAta,
		"attr.1.attribute_id":      "1",
		"attr.1.failure_rate":      float64(0),
		"attr.1.raw_string":        "108",
		"attr.1.raw_value":         int64(108),
		"attr.1.status":            int64(0),
		"attr.1.status_reason":     "",
		"attr.1.thresh":            int64(54),
		"attr.1.transformed_value": int64(0),
		"attr.1.value":             int64(135),
		"attr.1.when_failed":       "",
		"attr.1.worst":             int64(135),
		"power_cycle_count":        int64(10),
		"power_on_hours":           int64(10),
		"temp":                     int64(50),
	}

	//test
	smart, err := measurements.NewSmartFromInfluxDB(attrs)

	//assert
	require.NoError(t, err)
	require.Equal(t, &measurements.Smart{
		Date:            timeNow,
		DeviceWWN:       "test-wwn",
		DeviceProtocol:  "ATA",
		Temp:            50,
		PowerOnHours:    10,
		PowerCycleCount: 10,
		Attributes: map[string]measurements.SmartAttribute{
			"1": &measurements.SmartAtaAttribute{
				AttributeId: 1,
				Value:       135,
				Threshold:   54,
				Worst:       135,
				RawValue:    108,
				RawString:   "108",
				WhenFailed:  "",
			},
		}, Status: 0}, smart)
}

func TestNewSmartFromInfluxDB_NVMe(t *testing.T) {
	//setup
	timeNow := time.Now()
	attrs := map[string]interface{}{
		"_time":                                  timeNow,
		"device_wwn":                             "test-wwn",
		"device_protocol":                        pkg.DeviceProtocolNvme,
		"attr.available_spare.attribute_id":      "available_spare",
		"attr.available_spare.failure_rate":      float64(0),
		"attr.available_spare.status":            int64(0),
		"attr.available_spare.status_reason":     "",
		"attr.available_spare.thresh":            int64(0),
		"attr.available_spare.transformed_value": int64(0),
		"attr.available_spare.value":             int64(100),
		"power_cycle_count":                      int64(10),
		"power_on_hours":                         int64(10),
		"temp":                                   int64(50),
	}

	//test
	smart, err := measurements.NewSmartFromInfluxDB(attrs)

	//assert
	require.NoError(t, err)
	require.Equal(t, &measurements.Smart{
		Date:            timeNow,
		DeviceWWN:       "test-wwn",
		DeviceProtocol:  "NVMe",
		Temp:            50,
		PowerOnHours:    10,
		PowerCycleCount: 10,
		Attributes: map[string]measurements.SmartAttribute{
			"available_spare": &measurements.SmartNvmeAttribute{
				AttributeId: "available_spare",
				Value:       int64(100),
			},
		}, Status: 0}, smart)
}

func TestNewSmartFromInfluxDB_SCSI(t *testing.T) {
	//setup
	timeNow := time.Now()
	attrs := map[string]interface{}{
		"_time":           timeNow,
		"device_wwn":      "test-wwn",
		"device_protocol": pkg.DeviceProtocolScsi,
		"attr.read_errors_corrected_by_eccfast.attribute_id":      "read_errors_corrected_by_eccfast",
		"attr.read_errors_corrected_by_eccfast.failure_rate":      float64(0),
		"attr.read_errors_corrected_by_eccfast.status":            int64(0),
		"attr.read_errors_corrected_by_eccfast.status_reason":     "",
		"attr.read_errors_corrected_by_eccfast.thresh":            int64(0),
		"attr.read_errors_corrected_by_eccfast.transformed_value": int64(0),
		"attr.read_errors_corrected_by_eccfast.value":             int64(300357663),
		"power_cycle_count": int64(10),
		"power_on_hours":    int64(10),
		"temp":              int64(50),
	}

	//test
	smart, err := measurements.NewSmartFromInfluxDB(attrs)

	//assert
	require.NoError(t, err)
	require.Equal(t, &measurements.Smart{
		Date:            timeNow,
		DeviceWWN:       "test-wwn",
		DeviceProtocol:  "SCSI",
		Temp:            50,
		PowerOnHours:    10,
		PowerCycleCount: 10,
		Attributes: map[string]measurements.SmartAttribute{
			"read_errors_corrected_by_eccfast": &measurements.SmartScsiAttribute{
				AttributeId: "read_errors_corrected_by_eccfast",
				Value:       int64(300357663),
			},
		}, Status: 0}, smart)
}

func TestFromCollectorSmartInfo(t *testing.T) {
	//setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetIntSlice("failures.transient.ata").Return([]int{195}).AnyTimes()

	smartDataFile, err := os.Open("../testdata/smart-ata.json")
	require.NoError(t, err)
	defer smartDataFile.Close()

	var smartJson collector.SmartInfo

	smartDataBytes, err := ioutil.ReadAll(smartDataFile)
	require.NoError(t, err)
	err = json.Unmarshal(smartDataBytes, &smartJson)
	require.NoError(t, err)

	//test
	smartMdl := measurements.Smart{}
	err = smartMdl.FromCollectorSmartInfo(fakeConfig, "WWN-test", smartJson)

	//assert
	require.NoError(t, err)
	require.Equal(t, "WWN-test", smartMdl.DeviceWWN)
	require.Equal(t, pkg.DeviceStatusPassed, smartMdl.Status)
	require.Equal(t, 18, len(smartMdl.Attributes))

	//check that temperature was correctly parsed
	require.Equal(t, int64(163210330144), smartMdl.Attributes["194"].(*measurements.SmartAtaAttribute).RawValue)
	require.Equal(t, int64(32), smartMdl.Attributes["194"].(*measurements.SmartAtaAttribute).TransformedValue)

	//ensure that Scrutiny warning for a non critical attribute does not set device status to failed.
	require.Equal(t, pkg.AttributeStatusWarningScrutiny, smartMdl.Attributes["3"].GetStatus())

}

func TestFromCollectorSmartInfo_Fail_Smart(t *testing.T) {
	//setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetIntSlice("failures.transient.ata").Return([]int{195}).AnyTimes()

	smartDataFile, err := os.Open("../testdata/smart-fail.json")
	require.NoError(t, err)
	defer smartDataFile.Close()

	var smartJson collector.SmartInfo

	smartDataBytes, err := ioutil.ReadAll(smartDataFile)
	require.NoError(t, err)
	err = json.Unmarshal(smartDataBytes, &smartJson)
	require.NoError(t, err)

	//test
	smartMdl := measurements.Smart{}
	err = smartMdl.FromCollectorSmartInfo(fakeConfig, "WWN-test", smartJson)

	//assert
	require.NoError(t, err)
	require.Equal(t, "WWN-test", smartMdl.DeviceWWN)
	require.Equal(t, pkg.DeviceStatusFailedSmart, smartMdl.Status)
	require.Equal(t, 0, len(smartMdl.Attributes))
}

func TestFromCollectorSmartInfo_Fail_ScrutinySmart(t *testing.T) {
	//setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetIntSlice("failures.transient.ata").Return([]int{195}).AnyTimes()

	smartDataFile, err := os.Open("../testdata/smart-fail2.json")
	require.NoError(t, err)
	defer smartDataFile.Close()

	var smartJson collector.SmartInfo

	smartDataBytes, err := ioutil.ReadAll(smartDataFile)
	require.NoError(t, err)
	err = json.Unmarshal(smartDataBytes, &smartJson)
	require.NoError(t, err)

	//test
	smartMdl := measurements.Smart{}
	err = smartMdl.FromCollectorSmartInfo(fakeConfig, "WWN-test", smartJson)

	//assert
	require.NoError(t, err)
	require.Equal(t, "WWN-test", smartMdl.DeviceWWN)
	require.Equal(t, pkg.DeviceStatusFailedScrutiny|pkg.DeviceStatusFailedSmart, smartMdl.Status)
	require.Equal(t, 17, len(smartMdl.Attributes))
}

func TestFromCollectorSmartInfo_Fail_ScrutinyNonCriticalFailed(t *testing.T) {
	//setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetIntSlice("failures.transient.ata").Return([]int{195}).AnyTimes()

	smartDataFile, err := os.Open("../testdata/smart-ata-failed-scrutiny.json")
	require.NoError(t, err)
	defer smartDataFile.Close()

	var smartJson collector.SmartInfo

	smartDataBytes, err := ioutil.ReadAll(smartDataFile)
	require.NoError(t, err)
	err = json.Unmarshal(smartDataBytes, &smartJson)
	require.NoError(t, err)

	//test
	smartMdl := measurements.Smart{}
	err = smartMdl.FromCollectorSmartInfo(fakeConfig, "WWN-test", smartJson)

	//assert
	require.NoError(t, err)
	require.Equal(t, "WWN-test", smartMdl.DeviceWWN)
	require.Equal(t, pkg.DeviceStatusFailedScrutiny, smartMdl.Status)
	require.Equal(t, pkg.AttributeStatusFailedScrutiny, smartMdl.Attributes["199"].GetStatus(),
		"scrutiny should detect that %d failed (status: %d, %s)",
		smartMdl.Attributes["199"].(*measurements.SmartAtaAttribute).AttributeId,
		smartMdl.Attributes["199"].GetStatus(), smartMdl.Attributes["199"].(*measurements.SmartAtaAttribute).StatusReason,
	)

	require.Equal(t, 14, len(smartMdl.Attributes))
}

//TODO: Scrutiny Warn
//TODO: Smart + Scrutiny Warn

func TestFromCollectorSmartInfo_NVMe_Fail_Scrutiny(t *testing.T) {
	//setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetIntSlice("failures.transient.ata").Return([]int{195}).AnyTimes()

	smartDataFile, err := os.Open("../testdata/smart-nvme-failed.json")
	require.NoError(t, err)
	defer smartDataFile.Close()

	var smartJson collector.SmartInfo

	smartDataBytes, err := ioutil.ReadAll(smartDataFile)
	require.NoError(t, err)
	err = json.Unmarshal(smartDataBytes, &smartJson)
	require.NoError(t, err)

	//test
	smartMdl := measurements.Smart{}
	err = smartMdl.FromCollectorSmartInfo(fakeConfig, "WWN-test", smartJson)

	//assert
	require.NoError(t, err)
	require.Equal(t, "WWN-test", smartMdl.DeviceWWN)
	require.Equal(t, pkg.DeviceStatusFailedScrutiny, smartMdl.Status)
	require.Equal(t, pkg.AttributeStatusFailedScrutiny, smartMdl.Attributes["media_errors"].GetStatus(),
		"scrutiny should detect that %s failed (status: %d, %s)",
		smartMdl.Attributes["media_errors"].(*measurements.SmartNvmeAttribute).AttributeId,
		smartMdl.Attributes["media_errors"].GetStatus(),
		smartMdl.Attributes["media_errors"].(*measurements.SmartNvmeAttribute).StatusReason,
	)

	require.Equal(t, 16, len(smartMdl.Attributes))
}

func TestFromCollectorSmartInfo_Nvme(t *testing.T) {
	//setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetIntSlice("failures.transient.ata").Return([]int{195}).AnyTimes()

	smartDataFile, err := os.Open("../testdata/smart-nvme.json")
	require.NoError(t, err)
	defer smartDataFile.Close()

	var smartJson collector.SmartInfo

	smartDataBytes, err := ioutil.ReadAll(smartDataFile)
	require.NoError(t, err)
	err = json.Unmarshal(smartDataBytes, &smartJson)
	require.NoError(t, err)

	//test
	smartMdl := measurements.Smart{}
	err = smartMdl.FromCollectorSmartInfo(fakeConfig, "WWN-test", smartJson)

	//assert
	require.NoError(t, err)
	require.Equal(t, "WWN-test", smartMdl.DeviceWWN)
	require.Equal(t, pkg.DeviceStatusPassed, smartMdl.Status)
	require.Equal(t, 16, len(smartMdl.Attributes))

	require.Equal(t, int64(111303174), smartMdl.Attributes["host_reads"].(*measurements.SmartNvmeAttribute).Value)
	require.Equal(t, int64(83170961), smartMdl.Attributes["host_writes"].(*measurements.SmartNvmeAttribute).Value)
}

func TestFromCollectorSmartInfo_Scsi(t *testing.T) {
	//setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetIntSlice("failures.transient.ata").Return([]int{195}).AnyTimes()

	smartDataFile, err := os.Open("../testdata/smart-scsi.json")
	require.NoError(t, err)
	defer smartDataFile.Close()

	var smartJson collector.SmartInfo

	smartDataBytes, err := ioutil.ReadAll(smartDataFile)
	require.NoError(t, err)
	err = json.Unmarshal(smartDataBytes, &smartJson)
	require.NoError(t, err)

	//test
	smartMdl := measurements.Smart{}
	err = smartMdl.FromCollectorSmartInfo(fakeConfig, "WWN-test", smartJson)

	//assert
	require.NoError(t, err)
	require.Equal(t, "WWN-test", smartMdl.DeviceWWN)
	require.Equal(t, pkg.DeviceStatusPassed, smartMdl.Status)
	require.Equal(t, 14, len(smartMdl.Attributes))

	require.Equal(t, int64(56), smartMdl.Attributes["scsi_grown_defect_list"].(*measurements.SmartScsiAttribute).Value)
	require.Equal(t, int64(300357663), smartMdl.Attributes["read_errors_corrected_by_eccfast"].(*measurements.SmartScsiAttribute).Value) //total_errors_corrected
}

// TestFromCollectorSmartInfo_Scsi_SAS_EnvironmentalReports tests that for SAS drives
// where the standard temperature field is 0, the temperature is correctly parsed
// from scsi_environmental_reports.temperature_1.current (fixes GitHub issue #26)
func TestFromCollectorSmartInfo_Scsi_SAS_EnvironmentalReports(t *testing.T) {
	//setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetIntSlice("failures.transient.ata").Return([]int{195}).AnyTimes()

	smartDataFile, err := os.Open("../testdata/smart-scsi-sas-env-temp.json")
	require.NoError(t, err)
	defer smartDataFile.Close()

	var smartJson collector.SmartInfo

	smartDataBytes, err := ioutil.ReadAll(smartDataFile)
	require.NoError(t, err)
	err = json.Unmarshal(smartDataBytes, &smartJson)
	require.NoError(t, err)

	//test
	smartMdl := measurements.Smart{}
	err = smartMdl.FromCollectorSmartInfo(fakeConfig, "WWN-test", smartJson)

	//assert
	require.NoError(t, err)
	require.Equal(t, "WWN-test", smartMdl.DeviceWWN)
	require.Equal(t, pkg.DeviceStatusPassed, smartMdl.Status)

	// The temperature.current in the JSON is 0, but scsi_environmental_reports.temperature_1.current is 38
	// The fix should correctly extract 38 from scsi_environmental_reports
	require.Equal(t, int64(38), smartMdl.Temp, "Temperature should be parsed from scsi_environmental_reports when standard temperature is 0")
}

// TestFromCollectorSmartInfo_ATA_DeviceStatistics tests that ATA Device Statistics
// from GP Log 0x04 are correctly parsed, including enterprise SSD metrics like
// "Percentage Used Endurance Indicator" (devstat_7_8). Fixes GitHub issue #7 (SCR-11).
func TestFromCollectorSmartInfo_ATA_DeviceStatistics(t *testing.T) {
	//setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().GetIntSlice("failures.transient.ata").Return([]int{195}).AnyTimes()

	smartDataFile, err := os.Open("../testdata/smart-ata-full.json")
	require.NoError(t, err)
	defer smartDataFile.Close()

	var smartJson collector.SmartInfo

	smartDataBytes, err := ioutil.ReadAll(smartDataFile)
	require.NoError(t, err)
	err = json.Unmarshal(smartDataBytes, &smartJson)
	require.NoError(t, err)

	//test
	smartMdl := measurements.Smart{}
	err = smartMdl.FromCollectorSmartInfo(fakeConfig, "WWN-test", smartJson)

	//assert
	require.NoError(t, err)
	require.Equal(t, "WWN-test", smartMdl.DeviceWWN)
	require.Equal(t, pkg.DeviceProtocolAta, smartMdl.DeviceProtocol)

	// Check that device statistics are parsed with correct attribute IDs
	// devstat_7_8 is "Percentage Used Endurance Indicator" - critical for enterprise SSDs
	devstat, ok := smartMdl.Attributes["devstat_7_8"]
	require.True(t, ok, "devstat_7_8 (Percentage Used Endurance Indicator) should be present")

	// Verify it's the correct type
	devstatAttr, ok := devstat.(*measurements.SmartAtaDeviceStatAttribute)
	require.True(t, ok, "devstat_7_8 should be SmartAtaDeviceStatAttribute type")

	require.Equal(t, "devstat_7_8", devstatAttr.AttributeId)
	require.Equal(t, int64(19), devstatAttr.Value, "Percentage Used should be 19%")

	// Verify other device statistics are also present (page 1, offset 8 = Lifetime Power-On Resets)
	_, ok = smartMdl.Attributes["devstat_1_8"]
	require.True(t, ok, "devstat_1_8 (Lifetime Power-On Resets) should be present")
}

// TestSmart_Flatten_WithDeviceStatistics tests that device statistics are correctly
// flattened for storage in InfluxDB with proper string-based attribute IDs.
func TestSmart_Flatten_WithDeviceStatistics(t *testing.T) {
	//setup
	timeNow := time.Now()
	smart := measurements.Smart{
		Date:            timeNow,
		DeviceWWN:       "test-wwn",
		DeviceProtocol:  pkg.DeviceProtocolAta,
		Temp:            50,
		PowerOnHours:    10,
		PowerCycleCount: 10,
		Attributes: map[string]measurements.SmartAttribute{
			"devstat_7_8": &measurements.SmartAtaDeviceStatAttribute{
				AttributeId:      "devstat_7_8",
				Value:            25,
				TransformedValue: 25,
			},
		},
		Status: 0,
	}

	//test
	_, fields := smart.Flatten()

	//assert
	require.Equal(t, "devstat_7_8", fields["attr.devstat_7_8.attribute_id"])
	require.Equal(t, int64(25), fields["attr.devstat_7_8.value"])
	require.Equal(t, int64(25), fields["attr.devstat_7_8.transformed_value"])
}

// TestNewSmartFromInfluxDB_WithDeviceStatistics tests that device statistics are
// correctly inflated from InfluxDB data.
func TestNewSmartFromInfluxDB_WithDeviceStatistics(t *testing.T) {
	//setup
	timeNow := time.Now()
	attrs := map[string]interface{}{
		"_time":           timeNow,
		"device_wwn":      "test-wwn",
		"device_protocol": pkg.DeviceProtocolAta,
		"temp":            int64(50),
		"power_on_hours":  int64(10),
		// Device statistics with string-based attribute ID
		"attr.devstat_7_8.attribute_id":      "devstat_7_8",
		"attr.devstat_7_8.value":             int64(42),
		"attr.devstat_7_8.thresh":            int64(100),
		"attr.devstat_7_8.transformed_value": int64(42),
		"attr.devstat_7_8.status":            int64(0),
		"attr.devstat_7_8.status_reason":     "",
		"attr.devstat_7_8.failure_rate":      float64(0),
	}

	//test
	smart, err := measurements.NewSmartFromInfluxDB(attrs)

	//assert
	require.NoError(t, err)
	require.Contains(t, smart.Attributes, "devstat_7_8")

	devstatAttr, ok := smart.Attributes["devstat_7_8"].(*measurements.SmartAtaDeviceStatAttribute)
	require.True(t, ok, "devstat_7_8 should be SmartAtaDeviceStatAttribute type")
	require.Equal(t, "devstat_7_8", devstatAttr.AttributeId)
	require.Equal(t, int64(42), devstatAttr.Value)
	require.Equal(t, int64(100), devstatAttr.Threshold)
}
