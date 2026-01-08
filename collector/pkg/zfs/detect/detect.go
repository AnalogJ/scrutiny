package detect

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/analogj/scrutiny/collector/pkg/config"
	"github.com/analogj/scrutiny/collector/pkg/zfs/models"
	"github.com/sirupsen/logrus"
)

// Detect handles ZFS pool detection
type Detect struct {
	Logger *logrus.Entry
	Config config.Interface
}

// Start detects all ZFS pools on the system
func (d *Detect) Start() ([]models.ZFSPool, error) {
	// Check if zpool command exists
	zpoolPath, err := exec.LookPath("zpool")
	if err != nil {
		d.Logger.Warnf("zpool command not found: %v", err)
		return nil, fmt.Errorf("zpool command not found: %w", err)
	}
	d.Logger.Debugf("Found zpool at: %s", zpoolPath)

	// List all pools with basic properties
	pools, err := d.listPools()
	if err != nil {
		return nil, err
	}

	// Get detailed status for each pool (vdevs, scrub, errors)
	for i := range pools {
		if err := d.getPoolStatus(&pools[i]); err != nil {
			d.Logger.Warnf("Failed to get status for pool %s: %v", pools[i].Name, err)
		}
	}

	return pools, nil
}

// listPools lists all ZFS pools with their properties
func (d *Detect) listPools() ([]models.ZFSPool, error) {
	// zpool list -H -p -o name,guid,size,alloc,free,frag,cap,health,ashift
	cmd := exec.Command("zpool", "list", "-H", "-p", "-o",
		"name,guid,size,alloc,free,frag,cap,health,ashift")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list pools: %w", err)
	}

	var pools []models.ZFSPool
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		fields := strings.Split(line, "\t")
		if len(fields) < 9 {
			d.Logger.Warnf("Unexpected zpool list output: %s", line)
			continue
		}

		size, _ := strconv.ParseInt(fields[2], 10, 64)
		alloc, _ := strconv.ParseInt(fields[3], 10, 64)
		free, _ := strconv.ParseInt(fields[4], 10, 64)
		frag, _ := strconv.Atoi(fields[5])
		cap, _ := strconv.Atoi(fields[6])
		ashift, _ := strconv.Atoi(fields[8])

		pool := models.ZFSPool{
			Name:            fields[0],
			GUID:            fields[1],
			Size:            size,
			Allocated:       alloc,
			Free:            free,
			Fragmentation:   frag,
			CapacityPercent: float64(cap),
			Health:          fields[7],
			Status:          models.ZFSPoolStatus(fields[7]),
			Ashift:          ashift,
			ScrubState:      models.ZFSScrubStateNone,
		}

		// Set HostID if configured
		if d.Config != nil && d.Config.IsSet("host.id") {
			pool.HostID = d.Config.GetString("host.id")
		}

		pools = append(pools, pool)
	}

	return pools, nil
}

// getPoolStatus gets detailed status for a pool including vdevs and scrub
func (d *Detect) getPoolStatus(pool *models.ZFSPool) error {
	// zpool status -p <poolname>
	cmd := exec.Command("zpool", "status", "-p", pool.Name)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get pool status: %w", err)
	}

	statusStr := string(output)

	// Parse vdev tree
	pool.Vdevs = d.parseVdevTree(statusStr, pool.Name)

	// Calculate total errors from vdevs
	d.calculateTotalErrors(pool)

	// Parse scrub status
	d.parseScrubStatus(pool, statusStr)

	return nil
}

// parseVdevTree parses the vdev configuration from zpool status output
func (d *Detect) parseVdevTree(output string, poolName string) []models.ZFSVdev {
	var vdevs []models.ZFSVdev
	var currentParent *models.ZFSVdev
	var inConfig bool
	var baseIndent int

	// Pattern to match vdev lines with status and errors
	// Example: "  mirror-0  ONLINE       0     0     0"
	vdevPattern := regexp.MustCompile(`^(\s*)(\S+)\s+(ONLINE|DEGRADED|FAULTED|OFFLINE|REMOVED|UNAVAIL|AVAIL|INUSE)\s+(\d+)\s+(\d+)\s+(\d+)`)

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()

		// Look for config section
		if strings.Contains(line, "config:") {
			inConfig = true
			continue
		}

		// Skip header lines
		if strings.Contains(line, "NAME") && strings.Contains(line, "STATE") {
			continue
		}

		// End of config section
		if inConfig && (strings.Contains(line, "errors:") || strings.Contains(line, "scan:")) {
			break
		}

		if !inConfig {
			continue
		}

		matches := vdevPattern.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		indent := len(matches[1])
		name := matches[2]
		status := matches[3]
		readErr, _ := strconv.ParseInt(matches[4], 10, 64)
		writeErr, _ := strconv.ParseInt(matches[5], 10, 64)
		ckErr, _ := strconv.ParseInt(matches[6], 10, 64)

		// Skip the pool name line itself
		if name == poolName {
			baseIndent = indent
			continue
		}

		vdev := models.ZFSVdev{
			Name:           name,
			Status:         models.ZFSPoolStatus(status),
			ReadErrors:     readErr,
			WriteErrors:    writeErr,
			ChecksumErrors: ckErr,
			Type:           d.detectVdevType(name),
		}

		// Detect if this is a device path
		if strings.HasPrefix(name, "/dev/") || strings.Contains(name, "sd") ||
			strings.Contains(name, "nvme") || strings.Contains(name, "ada") ||
			strings.Contains(name, "da") || strings.Contains(name, "disk") {
			vdev.Type = models.ZFSVdevTypeDisk
			vdev.Path = d.resolveDevicePath(name)
		}

		// Determine hierarchy based on indentation
		if indent == baseIndent+2 {
			// Top-level vdev (mirror, raidz, disk directly under pool)
			vdevs = append(vdevs, vdev)
			currentParent = &vdevs[len(vdevs)-1]
		} else if indent > baseIndent+2 && currentParent != nil {
			// Child of current parent
			currentParent.Children = append(currentParent.Children, vdev)
		}
	}

	return vdevs
}

// detectVdevType determines the vdev type from its name
func (d *Detect) detectVdevType(name string) models.ZFSVdevType {
	nameLower := strings.ToLower(name)

	switch {
	case strings.HasPrefix(nameLower, "mirror"):
		return models.ZFSVdevTypeMirror
	case strings.HasPrefix(nameLower, "raidz3"):
		return models.ZFSVdevTypeRaidz3
	case strings.HasPrefix(nameLower, "raidz2"):
		return models.ZFSVdevTypeRaidz2
	case strings.HasPrefix(nameLower, "raidz1"), strings.HasPrefix(nameLower, "raidz"):
		return models.ZFSVdevTypeRaidz1
	case strings.HasPrefix(nameLower, "draid3"):
		return models.ZFSVdevTypeDraid3
	case strings.HasPrefix(nameLower, "draid2"):
		return models.ZFSVdevTypeDraid2
	case strings.HasPrefix(nameLower, "draid1"), strings.HasPrefix(nameLower, "draid"):
		return models.ZFSVdevTypeDraid1
	case nameLower == "spare" || nameLower == "spares":
		return models.ZFSVdevTypeSpare
	case nameLower == "log" || nameLower == "logs":
		return models.ZFSVdevTypeLog
	case nameLower == "cache":
		return models.ZFSVdevTypeCache
	case nameLower == "special":
		return models.ZFSVdevTypeSpecial
	case nameLower == "dedup":
		return models.ZFSVdevTypeDedup
	default:
		return models.ZFSVdevTypeDisk
	}
}

// resolveDevicePath resolves a device name to its full path
func (d *Detect) resolveDevicePath(name string) string {
	if strings.HasPrefix(name, "/dev/") {
		return name
	}
	// Common Linux device names
	if strings.HasPrefix(name, "sd") || strings.HasPrefix(name, "nvme") ||
		strings.HasPrefix(name, "hd") || strings.HasPrefix(name, "vd") {
		return "/dev/" + name
	}
	// FreeBSD device names
	if strings.HasPrefix(name, "ada") || strings.HasPrefix(name, "da") ||
		strings.HasPrefix(name, "nda") {
		return "/dev/" + name
	}
	return name
}

// calculateTotalErrors calculates total errors from all vdevs
func (d *Detect) calculateTotalErrors(pool *models.ZFSPool) {
	pool.TotalReadErrors = 0
	pool.TotalWriteErrors = 0
	pool.TotalChecksumErrors = 0

	var addErrors func(vdevs []models.ZFSVdev)
	addErrors = func(vdevs []models.ZFSVdev) {
		for _, vdev := range vdevs {
			pool.TotalReadErrors += vdev.ReadErrors
			pool.TotalWriteErrors += vdev.WriteErrors
			pool.TotalChecksumErrors += vdev.ChecksumErrors
			if len(vdev.Children) > 0 {
				addErrors(vdev.Children)
			}
		}
	}

	addErrors(pool.Vdevs)
}

// parseScrubStatus parses scrub information from zpool status output
func (d *Detect) parseScrubStatus(pool *models.ZFSPool, output string) {
	// Look for scan: line
	// Examples:
	// scan: scrub repaired 0B in 00:10:30 with 0 errors on Sun Jan  5 00:34:31 2026
	// scan: scrub in progress since Sun Jan  5 00:24:01 2026
	// scan: scrub canceled on Sun Jan  5 00:30:00 2026
	// scan: none requested

	scrubInProgress := regexp.MustCompile(`scan:\s+scrub in progress since (.+)`)
	scrubFinished := regexp.MustCompile(`scan:\s+scrub repaired \S+ in (\S+) with (\d+) errors on (.+)`)
	scrubCanceled := regexp.MustCompile(`scan:\s+scrub canceled on (.+)`)
	scrubProgress := regexp.MustCompile(`(\d+(?:\.\d+)?)\s*%\s+done`)

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()

		if matches := scrubInProgress.FindStringSubmatch(line); matches != nil {
			pool.ScrubState = models.ZFSScrubStateScanning
			if t, err := d.parseZFSDate(matches[1]); err == nil {
				pool.ScrubStartTime = &t
			}
			continue
		}

		if matches := scrubFinished.FindStringSubmatch(line); matches != nil {
			pool.ScrubState = models.ZFSScrubStateFinished
			pool.ScrubErrorsCount, _ = strconv.ParseInt(matches[2], 10, 64)
			if t, err := d.parseZFSDate(matches[3]); err == nil {
				pool.ScrubEndTime = &t
			}
			pool.ScrubPercentComplete = 100.0
			continue
		}

		if matches := scrubCanceled.FindStringSubmatch(line); matches != nil {
			pool.ScrubState = models.ZFSScrubStateCanceled
			if t, err := d.parseZFSDate(matches[1]); err == nil {
				pool.ScrubEndTime = &t
			}
			continue
		}

		// Parse progress percentage if scrubbing
		if pool.ScrubState == models.ZFSScrubStateScanning {
			if matches := scrubProgress.FindStringSubmatch(line); matches != nil {
				pool.ScrubPercentComplete, _ = strconv.ParseFloat(matches[1], 64)
			}
		}
	}
}

// parseZFSDate parses a date string from ZFS output
func (d *Detect) parseZFSDate(dateStr string) (time.Time, error) {
	// ZFS uses format like: "Sun Jan  5 00:34:31 2026"
	layouts := []string{
		"Mon Jan  2 15:04:05 2006",
		"Mon Jan 2 15:04:05 2006",
		time.ANSIC,
		time.UnixDate,
	}

	dateStr = strings.TrimSpace(dateStr)

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}
