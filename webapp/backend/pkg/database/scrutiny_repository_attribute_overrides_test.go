package database

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	mock_config "github.com/analogj/scrutiny/webapp/backend/pkg/config/mock"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/gofrs/uuid/v5"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// configWithDevices builds a real config from YAML, so this exercises the same parsing path the
// server uses at startup rather than a hand-stubbed decode.
func configWithDevices(t *testing.T, yaml string) config.Interface {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "scrutiny.yaml")
	require.NoError(t, os.WriteFile(path, []byte(yaml), 0600))

	c, err := config.Create()
	require.NoError(t, err)
	require.NoError(t, c.ReadConfig(path))
	return c
}

func Test_AttributeOverridesFor_NoneConfigured(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().IsSet("devices").Return(false).AnyTimes()

	repo := scrutinyRepository{appConfig: fakeConfig, logger: logrus.New()}
	overrides, err := repo.attributeOverridesFor(uuid.Must(uuid.NewV4()))
	require.NoError(t, err)
	require.Nil(t, overrides)
}

func Test_AttributeOverridesFor_ReturnsThisDevicesOverrides(t *testing.T) {
	t.Parallel()
	// The uuid is upper-cased in the config on purpose: it is hand-written, and must still match
	// the lower-cased uuid the collector reports.
	appConfig := configWithDevices(t, `
devices:
  - scrutiny_uuid: 8301E652-2817-5DF7-8BD4-49FE76DE25B0
    attribute_overrides:
      - protocol: ATA
        attribute_id: "199"
        fail_threshold: 500
  - scrutiny_uuid: 106089ed-2273-54e0-b498-ec4bdfc8ca6c
    attribute_overrides:
      - protocol: ATA
        attribute_id: "5"
        fail_threshold: 1
`)
	repo := scrutinyRepository{appConfig: appConfig, logger: logrus.New()}

	mine := uuid.Must(uuid.FromString("8301e652-2817-5df7-8bd4-49fe76de25b0"))
	overrides, err := repo.attributeOverridesFor(mine)
	require.NoError(t, err)
	require.Len(t, overrides, 1)
	match := overrides.Find(pkg.DeviceProtocolAta, "199")
	require.NotNil(t, match)
	require.Equal(t, int64(500), *match.FailThreshold)

	// A device with overrides configured for *other* devices gets none of them.
	other := uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000000"))
	overrides, err = repo.attributeOverridesFor(other)
	require.NoError(t, err)
	require.Empty(t, overrides)
}

// A `devices` block that cannot be decoded must surface as an error rather than being treated as
// "no overrides configured" -- silently ignoring it would apply stock analysis to a device the
// user believes they have tuned.
func Test_AttributeOverridesFor_ParseFailurePropagates(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().IsSet("devices").Return(true).AnyTimes()
	fakeConfig.EXPECT().UnmarshalKey("devices", gomock.Any()).
		Return(errors.New("boom: cannot decode devices")).AnyTimes()

	repo := scrutinyRepository{appConfig: fakeConfig, logger: logrus.New()}
	overrides, err := repo.attributeOverridesFor(uuid.Must(uuid.NewV4()))
	require.Error(t, err)
	require.Nil(t, overrides)
	require.Contains(t, err.Error(), "devices")
}

// ...and SaveSmartAttributes must refuse to persist rather than silently writing metrics graded by
// the wrong rules.
func Test_SaveSmartAttributes_RefusesOnBadOverrideConfig(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	fakeConfig := mock_config.NewMockInterface(mockCtrl)
	fakeConfig.EXPECT().IsSet("devices").Return(true).AnyTimes()
	fakeConfig.EXPECT().UnmarshalKey("devices", gomock.Any()).
		Return(errors.New("boom: cannot decode devices")).AnyTimes()

	repo := scrutinyRepository{appConfig: fakeConfig, logger: logrus.New()}
	smart, err := repo.SaveSmartAttributes(context.Background(), uuid.Must(uuid.NewV4()), collector.SmartInfo{})
	require.Error(t, err)
	require.Empty(t, smart.Attributes)
}
