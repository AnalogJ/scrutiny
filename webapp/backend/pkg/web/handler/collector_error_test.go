package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	mock_database "github.com/analogj/scrutiny/webapp/backend/pkg/database/mock"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/collector"
	"github.com/analogj/scrutiny/webapp/backend/pkg/notify"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

var collectorErrorDevice = models.Device{
	ScrutinyUUID: uuid.FromStringOrNil("ecfaaf20-d1f6-558b-b33a-3e8db19a6c2c"),
	HostId:       "host1",
	DeviceName:   "sda",
}

type collectorErrorFixture struct {
	repo     *mock_database.MockDeviceRepo
	tracker  *notify.MissedUploadTracker
	payloads chan notify.Payload
	report   func(scrutinyUUID string)
}

// notifications go to a webhook, which forwards them to the fixture's channel
func newCollectorErrorFixture(t *testing.T, onSmartctlError bool, uploadPeriod time.Duration) collectorErrorFixture {
	payloads := make(chan notify.Payload, 10)
	webhook := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload notify.Payload
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		payloads <- payload
	}))
	t.Cleanup(webhook.Close)

	appConfig, err := config.Create()
	require.NoError(t, err)
	appConfig.Set("notify.urls", []string{webhook.URL})
	appConfig.Set("notify.on_smartctl_error", onSmartctlError)
	appConfig.Set("notify.on_missed_upload", true)
	appConfig.Set("notify.upload_period", uploadPeriod)

	repo := mock_database.NewMockDeviceRepo(gomock.NewController(t))
	repo.EXPECT().GetDevices(gomock.Any()).AnyTimes().Return(nil, nil)

	logger := logrus.New()
	tracker := notify.NewMissedUploadTracker(logger, appConfig, repo)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/collector/error", func(c *gin.Context) {
		c.Set("LOGGER", logrus.NewEntry(logger))
		c.Set("CONFIG", appConfig)
		c.Set("DEVICE_REPOSITORY", repo)
		c.Set("MISSED_UPLOAD_TRACKER", tracker)
		CollectorError(c)
	})

	return collectorErrorFixture{
		repo:     repo,
		tracker:  tracker,
		payloads: payloads,
		report: func(scrutinyUUID string) {
			body, err := json.Marshal(collector.CollectorError{
				Error:        "smartctl could not open the device",
				HostId:       "host1",
				DeviceName:   "sda",
				ScrutinyUUID: scrutinyUUID,
			})
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/collector/error", bytes.NewReader(body))
			router.ServeHTTP(recorder, req)
			require.Equal(t, http.StatusOK, recorder.Code)
		},
	}
}

func requireNoNotification(t *testing.T, payloads chan notify.Payload, wait time.Duration) {
	select {
	case payload := <-payloads:
		require.FailNow(t, "unexpected notification", payload.FailureType)
	case <-time.After(wait):
	}
}

func requireNotificationType(t *testing.T, payloads chan notify.Payload, failureType string) {
	select {
	case payload := <-payloads:
		require.Equal(t, failureType, payload.FailureType)
	case <-time.After(2 * time.Second):
		require.FailNow(t, "no notification", failureType)
	}
}

func TestCollectorError_RestartsMissedUploadTimer(t *testing.T) {
	for name, onSmartctlError := range map[string]bool{
		"on_smartctl_error enabled":  true,
		"on_smartctl_error disabled": false,
	} {
		t.Run(name, func(t *testing.T) {
			fixture := newCollectorErrorFixture(t, onSmartctlError, 200*time.Millisecond)
			fixture.repo.EXPECT().GetDeviceDetails(gomock.Any(), collectorErrorDevice.ScrutinyUUID).AnyTimes().Return(collectorErrorDevice, nil)

			fixture.tracker.RecordUpload(collectorErrorDevice, time.Now())
			time.Sleep(100 * time.Millisecond)
			fixture.report(collectorErrorDevice.ScrutinyUUID.String())
			if onSmartctlError {
				requireNotificationType(t, fixture.payloads, notify.NotifyFailureTypeCollectorError)
			}

			// the upload's timer would fire 220ms after it
			requireNoNotification(t, fixture.payloads, 150*time.Millisecond)
			requireNotificationType(t, fixture.payloads, notify.NotifyFailureTypeMissedUpload)
		})
	}
}

func TestCollectorError_RecordsUploadWithoutLookup(t *testing.T) {
	// the mock fails the test on any GetDeviceDetails call
	fixture := newCollectorErrorFixture(t, false, time.Hour)

	fixture.report(collectorErrorDevice.ScrutinyUUID.String())

	requireNoNotification(t, fixture.payloads, 50*time.Millisecond)
}

func TestCollectorError_WithoutRegisteredDeviceDoesNotNotify(t *testing.T) {
	unknownUUID := uuid.FromStringOrNil("9a4d34b5-b2ee-51ef-8506-90eea09be417")
	for name, scrutinyUUID := range map[string]string{
		"missing":      "",
		"malformed":    "not-a-uuid",
		"nil":          uuid.Nil.String(),
		"unregistered": unknownUUID.String(),
	} {
		t.Run(name, func(t *testing.T) {
			fixture := newCollectorErrorFixture(t, false, 50*time.Millisecond)
			fixture.repo.EXPECT().GetDeviceDetails(gomock.Any(), unknownUUID).AnyTimes().Return(models.Device{}, gorm.ErrRecordNotFound)

			fixture.report(scrutinyUUID)

			requireNoNotification(t, fixture.payloads, 150*time.Millisecond)
		})
	}
}
