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
	require.Equal(t, smartMdl.DeviceWWN, "WWN-test")
	require.Equal(t, smartMdl.SmartStatus, "passed")

	//check that temperature was correctly parsed
	for _, attr := range smartMdl.SmartAttributes {
		if attr.AttributeId == 194 {
			require.Equal(t, int64(163210330144), attr.RawValue)
			require.Equal(t, int64(32), attr.TransformedValue)
		}
	}
}
