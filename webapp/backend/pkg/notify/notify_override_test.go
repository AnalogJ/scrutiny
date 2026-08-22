package notify

import (
	"testing"

	"github.com/analogj/scrutiny/webapp/backend/pkg"
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

	require.False(t, ShouldNotify(logrus.StandardLogger(), device, smartAttrs, uuid.Must(uuid.NewV4()),
		pkg.MetricsStatusThresholdBoth, pkg.MetricsStatusFilterAttributesAll, true, &gin.Context{}, fakeDatabase))
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

// With the critical-attributes filter on, an attribute that an override marked as passing must
// not be collected as a failing attribute. 199 is non-critical, so the clearest probe is a
// critical attribute (5) that an override cleared: nothing should be reported.
func TestShouldNotify_PassedByOverride_NotCountedAsFailingAttribute(t *testing.T) {
	t.Parallel()
	device := models.Device{DeviceStatus: pkg.DeviceStatusPassedOverride, DeviceProtocol: pkg.DeviceProtocolAta}
	smartAttrs := measurements.Smart{
		Attributes: map[string]measurements.SmartAttribute{
			"5": &measurements.SmartAtaAttribute{AttributeId: 5, Status: pkg.AttributeStatusPassedOverride},
		},
	}

	mockCtrl := gomock.NewController(t)
	fakeDatabase := mock_database.NewMockDeviceRepo(mockCtrl)

	require.False(t, ShouldNotify(logrus.StandardLogger(), device, smartAttrs, uuid.Must(uuid.NewV4()),
		pkg.MetricsStatusThresholdBoth, pkg.MetricsStatusFilterAttributesCritical, true, &gin.Context{}, fakeDatabase))
}
