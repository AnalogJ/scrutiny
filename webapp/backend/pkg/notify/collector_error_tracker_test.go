package notify_test

import (
	"testing"

	"github.com/analogj/scrutiny/webapp/backend/pkg/notify"
	"github.com/stretchr/testify/assert"
)

func TestCollectorErrorTracker_ShouldNotifyCollectorError(t *testing.T) {
	tracker := notify.NewCollectorErrorTracker()

	assert.True(t, tracker.ShouldNotifyCollectorError("host", "/dev/sda", "cannot open device", false))
	//the same failure on every collector run is not worth notifying about again
	assert.False(t, tracker.ShouldNotifyCollectorError("host", "/dev/sda", "cannot open device", false))
	//a different failure is
	assert.True(t, tracker.ShouldNotifyCollectorError("host", "/dev/sda", "smartctl exited with status 1", false))
	//and so is the same failure on a different device, or a different host
	assert.True(t, tracker.ShouldNotifyCollectorError("host", "/dev/sdb", "smartctl exited with status 1", false))
	assert.True(t, tracker.ShouldNotifyCollectorError("otherhost", "/dev/sda", "smartctl exited with status 1", false))
}

func TestCollectorErrorTracker_RepeatNotifications(t *testing.T) {
	tracker := notify.NewCollectorErrorTracker()

	assert.True(t, tracker.ShouldNotifyCollectorError("host", "/dev/sda", "cannot open device", true))
	assert.True(t, tracker.ShouldNotifyCollectorError("host", "/dev/sda", "cannot open device", true))
}

func TestCollectorErrorTracker_ResolveCollectorError(t *testing.T) {
	tracker := notify.NewCollectorErrorTracker()

	assert.True(t, tracker.ShouldNotifyCollectorError("host", "/dev/sda", "cannot open device", false))
	assert.False(t, tracker.ShouldNotifyCollectorError("host", "/dev/sda", "cannot open device", false))

	//the device reported successfully, so the next failure is news again
	tracker.ResolveCollectorError("host", "/dev/sda")
	assert.True(t, tracker.ShouldNotifyCollectorError("host", "/dev/sda", "cannot open device", false))
}
