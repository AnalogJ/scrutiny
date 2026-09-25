package measurements_test

import (
	"testing"

	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
	"github.com/stretchr/testify/require"
)

func TestSmartAtaAttribute_GetComparableValue(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		attr     measurements.SmartAtaAttribute
		expected int64
	}{
		{
			name:     "normalized display type uses Value",
			attr:     measurements.SmartAtaAttribute{AttributeId: 1, Value: 100, RawValue: 12345, TransformedValue: 7},
			expected: 100,
		},
		{
			name:     "raw display type uses RawValue",
			attr:     measurements.SmartAtaAttribute{AttributeId: 199, Value: 99, RawValue: 108},
			expected: 108,
		},
		{
			name:     "transformed display type uses TransformedValue",
			attr:     measurements.SmartAtaAttribute{AttributeId: 194, Value: 67, RawValue: 0x1A0021, TransformedValue: 33},
			expected: 33,
		},
		{
			name:     "attribute without metadata uses RawValue",
			attr:     measurements.SmartAtaAttribute{AttributeId: 999, Value: 100, RawValue: 42},
			expected: 42,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, tc.attr.GetComparableValue())
		})
	}
}

func TestSmartNvmeAttribute_GetComparableValue(t *testing.T) {
	t.Parallel()
	attr := measurements.SmartNvmeAttribute{AttributeId: "media_errors", Value: 5}
	require.Equal(t, int64(5), attr.GetComparableValue())
}

func TestSmartScsiAttribute_GetComparableValue(t *testing.T) {
	t.Parallel()
	attr := measurements.SmartScsiAttribute{AttributeId: "scsi_grown_defect_list", Value: 12}
	require.Equal(t, int64(12), attr.GetComparableValue())
}
