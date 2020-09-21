package detect_test

import (
	"github.com/analogj/scrutiny/collector/pkg/detect"
	"testing"
)

func TestDevicePrefix(t *testing.T) {
	//setup

	//test

	//assert
	require.Equal(t, "/dev/", detect.DevicePrefix())
}
