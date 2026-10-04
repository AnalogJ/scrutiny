package notify

import (
	"fmt"
	"sync"
	"time"

	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/gofrs/uuid/v5"
	"github.com/sirupsen/logrus"
)

// MissedUploadTracker notifies when a device misses an upload
type MissedUploadTracker struct {
	logger  logrus.FieldLogger
	config  config.Interface
	mu      sync.Mutex
	devices map[uuid.UUID]*uploadHistory
}

type uploadHistory struct {
	disabled bool          // `on_missed_upload: false`
	period   time.Duration // 0 means estimate it from uploads
	// circular buffer of the latest upload times, uploads[next] is the oldest
	uploads [3]time.Time
	next    int
	timer   *time.Timer
}

func NewMissedUploadTracker(logger logrus.FieldLogger, appConfig config.Interface) *MissedUploadTracker {
	tracker := &MissedUploadTracker{
		logger:  logger,
		config:  appConfig,
		devices: map[uuid.UUID]*uploadHistory{},
	}

	devices, err := appConfig.GetDevices()
	if err != nil {
		logger.Errorf("Could not read `devices`, missed upload notifications are disabled: %v", err)
	}
	for _, device := range devices {
		onMissedUpload := device.Notifications.OnMissedUpload
		// ValidateConfig already rejected invalid UUIDs
		tracker.devices[uuid.FromStringOrNil(device.ScrutinyUUID)] = &uploadHistory{
			disabled: onMissedUpload != nil && !*onMissedUpload,
			period:   device.Notifications.UploadPeriod,
		}
	}
	return tracker
}

// RecordUpload restarts the device's timer, which notifies if the next upload is more than 10% late
func (t *MissedUploadTracker) RecordUpload(device models.Device, uploadedAt time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()

	history, found := t.devices[device.ScrutinyUUID]
	if !found {
		history = &uploadHistory{}
		t.devices[device.ScrutinyUUID] = history
	}
	if history.disabled {
		return
	}
	if history.timer != nil {
		history.timer.Stop()
		history.timer = nil
	}
	history.uploads[history.next] = uploadedAt
	history.next = (history.next + 1) % len(history.uploads)

	period := history.period
	if period == 0 {
		oldest := history.uploads[history.next]
		// fewer than 3 uploads so far
		if oldest.IsZero() {
			return
		}
		middle := history.uploads[(history.next+1)%len(history.uploads)]
		firstDelta, secondDelta := middle.Sub(oldest), uploadedAt.Sub(middle)
		period = (firstDelta + secondDelta) / 2
		if (firstDelta - secondDelta).Abs() > period/10 {
			t.logger.Warnf("The last uploads from device %s were %s and %s apart, so its upload period is unknown until they settle. Set its `notifications.upload_period` to skip estimating it.",
				device.ScrutinyUUID, firstDelta, secondDelta)
			return
		}
	}
	// e.g. uploads with identical timestamps, which would notify immediately
	if period <= 0 {
		return
	}

	history.timer = time.AfterFunc(period+period/10, func() {
		missedUpload := NewMissedUpload(t.logger, t.config, device, uploadedAt, period)
		_ = missedUpload.Send() //we ignore error message when sending notifications.
	})
}

func NewMissedUpload(logger logrus.FieldLogger, appconfig config.Interface, device models.Device, lastUpload time.Time, uploadPeriod time.Duration) Notify {
	payload := NewPayload(device, false)
	payload.FailureType = NotifyFailureTypeMissedUpload
	payload.Subject = payload.GenerateMissedUploadSubject()
	payload.Message = payload.generateDeviceMessage("Scrutiny missed upload notification",
		fmt.Sprintf("Last Upload: %s", lastUpload.Format(time.RFC3339)),
		fmt.Sprintf("Upload Period: %s", uploadPeriod),
	)

	return Notify{
		Logger:  logger,
		Config:  appconfig,
		Payload: payload,
	}
}

func (p *Payload) GenerateMissedUploadSubject() string {
	if len(p.HostId) > 0 {
		return fmt.Sprintf("Scrutiny missed upload from [host]device: [%s]%s", p.HostId, p.DeviceName)
	}
	return fmt.Sprintf("Scrutiny missed upload from device: %s", p.DeviceName)
}
