package collector

import (
	"encoding/json"
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

func TestSmartInfo_CurrentTemperature(t *testing.T) {
	for _, tt := range []struct {
		name     string
		input    string
		expected int64
	}{
		{
			name:     "top-level temperature only",
			input:    `{"temperature": {"current": 28, "drive_trip": 65}}`,
			expected: 28,
		},
		{
			name:     "scsi environmental temperature only",
			input:    `{"scsi_environmental_reports": {"temperature_1": {"parameter_code": 0, "current": 25, "lifetime_maximum": 30}}}`,
			expected: 25,
		},
		{
			name:     "scsi environmental temperature with null top-level temperature",
			input:    `{"temperature": null, "scsi_environmental_reports": {"temperature_1": {"parameter_code": 0, "current": 25}}}`,
			expected: 25,
		},
		{
			name:     "both available prefers top-level temperature",
			input:    `{"temperature": {"current": 28}, "scsi_environmental_reports": {"temperature_1": {"parameter_code": 0, "current": 25}}}`,
			expected: 28,
		},
		{
			name:     "invalid top-level temperature falls back to scsi environmental temperature",
			input:    `{"temperature": {"current": 0}, "scsi_environmental_reports": {"temperature_1": {"parameter_code": 0, "current": 25}}}`,
			expected: 25,
		},
		{
			name:     "neither available",
			input:    `{}`,
			expected: 0,
		},
		{
			name:     "scsi environmental report without current temperature",
			input:    `{"scsi_environmental_reports": {"temperature_1": {"parameter_code": 0, "lifetime_maximum": 30}}}`,
			expected: 0,
		},
		{
			name:     "scsi environmental report without temperature_1",
			input:    `{"scsi_environmental_reports": {"humidity_1": {"parameter_code": 256, "current": 40}}}`,
			expected: 0,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var smartInfo SmartInfo
			require.NoError(t, json.Unmarshal([]byte(tt.input), &smartInfo))
			assert.Equal(t, tt.expected, smartInfo.CurrentTemperature())
		})
	}
}
