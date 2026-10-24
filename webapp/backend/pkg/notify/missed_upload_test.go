package notify

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	mock_database "github.com/analogj/scrutiny/webapp/backend/pkg/database/mock"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/gofrs/uuid/v5"
	"github.com/sirupsen/logrus"
	logrustest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var missedUploadDevice = models.Device{
	ScrutinyUUID: uuid.FromStringOrNil("ecfaaf20-d1f6-558b-b33a-3e8db19a6c2c"),
	HostId:       "host1",
	DeviceName:   "sda",
}

// the tracker sends its notifications to a webhook, which forwards them to the returned channel
func newTestMissedUploadTracker(t *testing.T, devices ...map[string]interface{}) (*MissedUploadTracker, *logrustest.Hook, chan Payload) {
	deviceRepo := mock_database.NewMockDeviceRepo(gomock.NewController(t))
	deviceRepo.EXPECT().GetDevices(gomock.Any()).AnyTimes().Return(nil, nil)
	deviceRepo.EXPECT().UpdateDeviceEstimatedUploadPeriod(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
	return newTestMissedUploadTrackerWithRepo(t, deviceRepo, devices...)
}

func newTestMissedUploadTrackerWithRepo(t *testing.T, deviceRepo *mock_database.MockDeviceRepo, devices ...map[string]interface{}) (*MissedUploadTracker, *logrustest.Hook, chan Payload) {
	payloads := make(chan Payload, 10)
	webhook := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload Payload
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		payloads <- payload
	}))
	t.Cleanup(webhook.Close)

	appConfig, err := config.Create()
	require.NoError(t, err)
	appConfig.Set("notify.urls", []string{webhook.URL})
	appConfig.Set("devices", devices)

	logger, hook := logrustest.NewNullLogger()
	return NewMissedUploadTracker(logger, appConfig, deviceRepo), hook, payloads
}

func requireNotification(t *testing.T, payloads chan Payload) Payload {
	select {
	case payload := <-payloads:
		return payload
	case <-time.After(2 * time.Second):
		require.FailNow(t, "no missed upload notification")
		return Payload{}
	}
}

func warnings(hook *logrustest.Hook) []string {
	var messages []string
	for _, entry := range hook.AllEntries() {
		if entry.Level == logrus.WarnLevel {
			messages = append(messages, entry.Message)
		}
	}
	return messages
}

func TestMissedUploadTracker_ConfiguredPeriod(t *testing.T) {
	tracker, _, payloads := newTestMissedUploadTracker(t, map[string]interface{}{
		"scrutiny_uuid": missedUploadDevice.ScrutinyUUID.String(),
		"notifications": map[string]interface{}{"upload_period": "100ms"},
	})

	upload := time.Now()
	tracker.RecordUpload(missedUploadDevice, upload)
	payload := requireNotification(t, payloads)

	require.GreaterOrEqual(t, time.Since(upload), 110*time.Millisecond)
	require.Equal(t, NotifyFailureTypeMissedUpload, payload.FailureType)
	require.Equal(t, "Scrutiny missed upload for sda from host1", payload.Subject)
	require.Contains(t, payload.Message, "Host Id: host1\nDevice Name: sda\n")
	require.Contains(t, payload.Message, "Scrutiny UUID: "+missedUploadDevice.ScrutinyUUID.String())
	require.Contains(t, payload.Message, "Last Upload: "+upload.Format(time.RFC3339))
	require.Contains(t, payload.Message, "Upload Period: 100ms")
}

func TestMissedUploadTracker_UploadRestartsTimer(t *testing.T) {
	tracker, _, payloads := newTestMissedUploadTracker(t, map[string]interface{}{
		"scrutiny_uuid": missedUploadDevice.ScrutinyUUID.String(),
		"notifications": map[string]interface{}{"upload_period": "200ms"},
	})

	tracker.RecordUpload(missedUploadDevice, time.Now())
	time.Sleep(100 * time.Millisecond)
	tracker.RecordUpload(missedUploadDevice, time.Now())

	// the first upload's timer would fire 220ms after it
	select {
	case <-payloads:
		require.FailNow(t, "notified although the second upload was on time")
	case <-time.After(150 * time.Millisecond):
	}
	requireNotification(t, payloads)
}

func TestMissedUploadTracker_EstimatesPeriodOnceUploadsSettle(t *testing.T) {
	// devices missing from the config are tracked too
	tracker, hook, payloads := newTestMissedUploadTracker(t)
	now := time.Now()

	tracker.RecordUpload(missedUploadDevice, now.Add(-3*time.Hour))
	tracker.RecordUpload(missedUploadDevice, now.Add(-100*time.Millisecond))
	require.Nil(t, tracker.devices[missedUploadDevice.ScrutinyUUID].timer, "estimated a period from 2 uploads")

	tracker.RecordUpload(missedUploadDevice, now.Add(-50*time.Millisecond))
	require.Nil(t, tracker.devices[missedUploadDevice.ScrutinyUUID].timer, "estimated a period from 3h and 50ms deltas")
	require.Len(t, warnings(hook), 1)
	hook.Reset()

	// overwrites the -3h upload, so the deltas settle at 50ms
	tracker.RecordUpload(missedUploadDevice, now)
	require.Empty(t, warnings(hook))
	payload := requireNotification(t, payloads)
	require.Contains(t, payload.Message, "Upload Period: 50ms")
}

func TestMissedUploadTracker_IrregularUploadKeepsEstimate(t *testing.T) {
	deviceRepo := mock_database.NewMockDeviceRepo(gomock.NewController(t))
	deviceRepo.EXPECT().GetDevices(gomock.Any()).Return(nil, nil)
	// only the settled estimate is saved
	deviceRepo.EXPECT().UpdateDeviceEstimatedUploadPeriod(gomock.Any(), gomock.Any(), 200*time.Millisecond).Return(nil)
	tracker, hook, payloads := newTestMissedUploadTrackerWithRepo(t, deviceRepo)
	now := time.Now()

	tracker.RecordUpload(missedUploadDevice, now.Add(-400*time.Millisecond))
	tracker.RecordUpload(missedUploadDevice, now.Add(-200*time.Millisecond))
	tracker.RecordUpload(missedUploadDevice, now)
	require.Empty(t, warnings(hook))

	// eg. a manual collector run
	irregular := time.Now()
	tracker.RecordUpload(missedUploadDevice, now.Add(50*time.Millisecond))
	require.Len(t, warnings(hook), 1)
	payload := requireNotification(t, payloads)
	require.GreaterOrEqual(t, time.Since(irregular), 220*time.Millisecond)
	require.Contains(t, payload.Message, "Upload Period: 200ms")
}

func TestMissedUploadTracker_SkipsDisabledDevices(t *testing.T) {
	tracker, _, _ := newTestMissedUploadTracker(t, map[string]interface{}{
		"scrutiny_uuid": missedUploadDevice.ScrutinyUUID.String(),
		"notifications": map[string]interface{}{"on_missed_upload": false, "upload_period": "10ms"},
	})

	tracker.RecordUpload(missedUploadDevice, time.Now())
	require.Nil(t, tracker.devices[missedUploadDevice.ScrutinyUUID].timer)
}

func TestMissedUploadTracker_SavesEstimatedPeriod(t *testing.T) {
	deviceRepo := mock_database.NewMockDeviceRepo(gomock.NewController(t))
	deviceRepo.EXPECT().GetDevices(gomock.Any()).Return(nil, nil)
	tracker, _, _ := newTestMissedUploadTrackerWithRepo(t, deviceRepo)
	now := time.Now()

	// the period is only known once the third upload settles it
	tracker.RecordUpload(missedUploadDevice, now.Add(-2*time.Hour))
	tracker.RecordUpload(missedUploadDevice, now.Add(-time.Hour))
	deviceRepo.EXPECT().UpdateDeviceEstimatedUploadPeriod(gomock.Any(), missedUploadDevice.ScrutinyUUID, time.Hour).Return(nil)
	tracker.RecordUpload(missedUploadDevice, now)
}

func TestMissedUploadTracker_SaveFailureIsLogged(t *testing.T) {
	deviceRepo := mock_database.NewMockDeviceRepo(gomock.NewController(t))
	deviceRepo.EXPECT().GetDevices(gomock.Any()).Return(nil, nil)
	deviceRepo.EXPECT().UpdateDeviceEstimatedUploadPeriod(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("db is locked"))
	tracker, hook, _ := newTestMissedUploadTrackerWithRepo(t, deviceRepo)
	now := time.Now()

	tracker.RecordUpload(missedUploadDevice, now.Add(-2*time.Hour))
	tracker.RecordUpload(missedUploadDevice, now.Add(-time.Hour))
	tracker.RecordUpload(missedUploadDevice, now)
	require.Contains(t, hook.LastEntry().Message, "db is locked")
	require.NotNil(t, tracker.devices[missedUploadDevice.ScrutinyUUID].timer)
}

func TestMissedUploadTracker_LoadsEstimatedPeriods(t *testing.T) {
	newDevice := func(last byte, period time.Duration) models.Device {
		return models.Device{ScrutinyUUID: uuid.UUID{15: last}, EstimatedUploadPeriod: period}
	}
	archived, deleted := newDevice(2, time.Hour), newDevice(3, time.Hour)
	archived.Archived = true
	deletedAt := time.Now()
	deleted.DeletedAt = &deletedAt
	configured, zeroConfigured := newDevice(4, time.Hour), newDevice(5, time.Hour)

	deviceRepo := mock_database.NewMockDeviceRepo(gomock.NewController(t))
	deviceRepo.EXPECT().GetDevices(gomock.Any()).Return([]models.Device{
		newDevice(1, time.Hour), archived, deleted, newDevice(6, 0), configured, zeroConfigured,
	}, nil)
	tracker, _, _ := newTestMissedUploadTrackerWithRepo(t, deviceRepo,
		map[string]interface{}{"scrutiny_uuid": configured.ScrutinyUUID.String(), "notifications": map[string]interface{}{"upload_period": "2h"}},
		map[string]interface{}{"scrutiny_uuid": zeroConfigured.ScrutinyUUID.String(), "notifications": map[string]interface{}{"upload_period": "0s"}},
	)

	// the periods are loaded, so a late upload is noticed without 3 uploads after the restart
	require.NotNil(t, tracker.devices[newDevice(1, 0).ScrutinyUUID].timer)
	require.NotNil(t, tracker.devices[zeroConfigured.ScrutinyUUID].timer)
	require.NotContains(t, tracker.devices, archived.ScrutinyUUID)
	require.NotContains(t, tracker.devices, deleted.ScrutinyUUID)
	require.NotContains(t, tracker.devices, newDevice(6, 0).ScrutinyUUID)
	require.Equal(t, 2*time.Hour, tracker.devices[configured.ScrutinyUUID].period)
	require.Nil(t, tracker.devices[configured.ScrutinyUUID].timer)
}

func TestMissedUploadTracker_LoadFailureIsLogged(t *testing.T) {
	deviceRepo := mock_database.NewMockDeviceRepo(gomock.NewController(t))
	deviceRepo.EXPECT().GetDevices(gomock.Any()).Return(nil, errors.New("db is locked"))
	_, hook, _ := newTestMissedUploadTrackerWithRepo(t, deviceRepo)
	require.Contains(t, hook.LastEntry().Message, "db is locked")
}
