package config

import (
	"strings"
	"testing"

	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/require"
)

func TestDeviceConfig_AttributeOverrides(t *testing.T) {
	c, err := loadTestConfig(t, `
devices:
  - scrutiny_uuid: 106089ed-2273-54e0-b498-ec4bdfc8ca6c
    attribute_overrides:
      - protocol: ATA
        attribute_id: "199"
        warn_threshold: 100
        fail_threshold: 200
      - protocol: ATA
        attribute_id: "188"
        fail_threshold: 0
  - scrutiny_uuid: 8301E652-2817-5DF7-8BD4-49FE76DE25B0
    attribute_overrides:
      - protocol: NVMe
        attribute_id: media_errors
        warn_threshold: 1
`)
	require.NoError(t, err)

	first := c.GetDeviceConfig(uuid.FromStringOrNil("106089ed-2273-54e0-b498-ec4bdfc8ca6c")).AttributeOverrides
	require.Len(t, first, 2)

	crc := first.Find(pkg.DeviceProtocolAta, "199")
	require.NotNil(t, crc)
	require.Equal(t, int64(100), *crc.WarnThreshold)
	require.Equal(t, int64(200), *crc.FailThreshold)

	// An explicit 0 must survive as a set threshold rather than collapsing to "unset".
	timeout := first.Find(pkg.DeviceProtocolAta, "188")
	require.NotNil(t, timeout)
	require.Nil(t, timeout.WarnThreshold)
	require.NotNil(t, timeout.FailThreshold)
	require.Equal(t, int64(0), *timeout.FailThreshold)

	// UUIDs are matched case-insensitively, since the config is hand-written.
	second := c.GetDeviceConfig(uuid.FromStringOrNil("8301e652-2817-5df7-8bd4-49fe76de25b0")).AttributeOverrides
	require.Len(t, second, 1)
	media := second.Find(pkg.DeviceProtocolNvme, "media_errors")
	require.NotNil(t, media)
	require.Equal(t, int64(1), *media.WarnThreshold)
	require.Nil(t, media.FailThreshold)
}

// A malformed override should stop the server with a clear message rather than being
// silently dropped -- a silently ignored override looks identical to a working one.
func TestValidateConfig_AttributeOverrides(t *testing.T) {
	tests := []struct {
		name      string
		yaml      string
		expectErr string
	}{
		{
			name: "valid config passes",
			yaml: `
devices:
  - scrutiny_uuid: 106089ed-2273-54e0-b498-ec4bdfc8ca6c
    attribute_overrides:
      - protocol: ATA
        attribute_id: "199"
        fail_threshold: 200
`,
		},
		{
			name: "unknown protocol is rejected",
			yaml: `
devices:
  - scrutiny_uuid: 106089ed-2273-54e0-b498-ec4bdfc8ca6c
    attribute_overrides:
      - protocol: SATA
        attribute_id: "199"
        fail_threshold: 200
`,
			expectErr: "protocol",
		},
		{
			name: "missing attribute_id is rejected",
			yaml: `
devices:
  - scrutiny_uuid: 106089ed-2273-54e0-b498-ec4bdfc8ca6c
    attribute_overrides:
      - protocol: ATA
        fail_threshold: 200
`,
			expectErr: "attribute_id",
		},
		{
			name: "an override with no thresholds is rejected",
			yaml: `
devices:
  - scrutiny_uuid: 106089ed-2273-54e0-b498-ec4bdfc8ca6c
    attribute_overrides:
      - protocol: ATA
        attribute_id: "199"
`,
			expectErr: "threshold",
		},
		{
			name: "missing scrutiny_uuid is rejected",
			yaml: `
devices:
  - attribute_overrides:
      - protocol: ATA
        attribute_id: "199"
        fail_threshold: 200
`,
			expectErr: "scrutiny_uuid",
		},
		{
			name: "non-numeric ATA attribute_id is rejected",
			yaml: `
devices:
  - scrutiny_uuid: 106089ed-2273-54e0-b498-ec4bdfc8ca6c
    attribute_overrides:
      - protocol: ATA
        attribute_id: devstat_7_8
        fail_threshold: 200
`,
			expectErr: "attribute_id",
		},
		{
			name: "warn above fail is rejected when larger is worse",
			yaml: `
devices:
  - scrutiny_uuid: 106089ed-2273-54e0-b498-ec4bdfc8ca6c
    attribute_overrides:
      - protocol: ATA
        attribute_id: "199"
        warn_threshold: 500
        fail_threshold: 100
`,
			expectErr: "warn_threshold",
		},
		{
			name: "warn below fail is rejected when smaller is worse",
			yaml: `
devices:
  - scrutiny_uuid: 106089ed-2273-54e0-b498-ec4bdfc8ca6c
    attribute_overrides:
      - protocol: NVMe
        attribute_id: available_spare
        warn_threshold: 10
        fail_threshold: 20
`,
			expectErr: "warn_threshold",
		},
		{
			name: "warn above fail is valid when smaller is worse",
			yaml: `
devices:
  - scrutiny_uuid: 106089ed-2273-54e0-b498-ec4bdfc8ca6c
    attribute_overrides:
      - protocol: NVMe
        attribute_id: available_spare
        warn_threshold: 20
        fail_threshold: 10
`,
		},
		{
			name:      "a devices block that isn't a list is rejected",
			yaml:      "devices: \"not a list\"\n",
			expectErr: "devices",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := loadTestConfig(t, tc.yaml)
			if tc.expectErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			require.Contains(t, strings.ToLower(err.Error()), tc.expectErr)
		})
	}
}
