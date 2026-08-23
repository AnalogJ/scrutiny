package measurements_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/require"
)

func i64(v int64) *int64 { return &v }

// ATA attribute 199 (UltraDMA CRC Error Count) is DisplayType "raw", Ideal "low",
// non-critical. A raw value of 116 lands in the observed 70..130 bucket, whose annual
// failure rate is 22.3% -- over the 20% non-critical fail line -- so with no override it
// fails. These are real values from a drive in my fleet whose CRC counter took a single
// step from 0 to 116 and has not moved since.
func TestSmartAtaAttribute_PopulateAttributeStatus_Override(t *testing.T) {
	tests := []struct {
		name     string
		attr     measurements.SmartAtaAttribute
		override *pkg.AttributeOverride
		expected pkg.AttributeStatus
	}{
		{
			name:     "no override, 199 raw 116 still fails on the observed bucket",
			attr:     measurements.SmartAtaAttribute{AttributeId: 199, Value: 200, RawValue: 116},
			override: nil,
			expected: pkg.AttributeStatusFailedScrutiny,
		},
		{
			name:     "fail threshold above the value clears the failure",
			attr:     measurements.SmartAtaAttribute{AttributeId: 199, Value: 200, RawValue: 116},
			override: &pkg.AttributeOverride{FailThreshold: i64(200)},
			expected: pkg.AttributeStatusPassedOverride,
		},
		{
			name:     "value between warn and fail thresholds warns",
			attr:     measurements.SmartAtaAttribute{AttributeId: 199, Value: 200, RawValue: 116},
			override: &pkg.AttributeOverride{WarnThreshold: i64(100), FailThreshold: i64(200)},
			expected: pkg.AttributeStatusWarningScrutiny | pkg.AttributeStatusWarningOverride,
		},
		{
			name:     "value above the fail threshold fails",
			attr:     measurements.SmartAtaAttribute{AttributeId: 199, Value: 200, RawValue: 116},
			override: &pkg.AttributeOverride{FailThreshold: i64(100)},
			expected: pkg.AttributeStatusFailedScrutiny | pkg.AttributeStatusFailedOverride,
		},
		{
			name:     "warn threshold alone, value below it, passes",
			attr:     measurements.SmartAtaAttribute{AttributeId: 199, Value: 200, RawValue: 5},
			override: &pkg.AttributeOverride{WarnThreshold: i64(100)},
			expected: pkg.AttributeStatusPassedOverride,
		},
		{
			// A zero threshold must be honoured rather than read as "unset", which is why
			// the threshold fields are pointers.
			name:     "explicit zero fail threshold is honoured",
			attr:     measurements.SmartAtaAttribute{AttributeId: 199, Value: 200, RawValue: 1},
			override: &pkg.AttributeOverride{FailThreshold: i64(0)},
			expected: pkg.AttributeStatusFailedScrutiny | pkg.AttributeStatusFailedOverride,
		},
		{
			// Attribute 188 is Critical with a first bucket of Low=0,High=100. The bucket
			// test is exclusive on Low, so a *healthy* value of exactly 0 matches nothing
			// and falls through to the "could not determine failure rate" warning. This
			// documents that upstream quirk; the next case shows an override fixing it.
			name:     "no override, healthy 188 of 0 still warns (upstream off-by-one)",
			attr:     measurements.SmartAtaAttribute{AttributeId: 188, Value: 100, RawValue: 0, TransformedValue: 0},
			override: nil,
			expected: pkg.AttributeStatusWarningScrutiny,
		},
		{
			name:     "override replaces the bucket lookup entirely for 188",
			attr:     measurements.SmartAtaAttribute{AttributeId: 188, Value: 100, RawValue: 0, TransformedValue: 0},
			override: &pkg.AttributeOverride{WarnThreshold: i64(50), FailThreshold: i64(100)},
			expected: pkg.AttributeStatusPassedOverride,
		},
		{
			// An override tunes Scrutiny's own heuristics. It must not mask the drive
			// reporting itself as failing -- that is the manufacturer's verdict, not a
			// statistical guess.
			name:     "override does not mask the drive's own SMART failure",
			attr:     measurements.SmartAtaAttribute{AttributeId: 199, Value: 200, RawValue: 116, WhenFailed: pkg.AttributeWhenFailedFailingNow},
			override: &pkg.AttributeOverride{FailThreshold: i64(9999)},
			expected: pkg.AttributeStatusFailedSmart,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			attr := tc.attr
			result := attr.PopulateAttributeStatus(tc.override)
			require.Equal(t, tc.expected, result.Status, "status mismatch (reason: %q)", result.StatusReason)
		})
	}
}

// NVMe available_spare has Ideal "high": a *lower* value is worse. A naive "warn when the
// value exceeds the threshold" override would be backwards here, so the comparison has to
// respect the attribute's Ideal.
func TestSmartNvmeAttribute_PopulateAttributeStatus_Override_IdealHigh(t *testing.T) {
	tests := []struct {
		name     string
		attr     measurements.SmartNvmeAttribute
		override *pkg.AttributeOverride
		expected pkg.AttributeStatus
	}{
		{
			name:     "no override, spare below the drive's own threshold fails",
			attr:     measurements.SmartNvmeAttribute{AttributeId: "available_spare", Value: 5, Threshold: 10},
			override: nil,
			expected: pkg.AttributeStatusFailedScrutiny,
		},
		{
			name:     "override replaces the threshold and a higher spare passes",
			attr:     measurements.SmartNvmeAttribute{AttributeId: "available_spare", Value: 5, Threshold: 10},
			override: &pkg.AttributeOverride{FailThreshold: i64(3)},
			expected: pkg.AttributeStatusPassedOverride,
		},
		{
			name:     "override fails when the value drops below it, not above",
			attr:     measurements.SmartNvmeAttribute{AttributeId: "available_spare", Value: 5, Threshold: -1},
			override: &pkg.AttributeOverride{FailThreshold: i64(10)},
			expected: pkg.AttributeStatusFailedScrutiny | pkg.AttributeStatusFailedOverride,
		},
		{
			name:     "ideal-low NVMe attribute still compares upwards",
			attr:     measurements.SmartNvmeAttribute{AttributeId: "media_errors", Value: 25, Threshold: -1},
			override: &pkg.AttributeOverride{FailThreshold: i64(10)},
			expected: pkg.AttributeStatusFailedScrutiny | pkg.AttributeStatusFailedOverride,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			attr := tc.attr
			result := attr.PopulateAttributeStatus(tc.override)
			require.Equal(t, tc.expected, result.Status, "status mismatch (reason: %q)", result.StatusReason)
		})
	}
}

// End-to-end through the collector parser: a device that fails Scrutiny's own analysis on
// attribute 199 (raw 108, in the 70..130 bucket at 22.3% AFR) can be brought back to passing
// with an override, and the device status records that an override was responsible.
func TestSmart_FromCollectorSmartInfo_Overrides(t *testing.T) {
	loadFixture := func(t *testing.T) collector.SmartInfo {
		t.Helper()
		raw, err := os.ReadFile("../testdata/smart-ata-failed-scrutiny.json")
		require.NoError(t, err)
		var info collector.SmartInfo
		require.NoError(t, json.Unmarshal(raw, &info))
		return info
	}

	tests := []struct {
		name           string
		overrides      pkg.AttributeOverrideSet
		expectedDevice pkg.DeviceStatus
		expectedAttr   pkg.AttributeStatus
	}{
		{
			name:           "no overrides, device fails as before",
			overrides:      nil,
			expectedDevice: pkg.DeviceStatusFailedScrutiny,
			expectedAttr:   pkg.AttributeStatusFailedScrutiny,
		},
		{
			name: "override clears the failure and marks the device as passing by override",
			overrides: pkg.AttributeOverrideSet{
				{Protocol: pkg.DeviceProtocolAta, AttributeID: "199", FailThreshold: i64(200)},
			},
			expectedDevice: pkg.DeviceStatusPassedOverride,
			expectedAttr:   pkg.AttributeStatusPassedOverride,
		},
		{
			name: "override can also fail a device, and says so",
			overrides: pkg.AttributeOverrideSet{
				{Protocol: pkg.DeviceProtocolAta, AttributeID: "199", FailThreshold: i64(50)},
			},
			expectedDevice: pkg.DeviceStatusFailedScrutiny | pkg.DeviceStatusFailedOverride,
			expectedAttr:   pkg.AttributeStatusFailedScrutiny | pkg.AttributeStatusFailedOverride,
		},
		{
			name: "an override for another attribute leaves 199 alone",
			overrides: pkg.AttributeOverrideSet{
				{Protocol: pkg.DeviceProtocolAta, AttributeID: "5", FailThreshold: i64(200)},
			},
			expectedDevice: pkg.DeviceStatusFailedScrutiny,
			expectedAttr:   pkg.AttributeStatusFailedScrutiny,
		},
		{
			name: "an override for another protocol does not apply to an ATA device",
			overrides: pkg.AttributeOverrideSet{
				{Protocol: pkg.DeviceProtocolNvme, AttributeID: "199", FailThreshold: i64(200)},
			},
			expectedDevice: pkg.DeviceStatusFailedScrutiny,
			expectedAttr:   pkg.AttributeStatusFailedScrutiny,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			smartMdl := measurements.Smart{}
			err := smartMdl.FromCollectorSmartInfo(uuid.Must(uuid.NewV4()), loadFixture(t), tc.overrides)
			require.NoError(t, err)
			require.Equal(t, tc.expectedDevice, smartMdl.Status)
			require.Equal(t, tc.expectedAttr, smartMdl.Attributes["199"].GetStatus())
		})
	}
}

// SCSI attributes carry an Ideal in ScsiMetadata just like the other protocols, so overrides
// resolve their comparison direction the same way.
func TestSmartScsiAttribute_PopulateAttributeStatus_Override(t *testing.T) {
	tests := []struct {
		name     string
		attr     measurements.SmartScsiAttribute
		override *pkg.AttributeOverride
		expected pkg.AttributeStatus
	}{
		{
			// Documents existing upstream behaviour rather than endorsing it: the non-override
			// path looks SCSI attributes up in thresholds.NmveMetadata, so a SCSI-only id such
			// as scsi_grown_defect_list matches nothing and the recommended-threshold check is
			// skipped entirely -- even with a value past its threshold. Left as-is here to keep
			// this change focused; the override path below uses ScsiMetadata, which is correct.
			name:     "no override, value past its threshold is not flagged (upstream metadata mismatch)",
			attr:     measurements.SmartScsiAttribute{AttributeId: "scsi_grown_defect_list", Value: 5, Threshold: 0},
			override: nil,
			expected: pkg.AttributeStatusPassed,
		},
		{
			name:     "override passes when the value is within it",
			attr:     measurements.SmartScsiAttribute{AttributeId: "scsi_grown_defect_list", Value: 5, Threshold: 0},
			override: &pkg.AttributeOverride{FailThreshold: i64(10)},
			expected: pkg.AttributeStatusPassedOverride,
		},
		{
			name:     "override fails when the value exceeds it",
			attr:     measurements.SmartScsiAttribute{AttributeId: "scsi_grown_defect_list", Value: 5, Threshold: 0},
			override: &pkg.AttributeOverride{FailThreshold: i64(3)},
			expected: pkg.AttributeStatusFailedScrutiny | pkg.AttributeStatusFailedOverride,
		},
		{
			name:     "override warns between the thresholds",
			attr:     measurements.SmartScsiAttribute{AttributeId: "read_total_uncorrected_errors", Value: 7, Threshold: -1},
			override: &pkg.AttributeOverride{WarnThreshold: i64(5), FailThreshold: i64(50)},
			expected: pkg.AttributeStatusWarningScrutiny | pkg.AttributeStatusWarningOverride,
		},
		{
			name:     "an unknown attribute id defaults to ideal-low",
			attr:     measurements.SmartScsiAttribute{AttributeId: "not_a_real_scsi_attribute", Value: 9, Threshold: -1},
			override: &pkg.AttributeOverride{FailThreshold: i64(4)},
			expected: pkg.AttributeStatusFailedScrutiny | pkg.AttributeStatusFailedOverride,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			attr := tc.attr
			result := attr.PopulateAttributeStatus(tc.override)
			require.Equal(t, tc.expected, result.Status, "status mismatch (reason: %q)", result.StatusReason)
		})
	}
}
