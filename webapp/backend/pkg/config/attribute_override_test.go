package config

import (
	"strings"
	"testing"

	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func configFromYaml(t *testing.T, yaml string) *configuration {
	t.Helper()
	c := configuration{Viper: viper.New()}
	c.SetConfigType("yaml")
	require.NoError(t, c.MergeConfig(strings.NewReader(yaml)))
	return &c
}

func TestDeviceAttributeOverrides(t *testing.T) {
	c := configFromYaml(t, `
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

	overrides, err := DeviceAttributeOverrides(c)
	require.NoError(t, err)
	require.Len(t, overrides, 2)

	first := overrides["106089ed-2273-54e0-b498-ec4bdfc8ca6c"]
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
	second := overrides["8301e652-2817-5df7-8bd4-49fe76de25b0"]
	require.Len(t, second, 1)
	media := second.Find(pkg.DeviceProtocolNvme, "media_errors")
	require.NotNil(t, media)
	require.Equal(t, int64(1), *media.WarnThreshold)
	require.Nil(t, media.FailThreshold)
}

func TestDeviceAttributeOverrides_NotConfigured(t *testing.T) {
	c := configFromYaml(t, "web:\n  listen:\n    port: 8080\n")
	overrides, err := DeviceAttributeOverrides(c)
	require.NoError(t, err)
	require.Empty(t, overrides)
}

// A malformed override should stop the server with a clear message rather than being
// silently dropped -- a silently ignored override looks identical to a working one.
func TestValidateConfig_DeviceOverrides(t *testing.T) {
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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := configFromYaml(t, tc.yaml).ValidateConfig()
			if tc.expectErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			require.Contains(t, strings.ToLower(err.Error()), tc.expectErr)
		})
	}
}
