package detect

import (
	"github.com/analogj/scrutiny/collector/pkg/models"
	"strings"
)

func DevicePrefix() string {
	return ""
}

func (d *Detect) Start() ([]models.Device, error) {
	// call the base/common functionality to get a list of devices
	detectedDevices, err := d.SmartctlScan()
	if err != nil {
		return nil, err
	}

	//inflate device info for detected devices.
	for ndx, _ := range detectedDevices {
		d.SmartCtlInfo(&detectedDevices[ndx]) //ignore errors.
	}

	return detectedDevices, nil
}

//WWN values NVMe and SCSI
func (d *Detect) wwnFallback(detectedDevice *models.Device) {

	//fallback to serial number
	if len(detectedDevice.WWN) == 0 {
		detectedDevice.WWN = detectedDevice.SerialNumber
	}

	//wwn must always be lowercase.
	detectedDevice.WWN = strings.ToLower(detectedDevice.WWN)
}
