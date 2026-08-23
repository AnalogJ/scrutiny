package collector

import (
	"errors"
	"net/url"
	"testing"

	mock_shell "github.com/analogj/scrutiny/collector/pkg/common/shell/mock"
	mock_config "github.com/analogj/scrutiny/collector/pkg/config/mock"
	"github.com/analogj/scrutiny/collector/pkg/detect"
	"github.com/analogj/scrutiny/collector/pkg/models"
	"github.com/gofrs/uuid/v5"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestApiEndpointParse(t *testing.T) {
	baseURL, _ := url.Parse("http://localhost:8080/")

	url1, _ := baseURL.Parse("d/e")
	require.Equal(t, "http://localhost:8080/d/e", url1.String())

	url2, _ := baseURL.Parse("/d/e")
	require.Equal(t, "http://localhost:8080/d/e", url2.String())
}

func TestApiEndpointParse_WithBasepathWithoutTrailingSlash(t *testing.T) {
	baseURL, _ := url.Parse("http://localhost:8080/scrutiny")

	//This testcase is unexpected and can cause issues. We need to ensure the apiEndpoint always has a trailing slash.
	url1, _ := baseURL.Parse("d/e")
	require.Equal(t, "http://localhost:8080/d/e", url1.String())

	url2, _ := baseURL.Parse("/d/e")
	require.Equal(t, "http://localhost:8080/d/e", url2.String())
}

func TestApiEndpointParse_WithBasepathWithTrailingSlash(t *testing.T) {
	baseURL, _ := url.Parse("http://localhost:8080/scrutiny/")

	url1, _ := baseURL.Parse("d/e")
	require.Equal(t, "http://localhost:8080/scrutiny/d/e", url1.String())

	url2, _ := baseURL.Parse("/d/e")
	require.Equal(t, "http://localhost:8080/d/e", url2.String())
}

func TestMetricsCollector_Collect_DeviceType(t *testing.T) {
	for _, tt := range []struct {
		name         string
		deviceType   string
		overrides    []models.ScanOverride
		expectedArgs []string
	}{
		{
			name:         "should drop a scanned scsi type",
			deviceType:   "scsi",
			overrides:    []models.ScanOverride{},
			expectedArgs: []string{"--xall", "--json"},
		},
		{
			name:       "should keep a configured scsi type",
			deviceType: "scsi",
			overrides: []models.ScanOverride{
				{Device: detect.DevicePrefix() + "sda", DeviceType: []string{"scsi"}},
			},
			expectedArgs: []string{"--xall", "--json", "--device", "scsi"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			const (
				someArgs       = "--xall --json"
				someDeviceName = "sda"
			)

			fullDeviceName := detect.DevicePrefix() + someDeviceName

			fakeConfig := mock_config.NewMockInterface(ctrl)
			fakeConfig.EXPECT().
				GetCommandMetricsSmartArgs(fullDeviceName).
				Return(someArgs)
			fakeConfig.EXPECT().
				GetString("commands.metrics_smartctl_bin").
				Return("smartctl")
			fakeConfig.EXPECT().
				GetDeviceOverrides().
				Return(tt.overrides).
				AnyTimes()

			someLogger := logrus.WithFields(logrus.Fields{})

			fakeShell := mock_shell.NewMockInterface(ctrl)
			fakeShell.EXPECT().
				Command(someLogger, "smartctl", append(tt.expectedArgs, fullDeviceName), "", gomock.Any()).
				Return("", errors.New("smartctl unavailable in test"))

			apiEndpoint, err := url.Parse("http://localhost:8080/")
			require.NoError(t, err)

			mc := MetricsCollector{
				config:        fakeConfig,
				BaseCollector: BaseCollector{logger: someLogger},
				apiEndpoint:   apiEndpoint,
				shell:         fakeShell,
			}

			mc.Collect(uuid.Must(uuid.NewV4()), someDeviceName, tt.deviceType)
		})
	}
}
