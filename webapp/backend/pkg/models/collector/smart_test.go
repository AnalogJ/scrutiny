package collector

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSmartInfo_Capacity(t *testing.T) {
	t.Run("should report nvme capacity", func(t *testing.T) {
		smartInfo := SmartInfo{
			UserCapacity: UserCapacity{
				Bytes: 1234,
			},
			NvmeTotalCapacity: 5678,
		}
		assert.Equal(t, int64(5678), smartInfo.Capacity())
	})

	t.Run("should report user capacity", func(t *testing.T) {
		smartInfo := SmartInfo{
			UserCapacity: UserCapacity{
				Bytes: 1234,
			},
		}
		assert.Equal(t, int64(1234), smartInfo.Capacity())
	})

	t.Run("should report 0 for unknown capacities", func(t *testing.T) {
		var smartInfo SmartInfo
		assert.Zero(t, smartInfo.Capacity())
	})
}

func TestSmartSupport_UnmarshalJSON(t *testing.T) {
	for _, tt := range []struct {
		name      string
		input     string
		available bool
		enabled   bool
		supported bool
	}{
		{
			name:      "should report legacy boolean support",
			input:     `true`,
			available: true,
			enabled:   true,
			supported: true,
		},
		{
			name:      "should report legacy boolean unsupported",
			input:     `false`,
			available: false,
			enabled:   false,
			supported: false,
		},
		{
			name:      "should report object support",
			input:     `{"available":true,"enabled":true}`,
			available: true,
			enabled:   true,
			supported: true,
		},
		{
			name:      "should require available and enabled object fields",
			input:     `{"available":true,"enabled":false}`,
			available: true,
			enabled:   false,
			supported: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var support SmartSupport
			require.NoError(t, json.Unmarshal([]byte(tt.input), &support))

			assert.Equal(t, tt.available, support.Available)
			assert.Equal(t, tt.enabled, support.Enabled)
			assert.Equal(t, tt.supported, support.Supported())
		})
	}
}

func TestSmartSupport_UnmarshalJSONInvalid(t *testing.T) {
	var support SmartSupport
	require.Error(t, json.Unmarshal([]byte(`"unsupported"`), &support))
}

func TestSmartInfo_HasInvalidData(t *testing.T) {
	newSmartInfo := func(exitStatus int) SmartInfo {
		var smartInfo SmartInfo
		smartInfo.Smartctl.ExitStatus = exitStatus
		return smartInfo
	}

	for _, tt := range []struct {
		name       string
		exitStatus int
		invalid    bool
	}{
		{name: "clean run", exitStatus: 0x00, invalid: false},
		{name: "commandline did not parse (bit 0)", exitStatus: 0x01, invalid: true},
		{name: "device open failed / standby (bit 1)", exitStatus: 0x02, invalid: true},
		{name: "checksum or command error (bit 2) keeps data", exitStatus: 0x04, invalid: false},
		{name: "failing disk (bit 3) keeps data", exitStatus: 0x08, invalid: false},
		{name: "prefail attributes (bit 4) keeps data", exitStatus: 0x10, invalid: false},
		{name: "health findings 0xD8 keeps data", exitStatus: 0xD8, invalid: false},
		{name: "open failed alongside health findings", exitStatus: 0x02 | 0x08, invalid: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			smartInfo := newSmartInfo(tt.exitStatus)
			assert.Equal(t, tt.invalid, smartInfo.HasInvalidData())
		})
	}
}

func TestSmartInfo_HasInvalidData_Fixtures(t *testing.T) {
	for _, tt := range []struct {
		fixture string
		invalid bool
	}{
		// clean read of a healthy drive
		{fixture: "smart-ata.json", invalid: false},
		// smartctl could not open the device (exit_status 2); payload is empty
		{fixture: "smart-fail.json", invalid: true},
		// genuinely failing drive with a complete payload (exit_status 216)
		{fixture: "smart-fail2.json", invalid: false},
		// megaraid passthrough sets bit 2 (exit_status 4) but returns valid data
		{fixture: "smart-megaraid0.json", invalid: false},
	} {
		t.Run(tt.fixture, func(t *testing.T) {
			contents, err := os.ReadFile("../testdata/" + tt.fixture)
			require.NoError(t, err)

			var smartInfo SmartInfo
			require.NoError(t, json.Unmarshal(contents, &smartInfo))

			assert.Equal(t, tt.invalid, smartInfo.HasInvalidData())
		})
	}
}
