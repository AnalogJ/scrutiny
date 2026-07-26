package collector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSmartctlExitStatus_IsFatal(t *testing.T) {
	for _, tt := range []struct {
		name     string
		exitCode int
		fatal    bool
	}{
		{"success", 0, false},
		{"commandline could not be parsed", 1, true},
		{"device open failed", 2, true},
		{"smart command failed", 4, true},
		{"disk failing", 8, false},
		{"prefail attribute below threshold", 16, false},
		{"attribute below threshold in the past", 32, false},
		{"errors in the error log", 64, false},
		{"errors in the self test log", 128, false},
		{"failing disk with errors in both logs", 8 | 64 | 128, false},
		{"device open failed while also reporting a failing disk", 2 | 8, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.fatal, SmartctlExitStatus(tt.exitCode).IsFatal())
		})
	}
}

func TestSmartctlExitStatus_Descriptions(t *testing.T) {
	assert.Empty(t, SmartctlExitStatus(0).Descriptions())

	// every set bit is reported, not just the lowest one
	assert.Len(t, SmartctlExitStatus(8|16|32).Descriptions(), 3)
}
