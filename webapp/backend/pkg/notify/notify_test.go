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
	notifyLevel := pkg.NotifyLevelFail
	notifyFilterAttributes := pkg.NotifyFilterAttributesAll

	//assert
	require.False(t, ShouldNotify(device, smartAttrs, notifyLevel, notifyFilterAttributes))
}

func TestShouldNotify_NotifyLevelFail_FailingSmartDevice(t *testing.T) {
	t.Parallel()
	//setup
	device := models.Device{
		DeviceStatus: pkg.DeviceStatusFailedSmart,
	}
	smartAttrs := measurements.Smart{}
	notifyLevel := pkg.NotifyLevelFail
	notifyFilterAttributes := pkg.NotifyFilterAttributesAll

	//assert
	require.True(t, ShouldNotify(device, smartAttrs, notifyLevel, notifyFilterAttributes))
}

func TestShouldNotify_NotifyLevelFailSmart_FailingSmartDevice(t *testing.T) {
	t.Parallel()
	//setup
	device := models.Device{
		DeviceStatus: pkg.DeviceStatusFailedSmart,
	}
	smartAttrs := measurements.Smart{}
	notifyLevel := pkg.NotifyLevelFailSmart
	notifyFilterAttributes := pkg.NotifyFilterAttributesAll

	//assert
	require.True(t, ShouldNotify(device, smartAttrs, notifyLevel, notifyFilterAttributes))
}

func TestShouldNotify_NotifyLevelFailScrutiny_FailingSmartDevice(t *testing.T) {
	t.Parallel()
	//setup
	device := models.Device{
		DeviceStatus: pkg.DeviceStatusFailedSmart,
	}
	smartAttrs := measurements.Smart{}
	notifyLevel := pkg.NotifyLevelFailScrutiny
	notifyFilterAttributes := pkg.NotifyFilterAttributesAll

	//assert
	require.False(t, ShouldNotify(device, smartAttrs, notifyLevel, notifyFilterAttributes))
}

func TestShouldNotify_NotifyFilterAttributesCritical_WithCriticalAttrs(t *testing.T) {
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
	notifyLevel := pkg.NotifyLevelFail
	notifyFilterAttributes := pkg.NotifyFilterAttributesCritical

	//assert
	require.True(t, ShouldNotify(device, smartAttrs, notifyLevel, notifyFilterAttributes))
}

func TestShouldNotify_NotifyFilterAttributesCritical_WithMultipleCriticalAttrs(t *testing.T) {
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
	notifyLevel := pkg.NotifyLevelFail
	notifyFilterAttributes := pkg.NotifyFilterAttributesCritical

	//assert
	require.True(t, ShouldNotify(device, smartAttrs, notifyLevel, notifyFilterAttributes))
}

func TestShouldNotify_NotifyFilterAttributesCritical_WithNoCriticalAttrs(t *testing.T) {
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
	notifyLevel := pkg.NotifyLevelFail
	notifyFilterAttributes := pkg.NotifyFilterAttributesCritical

	//assert
	require.False(t, ShouldNotify(device, smartAttrs, notifyLevel, notifyFilterAttributes))
}

func TestShouldNotify_NotifyFilterAttributesCritical_WithNoFailingCriticalAttrs(t *testing.T) {
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
	notifyLevel := pkg.NotifyLevelFail
	notifyFilterAttributes := pkg.NotifyFilterAttributesCritical

	//assert
	require.False(t, ShouldNotify(device, smartAttrs, notifyLevel, notifyFilterAttributes))
}

func TestShouldNotify_NotifyFilterAttributesCritical_NotifyLevelFailSmart_WithCriticalAttrsFailingScrutiny(t *testing.T) {
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
	notifyLevel := pkg.NotifyLevelFailSmart
	notifyFilterAttributes := pkg.NotifyFilterAttributesCritical

	//assert
	require.False(t, ShouldNotify(device, smartAttrs, notifyLevel, notifyFilterAttributes))
}
