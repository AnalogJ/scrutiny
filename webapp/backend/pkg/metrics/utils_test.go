package metrics

import (
	"testing"
	"time"

	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
)

func TestSanitizeMetricName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "convert dots to underscores",
			input:    "attr.5.raw_value",
			expected: "scrutiny_smart_attr_5_raw_value",
		},
		{
			name:     "convert hyphens to underscores",
			input:    "attr-5-raw-value",
			expected: "scrutiny_smart_attr_5_raw_value",
		},
		{
			name:     "convert spaces to underscores",
			input:    "attr 5 raw value",
			expected: "scrutiny_smart_attr_5_raw_value",
		},
		{
			name:     "convert to lowercase",
			input:    "Attr.5.Raw_Value",
			expected: "scrutiny_smart_attr_5_raw_value",
		},
		{
			name:     "already valid name",
			input:    "valid_metric_name",
			expected: "valid_metric_name",
		},
		{
			name:     "without attr prefix",
			input:    "some.metric.name",
			expected: "some_metric_name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeMetricName(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeMetricName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTryParseFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected float64
		shouldOk bool
	}{
		{
			name:     "parse int",
			input:    42,
			expected: 42.0,
			shouldOk: true,
		},
		{
			name:     "parse int64",
			input:    int64(12345),
			expected: 12345.0,
			shouldOk: true,
		},
		{
			name:     "parse float64",
			input:    3.14159,
			expected: 3.14159,
			shouldOk: true,
		},
		{
			name:     "parse float32",
			input:    float32(2.71),
			expected: float64(float32(2.71)), // Account for float32 precision
			shouldOk: true,
		},
		{
			name:     "parse string number",
			input:    "123.45",
			expected: 123.45,
			shouldOk: true,
		},
		{
			name:     "parse hexadecimal with 0x",
			input:    "0x1A",
			expected: 26.0,
			shouldOk: true,
		},
		{
			name:     "parse hexadecimal with 0X",
			input:    "0XFF",
			expected: 255.0,
			shouldOk: true,
		},
		{
			name:     "parse empty string",
			input:    "",
			expected: 0,
			shouldOk: false,
		},
		{
			name:     "parse invalid string",
			input:    "not_a_number",
			expected: 0,
			shouldOk: false,
		},
		{
			name:     "parse nil",
			input:    nil,
			expected: 0,
			shouldOk: false,
		},
		{
			name:     "parse bool",
			input:    true,
			expected: 0,
			shouldOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := TryParseFloat(tt.input)
			if ok != tt.shouldOk {
				t.Errorf("TryParseFloat(%v) ok = %v, want %v", tt.input, ok, tt.shouldOk)
			}
			if tt.shouldOk && result != tt.expected {
				t.Errorf("TryParseFloat(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSelectLatestSmartResult(t *testing.T) {
	now := time.Now()
	older := now.Add(-1 * time.Hour)
	oldest := now.Add(-2 * time.Hour)

	tests := []struct {
		name     string
		input    []measurements.Smart
		expected *time.Time
	}{
		{
			name:     "empty list",
			input:    []measurements.Smart{},
			expected: nil,
		},
		{
			name: "single result",
			input: []measurements.Smart{
				{Date: now},
			},
			expected: &now,
		},
		{
			name: "multiple results in order",
			input: []measurements.Smart{
				{Date: now},
				{Date: older},
				{Date: oldest},
			},
			expected: &now,
		},
		{
			name: "multiple results out of order",
			input: []measurements.Smart{
				{Date: older},
				{Date: now},
				{Date: oldest},
			},
			expected: &now,
		},
		{
			name: "multiple results reverse order",
			input: []measurements.Smart{
				{Date: oldest},
				{Date: older},
				{Date: now},
			},
			expected: &now,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SelectLatestSmartResult(tt.input)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("SelectLatestSmartResult() = %v, want nil", result)
				}
				return
			}

			if result == nil {
				t.Errorf("SelectLatestSmartResult() = nil, want non-nil")
				return
			}

			if !result.Date.Equal(*tt.expected) {
				t.Errorf("SelectLatestSmartResult() date = %v, want %v", result.Date, *tt.expected)
			}
		})
	}
}
