package notify

import (
	"testing"
	"time"

	"github.com/analogj/scrutiny/webapp/backend/pkg"
	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	mock_database "github.com/analogj/scrutiny/webapp/backend/pkg/database/mock"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/measurements"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// A device that passes only because of a user override must stay silent. Its DeviceStatus is
// non-zero (it carries the passed-by-override bit), so a bare `== DeviceStatusPassed` check
// would fall through and start evaluating attributes.
func TestShouldNotify_PassedByOverride_IsSilent(t *testing.T) {
	t.Parallel()
	device := models.Device{DeviceStatus: pkg.DeviceStatusPassedOverride}
	smartAttrs := measurements.Smart{
		Attributes: map[string]measurements.SmartAttribute{
			"199": &measurements.SmartAtaAttribute{AttributeId: 199, Status: pkg.AttributeStatusPassedOverride},
		},
	}

	mockCtrl := gomock.NewController(t)
	fakeDatabase := mock_database.NewMockDeviceRepo(mockCtrl)

	// repeat notifications off, so with no expectations on the mock any history lookup fails this
	require.False(t, ShouldNotify(logrus.StandardLogger(), device, smartAttrs, uuid.Must(uuid.NewV4()),
		pkg.MetricsStatusThresholdBoth, pkg.MetricsStatusFilterAttributesAll, false, &gin.Context{}, fakeDatabase))
}

// An override that *causes* a failure must still notify.
func TestShouldNotify_FailedByOverride_Notifies(t *testing.T) {
	t.Parallel()
	device := models.Device{DeviceStatus: pkg.DeviceStatusFailedScrutiny | pkg.DeviceStatusFailedOverride}
	smartAttrs := measurements.Smart{
		Attributes: map[string]measurements.SmartAttribute{
			"199": &measurements.SmartAtaAttribute{
				AttributeId: 199,
				Status:      pkg.AttributeStatusFailedScrutiny | pkg.AttributeStatusFailedOverride,
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	fakeDatabase := mock_database.NewMockDeviceRepo(mockCtrl)

	require.True(t, ShouldNotify(logrus.StandardLogger(), device, smartAttrs, uuid.Must(uuid.NewV4()),
		pkg.MetricsStatusThresholdBoth, pkg.MetricsStatusFilterAttributesAll, true, &gin.Context{}, fakeDatabase))
}

// An attribute an override marked as passing must not be looked up as a failing attribute.
func TestShouldNotify_PassedByOverride_NotCountedAsFailingAttribute(t *testing.T) {
	t.Parallel()
	device := models.Device{DeviceStatus: pkg.DeviceStatusFailedScrutiny, DeviceProtocol: pkg.DeviceProtocolAta}
	smartAttrs := measurements.Smart{
		Attributes: map[string]measurements.SmartAttribute{
			"5":   &measurements.SmartAtaAttribute{AttributeId: 5, Status: pkg.AttributeStatusPassedOverride},
			"199": &measurements.SmartAtaAttribute{AttributeId: 199, Status: pkg.AttributeStatusFailedScrutiny},
		},
	}
	scrutinyUUID := uuid.Must(uuid.NewV4())

	mockCtrl := gomock.NewController(t)
	fakeDatabase := mock_database.NewMockDeviceRepo(mockCtrl)
	fakeDatabase.EXPECT().GetSmartAttributeHistory(&gin.Context{}, scrutinyUUID, database.DURATION_KEY_FOREVER, time.Duration(0), 1, 1, []string{"199"}).Return(nil, nil).Times(1)

	require.True(t, ShouldNotify(logrus.StandardLogger(), device, smartAttrs, scrutinyUUID,
		pkg.MetricsStatusThresholdBoth, pkg.MetricsStatusFilterAttributesAll, false, &gin.Context{}, fakeDatabase))
}

// Editing an override can fail an attribute whose value hasn't changed.
func TestShouldNotify_FailedByEditedOverride_SameValue_Notifies(t *testing.T) {
	t.Parallel()
	device := models.Device{DeviceStatus: pkg.DeviceStatusFailedScrutiny, DeviceProtocol: pkg.DeviceProtocolAta}
	previous := measurements.Smart{Attributes: map[string]measurements.SmartAttribute{
		"199": &measurements.SmartAtaAttribute{AttributeId: 199, RawValue: 118, Status: pkg.AttributeStatusPassedOverride},
	}}
	smartAttrs := measurements.Smart{Attributes: map[string]measurements.SmartAttribute{
		"199": &measurements.SmartAtaAttribute{AttributeId: 199, RawValue: 118, Status: pkg.AttributeStatusFailedScrutiny | pkg.AttributeStatusFailedOverride},
	}}
	scrutinyUUID := uuid.Must(uuid.NewV4())

	mockCtrl := gomock.NewController(t)
	fakeDatabase := mock_database.NewMockDeviceRepo(mockCtrl)
	fakeDatabase.EXPECT().GetSmartAttributeHistory(&gin.Context{}, scrutinyUUID, database.DURATION_KEY_FOREVER, time.Duration(0), 1, 1, []string{"199"}).Return([]measurements.Smart{previous}, nil).Times(1)

	require.True(t, ShouldNotify(logrus.StandardLogger(), device, smartAttrs, scrutinyUUID,
		pkg.MetricsStatusThresholdBoth, pkg.MetricsStatusFilterAttributesAll, false, &gin.Context{}, fakeDatabase))
}
