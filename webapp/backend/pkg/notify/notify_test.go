package notify

import (
	"fmt"
	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestShouldNotify_MustSkipPassingDevices(t *testing.T) {
	t.Parallel()
	//setup
	device := models.Device{
		DeviceStatus: pkg.DeviceStatusPassed,
	}
	smartAttrs := measurements.Smart{}
	statusThreshold := pkg.MetricsStatusThresholdBoth
	notifyFilterAttributes := pkg.MetricsStatusFilterAttributesAll

	//assert
	require.False(t, ShouldNotify(device, smartAttrs, statusThreshold, notifyFilterAttributes))
}

func TestShouldNotify_MetricsStatusThresholdBoth_FailingSmartDevice(t *testing.T) {
	t.Parallel()
	//setup
	device := models.Device{
		DeviceStatus: pkg.DeviceStatusFailedSmart,
	}
	smartAttrs := measurements.Smart{}
	statusThreshold := pkg.MetricsStatusThresholdBoth
	notifyFilterAttributes := pkg.MetricsStatusFilterAttributesAll

	//assert
	require.True(t, ShouldNotify(device, smartAttrs, statusThreshold, notifyFilterAttributes))
}

func TestShouldNotify_MetricsStatusThresholdSmart_FailingSmartDevice(t *testing.T) {
	t.Parallel()
	//setup
	device := models.Device{
		DeviceStatus: pkg.DeviceStatusFailedSmart,
	}
	smartAttrs := measurements.Smart{}
	statusThreshold := pkg.MetricsStatusThresholdSmart
	notifyFilterAttributes := pkg.MetricsStatusFilterAttributesAll

	//assert
	require.True(t, ShouldNotify(device, smartAttrs, statusThreshold, notifyFilterAttributes))
}

func TestShouldNotify_MetricsStatusThresholdScrutiny_FailingSmartDevice(t *testing.T) {
	t.Parallel()
	//setup
	device := models.Device{
		DeviceStatus: pkg.DeviceStatusFailedSmart,
	}
	smartAttrs := measurements.Smart{}
	statusThreshold := pkg.MetricsStatusThresholdScrutiny
	notifyFilterAttributes := pkg.MetricsStatusFilterAttributesAll

	//assert
	require.False(t, ShouldNotify(device, smartAttrs, statusThreshold, notifyFilterAttributes))
}

func TestShouldNotify_MetricsStatusFilterAttributesCritical_WithCriticalAttrs(t *testing.T) {
	t.Parallel()
	//setup
	device := models.Device{
		DeviceStatus: pkg.DeviceStatusFailedSmart,
	}
	smartAttrs := measurements.Smart{Attributes: map[string]measurements.SmartAttribute{
		"5": &measurements.SmartAtaAttribute{
			Status: pkg.AttributeStatusFailedSmart,
		},
	}}
	statusThreshold := pkg.MetricsStatusThresholdBoth
	notifyFilterAttributes := pkg.MetricsStatusFilterAttributesCritical

	//assert
	require.True(t, ShouldNotify(device, smartAttrs, statusThreshold, notifyFilterAttributes))
}

func TestShouldNotify_MetricsStatusFilterAttributesCritical_WithMultipleCriticalAttrs(t *testing.T) {
	t.Parallel()
	//setup
	device := models.Device{
		DeviceStatus: pkg.DeviceStatusFailedSmart,
	}
	smartAttrs := measurements.Smart{Attributes: map[string]measurements.SmartAttribute{
		"5": &measurements.SmartAtaAttribute{
			Status: pkg.AttributeStatusPassed,
		},
		"10": &measurements.SmartAtaAttribute{
			Status: pkg.AttributeStatusFailedScrutiny,
		},
	}}
	statusThreshold := pkg.MetricsStatusThresholdBoth
	notifyFilterAttributes := pkg.MetricsStatusFilterAttributesCritical

	//assert
	require.True(t, ShouldNotify(device, smartAttrs, statusThreshold, notifyFilterAttributes))
}

func TestShouldNotify_MetricsStatusFilterAttributesCritical_WithNoCriticalAttrs(t *testing.T) {
	t.Parallel()
	//setup
	device := models.Device{
		DeviceStatus: pkg.DeviceStatusFailedSmart,
	}
	smartAttrs := measurements.Smart{Attributes: map[string]measurements.SmartAttribute{
		"1": &measurements.SmartAtaAttribute{
			Status: pkg.AttributeStatusFailedSmart,
		},
	}}
	statusThreshold := pkg.MetricsStatusThresholdBoth
	notifyFilterAttributes := pkg.MetricsStatusFilterAttributesCritical

	//assert
	require.False(t, ShouldNotify(device, smartAttrs, statusThreshold, notifyFilterAttributes))
}

func TestShouldNotify_MetricsStatusFilterAttributesCritical_WithNoFailingCriticalAttrs(t *testing.T) {
	t.Parallel()
	//setup
	device := models.Device{
		DeviceStatus: pkg.DeviceStatusFailedSmart,
	}
	smartAttrs := measurements.Smart{Attributes: map[string]measurements.SmartAttribute{
		"5": &measurements.SmartAtaAttribute{
			Status: pkg.AttributeStatusPassed,
		},
	}}
	statusThreshold := pkg.MetricsStatusThresholdBoth
	notifyFilterAttributes := pkg.MetricsStatusFilterAttributesCritical

	//assert
	require.False(t, ShouldNotify(device, smartAttrs, statusThreshold, notifyFilterAttributes))
}

func TestShouldNotify_MetricsStatusFilterAttributesCritical_MetricsStatusThresholdSmart_WithCriticalAttrsFailingScrutiny(t *testing.T) {
	t.Parallel()
	//setup
	device := models.Device{
		DeviceStatus: pkg.DeviceStatusFailedSmart,
	}
	smartAttrs := measurements.Smart{Attributes: map[string]measurements.SmartAttribute{
		"5": &measurements.SmartAtaAttribute{
			Status: pkg.AttributeStatusPassed,
		},
		"10": &measurements.SmartAtaAttribute{
			Status: pkg.AttributeStatusFailedScrutiny,
		},
	}}
	statusThreshold := pkg.MetricsStatusThresholdSmart
	notifyFilterAttributes := pkg.MetricsStatusFilterAttributesCritical

	//assert
	require.False(t, ShouldNotify(device, smartAttrs, statusThreshold, notifyFilterAttributes))
}

func TestNewPayload(t *testing.T) {
	t.Parallel()

	//setup
	device := models.Device{
		SerialNumber: "FAKEWDDJ324KSO",
		DeviceType:   pkg.DeviceProtocolAta,
		DeviceName:   "/dev/sda",
		DeviceStatus: pkg.DeviceStatusFailedScrutiny,
	}
	currentTime := time.Now()
	//test

	payload := NewPayload(device, false, currentTime)

	//assert
	require.Equal(t, "Scrutiny SMART error (ScrutinyFailure) detected on device: /dev/sda", payload.Subject)
	require.Equal(t, fmt.Sprintf(`Scrutiny SMART error notification for device: /dev/sda
Failure Type: ScrutinyFailure
Device Name: /dev/sda
Device Serial: FAKEWDDJ324KSO
Device Type: ATA

Date: %s`, currentTime.Format(time.RFC3339)), payload.Message)
}

func TestNewPayload_TestMode(t *testing.T) {
	t.Parallel()

	//setup
	device := models.Device{
		SerialNumber: "FAKEWDDJ324KSO",
		DeviceType:   pkg.DeviceProtocolAta,
		DeviceName:   "/dev/sda",
		DeviceStatus: pkg.DeviceStatusFailedScrutiny,
	}
	currentTime := time.Now()
	//test

	payload := NewPayload(device, true, currentTime)

	//assert
	require.Equal(t, "Scrutiny SMART error (EmailTest) detected on device: /dev/sda", payload.Subject)
	require.Equal(t, fmt.Sprintf(`TEST NOTIFICATION:
Scrutiny SMART error notification for device: /dev/sda
Failure Type: EmailTest
Device Name: /dev/sda
Device Serial: FAKEWDDJ324KSO
Device Type: ATA

Date: %s`, currentTime.Format(time.RFC3339)), payload.Message)
}

func TestNewPayload_WithHostId(t *testing.T) {
	t.Parallel()

	//setup
	device := models.Device{
		SerialNumber: "FAKEWDDJ324KSO",
		DeviceType:   pkg.DeviceProtocolAta,
		DeviceName:   "/dev/sda",
		DeviceStatus: pkg.DeviceStatusFailedScrutiny,
		HostId:       "custom-host",
	}
	currentTime := time.Now()
	//test

	payload := NewPayload(device, false, currentTime)

	//assert
	require.Equal(t, "Scrutiny SMART error (ScrutinyFailure) detected on [host]device: [custom-host]/dev/sda", payload.Subject)
	require.Equal(t, fmt.Sprintf(`Scrutiny SMART error notification for device: /dev/sda
Host Id: custom-host
Failure Type: ScrutinyFailure
Device Name: /dev/sda
Device Serial: FAKEWDDJ324KSO
Device Type: ATA

Date: %s`, currentTime.Format(time.RFC3339)), payload.Message)
}
