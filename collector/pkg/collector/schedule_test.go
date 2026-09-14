package collector

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestParseCronSchedule(t *testing.T) {
	from := time.Date(2026, 1, 1, 10, 30, 0, 0, time.Local)

	for spec, expected := range map[string]time.Time{
		"0 0 * * *":      time.Date(2026, 1, 2, 0, 0, 0, 0, time.Local),
		"\"0 0 * * *\"":  time.Date(2026, 1, 2, 0, 0, 0, 0, time.Local),
		"'*/15 * * * *'": time.Date(2026, 1, 1, 10, 45, 0, 0, time.Local),
		"@hourly":        time.Date(2026, 1, 1, 11, 0, 0, 0, time.Local),
		"@every 10s":     from.Add(10 * time.Second),
	} {
		schedule, err := ParseCronSchedule(spec)
		require.NoError(t, err, spec)
		require.Equal(t, expected, schedule.Next(from), spec)
	}
}

func TestParseCronSchedule_Invalid(t *testing.T) {
	for _, spec := range []string{"", "* * *", "0 0 * * * *", "'0 0 * * *\"", "not a schedule"} {
		_, err := ParseCronSchedule(spec)
		require.Error(t, err, spec)
	}
}

// fires every interval, so tests don't have to wait for a real cron tick
type intervalSchedule time.Duration

func (s intervalSchedule) Next(t time.Time) time.Time {
	return t.Add(time.Duration(s))
}

func TestRunOnSchedule_RunsRepeatedlyAndSurvivesErrors(t *testing.T) {
	var runs atomic.Int32
	stop := make(chan struct{})
	done := make(chan struct{})

	go func() {
		RunOnSchedule(logrus.WithField("test", t.Name()), intervalSchedule(10*time.Millisecond), false, 0, stop, func() error {
			runs.Add(1)
			return errors.New("collector failed")
		})
		close(done)
	}()

	require.Eventually(t, func() bool { return runs.Load() >= 3 }, 5*time.Second, 5*time.Millisecond)
	close(stop)
	require.Eventually(t, func() bool {
		select {
		case <-done:
			return true
		default:
			return false
		}
	}, 5*time.Second, 5*time.Millisecond)
}

func TestRunOnSchedule_RunStartup(t *testing.T) {
	for _, runStartup := range []bool{true, false} {
		var runs atomic.Int32
		stop := make(chan struct{})
		done := make(chan struct{})

		go func() {
			// the schedule never fires during the test, so any run comes from runStartup
			RunOnSchedule(logrus.WithField("test", t.Name()), intervalSchedule(time.Hour), runStartup, 10*time.Millisecond, stop, func() error {
				runs.Add(1)
				return nil
			})
			close(done)
		}()

		time.Sleep(200 * time.Millisecond)
		close(stop)
		<-done

		if runStartup {
			require.Equal(t, int32(1), runs.Load())
		} else {
			require.Equal(t, int32(0), runs.Load())
		}
	}
}
