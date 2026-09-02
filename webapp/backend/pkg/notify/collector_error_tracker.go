package notify

import (
	"strings"
	"sync"
)

// CollectorErrorTracker remembers the last collector failure reported for each device so that a
// permanently broken one — a USB enclosure that never answers `smartctl --info`, say — does not
// notify on every collector run. It mirrors what metrics.repeat_notifications does for SMART
// failures, and is deliberately in-memory: losing it on restart costs at most one extra
// notification.
type CollectorErrorTracker struct {
	mutex             sync.Mutex
	lastErrorMessages map[string]string
}

func NewCollectorErrorTracker() *CollectorErrorTracker {
	return &CollectorErrorTracker{lastErrorMessages: map[string]string{}}
}

// ShouldNotifyCollectorError records the failure and reports whether the user should hear about it:
// always when they asked for repeat notifications, otherwise only when it differs from the failure
// last reported for the same device.
func (t *CollectorErrorTracker) ShouldNotifyCollectorError(hostId string, deviceName string, errorMessage string, repeatNotifications bool) bool {
	key := collectorErrorKey(hostId, deviceName)

	t.mutex.Lock()
	defer t.mutex.Unlock()

	previousErrorMessage, seen := t.lastErrorMessages[key]
	t.lastErrorMessages[key] = errorMessage

	return repeatNotifications || !seen || previousErrorMessage != errorMessage
}

// ResolveCollectorError forgets a device's last failure, so that the next one notifies again.
func (t *CollectorErrorTracker) ResolveCollectorError(hostId string, deviceName string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	delete(t.lastErrorMessages, collectorErrorKey(hostId, deviceName))
}

// a null byte, because neither a host id nor a device name can contain one
func collectorErrorKey(hostId string, deviceName string) string {
	return strings.Join([]string{hostId, deviceName}, "\x00")
}
