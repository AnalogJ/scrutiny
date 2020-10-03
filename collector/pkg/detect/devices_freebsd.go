package detect

import (
	"github.com/analogj/scrutiny/collector/pkg/models"
	"github.com/jaypipes/ghw"
	"strings"
)

func DevicePrefix() string {
	return "/dev/"
}

func (d *Detect) Start() ([]models.Device, error) {
	// call the base/common functionality to get a list of devices
	detectedDevices, err := d.smartctlScan()
	if err != nil {
		return nil, err
	}

	//inflate device info for detected devices.
	for ndx, _ := range detectedDevices {
		d.smartCtlInfo(&detectedDevices[ndx]) //ignore errors.
	}

	return detectedDevices, nil
}

//WWN values NVMe and SCSI
func (d *Detect) wwnFallback(detectedDevice *models.Device) {
	block, err := ghw.Block()
	if err == nil {
		for _, disk := range block.Disks {
			if disk.Name == detectedDevice.DeviceName {
				d.Logger.Debugf("Found matching block device. WWN: %s", disk.WWN)
				detectedDevice.WWN = disk.WWN
				break
			}
		}
	}

	//no WWN found, or could not open Block devices. Either way, fallback to serial number
	if len(detectedDevice.WWN) == 0 {
		d.Logger.Debugf("WWN is empty, falling back to serial number: %s", detectedDevice.SerialNumber)
		detectedDevice.WWN = detectedDevice.SerialNumber
	}

	//wwn must always be lowercase.
	detectedDevice.WWN = strings.ToLower(detectedDevice.WWN)
}
