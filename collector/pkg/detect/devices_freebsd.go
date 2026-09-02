package detect

import (
	"strings"

	"github.com/analogj/scrutiny/collector/pkg/common/shell"
	"github.com/analogj/scrutiny/collector/pkg/models"
	"github.com/jaypipes/ghw"
)

func DevicePrefix() string {
	return "/dev/"
}

func (d *Detect) Scan() ([]models.Device, error) {
	d.Shell = shell.Create()
	// call the base/common functionality to get a list of devices
	return d.SmartctlScan()
}

func (d *Detect) Info(detectedDevices []models.Device) ([]models.Device, []DeviceInfoError) {
	//inflate device info for detected devices.
	infoErrors := []DeviceInfoError{}
	for ndx := range detectedDevices {
		if err := d.SmartCtlInfo(&detectedDevices[ndx]); err != nil {
			infoErrors = append(infoErrors, DeviceInfoError{Device: detectedDevices[ndx], Err: err})
		}
	}

	return detectedDevices, infoErrors
}

// WWN values NVMe and SCSI
func (d *Detect) wwnFallback(detectedDevice *models.Device) {
	block, err := ghw.Block()
	if err == nil {
		for _, disk := range block.Disks {
			if disk.Name == detectedDevice.DeviceName && strings.ToLower(disk.WWN) != "unknown" {
				d.Logger.Debugf("Found matching block device. WWN: %s", disk.WWN)
				detectedDevice.WWN = disk.WWN
				break
			}
		}
	}

	//wwn must always be lowercase.
	detectedDevice.WWN = strings.ToLower(detectedDevice.WWN)
}
