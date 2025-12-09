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

func TestSmartInfo_LargeLBAValues(t *testing.T) {
	// Test for GitHub issue #24 / upstream issue #800
	// LBA values can be large unsigned 64-bit integers that overflow signed int
	t.Run("should parse large LBA values in selective self-test log", func(t *testing.T) {
		// This JSON contains LBA values that exceed int64 max (9223372036854775807)
		// Value 18446743534724713985 is a valid uint64 but overflows int64
		jsonData := `{
			"ata_smart_selective_self_test_log": {
				"revision": 1,
				"table": [
					{
						"lba_min": 18446743534724713985,
						"lba_max": 7205816247684983039,
						"status": {
							"value": 0,
							"string": "Not_testing"
						}
					}
				],
				"flags": {
					"value": 0,
					"remainder_scan_enabled": false
				},
				"power_up_scan_resume_minutes": 0
			}
		}`

		var smartInfo SmartInfo
		err := json.Unmarshal([]byte(jsonData), &smartInfo)
		require.NoError(t, err, "should unmarshal large LBA values without error")

		// Verify the values were parsed correctly
		require.Len(t, smartInfo.AtaSmartSelectiveSelfTestLog.Table, 1)
		assert.Equal(t, uint64(18446743534724713985), smartInfo.AtaSmartSelectiveSelfTestLog.Table[0].LbaMin)
		assert.Equal(t, uint64(7205816247684983039), smartInfo.AtaSmartSelectiveSelfTestLog.Table[0].LbaMax)
	})

	t.Run("should parse large LBA values in error log", func(t *testing.T) {
		// LBA values in error logs can also be large
		jsonData := `{
			"ata_smart_error_log": {
				"summary": {
					"revision": 1,
					"count": 1,
					"logged_count": 1,
					"table": [
						{
							"error_number": 1,
							"lifetime_hours": 1000,
							"completion_registers": {
								"error": 0,
								"status": 0,
								"count": 0,
								"lba": 18446744073709551615,
								"device": 0
							},
							"error_description": "test",
							"previous_commands": [
								{
									"registers": {
										"command": 0,
										"features": 0,
										"count": 0,
										"lba": 18446744073709551615,
										"device": 0,
										"device_control": 0
									},
									"powerup_milliseconds": 0,
									"command_name": "test"
								}
							]
						}
					]
				}
			}
		}`

		var smartInfo SmartInfo
		err := json.Unmarshal([]byte(jsonData), &smartInfo)
		require.NoError(t, err, "should unmarshal large LBA values in error log without error")

		// Verify the values were parsed correctly
		require.Len(t, smartInfo.AtaSmartErrorLog.Summary.Table, 1)
		assert.Equal(t, uint64(18446744073709551615), smartInfo.AtaSmartErrorLog.Summary.Table[0].CompletionRegisters.Lba)
		require.Len(t, smartInfo.AtaSmartErrorLog.Summary.Table[0].PreviousCommands, 1)
		assert.Equal(t, uint64(18446744073709551615), smartInfo.AtaSmartErrorLog.Summary.Table[0].PreviousCommands[0].Registers.Lba)
	})
}
