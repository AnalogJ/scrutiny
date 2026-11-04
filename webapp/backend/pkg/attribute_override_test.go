package pkg_test

import (
	"testing"

	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/stretchr/testify/require"
)

func TestAttributeOverrideSet_Find(t *testing.T) {
	set := pkg.AttributeOverrideSet{
		{Protocol: pkg.DeviceProtocolAta, AttributeID: "199", FailThreshold: ptr(int64(200))},
		{Protocol: pkg.DeviceProtocolNvme, AttributeID: "media_errors", WarnThreshold: ptr(int64(1))},
	}

	tests := []struct {
		name        string
		protocol    string
		attributeID string
		found       bool
	}{
		{"matches protocol and id", pkg.DeviceProtocolAta, "199", true},
		{"matches a string attribute id", pkg.DeviceProtocolNvme, "media_errors", true},
		{"protocol match is case-insensitive", "ata", "199", true},
		{"same id under a different protocol does not match", pkg.DeviceProtocolNvme, "199", false},
		{"unknown attribute id does not match", pkg.DeviceProtocolAta, "5", false},
		{"unknown protocol does not match", pkg.DeviceProtocolScsi, "199", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := set.Find(tc.protocol, tc.attributeID)
			if !tc.found {
				require.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			require.Equal(t, tc.attributeID, got.AttributeID)
		})
	}
}

func TestAttributeOverrideSet_Find_EmptySet(t *testing.T) {
	var set pkg.AttributeOverrideSet
	require.Nil(t, set.Find(pkg.DeviceProtocolAta, "199"))
}

// The first matching entry wins, so a user can rely on document order.
func TestAttributeOverrideSet_Find_FirstMatchWins(t *testing.T) {
	set := pkg.AttributeOverrideSet{
		{Protocol: pkg.DeviceProtocolAta, AttributeID: "199", FailThreshold: ptr(int64(100))},
		{Protocol: pkg.DeviceProtocolAta, AttributeID: "199", FailThreshold: ptr(int64(200))},
	}
	got := set.Find(pkg.DeviceProtocolAta, "199")
	require.NotNil(t, got)
	require.Equal(t, int64(100), *got.FailThreshold)
}

func ptr[T any](v T) *T { return &v }
