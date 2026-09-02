package collector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsLowPowerExit(t *testing.T) {
	for _, tt := range []struct {
		name     string
		result   string
		expected bool
	}{
		{
			"device asleep",
			`{"smartctl":{"exit_status":2,"messages":[{"string":"Device is in STANDBY mode, exit(2)","severity":"information"}]}}`,
			true,
		},
		{
			//a device removed from a config `devices:` override still reaches Collect, and shares
			//the exit status with a sleeping one
			"device gone",
			`{"smartctl":{"exit_status":2,"messages":[{"string":"Smartctl open device: /dev/sdz failed: No such device","severity":"error"}]}}`,
			false,
		},
		{"no output at all", "", false},
		{"not json", "smartctl: command not found", false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isLowPowerExit([]byte(tt.result)))
		})
	}
}
