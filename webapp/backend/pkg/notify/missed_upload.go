package notify

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/analogj/scrutiny/webapp/backend/pkg/config"
	"github.com/analogj/scrutiny/webapp/backend/pkg/database"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/gofrs/uuid/v5"
	"github.com/sirupsen/logrus"
)

// MissedUploadTracker notifies when a device misses an upload
type MissedUploadTracker struct {
	logger     logrus.FieldLogger
	config     config.Interface
	deviceRepo database.DeviceRepo
	// the `notify` settings, for devices that aren't in `devices` and as the defaults for those that are
	defaults uploadHistory
	// requests are handled concurrently, so uploads race on the map and on each device's history
	mu      sync.Mutex
	devices map[uuid.UUID]*uploadHistory
}

type uploadHistory struct {
	enabled bool          // `on_missed_upload`
	period  time.Duration // 0 means estimate it from uploads
	// the last settled estimate, kept while irregular uploads settle again
	estimated time.Duration
	// circular buffer of the latest upload times, uploads[next] is the oldest
	uploads [3]time.Time
	next    int
	timer   *time.Timer
}

func NewMissedUploadTracker(logger logrus.FieldLogger, appConfig config.Interface, deviceRepo database.DeviceRepo) *MissedUploadTracker {
	tracker := &MissedUploadTracker{
		logger:     logger,
		config:     appConfig,
		deviceRepo: deviceRepo,
		devices:    map[uuid.UUID]*uploadHistory{},

		defaults: uploadHistory{
			enabled: appConfig.GetBool("notify.on_missed_upload"),
			period:  appConfig.GetDuration("notify.upload_period"),
		},
	}

	devices, err := appConfig.GetDevices()
	if err != nil {
		logger.Errorf("Could not read `devices`, missed upload notifications are disabled: %v", err)
	}
	for _, device := range devices {
		history := tracker.defaults
		if device.Notifications.OnMissedUpload != nil {
			history.enabled = *device.Notifications.OnMissedUpload
		}
		if device.Notifications.UploadPeriod != nil {
			history.period = *device.Notifications.UploadPeriod
		}
		// ValidateConfig already rejected invalid UUIDs
		tracker.devices[device.ScrutinyUUID] = &history
	}
	tracker.loadEstimatedPeriods()
	return tracker
}

// loadEstimatedPeriods resumes tracking devices that have an estimated period saved, unless their configured period overrides it
func (t *MissedUploadTracker) loadEstimatedPeriods() {
	devices, err := t.deviceRepo.GetDevices(context.Background())
	if err != nil {
		t.logger.Errorf("Could not load the estimated upload periods, they will be estimated again: %v", err)
		return
	}
	now := time.Now()
	for _, device := range devices {
		period := device.EstimatedUploadPeriod
		if device.Archived || device.DeletedAt != nil || period <= 0 {
			continue
		}
		history, found := t.devices[device.ScrutinyUUID]
		if !found {
			history = &t.defaults
		}
		if history.period != 0 {
			continue
		}
		// replays steady uploads that end now
		t.record(device, now.Add(-2*period))
		t.record(device, now.Add(-period))
		t.record(device, now)
	}
}

// RecordUpload restarts the device's timer, which notifies if the next upload is more than 10% late
func (t *MissedUploadTracker) RecordUpload(device models.Device, uploadedAt time.Time) {
	estimated := t.record(device, uploadedAt)
	if estimated <= 0 {
		return
	}
	// saved outside the lock so a slow write doesn't block other uploads
	if err := t.deviceRepo.UpdateDeviceEstimatedUploadPeriod(context.Background(), device.ScrutinyUUID, estimated); err != nil {
		t.logger.Errorf("Could not save the estimated upload period of device %s: %v", device.ScrutinyUUID, err)
	}
}

// record returns the upload period if it was estimated from the uploads, 0 otherwise
func (t *MissedUploadTracker) record(device models.Device, uploadedAt time.Time) (estimated time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()

	history, found := t.devices[device.ScrutinyUUID]
	if !found {
		defaults := t.defaults
		history = &defaults
		t.devices[device.ScrutinyUUID] = history
	}
	if !history.enabled {
		return 0
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
			return 0
		}
		middle := history.uploads[(history.next+1)%len(history.uploads)]
		firstDelta, secondDelta := middle.Sub(oldest), uploadedAt.Sub(middle)
		period = (firstDelta + secondDelta) / 2
		if (firstDelta - secondDelta).Abs() > period/10 {
			if history.estimated == 0 {
				t.logger.Warnf("The last uploads from device %s were %s and %s apart, so its upload period is unknown until they settle. Set its `notifications.upload_period` to skip estimating it.",
					device.ScrutinyUUID, firstDelta, secondDelta)
				return 0
			}
			t.logger.Warnf("The last uploads from device %s were %s and %s apart, so its estimated upload period of %s is kept until they settle. Set its `notifications.upload_period` to skip estimating it.",
				device.ScrutinyUUID, firstDelta, secondDelta, history.estimated)
			period = history.estimated
		} else {
			history.estimated = period
			estimated = period
		}
	}
	// e.g. uploads with identical timestamps, which would notify immediately
	if period <= 0 {
		return 0
	}

	history.timer = time.AfterFunc(period+period/10, func() {
		missedUpload := NewMissedUpload(t.logger, t.config, device, uploadedAt, period)
		_ = missedUpload.Send() //we ignore error message when sending notifications.
	})
	return estimated
}

func NewMissedUpload(logger logrus.FieldLogger, appconfig config.Interface, device models.Device, lastUpload time.Time, uploadPeriod time.Duration) Notify {
	payload := NewPayload(device, false)
	payload.FailureType = NotifyFailureTypeMissedUpload
	payload.Subject = payload.GenerateMissedUploadSubject()
	payload.Message = payload.generateDeviceMessage("Scrutiny missed upload notification",
		fmt.Sprintf("Scrutiny UUID: %s", device.ScrutinyUUID),
		fmt.Sprintf("Last Upload: %s", lastUpload.Format(time.RFC3339)),
		fmt.Sprintf("Upload Period: %s", uploadPeriod),
	)

	return Notify{
		Logger:  logger,
		Config:  appconfig,
		Payload: payload,

		ScrutinyUUID: device.ScrutinyUUID,
	}
}

func (p *Payload) GenerateMissedUploadSubject() string {
	if len(p.HostId) > 0 {
		return fmt.Sprintf("Scrutiny missed upload for %s from %s", p.DeviceName, p.HostId)
	}
	return fmt.Sprintf("Scrutiny missed upload from device: %s", p.DeviceName)
}
