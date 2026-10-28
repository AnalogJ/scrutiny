package collector

import (
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

// ParseCronSchedule parses a standard 5 field cron expression (eg. "0 0 * * *") or a descriptor such as
// "@daily" or "@every 6h". Surrounding quotes are stripped, matching the docker images, since
// COLLECTOR_CRON_SCHEDULE is often quoted twice in docker-compose files.
func ParseCronSchedule(spec string) (cron.Schedule, error) {
	spec = strings.TrimSpace(spec)
	if len(spec) >= 2 && (spec[0] == '"' || spec[0] == '\'') && spec[len(spec)-1] == spec[0] {
		spec = spec[1 : len(spec)-1]
	}
	return cron.ParseStandard(spec)
}

// RunOnSchedule calls run every time the schedule fires, and only returns once stop is closed.
// Runs never overlap, and a failed run is logged without stopping the schedule, like cron.
// If runStartup is true, run is also called once after waiting startupSleep.
func RunOnSchedule(logger *logrus.Entry, schedule cron.Schedule, runStartup bool, startupSleep time.Duration, stop <-chan struct{}, run func() error) {
	runAndLog := func() {
		if err := run(); err != nil {
			logger.Errorf("Collector run failed: %v", err)
		}
	}

	if runStartup {
		select {
		case <-stop:
			return
		case <-time.After(startupSleep):
		}
		logger.Infoln("Running collector on startup")
		runAndLog()
	}

	for {
		next := schedule.Next(time.Now())
		logger.Infof("Next collector run scheduled for %s", next.Format(time.RFC3339))

		select {
		case <-stop:
			return
		case <-time.After(time.Until(next)):
		}
		runAndLog()
	}
}
