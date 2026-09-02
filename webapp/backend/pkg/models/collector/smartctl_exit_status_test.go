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
		{"smart command failed", 4, false},
		{"disk failing", 8, false},
		{"prefail attribute below threshold", 16, false},
		{"attribute below threshold in the past", 32, false},
		{"errors in the error log", 64, false},
		{"errors in the self test log", 128, false},
		{"failing disk with errors in both logs", 8 | 64 | 128, false},
		{"device open failed while also reporting a failing disk", 2 | 8, true},
		//real megaraid/SAT output exits with FAILSMART while still reporting every attribute
		{"smart command failed on a disk in pre-fail", 4 | 16, false},
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

func TestSmartInfo_IsLowPowerExit(t *testing.T) {
	lowPowerMessage := struct {
		String   string `json:"string"`
		Severity string `json:"severity"`
	}{String: "Device is in STANDBY mode, exit(2)", Severity: "information"}
	openFailedMessage := lowPowerMessage
	openFailedMessage.String = "Smartctl open device: /dev/sdz failed: No such device"

	t.Run("device asleep", func(t *testing.T) {
		var smartInfo SmartInfo
		smartInfo.Smartctl.ExitStatus = SmartctlExitStatusFailDev
		smartInfo.Smartctl.Messages = append(smartInfo.Smartctl.Messages, lowPowerMessage)
		assert.True(t, smartInfo.IsLowPowerExit())
	})

	t.Run("device gone", func(t *testing.T) {
		var smartInfo SmartInfo
		smartInfo.Smartctl.ExitStatus = SmartctlExitStatusFailDev
		smartInfo.Smartctl.Messages = append(smartInfo.Smartctl.Messages, openFailedMessage)
		assert.False(t, smartInfo.IsLowPowerExit())
	})

	t.Run("no messages at all", func(t *testing.T) {
		var smartInfo SmartInfo
		smartInfo.Smartctl.ExitStatus = SmartctlExitStatusFailDev
		assert.False(t, smartInfo.IsLowPowerExit())
	})

	t.Run("asleep but also unreadable", func(t *testing.T) {
		var smartInfo SmartInfo
		smartInfo.Smartctl.ExitStatus = SmartctlExitStatusFailDev | SmartctlExitStatusFailCmd
		smartInfo.Smartctl.Messages = append(smartInfo.Smartctl.Messages, lowPowerMessage)
		assert.False(t, smartInfo.IsLowPowerExit())
	})
}
