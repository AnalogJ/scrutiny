package collector

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasPowerModeCheck(t *testing.T) {
	for _, tt := range []struct {
		args     string
		expected bool
	}{
		{"--xall --json", false},
		{"-a --json /dev/sda", false},
		{"-n standby --xall --json", true},
		{"-nstandby --xall --json", true},
		{"-n=standby --xall --json", true},
		{"--nocheck=standby --xall --json", true},
		{"--nocheck standby --xall --json", true},
	} {
		t.Run(tt.args, func(t *testing.T) {
			assert.Equal(t, tt.expected, hasPowerModeCheck(strings.Split(tt.args, " ")))
		})
	}
}
