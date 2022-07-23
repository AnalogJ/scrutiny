package notify

import (
	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
	"github.com/stretchr/testify/require"
	"testing"
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
