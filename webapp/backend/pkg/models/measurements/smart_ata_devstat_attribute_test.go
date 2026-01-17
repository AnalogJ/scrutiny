package measurements

import (
	"testing"

	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/stretchr/testify/require"
)

func TestSmartAtaDeviceStatAttribute_Flatten(t *testing.T) {
	// Test that device statistics are flattened with string-based attribute IDs
	attr := SmartAtaDeviceStatAttribute{
		AttributeId:      "devstat_7_8",
		Value:            19,
		Threshold:        100,
		TransformedValue: 19,
		Status:           pkg.AttributeStatusPassed,
	}

	flattened := attr.Flatten()

	require.Equal(t, "devstat_7_8", flattened["attr.devstat_7_8.attribute_id"])
	require.Equal(t, int64(19), flattened["attr.devstat_7_8.value"])
	require.Equal(t, int64(100), flattened["attr.devstat_7_8.thresh"])
	require.Equal(t, int64(19), flattened["attr.devstat_7_8.transformed_value"])
	require.Equal(t, int64(pkg.AttributeStatusPassed), flattened["attr.devstat_7_8.status"])
}

func TestSmartAtaDeviceStatAttribute_Inflate(t *testing.T) {
	// Test that device statistics can be inflated from InfluxDB data
	attr := SmartAtaDeviceStatAttribute{}

	attr.Inflate("attr.devstat_7_8.attribute_id", "devstat_7_8")
	attr.Inflate("attr.devstat_7_8.value", int64(25))
	attr.Inflate("attr.devstat_7_8.thresh", int64(100))
	attr.Inflate("attr.devstat_7_8.transformed_value", int64(25))
	attr.Inflate("attr.devstat_7_8.status", int64(pkg.AttributeStatusPassed))
	attr.Inflate("attr.devstat_7_8.status_reason", "")
	attr.Inflate("attr.devstat_7_8.failure_rate", float64(0))

	require.Equal(t, "devstat_7_8", attr.AttributeId)
	require.Equal(t, int64(25), attr.Value)
	require.Equal(t, int64(100), attr.Threshold)
	require.Equal(t, int64(25), attr.TransformedValue)
	require.Equal(t, pkg.AttributeStatusPassed, attr.Status)
}

func TestSmartAtaDeviceStatAttribute_FlattenInflateRoundtrip(t *testing.T) {
	// Test that flatten/inflate roundtrip preserves data
	original := SmartAtaDeviceStatAttribute{
		AttributeId:      "devstat_7_8",
		Value:            42,
		Threshold:        100,
		TransformedValue: 42,
		Status:           pkg.AttributeStatusWarningScrutiny,
		StatusReason:     "Test warning",
		FailureRate:      0.5,
	}

	flattened := original.Flatten()

	restored := SmartAtaDeviceStatAttribute{}
	for key, val := range flattened {
		restored.Inflate(key, val)
	}

	require.Equal(t, original.AttributeId, restored.AttributeId)
	require.Equal(t, original.Value, restored.Value)
	require.Equal(t, original.Threshold, restored.Threshold)
	require.Equal(t, original.TransformedValue, restored.TransformedValue)
	require.Equal(t, original.Status, restored.Status)
	require.Equal(t, original.StatusReason, restored.StatusReason)
	require.Equal(t, original.FailureRate, restored.FailureRate)
}

func TestSmartAtaDeviceStatAttribute_PopulateAttributeStatus_BelowThreshold(t *testing.T) {
	// Test that percentage used below threshold passes
	attr := SmartAtaDeviceStatAttribute{
		AttributeId: "devstat_7_8", // Percentage Used Endurance Indicator
		Value:       19,            // 19% used
		Threshold:   100,
	}

	attr.PopulateAttributeStatus()

	require.Equal(t, pkg.AttributeStatusPassed, attr.Status)
	require.Equal(t, int64(19), attr.TransformedValue)
}

func TestSmartAtaDeviceStatAttribute_PopulateAttributeStatus_AtThreshold(t *testing.T) {
	// Test that percentage used at threshold fails
	attr := SmartAtaDeviceStatAttribute{
		AttributeId: "devstat_7_8", // Percentage Used Endurance Indicator
		Value:       100,           // 100% used - device end of life
		Threshold:   100,
	}

	attr.PopulateAttributeStatus()

	require.True(t, pkg.AttributeStatusHas(attr.Status, pkg.AttributeStatusFailedScrutiny))
	require.NotEmpty(t, attr.StatusReason)
}

func TestSmartAtaDeviceStatAttribute_PopulateAttributeStatus_AboveThreshold(t *testing.T) {
	// Test that percentage used above threshold fails
	attr := SmartAtaDeviceStatAttribute{
		AttributeId: "devstat_7_8", // Percentage Used Endurance Indicator
		Value:       150,           // 150% used - past end of life
		Threshold:   100,
	}

	attr.PopulateAttributeStatus()

	require.True(t, pkg.AttributeStatusHas(attr.Status, pkg.AttributeStatusFailedScrutiny))
}

func TestSmartAtaDeviceStatAttribute_PopulateAttributeStatus_UnknownAttribute(t *testing.T) {
	// Test that unknown device statistics don't cause errors
	attr := SmartAtaDeviceStatAttribute{
		AttributeId: "devstat_99_99", // Unknown device statistic
		Value:       42,
	}

	attr.PopulateAttributeStatus()

	// Should pass since we don't have metadata for this attribute
	require.Equal(t, pkg.AttributeStatusPassed, attr.Status)
	require.Equal(t, int64(42), attr.TransformedValue)
}

func TestSmartAtaDeviceStatAttribute_GetTransformedValue(t *testing.T) {
	attr := SmartAtaDeviceStatAttribute{
		TransformedValue: 123,
	}
	require.Equal(t, int64(123), attr.GetTransformedValue())
}

func TestSmartAtaDeviceStatAttribute_GetStatus(t *testing.T) {
	attr := SmartAtaDeviceStatAttribute{
		Status: pkg.AttributeStatusWarningScrutiny,
	}
	require.Equal(t, pkg.AttributeStatusWarningScrutiny, attr.GetStatus())
}
