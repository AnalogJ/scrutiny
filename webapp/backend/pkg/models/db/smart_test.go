package db_test

import (
	"encoding/json"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/db"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func TestFromCollectorSmartInfo(t *testing.T) {
	//setup
	smartDataFile, err := os.Open("../testdata/smart-ata.json")
	require.NoError(t, err)
	defer smartDataFile.Close()

	var smartJson collector.SmartInfo

	smartDataBytes, err := ioutil.ReadAll(smartDataFile)
	require.NoError(t, err)
	err = json.Unmarshal(smartDataBytes, &smartJson)
	require.NoError(t, err)

	//test
	smartMdl := db.Smart{}
	err = smartMdl.FromCollectorSmartInfo("WWN-test", smartJson)

	//assert
	require.NoError(t, err)
	require.Equal(t, "WWN-test", smartMdl.DeviceWWN)
	require.Equal(t, "passed", smartMdl.SmartStatus)
	require.Equal(t, 18, len(smartMdl.AtaAttributes))
	require.Equal(t, 0, len(smartMdl.NvmeAttributes))
	require.Equal(t, 0, len(smartMdl.ScsiAttributes))

	//check that temperature was correctly parsed
	for _, attr := range smartMdl.AtaAttributes {
		if attr.AttributeId == 194 {
			require.Equal(t, int64(163210330144), attr.RawValue)
			require.Equal(t, int64(32), attr.TransformedValue)
		}
	}
}

func TestFromCollectorSmartInfo_Fail(t *testing.T) {
	//setup
	smartDataFile, err := os.Open("../testdata/smart-fail.json")
	require.NoError(t, err)
	defer smartDataFile.Close()

	var smartJson collector.SmartInfo

	smartDataBytes, err := ioutil.ReadAll(smartDataFile)
	require.NoError(t, err)
	err = json.Unmarshal(smartDataBytes, &smartJson)
	require.NoError(t, err)

	//test
	smartMdl := db.Smart{}
	err = smartMdl.FromCollectorSmartInfo("WWN-test", smartJson)

	//assert
	require.NoError(t, err)
	require.Equal(t, "WWN-test", smartMdl.DeviceWWN)
	require.Equal(t, "failed", smartMdl.SmartStatus)
	require.Equal(t, 0, len(smartMdl.AtaAttributes))
	require.Equal(t, 0, len(smartMdl.NvmeAttributes))
	require.Equal(t, 0, len(smartMdl.ScsiAttributes))
}

func TestFromCollectorSmartInfo_Fail2(t *testing.T) {
	//setup
	smartDataFile, err := os.Open("../testdata/smart-fail2.json")
	require.NoError(t, err)
	defer smartDataFile.Close()

	var smartJson collector.SmartInfo

	smartDataBytes, err := ioutil.ReadAll(smartDataFile)
	require.NoError(t, err)
	err = json.Unmarshal(smartDataBytes, &smartJson)
	require.NoError(t, err)

	//test
	smartMdl := db.Smart{}
	err = smartMdl.FromCollectorSmartInfo("WWN-test", smartJson)

	//assert
	require.NoError(t, err)
	require.Equal(t, "WWN-test", smartMdl.DeviceWWN)
	require.Equal(t, "failed", smartMdl.SmartStatus)
	require.Equal(t, 17, len(smartMdl.AtaAttributes))
	require.Equal(t, 0, len(smartMdl.NvmeAttributes))
	require.Equal(t, 0, len(smartMdl.ScsiAttributes))
}

func TestFromCollectorSmartInfo_Nvme(t *testing.T) {
	//setup
	smartDataFile, err := os.Open("../testdata/smart-nvme.json")
	require.NoError(t, err)
	defer smartDataFile.Close()

	var smartJson collector.SmartInfo

	smartDataBytes, err := ioutil.ReadAll(smartDataFile)
	require.NoError(t, err)
	err = json.Unmarshal(smartDataBytes, &smartJson)
	require.NoError(t, err)

	//test
	smartMdl := db.Smart{}
	err = smartMdl.FromCollectorSmartInfo("WWN-test", smartJson)

	//assert
	require.NoError(t, err)
	require.Equal(t, "WWN-test", smartMdl.DeviceWWN)
	require.Equal(t, "passed", smartMdl.SmartStatus)
	require.Equal(t, 0, len(smartMdl.AtaAttributes))
	require.Equal(t, 16, len(smartMdl.NvmeAttributes))
	require.Equal(t, 0, len(smartMdl.ScsiAttributes))

	require.Equal(t, 111303174, smartMdl.NvmeAttributes[6].Value)
	require.Equal(t, 83170961, smartMdl.NvmeAttributes[7].Value)
}

func TestFromCollectorSmartInfo_Scsi(t *testing.T) {
	//setup
	smartDataFile, err := os.Open("../testdata/smart-scsi.json")
	require.NoError(t, err)
	defer smartDataFile.Close()

	var smartJson collector.SmartInfo

	smartDataBytes, err := ioutil.ReadAll(smartDataFile)
	require.NoError(t, err)
	err = json.Unmarshal(smartDataBytes, &smartJson)
	require.NoError(t, err)

	//test
	smartMdl := db.Smart{}
	err = smartMdl.FromCollectorSmartInfo("WWN-test", smartJson)

	//assert
	require.NoError(t, err)
	require.Equal(t, "WWN-test", smartMdl.DeviceWWN)
	require.Equal(t, "passed", smartMdl.SmartStatus)
	require.Equal(t, 0, len(smartMdl.AtaAttributes))
	require.Equal(t, 0, len(smartMdl.NvmeAttributes))
	require.Equal(t, 13, len(smartMdl.ScsiAttributes))

	require.Equal(t, 56, smartMdl.ScsiAttributes[0].Value)
	require.Equal(t, 300357663, smartMdl.ScsiAttributes[4].Value) //total_errors_corrected
}
