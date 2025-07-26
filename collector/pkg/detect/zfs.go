package detect

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/analogj/scrutiny/collector/pkg/common/shell"
	"github.com/analogj/scrutiny/collector/pkg/config"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models"
	"github.com/sirupsen/logrus"
)

type ZfsDetect struct {
	Logger *logrus.Entry
	Config config.Interface
	Shell  shell.Interface
}

// ZfsStatus represents the JSON structure returned by `zpool status --json`
type ZfsStatus struct {
	OutputVersion struct {
		Command   string `json:"command"`
		VersMajor int    `json:"vers_major"`
		VersMinor int    `json:"vers_minor"`
	} `json:"output_version"`
	Pools map[string]ZfsPoolStatus `json:"pools"`
}

type ZfsPoolStatus struct {
	Name       string `json:"name"`
	State      string `json:"state"`
	PoolGuid   string `json:"pool_guid"`
	Txg        string `json:"txg"`
	SpaVersion string `json:"spa_version"`
	ZplVersion string `json:"zpl_version"`
	Status     string `json:"status,omitempty"`
	Action     string `json:"action,omitempty"`
	ErrorCount string `json:"error_count"`
	ScanStats  *struct {
		Function          string `json:"function"`
		State             string `json:"state"`
		StartTime         string `json:"start_time"`
		EndTime           string `json:"end_time"`
		ToExamine         string `json:"to_examine"`
		Examined          string `json:"examined"`
		Skipped           string `json:"skipped"`
		Processed         string `json:"processed"`
		Errors            string `json:"errors"`
		BytesPerScan      string `json:"bytes_per_scan"`
		PassStart         string `json:"pass_start"`
		ScrubPause        string `json:"scrub_pause"`
		ScrubSpentPaused  string `json:"scrub_spent_paused"`
		IssuedBytesPerScan string `json:"issued_bytes_per_scan"`
		Issued            string `json:"issued"`
	} `json:"scan_stats,omitempty"`
	Vdevs map[string]ZfsVdevStatus `json:"vdevs"`
}

type ZfsVdevStatus struct {
	Name           string                    `json:"name"`
	VdevType       string                    `json:"vdev_type"`
	Guid           string                    `json:"guid"`
	Class          string                    `json:"class"`
	State          string                    `json:"state"`
	AllocSpace     string                    `json:"alloc_space,omitempty"`
	TotalSpace     string                    `json:"total_space,omitempty"`
	DefSpace       string                    `json:"def_space,omitempty"`
	RepDevSize     string                    `json:"rep_dev_size,omitempty"`
	PhysSpace      string                    `json:"phys_space,omitempty"`
	ReadErrors     string                    `json:"read_errors"`
	WriteErrors    string                    `json:"write_errors"`
	ChecksumErrors string                    `json:"checksum_errors"`
	SlowIos        string                    `json:"slow_ios,omitempty"`
	Path           string                    `json:"path,omitempty"`
	PhysPath       string                    `json:"phys_path,omitempty"`
	DevId          string                    `json:"devid,omitempty"`
	ScanProcessed  string                    `json:"scan_processed,omitempty"`
	Vdevs          map[string]ZfsVdevStatus  `json:"vdevs,omitempty"`
}

// DetectZfsPools scans for ZFS pools using `zpool status --json`
func (z *ZfsDetect) DetectZfsPools() ([]models.ZfsPool, error) {
	// Check if ZFS is available and enabled in config
	if !z.Config.GetBool("zfs.enabled") {
		z.Logger.Debug("ZFS monitoring is disabled in configuration")
		return []models.ZfsPool{}, nil
	}

	// Log container environment information for debugging
	z.logEnvironmentInfo()

	// Try to execute zpool status --json
	zpoolBin := z.Config.GetString("commands.zpool_bin")
	if zpoolBin == "" {
		zpoolBin = "zpool"
	}

	args := strings.Split("status --json", " ")
	zpoolOutput, err := z.Shell.Command(z.Logger, zpoolBin, args, "", os.Environ())
	if err != nil {
		z.Logger.Debugf("Error running zpool status: %v", err)
		
		// Check if it's a hostid mismatch issue
		if strings.Contains(err.Error(), "hostid") || strings.Contains(zpoolOutput, "hostid") {
			return nil, fmt.Errorf("ZFS hostid mismatch detected. To resolve this issue:\n" +
				"1. Mount /etc/hostid as read-only in the container: -v /etc/hostid:/etc/hostid:ro\n" +
				"2. For Docker Compose, add to volumes: - /etc/hostid:/etc/hostid:ro\n" +
				"3. For Kubernetes, mount it as a hostPath volume\n" +
				"Original error: %v", err)
		}
		
		return []models.ZfsPool{}, nil // Return empty slice for other errors
	}

	// Log raw JSON output for debugging
	z.Logger.Debugf("Raw zpool status output: %s", zpoolOutput)

	// Check for hostid mismatch warnings in the output
	if strings.Contains(zpoolOutput, "hostid") && strings.Contains(zpoolOutput, "mismatch") {
		return nil, fmt.Errorf("ZFS hostid mismatch detected in pool status. To resolve this issue:\n" +
			"1. Mount /etc/hostid as read-only in the container: -v /etc/hostid:/etc/hostid:ro\n" +
			"2. For Docker Compose, add to volumes: - /etc/hostid:/etc/hostid:ro\n" +
			"3. For Kubernetes, mount it as a hostPath volume\n" +
			"This ensures the container uses the same hostid as your ZFS host system.")
	}

	// Basic JSON validation
	if !json.Valid([]byte(zpoolOutput)) {
		z.Logger.Errorf("Invalid JSON received from zpool command")
		return nil, fmt.Errorf("invalid JSON from zpool status")
	}

	var zfsStatus ZfsStatus
	err = json.Unmarshal([]byte(zpoolOutput), &zfsStatus)
	if err != nil {
		z.Logger.Errorf("Error parsing zpool status JSON: %v", err)
		z.Logger.Errorf("Raw JSON that failed to parse: %s", zpoolOutput)
		
		// Try to identify the specific issue
		if strings.Contains(err.Error(), "unknown field") {
			z.Logger.Errorf("JSON structure mismatch - the zpool output contains fields not defined in our struct")
		}
		if strings.Contains(err.Error(), "cannot unmarshal") {
			z.Logger.Errorf("JSON type mismatch - field types don't match expected Go types")
		}
		
		return nil, err
	}

	// Validate required fields are present
	if zfsStatus.Pools == nil {
		z.Logger.Errorf("No 'pools' field found in zpool JSON output")
		return nil, fmt.Errorf("missing pools field in zpool output")
	}

	z.Logger.Debugf("Successfully parsed %d ZFS pools from JSON", len(zfsStatus.Pools))

	var pools []models.ZfsPool
	hostId := z.Config.GetString("host.id")

	for poolName, poolStatus := range zfsStatus.Pools {
		z.Logger.Debugf("Processing pool '%s': state=%s, status='%s', action='%s'", 
			poolName, poolStatus.State, poolStatus.Status, poolStatus.Action)

		pool := models.ZfsPool{
			PoolGuid:   poolStatus.PoolGuid,
			Name:       poolStatus.Name,
			HostId:     hostId,
			State:      poolStatus.State,
			Txg:        poolStatus.Txg,
			SpaVersion: poolStatus.SpaVersion,
			ZplVersion: poolStatus.ZplVersion,
			Status:     strings.TrimSpace(poolStatus.Status),
			Action:     strings.TrimSpace(poolStatus.Action),
			ErrorCount: poolStatus.ErrorCount,
		}

		// Set default values for empty status/action fields on healthy pools
		if pool.Status == "" && pool.State == "ONLINE" {
			pool.Status = "Pool is healthy"
			z.Logger.Debugf("Set default status for healthy pool '%s'", poolName)
		}
		if pool.Action == "" && pool.State == "ONLINE" {
			pool.Action = "No action required"
			z.Logger.Debugf("Set default action for healthy pool '%s'", poolName)
		}

		// Add scan information if available
		if poolStatus.ScanStats != nil {
			z.Logger.Debugf("Found scan stats for pool '%s': function=%s, state=%s", 
				poolName, poolStatus.ScanStats.Function, poolStatus.ScanStats.State)
			pool.ScanFunction = poolStatus.ScanStats.Function
			pool.ScanState = poolStatus.ScanStats.State
			pool.ScanStartTime = poolStatus.ScanStats.StartTime
			pool.ScanEndTime = poolStatus.ScanStats.EndTime
			pool.ScanToExamine = poolStatus.ScanStats.ToExamine
			pool.ScanExamined = poolStatus.ScanStats.Examined
			pool.ScanProcessed = poolStatus.ScanStats.Processed
			pool.ScanErrors = poolStatus.ScanStats.Errors
			pool.ScanIssued = poolStatus.ScanStats.Issued
		} else {
			z.Logger.Debugf("No scan stats available for pool '%s'", poolName)
		}

		// Process vdevs (typically there's a root vdev with the pool name)
		z.Logger.Debugf("Processing %d vdevs for pool '%s'", len(poolStatus.Vdevs), poolName)
		for vdevName, vdevStatus := range poolStatus.Vdevs {
			z.Logger.Debugf("Processing vdev '%s' of type '%s' for pool '%s'", 
				vdevName, vdevStatus.VdevType, poolName)
			if vdevStatus.VdevType == "root" {
				// Copy root vdev information to pool
				pool.AllocSpace = vdevStatus.AllocSpace
				pool.TotalSpace = vdevStatus.TotalSpace
				pool.DefSpace = vdevStatus.DefSpace
				pool.ReadErrors = vdevStatus.ReadErrors
				pool.WriteErrors = vdevStatus.WriteErrors
				pool.ChecksumErrors = vdevStatus.ChecksumErrors

				// Process child vdevs
				z.Logger.Debugf("Processing %d child vdevs for root vdev of pool '%s'", 
					len(vdevStatus.Vdevs), poolName)
				pool.Vdevs = z.processVdevs(vdevStatus.Vdevs, poolStatus.PoolGuid, nil)
			}
		}

		pools = append(pools, pool)
		z.Logger.Infof("Detected ZFS pool: %s (%s)", pool.Name, pool.State)
	}

	return pools, nil
}

// processVdevs recursively processes vdev hierarchy
func (z *ZfsDetect) processVdevs(vdevs map[string]ZfsVdevStatus, poolGuid string, parentId *uint) []models.ZfsVdev {
	var result []models.ZfsVdev

	for _, vdevStatus := range vdevs {
		vdev := models.ZfsVdev{
			PoolGuid:       poolGuid,
			ParentId:       parentId,
			Guid:           vdevStatus.Guid,
			Name:           vdevStatus.Name,
			VdevType:       vdevStatus.VdevType,
			Class:          vdevStatus.Class,
			State:          vdevStatus.State,
			AllocSpace:     vdevStatus.AllocSpace,
			TotalSpace:     vdevStatus.TotalSpace,
			DefSpace:       vdevStatus.DefSpace,
			RepDevSize:     vdevStatus.RepDevSize,
			PhysSpace:      vdevStatus.PhysSpace,
			ReadErrors:     vdevStatus.ReadErrors,
			WriteErrors:    vdevStatus.WriteErrors,
			ChecksumErrors: vdevStatus.ChecksumErrors,
			SlowIos:        vdevStatus.SlowIos,
			Path:           vdevStatus.Path,
			PhysPath:       vdevStatus.PhysPath,
			DevId:          vdevStatus.DevId,
			ScanProcessed:  vdevStatus.ScanProcessed,
		}

		result = append(result, vdev)

		// Process child vdevs recursively
		if len(vdevStatus.Vdevs) > 0 {
			// Note: In a real implementation, we'd need to save the parent first to get its ID
			// For now, we'll process children with nil parent ID and handle hierarchy later
			childVdevs := z.processVdevs(vdevStatus.Vdevs, poolGuid, nil)
			result = append(result, childVdevs...)
		}
	}

	return result
}

// logEnvironmentInfo logs container and system information for debugging
func (z *ZfsDetect) logEnvironmentInfo() {
	// Check if running in container
	if _, err := os.Stat("/.dockerenv"); err == nil {
		z.Logger.Debug("Running inside Docker container")
	}

	// Check for container runtime environment variables
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		z.Logger.Debug("Running in Kubernetes environment")
	}

	// Log ZFS-related kernel modules
	if moduleInfo, err := z.Shell.Command(z.Logger, "lsmod", []string{}, "", os.Environ()); err == nil {
		if strings.Contains(moduleInfo, "zfs") {
			z.Logger.Debug("ZFS kernel module is loaded")
		} else {
			z.Logger.Warn("ZFS kernel module not found in lsmod output")
		}
	}

	// Check if ZFS filesystem is available
	if _, err := os.Stat("/sys/module/zfs"); err == nil {
		z.Logger.Debug("ZFS sysfs entries found at /sys/module/zfs")
	} else {
		z.Logger.Debug("No ZFS sysfs entries found")
	}

	// Log host ID for debugging
	hostId := z.Config.GetString("host.id")
	z.Logger.Debugf("Using host ID: %s", hostId)
	
	// Check system hostid for ZFS troubleshooting
	if systemHostId, err := z.Shell.Command(z.Logger, "hostid", []string{}, "", os.Environ()); err == nil {
		z.Logger.Debugf("System hostid: %s", strings.TrimSpace(systemHostId))
		
		// Check if /etc/hostid exists
		if _, err := os.Stat("/etc/hostid"); err == nil {
			z.Logger.Debug("Found /etc/hostid file")
		} else {
			z.Logger.Debug("No /etc/hostid file found")
		}
	}
	
	// Check for ZFS_HOSTID_CONFIGURED environment variable
	if os.Getenv("ZFS_HOSTID_CONFIGURED") == "true" {
		z.Logger.Debug("ZFS hostid has been configured by container initialization")
	}
}

// IsZfsAvailable checks if ZFS tools are available on the system
func (z *ZfsDetect) IsZfsAvailable() bool {
	zpoolBin := z.Config.GetString("commands.zpool_bin")
	if zpoolBin == "" {
		zpoolBin = "zpool"
	}

	_, err := z.Shell.Command(z.Logger, "which", []string{zpoolBin}, "", os.Environ())
	return err == nil
}