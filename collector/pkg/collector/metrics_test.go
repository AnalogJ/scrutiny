package collector

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os/exec"
	"testing"

	mock_shell "github.com/analogj/scrutiny/collector/pkg/common/shell/mock"
	mock_config "github.com/analogj/scrutiny/collector/pkg/config/mock"
	"github.com/analogj/scrutiny/collector/pkg/detect"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/gofrs/uuid/v5"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// helperCollectorErrorApi stands in for the webapp, capturing the collector errors reported to it.
func helperCollectorErrorApi(t *testing.T) (*url.URL, chan collector.CollectorError) {
	reported := make(chan collector.CollectorError, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var errorPayload collector.CollectorError
		require.NoError(t, json.NewDecoder(r.Body).Decode(&errorPayload))
		reported <- errorPayload
		_, _ = w.Write([]byte(`{"success": true}`))
	}))
	t.Cleanup(server.Close)

	apiEndpoint, err := url.Parse(server.URL + "/")
	require.NoError(t, err)

	return apiEndpoint, reported
}

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

func TestMetricsCollector_Collect_ConfiguredDeviceType(t *testing.T) {
	for _, tt := range []struct {
		name                  string
		deviceType            string
		hasDeviceTypeOverride bool
		expectedArgs          []string
	}{
		{"should drop a scanned scsi type", "scsi", false, []string{"--xall", "--json"}},
		{"should keep a configured scsi type", "scsi", true, []string{"--xall", "--json", "--device", "scsi"}},
		{"should keep a scanned non-standard type", "megaraid,14", false, []string{"--xall", "--json", "--device", "megaraid,14"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			const someDeviceName = "sda"

			fullDeviceName := detect.DevicePrefix() + someDeviceName

			fakeConfig := mock_config.NewMockInterface(ctrl)
			fakeConfig.EXPECT().GetCommandMetricsSmartArgs(fullDeviceName).Return("--xall --json")
			fakeConfig.EXPECT().GetString("commands.metrics_smartctl_bin").Return("smartctl")
			//only consulted for a scsi/ata type; the argv below is what the test pins
			fakeConfig.EXPECT().HasDeviceTypeOverride(fullDeviceName).AnyTimes().Return(tt.hasDeviceTypeOverride)
			fakeConfig.EXPECT().GetString("host.id").AnyTimes().Return("testhost")

			someLogger := logrus.WithFields(logrus.Fields{})

			fakeShell := mock_shell.NewMockInterface(ctrl)
			fakeShell.EXPECT().
				Command(someLogger, "smartctl", append(tt.expectedArgs, fullDeviceName), "", gomock.Any()).
				Return("", errors.New("smartctl is not available in tests"))

			apiEndpoint, reported := helperCollectorErrorApi(t)

			mc := MetricsCollector{
				config:        fakeConfig,
				BaseCollector: BaseCollector{logger: someLogger},
				apiEndpoint:   apiEndpoint,
				shell:         fakeShell,
			}

			mc.Collect(uuid.Must(uuid.NewV4()), someDeviceName, tt.deviceType)

			//smartctl could not be run at all, so the failure is reported rather than swallowed
			require.Equal(t, someDeviceName, (<-reported).DeviceName)
		})
	}
}

func TestMetricsCollector_Collect_SignalledSmartctl(t *testing.T) {
	ctrl := gomock.NewController(t)

	const someDeviceName = "sda"
	fullDeviceName := detect.DevicePrefix() + someDeviceName

	fakeConfig := mock_config.NewMockInterface(ctrl)
	fakeConfig.EXPECT().GetCommandMetricsSmartArgs(fullDeviceName).Return("--xall --json")
	fakeConfig.EXPECT().GetString("commands.metrics_smartctl_bin").Return("smartctl")
	fakeConfig.EXPECT().HasDeviceTypeOverride(fullDeviceName).AnyTimes().Return(false)
	fakeConfig.EXPECT().GetString("host.id").AnyTimes().Return("testhost")

	someLogger := logrus.WithFields(logrus.Fields{})

	//a process that never produced an exit status - killed by a signal, say - reports exit code -1,
	//which would otherwise decode as every smartctl failure bit being set
	signalled := &exec.ExitError{}
	require.Equal(t, -1, signalled.ExitCode())

	fakeShell := mock_shell.NewMockInterface(ctrl)
	fakeShell.EXPECT().
		Command(someLogger, "smartctl", []string{"--xall", "--json", fullDeviceName}, "", gomock.Any()).
		Return("", signalled)

	apiEndpoint, reported := helperCollectorErrorApi(t)

	mc := MetricsCollector{
		config:        fakeConfig,
		BaseCollector: BaseCollector{logger: someLogger},
		apiEndpoint:   apiEndpoint,
		shell:         fakeShell,
	}

	mc.Collect(uuid.Must(uuid.NewV4()), someDeviceName, "")

	require.Equal(t, someDeviceName, (<-reported).DeviceName)
}
