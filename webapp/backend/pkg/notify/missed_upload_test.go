package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/gofrs/uuid/v5"
	"github.com/sirupsen/logrus"
	logrustest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var missedUploadDevice = models.Device{
	ScrutinyUUID: uuid.FromStringOrNil("ecfaaf20-d1f6-558b-b33a-3e8db19a6c2c"),
	HostId:       "host1",
	DeviceName:   "sda",
}

// the tracker sends its notifications to a webhook, which forwards them to the returned channel
func newTestMissedUploadTracker(t *testing.T, devices ...map[string]interface{}) (*MissedUploadTracker, *logrustest.Hook, chan Payload) {
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
	return NewMissedUploadTracker(logger, appConfig), hook, payloads
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
	require.Equal(t, "Scrutiny missed upload from [host]device: [host1]sda", payload.Subject)
	require.Contains(t, payload.Message, "Host Id: host1\nDevice Name: sda\n")
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

func TestMissedUploadTracker_IrregularUploadStopsTimer(t *testing.T) {
	tracker, hook, payloads := newTestMissedUploadTracker(t)
	now := time.Now()

	tracker.RecordUpload(missedUploadDevice, now.Add(-400*time.Millisecond))
	tracker.RecordUpload(missedUploadDevice, now.Add(-200*time.Millisecond))
	tracker.RecordUpload(missedUploadDevice, now)
	require.NotNil(t, tracker.devices[missedUploadDevice.ScrutinyUUID].timer)

	// eg. a manual collector run
	tracker.RecordUpload(missedUploadDevice, now.Add(50*time.Millisecond))
	require.Nil(t, tracker.devices[missedUploadDevice.ScrutinyUUID].timer)
	require.Len(t, warnings(hook), 1)

	// the 200ms timer would fire 220ms after the third upload
	select {
	case <-payloads:
		require.FailNow(t, "notified although the device uploaded")
	case <-time.After(400 * time.Millisecond):
	}
}

func TestMissedUploadTracker_SkipsDisabledDevices(t *testing.T) {
	tracker, _, _ := newTestMissedUploadTracker(t, map[string]interface{}{
		"scrutiny_uuid": missedUploadDevice.ScrutinyUUID.String(),
		"notifications": map[string]interface{}{"on_missed_upload": false, "upload_period": "10ms"},
	})

	tracker.RecordUpload(missedUploadDevice, time.Now())
	require.Nil(t, tracker.devices[missedUploadDevice.ScrutinyUUID].timer)
}
